package repository

import (
	"time"

	"khan/internal/database"
	"khan/internal/models"
)

// ReadRepo handles read receipts (رسید خواندن)
type ReadRepo struct {
	store *database.Store
}

func NewReadRepo(store *database.Store) *ReadRepo { return &ReadRepo{store: store} }

// SetRead marks a room as read up to messageID for a user
func (r *ReadRepo) SetRead(roomID, userID, messageID int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()

	for i, rec := range r.store.Data().Reads {
		if rec.RoomID == roomID && rec.UserID == userID {
			if messageID > rec.LastRead {
				r.store.Data().Reads[i].LastRead = messageID
				r.store.Data().Reads[i].UpdatedAt = time.Now().Format(time.RFC3339)
				return r.store.SaveLocked()
			}
			return nil
		}
	}
	r.store.Data().Reads = append(r.store.Data().Reads, database.ReadRecord{
		ID:        r.store.NextID("reads"),
		RoomID:    roomID,
		UserID:    userID,
		LastRead:  messageID,
		UpdatedAt: time.Now().Format(time.RFC3339),
	})
	return r.store.SaveLocked()
}

// GetRead returns the last read message id for a user in a room
func (r *ReadRepo) GetRead(roomID, userID int64) (int64, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, rec := range r.store.Data().Reads {
		if rec.RoomID == roomID && rec.UserID == userID {
			return rec.LastRead, nil
		}
	}
	return 0, nil
}

// ListReadsForRoom returns read positions of all members in a room
func (r *ReadRepo) ListReadsForRoom(roomID int64) ([]models.ReadReceipt, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	var out = make([]models.ReadReceipt, 0)
	for _, rec := range r.store.Data().Reads {
		if rec.RoomID == roomID {
			out = append(out, models.ReadReceipt{
				ID: rec.ID, RoomID: rec.RoomID, UserID: rec.UserID, LastRead: rec.LastRead,
			})
		}
	}
	return out, nil
}

// UnreadCount returns how many messages a user hasn't read in a room
func (r *ReadRepo) UnreadCount(roomID, userID int64) (int64, error) {
	lastRead, _ := r.GetRead(roomID, userID)
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	var count int64
	for _, rec := range r.store.Data().Messages {
		if rec.RoomID == roomID && rec.ID > lastRead && rec.DeletedAt == "" {
			count++
		}
	}
	return count, nil
}
