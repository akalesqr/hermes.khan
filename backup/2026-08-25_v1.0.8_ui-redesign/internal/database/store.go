package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Store is a lightweight JSON-file datastore (replaces SQLite for LAN scale).
// v1.0.3: one file per table so chat volume never slows down user operations.
//
//	users.json          — user accounts
//	rooms.json          — rooms/channels
//	room_members.json   — membership
//	messages.json       — messages (largest table, isolated)
//	reactions.json      — emoji reactions
//	sessions.json       — login sessions
//	files.json          — uploaded file metadata
//	reads.json          — read receipts (per user, per room, last_read_id)
//	polls.json          — polls
//	poll_votes.json     — poll votes
//	invites.json        — private room invites
//	pins.json           — pinned messages
//	seq.json            — per-table id counters
type Store struct {
	dir  string
	mu   sync.RWMutex
	data *Data
}

// Data holds all tables in memory, persisted to separate files on disk.
type Data struct {
	Users       []UserRecord       `json:"users"`
	Rooms       []RoomRecord       `json:"rooms"`
	RoomMembers []RoomMemberRecord `json:"room_members"`
	Messages    []MessageRecord    `json:"messages"`
	Reactions   []ReactionRecord   `json:"reactions"`
	Sessions    []SessionRecord    `json:"sessions"`
	Files       []FileRecord       `json:"files"`
	Reads       []ReadRecord       `json:"reads"`
	Polls       []PollRecord       `json:"polls"`
	PollVotes   []PollVoteRecord   `json:"poll_votes"`
	Invites     []InviteRecord     `json:"invites"`
	Pins        []PinRecord        `json:"pins"`
	Departments []DepartmentRecord `json:"departments"`
	Seq         map[string]int64   `json:"seq"`
}

// UserRecord mirrors the users table
type UserRecord struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	PasswordHash  string `json:"password_hash"`
	DisplayName   string `json:"display_name"`
	AvatarPath    string `json:"avatar_path,omitempty"`
	Role          string `json:"role"`
	Hidden        bool   `json:"hidden"`
	Active        bool   `json:"active"`
	MustChangePwd bool   `json:"must_change_pwd"`
	CreatedAt     string `json:"created_at"`
	LastSeen      string `json:"last_seen,omitempty"`
	Status        string `json:"status,omitempty"` // "online" | "offline" | "away"
	StatusText    string `json:"status_text,omitempty"`
}

type RoomRecord struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"` // dm | group | public | private | channel
	CreatorID   int64  `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
	Department  int64  `json:"department,omitempty"` // department id (0 = none)
	Topic       string `json:"topic,omitempty"`      // channel topic
	Description string `json:"description,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Archived    bool   `json:"archived,omitempty"`
}

type RoomMemberRecord struct {
	ID       int64  `json:"id"`
	RoomID   int64  `json:"room_id"`
	UserID   int64  `json:"user_id"`
	Role     string `json:"role"` // "owner" | "admin" | "member"
	JoinedAt string `json:"joined_at"`
}

type MessageRecord struct {
	ID         int64   `json:"id"`
	RoomID     int64   `json:"room_id"`
	SenderID   int64   `json:"sender_id"` // -1 = deleted user
	Body       string  `json:"body"`      // base64 AES-GCM
	FileID     *int64  `json:"file_id,omitempty"`
	CreatedAt  string  `json:"created_at"`
	EditedAt   string  `json:"edited_at,omitempty"`
	DeletedAt  string  `json:"deleted_at,omitempty"`
	ReplyTo    *int64  `json:"reply_to,omitempty"`
	Forwarded  *int64  `json:"forwarded_from,omitempty"` // original message id
	Mentions   []int64 `json:"mentions,omitempty"`       // mentioned user ids
	Urgent     bool    `json:"urgent,omitempty"`         // important/urgent alert
	PollID     *int64  `json:"poll_id,omitempty"`
	Attachment string  `json:"attachment,omitempty"` // sticker id / emoji
}

