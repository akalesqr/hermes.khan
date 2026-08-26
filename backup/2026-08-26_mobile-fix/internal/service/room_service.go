package service

import (
	"errors"
	"time"

	"khan/internal/models"
	"khan/internal/repository"
)

// RoomService handles room/membership operations
type RoomService struct {
	rooms  *repository.RoomRepo
	users  *repository.UserRepo
	invites *repository.InviteRepo
}

func NewRoomService(rooms *repository.RoomRepo, users *repository.UserRepo, invites *repository.InviteRepo) *RoomService {
	return &RoomService{rooms: rooms, users: users, invites: invites}
}

// CreateRoom creates a new room (user+)
func (s *RoomService) CreateRoom(actor *models.User, name, roomType string, deptID int64) (*models.Room, error) {
	switch roomType {
	case models.RoomGroup, models.RoomPublic, models.RoomPrivate, models.RoomChannel:
		if name == "" {
			return nil, errors.New("نام گروه الزامی است")
		}
	case models.RoomDM:
		return nil, errors.New("اتاق خصوصی به صورت خودکار ساخته می‌شود")
	default:
		return nil, errors.New("نوع اتاق نامعتبر است")
	}

	room := &models.Room{
		Name:        name,
		Type:        roomType,
		CreatorID:   actor.ID,
		Department:  deptID,
		CreatedAt:   time.Now(),
	}
	id, err := s.rooms.Create(room)
	if err != nil {
		return nil, err
	}
	room.ID = id
	// Creator becomes owner
	if err := s.rooms.AddMember(id, actor.ID, "owner"); err != nil {
		return nil, err
	}
	return room, nil
}

// ListRooms returns rooms for a user (excludes archived by default)
func (s *RoomService) ListRooms(actor *models.User) ([]models.Room, error) {
	return s.rooms.ListForUser(actor.ID)
}

// JoinRoom adds user to a public room
func (s *RoomService) JoinRoom(actor *models.User, roomID int64) error {
	room, err := s.rooms.GetByID(roomID)
	if err != nil || room == nil {
		return errors.New("اتاق یافت نشد")
	}
	if room.Type == models.RoomPrivate {
		return errors.New("این اتاق خصوصی است — برای عضویت دعوت لازم است")
	}
	if room.Type == models.RoomDM {
		return errors.New("اتاق پیام‌رسان خصوصی است")
	}
	return s.rooms.AddMember(roomID, actor.ID, "member")
}

// InviteUser invites a user to a private room (owner/admin only)
func (s *RoomService) InviteUser(actor *models.User, roomID, userID int64) error {
	room, err := s.rooms.GetByID(roomID)
	if err != nil || room == nil {
		return errors.New("اتاق یافت نشد")
	}
	if room.Type == models.RoomDM {
		return errors.New("اتاق پیام‌رسان خصوصی است")
	}
	// Must be member with owner/admin role
	role, err := s.rooms.MemberRole(roomID, actor.ID)
	if err != nil || (role != "owner" && role != "admin") {
		// admin+ can invite to any room
		if !actor.CanManageUsers() {
			return errors.New("فقط مالک/ادمین اتاق می‌تواند دعوت کند")
		}
	}

	target, err := s.users.GetByID(userID)
	if err != nil || target == nil || target.Hidden {
		return errors.New("کاربر یافت نشد")
	}

	// Add invite record
	if err := s.invites.Invite(roomID, userID, actor.ID); err != nil {
		return err
	}
	// Auto-add as member
	return s.rooms.AddMember(roomID, userID, "member")
}

// AcceptInvite accepts an invite and adds user to the room
func (s *RoomService) AcceptInvite(actor *models.User, roomID int64) error {
	return s.rooms.AddMember(roomID, actor.ID, "member")
}

// ListPrivateRooms returns private rooms the user can be invited to
func (s *RoomService) ListPrivateRooms(actor *models.User) ([]models.Room, error) {
	return s.rooms.ListPrivateRooms(actor.ID)
}

