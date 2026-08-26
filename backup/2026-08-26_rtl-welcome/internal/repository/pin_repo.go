package repository

import (
	"time"

	"khan/internal/database"
	"khan/internal/models"
)

// PinRepo handles pinned messages (پین پیام بالای اتاق)
type PinRepo struct {
	store *database.Store
}

func NewPinRepo(store *database.Store) *PinRepo { return &PinRepo{store: store} }

// Pin pins a message to the top of a room
func (p *PinRepo) Pin(roomID, messageID, userID int64) error {
	p.store.Mu().Lock()
	defer p.store.Mu().Unlock()

	// Check not already pinned
	for _, rec := range p.store.Data().Pins {
		if rec.RoomID == roomID && rec.MessageID == messageID {
			return nil // already pinned
		}
	}

	p.store.Data().Pins = append(p.store.Data().Pins, database.PinRecord{
		ID:        p.store.NextID("pins"),
		RoomID:    roomID,
		MessageID: messageID,
		UserID:    userID,
		CreatedAt: time.Now().Format(time.RFC3339),
	})
	return p.store.SaveLocked()
}

// Unpin removes a pinned message
func (p *PinRepo) Unpin(roomID, messageID int64) error {
	p.store.Mu().Lock()
	defer p.store.Mu().Unlock()
	var kept []database.PinRecord
	for _, rec := range p.store.Data().Pins {
		if !(rec.RoomID == roomID && rec.MessageID == messageID) {
			kept = append(kept, rec)
		}
	}
	p.store.Data().Pins = kept
	return p.store.SaveLocked()
}

// ListForRoom returns pinned messages in a room
func (p *PinRepo) ListForRoom(roomID int64) ([]models.Pin, error) {
	p.store.Mu().RLock()
	defer p.store.Mu().RUnlock()
	var out = make([]models.Pin, 0)
	for _, rec := range p.store.Data().Pins {
		if rec.RoomID == roomID {
			pin := models.Pin{ID: rec.ID, RoomID: rec.RoomID, MessageID: rec.MessageID, UserID: rec.UserID}
			if t, err := time.Parse(time.RFC3339, rec.CreatedAt); err == nil {
				pin.CreatedAt = t
			}
			out = append(out, pin)
		}
	}
	return out, nil
}

// IsPinned checks if a message is pinned
func (p *PinRepo) IsPinned(messageID int64) bool {
	p.store.Mu().RLock()
	defer p.store.Mu().RUnlock()
	for _, rec := range p.store.Data().Pins {
		if rec.MessageID == messageID {
			return true
		}
	}
	return false
}