type ReactionRecord struct {
	ID        int64  `json:"id"`
	MessageID int64  `json:"message_id"`
	UserID    int64  `json:"user_id"`
	Emoji     string `json:"emoji"`
}

type SessionRecord struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

type FileRecord struct {
	ID        int64  `json:"id"`
	OwnerID   int64  `json:"owner_id"`
	RoomID    int64  `json:"room_id"`
	FileName  string `json:"file_name"`
	StoredAs  string `json:"stored_as"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mime_type"`
	CreatedAt string `json:"created_at"`
}

// ReadRecord tracks per-room read position (read receipts)
type ReadRecord struct {
	ID        int64  `json:"id"`
	RoomID    int64  `json:"room_id"`
	UserID    int64  `json:"user_id"`
	LastRead  int64  `json:"last_read"` // last message id read
	UpdatedAt string `json:"updated_at"`
}

// PollRecord is a poll attached to a room
type PollRecord struct {
	ID        int64    `json:"id"`
	RoomID    int64    `json:"room_id"`
	CreatorID int64    `json:"creator_id"`
	Question  string   `json:"question"`
	Options   []string `json:"options"`
	CreatedAt string   `json:"created_at"`
	Closed    bool     `json:"closed"`
	ClosedAt  string   `json:"closed_at,omitempty"`
}

// PollVoteRecord is one user's vote on a poll option
type PollVoteRecord struct {
	ID       int64  `json:"id"`
	PollID   int64  `json:"poll_id"`
	UserID   int64  `json:"user_id"`
	Option   int    `json:"option"` // index into Options
	VotedAt  string `json:"voted_at"`
}

// InviteRecord is an invitation to a private room
type InviteRecord struct {
	ID       int64  `json:"id"`
	RoomID   int64  `json:"room_id"`
	UserID   int64  `json:"user_id"`   // invitee
	ByUserID int64  `json:"by_user_id"` // inviter
	CreatedAt string `json:"created_at"`
}