// AddMember adds a user to a room (supervisor+ in their groups, admin+ anywhere)
func (s *RoomService) AddMember(actor *models.User, roomID, userID int64) error {
	room, err := s.rooms.GetByID(roomID)
	if err != nil || room == nil {
		return errors.New("اتاق یافت نشد")
	}

	// DM rooms can't be extended
	if room.Type == models.RoomDM {
		return errors.New("این اتاق خصوصی است")
	}

	// Permission: admin can manage all groups; supervisor only own groups
	if actor.CanManageUsers() {
		// admin+ can manage any room
	} else if actor.CanManageRooms() {
		role, err := s.rooms.MemberRole(roomID, actor.ID)
		if err != nil {
			return err
		}
		if role == "" {
			return errors.New("شما عضو این گروه نیستید")
		}
	} else {
		return errors.New("دسترسی غیرمجاز")
	}

	target, err := s.users.GetByID(userID)
	if err != nil || target == nil {
		return errors.New("کاربر یافت نشد")
	}
	if target.Hidden {
		return errors.New("کاربر یافت نشد")
	}
	if !target.Active {
		return errors.New("کاربر غیرفعال است")
	}
	return s.rooms.AddMember(roomID, userID, "member")
}

// RemoveMember removes a user from a room (supervisor+ in their groups, admin+ anywhere)
func (s *RoomService) RemoveMember(actor *models.User, roomID, userID int64) error {
	room, err := s.rooms.GetByID(roomID)
	if err != nil || room == nil {
		return errors.New("اتاق یافت نشد")
	}
	if room.Type == models.RoomDM {
		return errors.New("این اتاق خصوصی است")
	}

	if actor.CanManageUsers() {
		// admin+ can manage any room
	} else if actor.CanManageRooms() {
		role, err := s.rooms.MemberRole(roomID, actor.ID)
		if err != nil {
			return err
		}
		if role == "" {
			return errors.New("شما عضو این گروه نیستید")
		}
	} else {
		return errors.New("دسترسی غیرمجاز")
	}

	target, err := s.users.GetByID(userID)
	if err != nil || target == nil {
		return errors.New("کاربر یافت نشد")
	}
	if target.Hidden {
		return errors.New("کاربر یافت نشد")
	}

	// Can't remove room owner
	role, err := s.rooms.MemberRole(roomID, userID)
	if err != nil {
		return err
	}
	if role == "owner" && actor.ID != userID {
		return errors.New("نمی‌توانید سازنده گروه را حذف کنید")
	}
	return s.rooms.RemoveMember(roomID, userID)
}

// RenameRoom updates group name (supervisor+)
func (s *RoomService) RenameRoom(actor *models.User, roomID int64, name string) error {
	if name == "" {
		return errors.New("نام گروه الزامی است")
	}
	room, err := s.rooms.GetByID(roomID)
	if err != nil || room == nil {
		return errors.New("اتاق یافت نشد")
	}
	if room.Type == models.RoomDM {
		return errors.New("اتاق خصوصی نام ندارد")
	}

	if actor.CanManageUsers() {
		// admin+
	} else if actor.CanManageRooms() {
		role, err := s.rooms.MemberRole(roomID, actor.ID)
		if err != nil || role == "" {
			return errors.New("شما عضو این گروه نیستید")
		}
	} else {
		return errors.New("دسترسی غیرمجاز")
	}

	room.Name = name
	return s.rooms.UpdateRoom(room)
}

// Members lists visible members of a room
func (s *RoomService) Members(actor *models.User, roomID int64) ([]models.User, error) {
	isMember, err := s.rooms.IsMember(roomID, actor.ID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("شما عضو این اتاق نیستید")
	}
	return s.rooms.ListMembers(roomID)
}

// FindOrCreateDM finds or creates a DM room between two users
func (s *RoomService) FindOrCreateDM(actor *models.User, otherID int64) (int64, error) {
	target, err := s.users.GetByID(otherID)
	if err != nil || target == nil {
		return 0, errors.New("کاربر یافت نشد")
	}
	if target.Hidden {
		return 0, errors.New("کاربر یافت نشد")
	}
	if !target.Active {
		return 0, errors.New("کاربر غیرفعال است")
	}
	if actor.ID == otherID {
		return 0, errors.New("نمی‌توانید با خودتان چت کنید")
	}
	return s.rooms.FindOrCreateDM(actor.ID, otherID)
}
