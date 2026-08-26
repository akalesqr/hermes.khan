package repository

import (
	"errors"
	"sort"
	"time"

	"khan/internal/database"
	"khan/internal/models"
)

// RoomRepo handles room persistence
type RoomRepo struct {
	store *database.Store
}

func NewRoomRepo(store *database.Store) *RoomRepo { return &RoomRepo{store: store} }

// ListDepartments returns all departments
func (r *RoomRepo) ListDepartments() ([]models.Department, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	depts := make([]models.Department, 0)
	for _, rec := range r.store.Data().Departments {
		depts = append(depts, models.Department{ID: rec.ID, Name: rec.Name})
	}
	return depts, nil
}

// CreateDepartment creates a department
func (r *RoomRepo) CreateDepartment(name string) (*models.Department, error) {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	id := r.store.NextID("departments")
	r.store.Data().Departments = append(r.store.Data().Departments, database.DepartmentRecord{
		ID: id, Name: name,
	})
	if err := r.store.SaveLocked(); err != nil {
		return nil, err
	}
	return &models.Department{ID: id, Name: name}, nil
}

// DeleteDepartment removes a department
func (r *RoomRepo) DeleteDepartment(id int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	data := r.store.Data()
	var kept []database.DepartmentRecord
	for _, d := range data.Departments {
		if d.ID != id {
			kept = append(kept, d)
		}
	}
	data.Departments = kept
	for i := range data.Rooms {
		if data.Rooms[i].Department == id {
			data.Rooms[i].Department = 0
		}
	}
	return r.store.SaveLocked()
}

func recToRoom(r database.RoomRecord) *models.Room {
	room := &models.Room{
		ID:          r.ID,
		Name:        r.Name,
		Type:        r.Type,
		CreatorID:   r.CreatorID,
		Department:  r.Department,
		Topic:       r.Topic,
		Description: r.Description,
		Avatar:      r.Avatar,
		Archived:    r.Archived,
	}
	if t, err := time.Parse(time.RFC3339, r.CreatedAt); err == nil {
		room.CreatedAt = t
	}
	return room
}

// Create inserts a new room and returns its id
func (r *RoomRepo) Create(room *models.Room) (int64, error) {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()

	room.ID = r.store.NextID("rooms")
	rec := database.RoomRecord{
		ID:          room.ID,
		Name:        room.Name,
		Type:        room.Type,
		CreatorID:   room.CreatorID,
		Department:  room.Department,
		Topic:       room.Topic,
		Description: room.Description,
		CreatedAt:   room.CreatedAt.Format(time.RFC3339),
	}
	r.store.Data().Rooms = append(r.store.Data().Rooms, rec)
	return room.ID, r.store.SaveLocked()
}

// GetByID finds a room
func (r *RoomRepo) GetByID(id int64) (*models.Room, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, rec := range r.store.Data().Rooms {
		if rec.ID == id {
			return recToRoom(rec), nil
		}
	}
	return nil, nil
}

// ListForUser returns rooms the user belongs to
func (r *RoomRepo) ListForUser(userID int64) ([]models.Room, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	rooms := make([]models.Room, 0)
	for _, m := range r.store.Data().RoomMembers {
		if m.UserID == userID {
			for _, rec := range r.store.Data().Rooms {
				if rec.ID == m.RoomID {
					rooms = append(rooms, *recToRoom(rec))
					break
				}
			}
		}
	}
	sort.Slice(rooms, func(i, j int) bool { return rooms[i].ID > rooms[j].ID })
	return rooms, nil
}

// ListPublic returns all public rooms
func (r *RoomRepo) ListPublic() ([]models.Room, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	rooms := make([]models.Room, 0)
	for _, rec := range r.store.Data().Rooms {
		if rec.Type == models.RoomPublic {
			rooms = append(rooms, *recToRoom(rec))
		}
	}
	return rooms, nil
}

// ListPrivateRooms returns private rooms (for admin browsing / invites)
func (r *RoomRepo) ListPrivateRooms(userID int64) ([]models.Room, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	rooms := make([]models.Room, 0)
	for _, rec := range r.store.Data().Rooms {
		if rec.Type == models.RoomPrivate {
			rooms = append(rooms, *recToRoom(rec))
		}
	}
	return rooms, nil
}

