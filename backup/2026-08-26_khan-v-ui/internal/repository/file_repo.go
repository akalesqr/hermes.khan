package repository

import (
	"errors"

	"khan/internal/database"
	"khan/internal/models"
)

// ErrNotFound is returned when a record doesn't exist
var ErrNotFound = errors.New("not found")

var errNotFound = ErrNotFound

// FileRepo handles file metadata persistence
type FileRepo struct {
	store *database.Store
}

func NewFileRepo(store *database.Store) *FileRepo { return &FileRepo{store: store} }

// Create inserts file metadata
func (r *FileRepo) Create(f *models.File) (int64, error) {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	f.ID = r.store.NextID("files")
	rec := database.FileRecord{
		ID:        f.ID,
		OwnerID:   f.OwnerID,
		RoomID:    f.RoomID,
		FileName:  f.FileName,
		StoredAs:  f.StoredAs,
		Size:      f.Size,
		MimeType:  f.MimeType,
		CreatedAt: f.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	r.store.Data().Files = append(r.store.Data().Files, rec)
	return f.ID, r.store.SaveLocked()
}

// GetByID finds a file
func (r *FileRepo) GetByID(id int64) (*models.File, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, rec := range r.store.Data().Files {
		if rec.ID == id {
			return &models.File{
				ID: rec.ID, OwnerID: rec.OwnerID, RoomID: rec.RoomID,
				FileName: rec.FileName, StoredAs: rec.StoredAs,
				Size: rec.Size, MimeType: rec.MimeType,
			}, nil
		}
	}
	return nil, nil
}

// CanAccess checks if user has access to a file (member of its room)
func (r *FileRepo) CanAccess(fileID, userID int64) (bool, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, rec := range r.store.Data().Files {
		if rec.ID != fileID {
			continue
		}
		for _, m := range r.store.Data().RoomMembers {
			if m.RoomID == rec.RoomID && m.UserID == userID {
				return true, nil
			}
		}
		return false, nil
	}
	return false, nil
}
