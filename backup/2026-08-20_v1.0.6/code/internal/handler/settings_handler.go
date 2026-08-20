package handler

import (
	"archive/zip"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"khan/config"
	"khan/internal/database"
	"khan/internal/service"
)

// SettingsHandler serves /api/settings/* (admin panel + public server info)
type SettingsHandler struct {
	cfg      *config.Config
	auth     *service.AuthService
	storeRef *database.Store
}

func NewSettingsHandler(cfg *config.Config, auth *service.AuthService) *SettingsHandler {
	return &SettingsHandler{cfg: cfg, auth: auth}
}

// SetStore injects the datastore for backup operations
func (h *SettingsHandler) SetStore(s *database.Store) { h.storeRef = s }

// Info returns public server info (no auth needed — used by login page)
func (h *SettingsHandler) Info(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":         "خان",
		"version": "1.0.6",
		"address_type": h.cfg.Server.AddressType,
		"ip":           h.cfg.Server.IP,
		"dns":          h.cfg.Server.DNS,
		"port":         h.cfg.Server.Port,
	})
}

// LicenseStatus returns license state for admin panel
func (h *SettingsHandler) LicenseStatus(w http.ResponseWriter, r *http.Request) {
	state, to, errMsg, maxUsers, expiry := service.LicenseState()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"state":      state,
		"licensed_to": to,
		"error":      errMsg,
		"max_users":  maxUsers,
		"expiry":     expiry.Format(time.RFC3339),
	})
}

// LicenseApply uploads a license.key file and verifies it
func (h *SettingsHandler) LicenseApply(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	file, _, err := r.FormFile("license")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "فایل لایسنس انتخاب نشده است")
		return
	}
	defer file.Close()

	data := make([]byte, 1<<20)
	n, err := file.Read(data)
	if err != nil && err.Error() != "EOF" {
		writeErr(w, http.StatusBadRequest, "خطای خواندن فایل")
		return
	}

	// save to data dir then verify
	path := filepath.Join(h.cfg.Server.DataDir, "license.key")
	if err := os.WriteFile(path, data[:n], 0o644); err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای ذخیره فایل")
		return
	}

	if err := service.LoadLicense(path); err != nil {
		// tampered → penalty applied, still respond with state
		state, to, errMsg, maxUsers, expiry := service.LicenseState()
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"applied": false, "state": state, "licensed_to": to,
			"error": errMsg, "max_users": maxUsers, "expiry": expiry.Format(time.RFC3339),
		})
		return
	}

	state, to, _, maxUsers, expiry := service.LicenseState()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"applied": true, "state": state, "licensed_to": to,
		"max_users": maxUsers, "expiry": expiry.Format(time.RFC3339),
	})
}

// LicenseRemove deletes the license → back to 20 free users
func (h *SettingsHandler) LicenseRemove(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.cfg.Server.DataDir, "license.key")
	if err := service.RemoveLicense(path); err != nil {
		writeErr(w, http.StatusInternalServerError, "خطا در حذف لایسنس")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"removed": true, "state": "free", "max_users": 20,
	})
}

// NetworkInfo returns current server network config
func (h *SettingsHandler) NetworkInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"address_type": h.cfg.Server.AddressType,
		"ip":           h.cfg.Server.IP,
		"dns":          h.cfg.Server.DNS,
		"port":         h.cfg.Server.Port,
	})
}

