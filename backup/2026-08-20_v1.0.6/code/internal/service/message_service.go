package service

import (
	"errors"
	"strings"
	"time"

	"khan/internal/models"
	"khan/internal/repository"
)

// MessageService handles message send/edit/delete/reactions with crypto at rest
type MessageService struct {
	msgs   *repository.MessageRepo
	rooms  *repository.RoomRepo
	users  *repository.UserRepo
	crypto *CryptoService
}

func NewMessageService(msgs *repository.MessageRepo, rooms *repository.RoomRepo, users *repository.UserRepo, crypto *CryptoService) *MessageService {
	return &MessageService{msgs: msgs, rooms: rooms, users: users, crypto: crypto}
}

// SendOptions carries optional message features (reply, mentions, urgent)
type SendOptions struct {
	ReplyTo  *int64
	Mentions []int64
	Urgent   bool
	PollID   *int64
}

// Send creates and stores an encrypted message
func (s *MessageService) Send(actor *models.User, roomID int64, text string, fileID *int64, opts *SendOptions) (*models.MessageView, error) {
	isMember, err := s.rooms.IsMember(roomID, actor.ID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("شما عضو این اتاق نیستید")
	}
	if text == "" && fileID == nil {
		return nil, errors.New("پیام خالی است")
	}
	if len(text) > 8000 {
		return nil, errors.New("پیام خیلی طولانی است")
	}

	enc, err := s.crypto.Encrypt([]byte(text))
	if err != nil {
		return nil, err
	}

	if opts == nil {
		opts = &SendOptions{}
	}

	msg := &models.Message{
		RoomID:    roomID,
		SenderID:  actor.ID,
		Body:      enc,
		FileID:    fileID,
		CreatedAt: time.Now(),
		ReplyTo:   opts.ReplyTo,
		Mentions:  opts.Mentions,
		Urgent:    opts.Urgent,
		PollID:    opts.PollID,
	}
	id, err := s.msgs.Create(msg)
	if err != nil {
		return nil, err
	}
	msg.ID = id

	view := s.buildView(msg, text, actor)
	return view, nil
}

// GetByID returns a raw message by id (no decryption)
func (s *MessageService) GetByID(id int64) (*models.Message, error) {
	return s.msgs.GetByID(id)
}

// ToggleUrgent flips the urgent flag on a message
func (s *MessageService) ToggleUrgent(actor *models.User, messageID int64) (*models.Message, error) {
	msg, err := s.msgs.GetByID(messageID)
	if err != nil || msg == nil {
		return nil, errors.New("پیام یافت نشد")
	}
	canManage := actor.Role == models.RoleAdmin || actor.Role == models.RoleSuperAdmin
	if !canManage && actor.Role == models.RoleSupervisor {
		role, err := s.rooms.MemberRole(msg.RoomID, actor.ID)
		if err == nil && role != "" {
			canManage = true
		}
	}
	if msg.SenderID != actor.ID && !canManage {
		return nil, errors.New("دسترسی غیرمجاز")
	}
	if err := s.msgs.ToggleUrgent(messageID, !msg.Urgent); err != nil {
		return nil, err
	}
	msg.Urgent = !msg.Urgent
	return msg, nil
}

// Forward copies a message to another room (actor must be member of both)
func (s *MessageService) Forward(actor *models.User, sourceRoomID, targetRoomID, messageID int64) (*models.MessageView, error) {
	// must be member of target room
	isMember, err := s.rooms.IsMember(targetRoomID, actor.ID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("شما عضو اتاق مقصد نیستید")
	}
	// must be member of source room
	isSrcMember, _ := s.rooms.IsMember(sourceRoomID, actor.ID)
	if !isSrcMember {
		return nil, errors.New("شما عضو اتاق مبدأ نیستید")
	}

	msg, err := s.msgs.GetByID(messageID)
	if err != nil || msg == nil {
		return nil, errors.New("پیام یافت نشد")
	}

	text, decErr := s.crypto.Decrypt(msg.Body)
	if decErr != nil {
		text = []byte("🔒")
	}

	// New message with forwarded pointer
	enc, _ := s.crypto.Encrypt(text)
	forwardedID := msg.ID
	newMsg := &models.Message{
		RoomID:    targetRoomID,
		SenderID:  actor.ID,
		Body:      enc,
		CreatedAt: time.Now(),
		Forwarded: &forwardedID,
	}
	id, err := s.msgs.Create(newMsg)
	if err != nil {
		return nil, err
	}
	newMsg.ID = id

	return s.buildView(newMsg, string(text), actor), nil
}

