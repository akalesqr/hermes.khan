package handler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"khan/config"
	"khan/internal/models"
	"khan/internal/repository"
	"khan/internal/service"
)

// FileHandler serves /api/files/*
type FileHandler struct {
	files *repository.FileRepo
	rooms *repository.RoomRepo
	cfg   *config.Config
	crypto *service.CryptoService
}

func NewFileHandler(files *repository.FileRepo, rooms *repository.RoomRepo, cfg *config.Config, crypto *service.CryptoService) *FileHandler {
	return &FileHandler{files: files, rooms: rooms, cfg: cfg, crypto: crypto}
}

// Upload handles multipart file upload (member only)
func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	roomID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}

	// membership check
	isMember, err := h.rooms.IsMember(roomID, u.ID)
	if err != nil || !isMember {
		writeErr(w, http.StatusForbidden, "شما عضو این اتاق نیستید")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(h.cfg.Server.MaxUploadMB)*1024*1024)
	if err := r.ParseMultipartForm(int64(h.cfg.Server.MaxUploadMB) * 1024 * 1024); err != nil {
		writeErr(w, http.StatusRequestEntityTooLarge, "حجم فایل بیشتر از حد مجاز است")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "فایل ارسال نشده است")
		return
	}
	defer file.Close()

	// sanitize filename
	name := filepath.Base(header.Filename)
	name = strings.Map(func(rn rune) rune {
		if rn == '/' || rn == '\\' || rn == 0 {
			return '_'
		}
		return rn
	}, name)

	// encrypt content at rest
	data, err := io.ReadAll(file)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای خواندن فایل")
		return
	}
	enc, err := h.crypto.Encrypt(data)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای رمزنگاری")
		return
	}

	// store under data/uploads/<yyyy-mm>/<id>.enc
	dataDir := h.cfg.Server.DataDir
	monthDir := filepath.Join(dataDir, "uploads", time.Now().Format("2006-01"))
	if err := os.MkdirAll(monthDir, 0o755); err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای ایجاد پوشه")
		return
	}

	f := &models.File{
		OwnerID:   u.ID,
		RoomID:    roomID,
		FileName:  name,
		StoredAs:  "", // set after id
		Size:      int64(len(data)),
		MimeType:  header.Header.Get("Content-Type"),
		CreatedAt: time.Now(),
	}
	id, err := h.files.Create(f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای ذخیره فایل")
		return
	}
	stored := filepath.Join(monthDir, time.Now().Format("150405")+"-"+itoa(id)+".enc")
	if err := os.WriteFile(stored, enc, 0o600); err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای ذخیره فایل")
		return
	}
	f.ID = id
	f.StoredAs = stored
	f.Size = int64(len(data))
	f.MimeType = header.Header.Get("Content-Type")

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id": id, "file_name": name, "size": len(data), "mime_type": f.MimeType,
	})
}

// Download serves a decrypted file (member only)
func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	fileID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}

	can, err := h.files.CanAccess(fileID, u.ID)
	if err != nil || !can {
		writeErr(w, http.StatusForbidden, "دسترسی غیرمجاز")
		return
	}

	f, err := h.files.GetByID(fileID)
	if err != nil || f == nil {
		writeErr(w, http.StatusNotFound, "فایل یافت نشد")
		return
	}

	enc, err := os.ReadFile(f.StoredAs)
	if err != nil {
		writeErr(w, http.StatusNotFound, "فایل یافت نشد")
		return
	}
	data, err := h.crypto.Decrypt(enc)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای رمزگشایی")
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+f.FileName+"\"")
	w.Header().Set("Content-Type", f.MimeType)
	w.Header().Set("Content-Length", itoa(int64(len(data))))
	_, _ = w.Write(data)
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}
