package handler

import (
	"encoding/json"
	"strings"
	"net/http"

	"khan/internal/repository"
	"khan/internal/service"
)

// RoomHandler serves /api/rooms/*
type RoomHandler struct {
	rooms  *service.RoomService
	invites *repository.InviteRepo
	roomRep *repository.RoomRepo
}

func NewRoomHandler(rooms *service.RoomService) *RoomHandler {
	return &RoomHandler{rooms: rooms}
}

// SetRepos wires extra repositories (called from bootstrap)
func (h *RoomHandler) SetRepos(invites *repository.InviteRepo, roomRep *repository.RoomRepo) {
	h.invites = invites
	h.roomRep = roomRep
}

// PublicRooms lists all public rooms
func (h *RoomHandler) PublicRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomRep.ListPublic()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای داخلی")
		return
	}
	writeJSON(w, http.StatusOK, rooms)
}

// PrivateRooms lists all private rooms
func (h *RoomHandler) PrivateRooms(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	rooms, err := h.roomRep.ListPrivateRooms(u.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای داخلی")
		return
	}
	writeJSON(w, http.StatusOK, rooms)
}

// Invite invites a user to a private room
func (h *RoomHandler) Invite(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	id, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if err := h.rooms.InviteUser(u, id, req.UserID); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// Departments lists departments
func (h *RoomHandler) Departments(w http.ResponseWriter, r *http.Request) {
	if h.roomRep == nil {
		writeJSON(w, http.StatusOK, []interface{}{})
		return
	}
	depts, _ := h.roomRep.ListDepartments()
	writeJSON(w, http.StatusOK, depts)
}

// CreateDepartment creates a department (admin+)
func (h *RoomHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	if !u.CanManageUsers() {
		writeErr(w, http.StatusForbidden, "دسترسی غیرمجاز")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if req.Name == "" {
		writeErr(w, http.StatusBadRequest, "نام بخش الزامی است")
		return
	}
	dept, err := h.roomRep.CreateDepartment(req.Name)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, dept)
}

// DeleteDepartment removes a department (admin+)
func (h *RoomHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	if !u.CanManageUsers() {
		writeErr(w, http.StatusForbidden, "دسترسی غیرمجاز")
		return
	}
	id, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	if err := h.roomRep.DeleteDepartment(id); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	rooms, err := h.rooms.ListRooms(u)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "خطای داخلی")
		return
	}
	writeJSON(w, http.StatusOK, rooms)
}

type createRoomReq struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Department  int64  `json:"department"`
	Description string `json:"description,omitempty"`
	Topic       string `json:"topic,omitempty"`
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	var req createRoomReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		writeErr(w, http.StatusBadRequest, "نام اتاق الزامی است")
		return
	}
	if len(req.Name) > 100 {
		writeErr(w, http.StatusBadRequest, "نام اتاق حداکثر ۱۰۰ کاراکتر")
		return
	}
	// Basic XSS filter for room name
	req.Name = strings.ReplaceAll(req.Name, "<", "")
	req.Name = strings.ReplaceAll(req.Name, ">", "")
	if len(req.Description) > 500 {
		writeErr(w, http.StatusBadRequest, "توضیحات حداکثر ۵۰۰ کاراکتر")
		return
	}
	if len(req.Topic) > 200 {
		writeErr(w, http.StatusBadRequest, "موضوع حداکثر ۲۰۰ کاراکتر")
		return
	}
	room, err := h.rooms.CreateRoom(u, req.Name, req.Type, req.Department)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, room)
}

func (h *RoomHandler) Join(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	id, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	if err := h.rooms.JoinRoom(u, id); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *RoomHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	id, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if err := h.rooms.AddMember(u, id, req.UserID); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *RoomHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	roomID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	userID, err := pathInt64(r, "uid")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	if err := h.rooms.RemoveMember(u, roomID, userID); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *RoomHandler) Members(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	id, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	members, err := h.rooms.Members(u, id)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	sanitized := make([]interface{}, 0, len(members))
	for i := range members {
		sanitized = append(sanitized, members[i].Sanitized())
	}
	writeJSON(w, http.StatusOK, sanitized)
}

func (h *RoomHandler) Rename(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	id, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if err := h.rooms.RenameRoom(u, id, req.Name); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// StartDM creates or finds a DM room with another user
func (h *RoomHandler) StartDM(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	otherID, err := pathInt64(r, "uid")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	roomID, err := h.rooms.FindOrCreateDM(u, otherID)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"room_id": roomID})
}