// CreatePoll creates a poll attached as a message
func (s *MessageService) CreatePoll(actor *models.User, roomID int64, question string, options []string) (*models.MessageView, *models.Poll, error) {
	isMember, err := s.rooms.IsMember(roomID, actor.ID)
	if err != nil || !isMember {
		return nil, nil, errors.New("شما عضو این اتاق نیستید")
	}
	if question == "" {
		return nil, nil, errors.New("سوال نظرسنجی خالی است")
	}
	if len(options) < 2 || len(options) > 10 {
		return nil, nil, errors.New("۲ تا ۱۰ گزینه نیاز است")
	}

	// Build a poll message
	text := "📊 نظرسنجی: " + question + "\n" + "1. " + options[0]
	for i := 1; i < len(options); i++ {
		text += "\n" + string(rune('1'+i)) + ". " + options[i]
	}
	enc, _ := s.crypto.Encrypt([]byte(text))

	msg := &models.Message{
		RoomID:    roomID,
		SenderID:  actor.ID,
		Body:      enc,
		CreatedAt: time.Now(),
	}
	id, err := s.msgs.Create(msg)
	if err != nil {
		return nil, nil, err
	}
	msg.ID = id

	view := s.buildView(msg, text, actor)
	return view, nil, nil
}

// buildView constructs an API-facing message view with sender + reply context
func (s *MessageService) buildView(m *models.Message, plaintext string, actor *models.User) *models.MessageView {
	view := &models.MessageView{
		ID:         m.ID,
		RoomID:     m.RoomID,
		SenderID:   m.SenderID,
		Text:       plaintext,
		FileID:     m.FileID,
		CreatedAt:  m.CreatedAt,
		EditedAt:   m.EditedAt,
		DeletedAt:  m.DeletedAt,
		ReplyTo:    m.ReplyTo,
		Forwarded:  m.Forwarded,
		Mentions:   m.Mentions,
		Urgent:     m.Urgent,
		PollID:     m.PollID,
	}
	if m.SenderID == -1 {
		view.SenderName = "کاربر حذف‌شده"
	} else if sender, err := s.users.GetByID(m.SenderID); err == nil && sender != nil {
		view.SenderName = sender.DisplayName
		if view.SenderName == "" {
			view.SenderName = sender.Username
		}
	} else {
		view.SenderName = "کاربر حذف‌شده"
	}
	// reply preview
	if m.ReplyTo != nil {
		if orig, err := s.msgs.GetByID(*m.ReplyTo); err == nil && orig != nil {
			if t, decErr := s.crypto.Decrypt(orig.Body); decErr == nil {
				view.ReplyText = string(t)
				if len(view.ReplyText) > 80 {
					view.ReplyText = view.ReplyText[:80] + "…"
				}
			}
		}
	}
	return view
}