// AddMember adds a user to a room
func (r *RoomRepo) AddMember(roomID, userID int64, role string) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	for _, m := range r.store.Data().RoomMembers {
		if m.RoomID == roomID && m.UserID == userID {
			return nil // already member
		}
	}
	rec := database.RoomMemberRecord{
		ID:       r.store.NextID("room_members"),
		RoomID:   roomID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now().Format(time.RFC3339),
	}
	r.store.Data().RoomMembers = append(r.store.Data().RoomMembers, rec)
	return r.store.SaveLocked()
}

// RemoveMember removes a user from a room
func (r *RoomRepo) RemoveMember(roomID, userID int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	var kept []database.RoomMemberRecord
	for _, m := range r.store.Data().RoomMembers {
		if !(m.RoomID == roomID && m.UserID == userID) {
			kept = append(kept, m)
		}
	}
	r.store.Data().RoomMembers = kept
	return r.store.SaveLocked()
}

// IsMember checks membership
func (r *RoomRepo) IsMember(roomID, userID int64) (bool, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, m := range r.store.Data().RoomMembers {
		if m.RoomID == roomID && m.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

// MemberRole returns the member's role in the room
func (r *RoomRepo) MemberRole(roomID, userID int64) (string, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, m := range r.store.Data().RoomMembers {
		if m.RoomID == roomID && m.UserID == userID {
			return m.Role, nil
		}
	}
	return "", nil
}

// ListMembers returns visible members of a room (hidden excluded)
func (r *RoomRepo) ListMembers(roomID int64) ([]models.User, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	users := make([]models.User, 0)
	for _, m := range r.store.Data().RoomMembers {
		if m.RoomID != roomID {
			continue
		}
		for _, rec := range r.store.Data().Users {
			if rec.ID == m.UserID && !rec.Hidden {
				users = append(users, *recToUser(rec))
				break
			}
		}
	}
	sort.Slice(users, func(i, j int) bool { return users[i].DisplayName < users[j].DisplayName })
	return users, nil
}

// FindOrCreateDM finds an existing DM between two users or creates one
func (r *RoomRepo) FindOrCreateDM(userA, userB int64) (int64, error) {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()

	for _, m := range r.store.Data().RoomMembers {
		if m.UserID != userA {
			continue
		}
		for _, m2 := range r.store.Data().RoomMembers {
			if m2.UserID == userB && m2.RoomID == m.RoomID {
				for _, rec := range r.store.Data().Rooms {
					if rec.ID == m.RoomID && rec.Type == "dm" {
						return rec.ID, nil
					}
				}
			}
		}
	}

	// Create new DM room
	roomID := r.store.NextID("rooms")
	now := time.Now().Format(time.RFC3339)
	r.store.Data().Rooms = append(r.store.Data().Rooms, database.RoomRecord{
		ID: roomID, Name: "", Type: "dm", CreatorID: userA, CreatedAt: now,
	})
	for _, uid := range []int64{userA, userB} {
		r.store.Data().RoomMembers = append(r.store.Data().RoomMembers, database.RoomMemberRecord{
			ID: r.store.NextID("room_members"), RoomID: roomID, UserID: uid,
			Role: "member", JoinedAt: now,
		})
	}
	return roomID, r.store.SaveLocked()
}

// UpdateRoom updates room name
func (r *RoomRepo) UpdateRoom(room *models.Room) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	for i, rec := range r.store.Data().Rooms {
		if rec.ID == room.ID {
			r.store.Data().Rooms[i].Name = room.Name
			return r.store.SaveLocked()
		}
	}
	return errors.New("room not found")
}

// DeleteRoom removes a room and its members/messages
func (r *RoomRepo) DeleteRoom(roomID int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	data := r.store.Data()

	var rooms []database.RoomRecord
	for _, rec := range data.Rooms {
		if rec.ID != roomID {
			rooms = append(rooms, rec)
		}
	}
	data.Rooms = rooms

	var members []database.RoomMemberRecord
	for _, m := range data.RoomMembers {
		if m.RoomID != roomID {
			members = append(members, m)
		}
	}
	data.RoomMembers = members

	var msgs []database.MessageRecord
	for _, msg := range data.Messages {
		if msg.RoomID != roomID {
			msgs = append(msgs, msg)
		}
	}
	data.Messages = msgs

	return r.store.SaveLocked()
}