// NetworkUpdate lets admin change IP/DNS/port (port limited to 1700-1799)
func (h *SettingsHandler) NetworkUpdate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AddressType string `json:"address_type"`
		IP          string `json:"ip"`
		DNS         string `json:"dns"`
		Port        int    `json:"port"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}

	if req.AddressType != "ip" && req.AddressType != "dns" {
		writeErr(w, http.StatusBadRequest, "نوع آدرس باید ip یا dns باشد")
		return
	}
	if req.Port < 1700 || req.Port > 1799 {
		writeErr(w, http.StatusBadRequest, "پورت باید بین 1700 تا 1799 باشد (دو رقم آخر قابل تغییر)")
		return
	}
	if req.AddressType == "dns" && req.DNS == "" {
		writeErr(w, http.StatusBadRequest, "آدرس DNS الزامی است")
		return
	}
	if req.AddressType == "ip" && req.IP == "" {
		writeErr(w, http.StatusBadRequest, "آدرس IP الزامی است")
		return
	}

	h.cfg.Server.AddressType = req.AddressType
	h.cfg.Server.IP = req.IP
	h.cfg.Server.DNS = req.DNS
	h.cfg.Server.Port = req.Port
	if err := h.cfg.Save(); err != nil {
		writeErr(w, http.StatusInternalServerError, "خطا در ذخیره تنظیمات")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"ok": true, "port": req.Port})
}

// Backup creates a backup of the data store
func (h *SettingsHandler) Backup(w http.ResponseWriter, r *http.Request) {
	store := h.store()
	// BackupPath is an absolute OS-appropriate path (e.g. %PROGRAMDATA%\Khan\backups).
	// If it's relative (legacy), fall back to data_dir-relative.
	backupDir := h.cfg.Backup.BackupPath
	if !filepath.IsAbs(backupDir) {
		backupDir = filepath.Join(h.cfg.Server.DataDir, backupDir)
	}
	dest, err := store.Backup(backupDir)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطا در بکاپ‌گیری")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"ok": true, "path": dest})
}

// backupDir resolves the configured backup directory (absolute OS path
// or relative to the data dir for legacy configs).
func (h *SettingsHandler) backupDir() string {
	dir := h.cfg.Backup.BackupPath
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(h.cfg.Server.DataDir, dir)
	}
	return dir
}

// ListBackups returns the most recent backup folders
func (h *SettingsHandler) ListBackups(w http.ResponseWriter, r *http.Request) {
	dir := h.backupDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		// No backups yet — empty list is fine
		writeJSON(w, http.StatusOK, []interface{}{})
		return
	}
	type backupInfo struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
		Time string `json:"time"`
	}
	var backups = make([]backupInfo, 0)
	for _, e := range entries {
		if !e.IsDir() {
			continue // only backup directories
		}
		// Calculate total size of backup directory
		var totalSize int64
		filepath.Walk(filepath.Join(dir, e.Name()), func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				totalSize += info.Size()
			}
			return nil
		})
		info, _ := e.Info()
		backups = append(backups, backupInfo{
			Name: e.Name(),
			Size: totalSize,
			Time: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}
	writeJSON(w, http.StatusOK, backups)
}

// RestoreBackup loads a backup folder (by name) back into the store.
func (h *SettingsHandler) RestoreBackup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if req.Name == "" || strings.Contains(req.Name, "..") {
		writeErr(w, http.StatusBadRequest, "نام بکاپ نامعتبر است")
		return
	}
	dir := filepath.Join(h.backupDir(), filepath.Base(req.Name))
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		writeErr(w, http.StatusNotFound, "بکاپ یافت نشد")
		return
	}
	if err := h.store().RestoreFrom(dir); err != nil {
		writeErr(w, http.StatusInternalServerError, "خطا در بازیابی: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"ok": true, "restored": req.Name})
}

// DownloadBackup zips a backup folder and streams it as an attachment.
func (h *SettingsHandler) DownloadBackup(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" || strings.Contains(name, "..") {
		writeErr(w, http.StatusBadRequest, "نام بکاپ نامعتبر است")
		return
	}
	dir := filepath.Join(h.backupDir(), filepath.Base(name))
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		writeErr(w, http.StatusNotFound, "بکاپ یافت نشد")
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+name+".zip")
	zw := zip.NewWriter(w)
	defer zw.Close()
	_ = filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		f, err := zw.Create(rel)
		if err != nil {
			return err
		}
		content, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		_, _ = f.Write(content)
		return nil
	})
}

// store returns the active store (set via SetStore)
func (h *SettingsHandler) store() *database.Store { return h.storeRef }

// Logs returns server logs
func (h *SettingsHandler) Logs(w http.ResponseWriter, r *http.Request) {
	logs := "=== Khan Server Logs ===\nServer started on port " + strconv.Itoa(h.cfg.Server.Port) + "\n"
	logs += "Version: 1.0.6\n"
	logs += "Data dir: " + h.cfg.Server.DataDir + "\n"
	logs += "Logs can be viewed via journalctl -u khan or in the terminal running the server.\n"
	writeJSON(w, http.StatusOK, map[string]string{"logs": logs})
}

var _ = service.LicenseState