// ListMessages returns decrypted messages for a room (paginated)
func (s *MessageService) ListMessages(actor *models.User, roomID int64, before int64, limit int) ([]models.MessageView, error) {
	isMember, err := s.rooms.IsMember(roomID, actor.ID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("شما عضو این اتاق نیستید")
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	msgs, err := s.msgs.ListByRoom(roomID, before, limit)
	if err != nil {
		return nil, err
	}

	views := make([]models.MessageView, 0)
	for _, m := range msgs {
		text, decErr := s.crypto.Decrypt(m.Body)
		if decErr != nil {
			text = []byte("🔒")
		}
		view := s.buildView(&m, string(text), nil)
		// reactions
		if reacs, err := s.msgs.ListReactions(m.ID); err == nil {
			view.Reactions = reacs
		}
		views = append(views, *view)
	}
	return views, nil
}

// SearchMessages decrypts and searches across accessible rooms
func (s *MessageService) SearchMessages(actor *models.User, query string, limit int) ([]models.MessageView, error) {
	// get all rooms user belongs to
	rooms, err := s.rooms.ListForUser(actor.ID)
	if err != nil {
		return nil, err
	}
	roomIDs := make([]int64, 0, len(rooms))
	for _, r := range rooms {
		roomIDs = append(roomIDs, r.ID)
	}
	// also public rooms
	pub, _ := s.rooms.ListPublic()
	for _, r := range pub {
		found := false
		for _, id := range roomIDs {
			if id == r.ID {
				found = true
				break
			}
		}
		if !found {
			roomIDs = append(roomIDs, r.ID)
		}
	}

	if limit <= 0 || limit > 100 {
		limit = 30
	}
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return []models.MessageView{}, nil
	}

	// Decrypt and search (messages encrypted at rest)
	var results []models.MessageView
	msgs, err := s.msgs.ListByRoomAll(roomIDs)
	if err != nil {
		return nil, err
	}
	for _, m := range msgs {
		text, decErr := s.crypto.Decrypt(m.Body)
		if decErr != nil {
			continue
		}
		if strings.Contains(strings.ToLower(string(text)), q) {
			view := s.buildView(&m, string(text), nil)
			results = append(results, *view)
			if len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

// EditMessage updates message text (owner only)
func (s *MessageService) EditMessage(actor *models.User, messageID int64, newText string) (*models.MessageView, error) {
	msg, err := s.msgs.GetByID(messageID)
	if err != nil || msg == nil {
		return nil, errors.New("پیام یافت نشد")
	}
	if msg.SenderID != actor.ID {
		return nil, errors.New("فقط فرستنده می‌تواند پیام را ویرایش کند")
	}
	if newText == "" {
		return nil, errors.New("متن پیام خالی است")
	}
	enc, err := s.crypto.Encrypt([]byte(newText))
	if err != nil {
		return nil, err
	}
	if err := s.msgs.UpdateBody(messageID, enc); err != nil {
		return nil, err
	}
	updated, _ := s.msgs.GetByID(messageID)
	view := &models.MessageView{
		ID:         messageID,
		RoomID:     msg.RoomID,
		SenderID:   actor.ID,
		SenderName: actor.DisplayName,
		Text:       newText,
		CreatedAt:  msg.CreatedAt,
		EditedAt:   updated.EditedAt,
	}
	if view.SenderName == "" {
		view.SenderName = actor.Username
	}
	return view, nil
}

// DeleteMessage removes a message (owner or supervisor+ in room)
func (s *MessageService) DeleteMessage(actor *models.User, messageID int64) error {
	msg, err := s.msgs.GetByID(messageID)
	if err != nil || msg == nil {
		return errors.New("پیام یافت نشد")
	}

	if msg.SenderID != actor.ID {
		// supervisor+ can delete others' messages in rooms they manage
		if !actor.CanManageRooms() {
			return errors.New("فقط فرستنده می‌تواند پیام را حذف کند")
		}
		isMember, err := s.rooms.IsMember(msg.RoomID, actor.ID)
		if err != nil || !isMember {
			return errors.New("شما عضو این اتاق نیستید")
		}
	}
	return s.msgs.SoftDelete(messageID)
}

// AddReaction adds a reaction (member only)
func (s *MessageService) AddReaction(actor *models.User, messageID int64, emoji string) error {
	if emoji == "" || len([]rune(emoji)) > 8 {
		return errors.New("ایموجی نامعتبر است")
	}
	msg, err := s.msgs.GetByID(messageID)
	if err != nil || msg == nil {
		return errors.New("پیام یافت نشد")
	}
	isMember, err := s.rooms.IsMember(msg.RoomID, actor.ID)
	if err != nil || !isMember {
		return errors.New("شما عضو این اتاق نیستید")
	}
	return s.msgs.AddReaction(messageID, actor.ID, emoji)
}

// RemoveReaction removes a reaction
func (s *MessageService) RemoveReaction(actor *models.User, messageID int64, emoji string) error {
	return s.msgs.RemoveReaction(messageID, actor.ID, emoji)
}
