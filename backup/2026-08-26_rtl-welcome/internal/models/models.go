package models

import "time"

// Role constants (obfuscated storage — short codes)
const (
	RoleUser       = "user"
	RoleSupervisor = "sup"
	RoleAdmin      = "adm"
	RoleSuperAdmin = "sadm"
)

// Room types
const (
	RoomDM      = "dm"
	RoomGroup   = "group"
	RoomPublic  = "public"
	RoomPrivate = "private"
	RoomChannel = "channel"
)

// User represents a chat user
type User struct {
	ID            int64      `json:"id" db:"id"`
	Username      string     `json:"username" db:"username"`
	PasswordHash  string     `json:"-" db:"password_hash"`
	DisplayName   string     `json:"display_name" db:"display_name"`
	AvatarPath    string     `json:"avatar_path,omitempty" db:"avatar_path"`
	Role          string     `json:"role" db:"role"`
	Hidden        bool       `json:"-" db:"hidden"`
	Active        bool       `json:"active" db:"active"`
	MustChangePwd bool       `json:"must_change_pwd" db:"must_change_pwd"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	LastSeen      *time.Time `json:"last_seen,omitempty" db:"last_seen"`
	Status        string     `json:"status,omitempty" db:"status"`
	StatusText    string     `json:"status_text,omitempty" db:"status_text"`
}

// Sanitized returns a user safe for API responses (no sensitive fields)
func (u *User) Sanitized() *User {
	c := *u
	c.PasswordHash = ""
	c.Hidden = false
	return &c
}

// IsHidden reports whether this user should be invisible to others
func (u *User) IsHidden() bool {
	return u.Hidden
}

// IsOnline reports online status (based on last_seen within 2 min)
func (u *User) IsOnline() bool {
	if u.LastSeen == nil {
		return false
	}
	return time.Since(*u.LastSeen) < 2*time.Minute
}

// CanManageRooms checks supervisor+ level
func (u *User) CanManageRooms() bool {
	return u.Role == RoleSupervisor || u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

// CanManageUsers checks admin+ level
func (u *User) CanManageUsers() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

// CanManageAdmins checks super admin only
func (u *User) CanManageAdmins() bool {
	return u.Role == RoleSuperAdmin
}

// Room represents a chat room
type Room struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Type        string    `json:"type" db:"type"`
	CreatorID   int64     `json:"creator_id" db:"creator_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Department  int64     `json:"department,omitempty" db:"department"`
	Topic       string    `json:"topic,omitempty" db:"topic"`
	Description string    `json:"description,omitempty" db:"description"`
	Avatar      string    `json:"avatar,omitempty" db:"avatar"`
	Archived    bool      `json:"archived,omitempty" db:"archived"`
}

// RoomMember represents membership of a user in a room
type RoomMember struct {
	ID       int64     `json:"id" db:"id"`
	RoomID   int64     `json:"room_id" db:"room_id"`
	UserID   int64     `json:"user_id" db:"user_id"`
	Role     string    `json:"role" db:"role"` // "owner" | "admin" | "member"
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

// Message represents a chat message (body is encrypted at rest)
type Message struct {
	ID        int64      `json:"id" db:"id"`
	RoomID    int64      `json:"room_id" db:"room_id"`
	SenderID  int64      `json:"sender_id" db:"sender_id"`
	Body      []byte     `json:"body" db:"body"`
	IV        []byte     `json:"-" db:"iv"`
	FileID    *int64     `json:"file_id,omitempty" db:"file_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	EditedAt  *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	ReplyTo   *int64     `json:"reply_to,omitempty" db:"reply_to"`
	Forwarded *int64     `json:"forwarded_from,omitempty" db:"forwarded_from"`
	Mentions  []int64    `json:"mentions,omitempty" db:"mentions"`
	Urgent    bool       `json:"urgent,omitempty" db:"urgent"`
	PollID    *int64     `json:"poll_id,omitempty" db:"poll_id"`
}

// MessageView is the API-facing message with sender info
type MessageView struct {
	ID         int64      `json:"id" db:"id"`
	RoomID     int64      `json:"room_id" db:"room_id"`
	SenderID   int64      `json:"sender_id" db:"sender_id"`
	SenderName string     `json:"sender_name" db:"sender_name"`
	Text       string     `json:"text" db:"text"`
	FileID     *int64     `json:"file_id,omitempty" db:"file_id"`
	FileName   string     `json:"file_name,omitempty" db:"file_name"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	EditedAt   *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	Reactions  []Reaction `json:"reactions,omitempty" db:"reactions"`
	ReplyTo    *int64     `json:"reply_to,omitempty" db:"reply_to"`
	ReplyText  string     `json:"reply_text,omitempty" db:"reply_text"`
	Forwarded  *int64     `json:"forwarded_from,omitempty" db:"forwarded_from"`
	Mentions   []int64    `json:"mentions,omitempty" db:"mentions"`
	Urgent     bool       `json:"urgent,omitempty" db:"urgent"`
	PollID     *int64     `json:"poll_id,omitempty" db:"poll_id"`
	Poll       *Poll      `json:"poll,omitempty" db:"poll"`
}

// Reaction represents an emoji reaction on a message
type Reaction struct {
	ID        int64     `json:"id" db:"id"`
	MessageID int64     `json:"message_id" db:"message_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Emoji     string    `json:"emoji" db:"emoji"`
}

// File represents an uploaded file
type File struct {
	ID        int64     `json:"id" db:"id"`
	OwnerID   int64     `json:"owner_id" db:"owner_id"`
	RoomID    int64     `json:"room_id" db:"room_id"`
	FileName  string    `json:"file_name" db:"file_name"`
	StoredAs  string    `json:"-" db:"stored_as"`
	Size      int64     `json:"size" db:"size"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Session represents an active login session
type Session struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Token     string    `json:"-" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ReadReceipt tracks a user's read position in a room
type ReadReceipt struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	UserID    int64     `json:"user_id"`
	LastRead  int64     `json:"last_read"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Poll is a poll in a room
type Poll struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	CreatorID int64     `json:"creator_id"`
	Question  string    `json:"question"`
	Options   []string  `json:"options"`
	CreatedAt time.Time `json:"created_at"`
	Closed    bool      `json:"closed"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	Votes     []PollVote `json:"votes,omitempty"`
}

// PollVote is one user's vote
type PollVote struct {
	ID      int64     `json:"id"`
	PollID  int64     `json:"poll_id"`
	UserID  int64     `json:"user_id"`
	Option  int       `json:"option"`
	VotedAt time.Time `json:"voted_at"`
}

// Invite is an invitation to a private room
type Invite struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	UserID    int64     `json:"user_id"`
	ByUserID  int64     `json:"by_user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Pin pins a message to the top of a room
type Pin struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	MessageID int64     `json:"message_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Department is an organizational section that groups rooms
type Department struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
