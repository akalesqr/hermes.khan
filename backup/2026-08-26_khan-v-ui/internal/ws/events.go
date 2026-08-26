package ws

import "time"

// EventType constants (server → client)
const (
	EvMessage        = "message"
	EvMessageEdit    = "message_edited"
	EvMessageDelete  = "message_deleted"
	EvReaction       = "reaction"
	EvTyping         = "typing"
	EvPresence       = "presence"
	EvUserAdded      = "user_added"
	EvUserRemoved    = "user_removed"
	EvRoleChanged    = "role_changed"
	EvUserDeleted    = "user_deleted"
	EvForceLogout    = "force_logout"
	EvReadReceipt    = "read_receipt"
	EvPin            = "pin"
	EvUnpin          = "unpin"
	EvPoll           = "poll"
	EvPollUpdate     = "poll_update"
	EvInvite         = "invite"
	EvOfflineMsg     = "offline_message"
	EvDepartment     = "department"
)

// Event is a server → client message
type Event struct {
	Type    string      `json:"type"`
	RoomID  int64       `json:"room_id,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// NewEvent builds an event
func NewEvent(typ string, roomID int64, payload interface{}) Event {
	return Event{Type: typ, RoomID: roomID, Payload: payload}
}

// ClientEventType constants (client → server)
const (
	CevSendMessage = "send_message"
	CevEditMessage = "edit_message"
	CevDeleteMsg   = "delete_message"
	CevAddReaction = "add_reaction"
	CevDelReaction = "remove_reaction"
	CevTyping      = "typing"
	CevMarkRead    = "mark_read"
	CevJoinRoom    = "join_room"
	CevForward     = "forward_message"
	CevReply       = "reply_message"
	CevPin         = "pin_message"
	CevUnpin       = "unpin_message"
	CevPollCreate  = "poll_create"
	CevPollVote    = "poll_vote"
	CevPollClose   = "poll_close"
	CevPresence    = "presence_set"
)

// ClientEvent is a client → server message
type ClientEvent struct {
	Type     string   `json:"type"`
	RoomID   int64    `json:"room_id"`
	MessageID int64   `json:"message_id,omitempty"`
	Text     string   `json:"text,omitempty"`
	Emoji    string   `json:"emoji,omitempty"`
	FileID   *int64   `json:"file_id,omitempty"`
	LastID   int64    `json:"last_message_id,omitempty"`
	UserID   int64    `json:"-"` // injected server-side
	ReplyTo  *int64   `json:"reply_to,omitempty"`
	Mentions []int64  `json:"mentions,omitempty"`
	Urgent   bool     `json:"urgent,omitempty"`
	ForwardMessageID int64 `json:"forward_message_id,omitempty"`
	TargetRoomID int64 `json:"target_room_id,omitempty"`
	PollQuestion string `json:"poll_question,omitempty"`
	PollOptions  []string `json:"poll_options,omitempty"`
	PollOption   int    `json:"poll_option,omitempty"`
	PollID       int64  `json:"poll_id,omitempty"`
	Status       string `json:"status,omitempty"`
}

// MessagePayload for EvMessage (v1.0.3: full message view)
type MessagePayload struct {
	ID         int64     `json:"id"`
	RoomID     int64     `json:"room_id"`
	SenderID   int64     `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Text       string    `json:"text"`
	FileID     *int64    `json:"file_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	EditedAt   *time.Time `json:"edited_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	ReplyTo    *int64    `json:"reply_to,omitempty"`
	ReplyText  string    `json:"reply_text,omitempty"`
	Forwarded  *int64    `json:"forwarded_from,omitempty"`
	Mentions   []int64   `json:"mentions,omitempty"`
	Urgent     bool      `json:"urgent,omitempty"`
	PollID     *int64    `json:"poll_id,omitempty"`
}

// PresencePayload for EvPresence
type PresencePayload struct {
	UserID int64  `json:"user_id"`
	Online bool   `json:"online"`
	Name   string `json:"name,omitempty"`
}

// TypingPayload for EvTyping
type TypingPayload struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
}

// ReadReceiptPayload for EvReadReceipt
type ReadReceiptPayload struct {
	RoomID int64 `json:"room_id"`
	UserID int64 `json:"user_id"`
	LastID int64 `json:"last_message_id"`
}

// PinPayload for EvPin/EvUnpin
type PinPayload struct {
	RoomID    int64     `json:"room_id"`
	MessageID int64     `json:"message_id"`
	UserID    int64     `json:"user_id,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// PollPayload for EvPoll/EvPollUpdate
type PollPayload struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	CreatorID int64     `json:"creator_id"`
	Question  string    `json:"question"`
	Options   []string  `json:"options"`
	CreatedAt time.Time `json:"created_at"`
	Closed    bool      `json:"closed"`
}
