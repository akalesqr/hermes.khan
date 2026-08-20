package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"khan/internal/models"
	"khan/internal/repository"
	"khan/internal/service"
)

// MessageHandler serves /api/messages/*
type MessageHandler struct {
	msgs   *service.MessageService
	reads  *repository.ReadRepo
	pins   *repository.PinRepo
	polls  *repository.PollRepo
}

func NewMessageHandler(msgs *service.MessageService) *MessageHandler {
	return &MessageHandler{msgs: msgs}
}

// SetRepos wires extra repositories (called from bootstrap)
func (h *MessageHandler) SetRepos(reads *repository.ReadRepo, pins *repository.PinRepo, polls *repository.PollRepo) {
	h.reads = reads
	h.pins = pins
	h.polls = polls
}

// List returns messages for a room (paginated by ?before=X&limit=N)
func (h *MessageHandler) List(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	roomID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	before, _ := strconv.ParseInt(r.URL.Query().Get("before"), 10, 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	msgs, err := h.msgs.ListMessages(u, roomID, before, limit)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, msgs)
}

// Search searches messages across the user's rooms (?q=...)
func (h *MessageHandler) Search(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	q := r.URL.Query().Get("q")
	results, err := h.msgs.SearchMessages(u, q, 30)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if results == nil {
		results = []models.MessageView{}
	}
	writeJSON(w, http.StatusOK, results)
}

func (h *MessageHandler) Edit(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	msgID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	view, err := h.msgs.EditMessage(u, msgID, req.Text)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (h *MessageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	msgID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	if err := h.msgs.DeleteMessage(u, msgID); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// Forward forwards a message to another room
func (h *MessageHandler) Forward(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	msgID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req struct {
		RoomID int64 `json:"room_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if req.RoomID == 0 {
		writeErr(w, http.StatusBadRequest, "اتاق مقصد الزامی است")
		return
	}
	// Need original message to know its room
	orig, err := h.msgs.GetByID(msgID)
	if err != nil || orig == nil {
		writeErr(w, http.StatusNotFound, "پیام یافت نشد")
		return
	}
	view, err := h.msgs.Forward(u, orig.RoomID, req.RoomID, msgID)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, view)
}

// ToggleUrgent marks a message as urgent/important
func (h *MessageHandler) ToggleUrgent(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	msgID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	msg, err := h.msgs.ToggleUrgent(u, msgID)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"urgent": msg.Urgent})
}

// Pins lists pinned messages for a room
func (h *MessageHandler) Pins(w http.ResponseWriter, r *http.Request) {
	if h.pins == nil {
		writeJSON(w, http.StatusOK, []interface{}{})
		return
	}
	roomID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	pins, _ := h.pins.ListForRoom(roomID)
	writeJSON(w, http.StatusOK, pins)
}

type reactionReq struct {
	Emoji string `json:"emoji"`
}

func (h *MessageHandler) AddReaction(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	msgID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req reactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if err := h.msgs.AddReaction(u, msgID, req.Emoji); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *MessageHandler) RemoveReaction(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	msgID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req reactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if err := h.msgs.RemoveReaction(u, msgID, req.Emoji); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