// PinRecord pins a message to the top of a room
type PinRecord struct {
	ID        int64  `json:"id"`
	RoomID    int64  `json:"room_id"`
	MessageID int64  `json:"message_id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// DepartmentRecord is an organizational section that groups rooms
type DepartmentRecord struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// Open loads (or creates) the store in dir. Migrates legacy single-file
// khan.db.json into the per-table layout automatically.
func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Store{dir: dir, data: &Data{
		Seq:         map[string]int64{},
		Users:       []UserRecord{},
		Rooms:       []RoomRecord{},
		RoomMembers: []RoomMemberRecord{},
		Messages:    []MessageRecord{},
		Reactions:   []ReactionRecord{},
		Sessions:    []SessionRecord{},
		Files:       []FileRecord{},
		Reads:       []ReadRecord{},
		Polls:       []PollRecord{},
		PollVotes:   []PollVoteRecord{},
		Invites:     []InviteRecord{},
		Pins:        []PinRecord{},
		Departments: []DepartmentRecord{},
	}}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) file(name string) string { return filepath.Join(s.dir, name) }

// tableFiles maps table -> filename
var tableFiles = map[string]string{
	"users":        "users.json",
	"rooms":        "rooms.json",
	"room_members": "room_members.json",
	"messages":     "messages.json",
	"reactions":    "reactions.json",
	"sessions":     "sessions.json",
	"files":        "files.json",
	"reads":        "reads.json",
	"polls":        "polls.json",
	"poll_votes":   "poll_votes.json",
	"invites":      "invites.json",
	"pins":         "pins.json",
	"departments":  "departments.json",
}

func (s *Store) load() error {
	legacy := s.file("khan.db.json")
	if _, err := os.Stat(legacy); err == nil {
		// Legacy single-file store exists → migrate to per-table files
		if err := s.migrateLegacy(legacy); err != nil {
			return err
		}
		return nil
	}

	// Load per-table files
	var firstErr error
	for table, fname := range tableFiles {
		data, err := os.ReadFile(s.file(fname))
		if err != nil {
			if os.IsNotExist(err) {
				continue // fresh table
			}
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if err := json.Unmarshal(data, tablePtr(s.data, table)); err != nil {
			return fmt.Errorf("corrupt %s: %w", fname, err)
		}
	}

	// Load seq.json
	if data, err := os.ReadFile(s.file("seq.json")); err == nil {
		var seq map[string]int64
		if err := json.Unmarshal(data, &seq); err == nil && seq != nil {
			s.data.Seq = seq
		}
	}
	if firstErr != nil {
		return firstErr
	}
	s.normalize()
	return nil
}

// tablePtr returns a pointer to the slice field for a table name
func tablePtr(d *Data, table string) any {
	switch table {
	case "users":
		return &d.Users
	case "rooms":
		return &d.Rooms
	case "room_members":
		return &d.RoomMembers
	case "messages":
		return &d.Messages
	case "reactions":
		return &d.Reactions
	case "sessions":
		return &d.Sessions
	case "files":
		return &d.Files
	case "reads":
		return &d.Reads
	case "polls":
		return &d.Polls
	case "poll_votes":
		return &d.PollVotes
	case "invites":
		return &d.Invites
	case "pins":
		return &d.Pins
	case "departments":
		return &d.Departments
	}
	return nil
}

// migrateLegacy converts khan.db.json → per-table files, then renames it to .bak
func (s *Store) migrateLegacy(legacy string) error {
	data, err := os.ReadFile(legacy)
	if err != nil {
		return err
	}
	var old Data
	if err := json.Unmarshal(data, &old); err != nil {
		return fmt.Errorf("corrupt legacy store: %w", err)
	}
	// Copy old data into new structure
	s.data.Users = old.Users
	s.data.Rooms = old.Rooms
	s.data.RoomMembers = old.RoomMembers
	s.data.Messages = old.Messages
	s.data.Reactions = old.Reactions
	s.data.Sessions = old.Sessions
	s.data.Files = old.Files
	if old.Seq != nil {
		s.data.Seq = old.Seq
	}
	// Empty new tables
	s.data.Reads = []ReadRecord{}
	s.data.Polls = []PollRecord{}
	s.data.PollVotes = []PollVoteRecord{}
	s.data.Invites = []InviteRecord{}
	s.data.Pins = []PinRecord{}
	s.data.Departments = []DepartmentRecord{}
	s.normalize()

	// Write all per-table files
	if err := s.saveAllLocked(); err != nil {
		return err
	}
	// Rename legacy out of the way (keep as backup)
	return os.Rename(legacy, legacy+".bak")
}

// normalize ensures no nil slices and seq map exists
func (s *Store) normalize() {
	if s.data.Seq == nil {
		s.data.Seq = map[string]int64{}
	}
	if s.data.Users == nil { s.data.Users = []UserRecord{} }
	if s.data.Rooms == nil { s.data.Rooms = []RoomRecord{} }
	if s.data.RoomMembers == nil { s.data.RoomMembers = []RoomMemberRecord{} }
	if s.data.Messages == nil { s.data.Messages = []MessageRecord{} }
	if s.data.Reactions == nil { s.data.Reactions = []ReactionRecord{} }
	if s.data.Sessions == nil { s.data.Sessions = []SessionRecord{} }
	if s.data.Files == nil { s.data.Files = []FileRecord{} }
	if s.data.Reads == nil { s.data.Reads = []ReadRecord{} }
	if s.data.Polls == nil { s.data.Polls = []PollRecord{} }
	if s.data.PollVotes == nil { s.data.PollVotes = []PollVoteRecord{} }
	if s.data.Invites == nil { s.data.Invites = []InviteRecord{} }
	if s.data.Pins == nil { s.data.Pins = []PinRecord{} }
	if s.data.Departments == nil { s.data.Departments = []DepartmentRecord{} }
}

// save persists every table to its own file (caller holds lock)
func (s *Store) saveAllLocked() error {
	for table, fname := range tableFiles {
		data, err := json.MarshalIndent(tablePtr(s.data, table), "", "  ")
		if err != nil {
			return err
		}
		tmp := s.file(fname + ".tmp")
		if err := os.WriteFile(tmp, data, 0o600); err != nil {
			return err
		}
		if err := os.Rename(tmp, s.file(fname)); err != nil {
			return err
		}
	}
	seqData, err := json.MarshalIndent(s.data.Seq, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.file("seq.json.tmp")
	if err := os.WriteFile(tmp, seqData, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.file("seq.json"))
}

// save atomically persists all tables
func (s *Store) save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.saveAllLocked()
}

// SaveNow forces a flush to disk
func (s *Store) SaveNow() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveAllLocked()
}

func (s *Store) saveLocked() error {
	return s.saveAllLocked()
}

// nextID returns the next auto-increment id for a table
func (s *Store) nextID(table string) int64 {
	s.data.Seq[table]++
	return s.data.Seq[table]
}

// Backup copies all data files to destDir (creates a timestamped folder)
func (s *Store) Backup(destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}
	sub := filepath.Join(destDir, fmt.Sprintf("khan-backup-%s", time.Now().Format("20060102-150405")))
	if err := os.MkdirAll(sub, 0o755); err != nil {
		return "", err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for table, fname := range tableFiles {
		data, err := json.MarshalIndent(tablePtr(s.data, table), "", "  ")
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(filepath.Join(sub, fname), data, 0o600); err != nil {
			return "", err
		}
	}
	seqData, err := json.MarshalIndent(s.data.Seq, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(sub, "seq.json"), seqData, 0o600); err != nil {
		return "", err
	}
	return sub, nil
}

// RestoreFrom replaces the in-memory data with the tables found in a backup
// directory (the timestamped folder created by Backup), then persists.
func (s *Store) RestoreFrom(backupDir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Build a fresh Data, load each per-table file from the backup dir.
	kept := Data{
		Seq: map[string]int64{},
		Users: []UserRecord{}, Rooms: []RoomRecord{},
		RoomMembers: []RoomMemberRecord{}, Messages: []MessageRecord{},
		Reactions: []ReactionRecord{}, Sessions: []SessionRecord{},
		Files: []FileRecord{}, Reads: []ReadRecord{},
		Polls: []PollRecord{}, PollVotes: []PollVoteRecord{},
		Invites: []InviteRecord{}, Pins: []PinRecord{},
		Departments: []DepartmentRecord{},
	}
	// Copy current seq so a partial restore can't collide with new ids.
	for k, v := range s.data.Seq {
		kept.Seq[k] = v
	}

	for table, fname := range tableFiles {
		data, err := os.ReadFile(filepath.Join(backupDir, fname))
		if err != nil {
			continue // missing file → keep empty table for that file
		}
		if err := json.Unmarshal(data, tablePtr(&kept, table)); err != nil {
			return fmt.Errorf("corrupt backup %s: %w", fname, err)
		}
	}
	// seq file too
	if data, err := os.ReadFile(filepath.Join(backupDir, "seq.json")); err == nil {
		var seq map[string]int64
		if err := json.Unmarshal(data, &seq); err == nil && seq != nil {
			kept.Seq = seq
		}
	}

	s.data = &kept
	return s.saveAllLocked()
}

// ---------- helpers ----------

func nowStr() string { return time.Now().Format(time.RFC3339) }

// Mu exposes the store mutex for repo transactions
func (s *Store) Mu() *sync.RWMutex { return &s.mu }

// Data exposes the in-memory dataset for repo transactions
func (s *Store) Data() *Data { return s.data }

// NextID returns the next auto-increment id for a table (caller holds lock)
func (s *Store) NextID(table string) int64 { return s.nextID(table) }

// SaveLocked persists while the caller holds the write lock
func (s *Store) SaveLocked() error { return s.saveLocked() }
