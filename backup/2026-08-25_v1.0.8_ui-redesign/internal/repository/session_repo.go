package repository

import (
	"time"

	"khan/internal/database"
	"khan/internal/models"
)

// SessionRepo handles session persistence
type SessionRepo struct {
	store *database.Store
}

func NewSessionRepo(store *database.Store) *SessionRepo { return &SessionRepo{store: store} }

// Create stores a new session token
func (r *SessionRepo) Create(s *models.Session) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	rec := database.SessionRecord{
		ID:        r.store.NextID("sessions"),
		UserID:    s.UserID,
		Token:     s.Token,
		ExpiresAt: s.ExpiresAt.Format(time.RFC3339),
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
	}
	r.store.Data().Sessions = append(r.store.Data().Sessions, rec)
	return r.store.SaveLocked()
}

// GetByToken finds a session by token
func (r *SessionRepo) GetByToken(token string) (*models.Session, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, rec := range r.store.Data().Sessions {
		if rec.Token == token {
			s := &models.Session{ID: rec.ID, UserID: rec.UserID, Token: rec.Token}
			if t, err := time.Parse(time.RFC3339, rec.ExpiresAt); err == nil {
				s.ExpiresAt = t
			}
			if t, err := time.Parse(time.RFC3339, rec.CreatedAt); err == nil {
				s.CreatedAt = t
			}
			return s, nil
		}
	}
	return nil, errNotFound
}

// Delete removes a session (logout)
func (r *SessionRepo) Delete(token string) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	var kept []database.SessionRecord
	for _, rec := range r.store.Data().Sessions {
		if rec.Token != token {
			kept = append(kept, rec)
		}
	}
	r.store.Data().Sessions = kept
	return r.store.SaveLocked()
}

// DeleteForUser removes all sessions for a user (force logout)
func (r *SessionRepo) DeleteForUser(userID int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	var kept []database.SessionRecord
	for _, rec := range r.store.Data().Sessions {
		if rec.UserID != userID {
			kept = append(kept, rec)
		}
	}
	r.store.Data().Sessions = kept
	return r.store.SaveLocked()
}

// CleanupExpired removes expired sessions
func (r *SessionRepo) CleanupExpired() error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	now := time.Now()
	var kept []database.SessionRecord
	for _, rec := range r.store.Data().Sessions {
		if t, err := time.Parse(time.RFC3339, rec.ExpiresAt); err == nil && t.After(now) {
			kept = append(kept, rec)
		}
	}
	r.store.Data().Sessions = kept
	return r.store.SaveLocked()
}
