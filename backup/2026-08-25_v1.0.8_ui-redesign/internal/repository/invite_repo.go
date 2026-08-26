package repository

import (
	"time"

	"khan/internal/database"
	"khan/internal/models"
)

// InviteRepo handles private room invitations (دعوت به اتاق خصوصی)
type InviteRepo struct {
	store *database.Store
}

func NewInviteRepo(store *database.Store) *InviteRepo { return &InviteRepo{store: store} }

// Invite adds a user to a private room
func (r *InviteRepo) Invite(roomID, userID, byUserID int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()

	r.store.Data().Invites = append(r.store.Data().Invites, database.InviteRecord{
		ID:        r.store.NextID("invites"),
		RoomID:    roomID,
		UserID:    userID,
		ByUserID:  byUserID,
		CreatedAt: time.Now().Format(time.RFC3339),
	})
	return r.store.SaveLocked()
}

// ListForUser returns private room invites for a user
func (r *InviteRepo) ListForUser(userID int64) ([]models.Invite, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	var out = make([]models.Invite, 0)
	for _, rec := range r.store.Data().Invites {
		if rec.UserID == userID {
			inv := models.Invite{ID: rec.ID, RoomID: rec.RoomID, UserID: rec.UserID, ByUserID: rec.ByUserID}
			if t, err := time.Parse(time.RFC3339, rec.CreatedAt); err == nil {
				inv.CreatedAt = t
			}
			out = append(out, inv)
		}
	}
	return out, nil
}

// Remove deletes an invite
func (r *InviteRepo) Remove(roomID, userID int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	var kept []database.InviteRecord
	for _, inv := range r.store.Data().Invites {
		if !(inv.RoomID == roomID && inv.UserID == userID) {
			kept = append(kept, inv)
		}
	}
	r.store.Data().Invites = kept
	return r.store.SaveLocked()
}
