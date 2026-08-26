package repository

import (
	"errors"
	"sort"
	"strings"
	"time"

	"khan/internal/database"
	"khan/internal/models"
)

// UserRepo handles user persistence on top of Store
type UserRepo struct {
	store *database.Store
}

func NewUserRepo(store *database.Store) *UserRepo { return &UserRepo{store: store} }

func recToUser(r database.UserRecord) *models.User {
	u := &models.User{
		ID:            r.ID,
		Username:      r.Username,
		PasswordHash:  r.PasswordHash,
		DisplayName:   r.DisplayName,
		AvatarPath:    r.AvatarPath,
		Role:          r.Role,
		Hidden:        r.Hidden,
		Active:        r.Active,
		MustChangePwd: r.MustChangePwd,
	}
	if r.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, r.CreatedAt); err == nil {
			u.CreatedAt = t
		}
	}
	if r.LastSeen != "" {
		if t, err := time.Parse(time.RFC3339, r.LastSeen); err == nil {
			u.LastSeen = &t
		}
	}
	return u
}

func userToRec(u *models.User) database.UserRecord {
	rec := database.UserRecord{
		ID:            u.ID,
		Username:      u.Username,
		PasswordHash:  u.PasswordHash,
		DisplayName:   u.DisplayName,
		AvatarPath:    u.AvatarPath,
		Role:          u.Role,
		Hidden:        u.Hidden,
		Active:        u.Active,
		MustChangePwd: u.MustChangePwd,
		CreatedAt:     u.CreatedAt.Format(time.RFC3339),
	}
	if u.LastSeen != nil {
		rec.LastSeen = u.LastSeen.Format(time.RFC3339)
	}
	return rec
}

// Create inserts a new user
func (r *UserRepo) Create(u *models.User) (int64, error) {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()

	u.ID = r.store.NextID("users")
	rec := userToRec(u)
	r.store.Data().Users = append(r.store.Data().Users, rec)
	return u.ID, r.store.SaveLocked()
}

// GetByID finds a user by id
func (r *UserRepo) GetByID(id int64) (*models.User, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, rec := range r.store.Data().Users {
		if rec.ID == id {
			return recToUser(rec), nil
		}
	}
	return nil, nil
}

// GetByUsername finds a user by username
func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, rec := range r.store.Data().Users {
		if strings.EqualFold(rec.Username, username) {
			return recToUser(rec), nil
		}
	}
	return nil, nil
}

// ListVisible returns all users except hidden super admins
func (r *UserRepo) ListVisible() ([]models.User, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	users := make([]models.User, 0)
	for _, rec := range r.store.Data().Users {
		if !rec.Hidden {
			users = append(users, *recToUser(rec))
		}
	}
	sort.Slice(users, func(i, j int) bool { return users[i].DisplayName < users[j].DisplayName })
	return users, nil
}

// ListAll returns every user (super admin use only)
func (r *UserRepo) ListAll() ([]models.User, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	users := make([]models.User, 0)
	for _, rec := range r.store.Data().Users {
		users = append(users, *recToUser(rec))
	}
	return users, nil
}

// CountVisible counts non-hidden users (for license limit)
func (r *UserRepo) CountVisible() (int, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	n := 0
	for _, rec := range r.store.Data().Users {
		if !rec.Hidden {
			n++
		}
	}
	return n, nil
}

// Update updates a user's mutable fields
func (r *UserRepo) Update(u *models.User) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	rec := userToRec(u)
	for i, existing := range r.store.Data().Users {
		if existing.ID == u.ID {
			r.store.Data().Users[i] = rec
			return r.store.SaveLocked()
		}
	}
	return errors.New("user not found")
}

// UpdatePassword only updates the password hash
func (r *UserRepo) UpdatePassword(id int64, hash string) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	for i, rec := range r.store.Data().Users {
		if rec.ID == id {
			r.store.Data().Users[i].PasswordHash = hash
			r.store.Data().Users[i].MustChangePwd = false
			return r.store.SaveLocked()
		}
	}
	return errors.New("user not found")
}

// Delete removes a user entirely (anonymizes messages)
func (r *UserRepo) Delete(id int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()

	data := r.store.Data()
	var kept []database.UserRecord
	for _, rec := range data.Users {
		if rec.ID != id {
			kept = append(kept, rec)
		}
	}
	data.Users = kept

	for i := range data.Messages {
		if data.Messages[i].SenderID == id {
			data.Messages[i].SenderID = -1 // deleted user
		}
	}

	var members []database.RoomMemberRecord
	for _, m := range data.RoomMembers {
		if m.UserID != id {
			members = append(members, m)
		}
	}
	data.RoomMembers = members

	var sessions []database.SessionRecord
	for _, s := range data.Sessions {
		if s.UserID != id {
			sessions = append(sessions, s)
		}
	}
	data.Sessions = sessions

	return r.store.SaveLocked()
}

// TouchLastSeen updates the presence timestamp
func (r *UserRepo) TouchLastSeen(id int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	for i, rec := range r.store.Data().Users {
		if rec.ID == id {
			r.store.Data().Users[i].LastSeen = time.Now().Format(time.RFC3339)
			return r.store.SaveLocked()
		}
	}
	return nil
}

// SetOnline marks a user online/offline and touches LastSeen
func (r *UserRepo) SetOnline(id int64, online bool) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	now := time.Now().Format(time.RFC3339)
	for i, rec := range r.store.Data().Users {
		if rec.ID == id {
			r.store.Data().Users[i].LastSeen = now
			if online {
				r.store.Data().Users[i].Status = "online"
			} else {
				r.store.Data().Users[i].Status = "offline"
			}
			return r.store.SaveLocked()
		}
	}
	return nil
}

// SetStatus updates a user's custom status
func (r *UserRepo) SetStatus(id int64, status string) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	for i, rec := range r.store.Data().Users {
		if rec.ID == id {
			r.store.Data().Users[i].Status = status
			return r.store.SaveLocked()
		}
	}
	return nil
}

// Search finds visible users by username/display name
func (r *UserRepo) Search(q string) ([]models.User, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	users := make([]models.User, 0)
	for _, rec := range r.store.Data().Users {
		if rec.Hidden {
			continue
		}
		if strings.Contains(strings.ToLower(rec.Username), strings.ToLower(q)) ||
			strings.Contains(strings.ToLower(rec.DisplayName), strings.ToLower(q)) {
			users = append(users, *recToUser(rec))
		}
		if len(users) >= 50 {
			break
		}
	}
	return users, nil
}
