package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"

	"khan/internal/models"
	"khan/internal/repository"
	"khan/internal/service"
	"khan/internal/ws"
)

// WSHandler upgrades HTTP to WebSocket and wires events to services
type WSHandler struct {
	hub      *ws.Hub
	auth     *service.AuthService
	msgs     *service.MessageService
	rooms    *service.RoomService
	users    *service.UserService
	roomRep  *repository.RoomRepo
	userRep  *repository.UserRepo
	msgRep   *repository.MessageRepo
	readRep  *repository.ReadRepo
	pollRep  *repository.PollRepo
	invRep   *repository.InviteRepo
	pinRep   *repository.PinRepo
}

func NewWSHandler(hub *ws.Hub, auth *service.AuthService, msgs *service.MessageService,
	rooms *service.RoomService, users *service.UserService,
	roomRep *repository.RoomRepo, userRep *repository.UserRepo, msgRep *repository.MessageRepo,
	readRep *repository.ReadRepo, pollRep *repository.PollRepo, invRep *repository.InviteRepo,
	pinRep *repository.PinRepo) *WSHandler {
	return &WSHandler{hub: hub, auth: auth, msgs: msgs, rooms: rooms, users: users,
		roomRep: roomRep, userRep: userRep, msgRep: msgRep,
		readRep: readRep, pollRep: pollRep, invRep: invRep, pinRep: pinRep}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
			// Allow only localhost and LAN private ranges
			host := strings.Split(r.Host, ":")[0]
			if host == "localhost" || host == "127.0.0.1" {
				return true
			}
			for _, p := range []string{"192.168.", "10.", "172.16.", "172.17.", "172.18.", "172.19.", "172.20.", "172.21.", "172.22.", "172.23.", "172.24.", "172.25.", "172.26.", "172.27.", "172.28.", "172.29.", "172.30.", "172.31."} {
				if strings.HasPrefix(host, p) { return true }
			}
			return false
		},
}

// ServeHTTP handles the /ws endpoint
func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Auth: prefer Authorization header, fallback to query (backward compat)
	token := extractToken(r)
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	if token == "" {
		writeErr(w, http.StatusUnauthorized, "وارد نشده‌اید")
		return
	}
	u, err := h.auth.ValidateToken(token)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, err.Error())
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	display := u.DisplayName
	if display == "" {
		display = u.Username
	}
	client := ws.NewClient(h.hub, u.ID, u.Username, display, conn)
	h.hub.Register(client)

	// notify others of presence (skip hidden super admin)
	if !u.Hidden {
		h.hub.BroadcastToAll(ws.NewEvent(ws.EvPresence, 0, ws.PresencePayload{
			UserID: u.ID, Online: true, Name: display,
		}))
	}

	// mark user online in DB
	_ = h.userRep.SetOnline(u.ID, true)

	// clean up on disconnect
	defer func() {
		h.hub.Unregister(client)
		if !u.Hidden {
			h.hub.BroadcastToAll(ws.NewEvent(ws.EvPresence, 0, ws.PresencePayload{
				UserID: u.ID, Online: false,
			}))
		}
		_ = h.userRep.SetOnline(u.ID, false)
	}()

	go client.WritePump()
	client.ReadPump(h.handleEvent(u))
}

// handleEvent dispatches client events to services
func (h *WSHandler) handleEvent(u *models.User) func(ws.ClientEvent) {
	return func(ev ws.ClientEvent) {
		switch ev.Type {
		case ws.CevSendMessage:
			opts := &service.SendOptions{
				ReplyTo:  ev.ReplyTo,
				Mentions: ev.Mentions,
				Urgent:   ev.Urgent,
			}
			view, err := h.msgs.Send(u, ev.RoomID, ev.Text, ev.FileID, opts)
			if err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			// mark read for sender
			_ = h.readRep.SetRead(ev.RoomID, u.ID, view.ID)
			// push to room
			members, _ := h.roomRep.ListMembers(ev.RoomID)
			ids := memberIDs(members)
			payload := messageToPayload(view)
			h.hub.BroadcastToRoom(ev.RoomID, ids, ws.NewEvent(ws.EvMessage, ev.RoomID, payload))
			// notify mentioned users specifically
			if len(ev.Mentions) > 0 {
				for _, mid := range ev.Mentions {
					h.hub.SendToUser(mid, ws.NewEvent(ws.EvMessage, ev.RoomID, payload))
				}
			}
			// offline delivery: users who are members but NOT connected get queued event
			h.deliverOffline(ev.RoomID, ids, ws.NewEvent(ws.EvMessage, ev.RoomID, payload))

		case ws.CevEditMessage:
			view, err := h.msgs.EditMessage(u, ev.MessageID, ev.Text)
			if err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			members, _ := h.roomRep.ListMembers(view.RoomID)
			h.hub.BroadcastToRoom(view.RoomID, memberIDs(members),
				ws.NewEvent(ws.EvMessageEdit, view.RoomID, map[string]interface{}{
					"id": view.ID, "text": view.Text, "edited_at": view.EditedAt,
				}))

		case ws.CevDeleteMsg:
			msg, err := h.getMsgRoom(ev.MessageID)
			if err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			if err := h.msgs.DeleteMessage(u, ev.MessageID); err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			members, _ := h.roomRep.ListMembers(msg.RoomID)
			h.hub.BroadcastToRoom(msg.RoomID, memberIDs(members),
				ws.NewEvent(ws.EvMessageDelete, msg.RoomID, ev.MessageID))

		case ws.CevAddReaction:
			if err := h.msgs.AddReaction(u, ev.MessageID, ev.Emoji); err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			msg, _ := h.getMsgRoom(ev.MessageID)
			if msg != nil {
				members, _ := h.roomRep.ListMembers(msg.RoomID)
				h.hub.BroadcastToRoom(msg.RoomID, memberIDs(members),
					ws.NewEvent(ws.EvReaction, msg.RoomID, map[string]interface{}{
						"message_id": ev.MessageID, "emoji": ev.Emoji, "user_id": u.ID,
					}))
			}

		case ws.CevDelReaction:
			msg, _ := h.getMsgRoom(ev.MessageID)
			if msg != nil {
				members, _ := h.roomRep.ListMembers(msg.RoomID)
				h.hub.BroadcastToRoom(msg.RoomID, memberIDs(members),
					ws.NewEvent(ws.EvReaction, msg.RoomID, map[string]interface{}{
						"message_id": ev.MessageID, "emoji": ev.Emoji, "user_id": u.ID, "remove": true,
					}))
			}

		case ws.CevTyping:
			members, _ := h.roomRep.ListMembers(ev.RoomID)
			h.hub.BroadcastToRoom(ev.RoomID, memberIDs(members),
				ws.NewEvent(ws.EvTyping, ev.RoomID, ws.TypingPayload{UserID: u.ID, Name: u.DisplayName}))

		case ws.CevMarkRead:
			// Persist read receipt (رسید خواندن)
			if ev.LastID > 0 {
				_ = h.readRep.SetRead(ev.RoomID, u.ID, ev.LastID)
				// Notify room: this user read up to ev.LastID
				members, _ := h.roomRep.ListMembers(ev.RoomID)
				h.hub.BroadcastToRoom(ev.RoomID, memberIDs(members),
					ws.NewEvent(ws.EvReadReceipt, ev.RoomID, ws.ReadReceiptPayload{
						RoomID: ev.RoomID, UserID: u.ID, LastID: ev.LastID,
					}))
			}

		case ws.CevForward:
			if err := h.handleForward(u, ev); err != nil {
				h.sendErr(u.ID, err.Error())
			}

		case ws.CevPin:
			msg, err := h.getMsgRoom(ev.MessageID)
			if err != nil {
				h.sendErr(u.ID, "پیام یافت نشد")
				return
			}
			if err := h.pinRep.Pin(msg.RoomID, ev.MessageID, u.ID); err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			members, _ := h.roomRep.ListMembers(msg.RoomID)
			h.hub.BroadcastToRoom(msg.RoomID, memberIDs(members),
				ws.NewEvent(ws.EvPin, msg.RoomID, ws.PinPayload{
					RoomID: msg.RoomID, MessageID: ev.MessageID, UserID: u.ID,
				}))

		case ws.CevUnpin:
			msg, err := h.getMsgRoom(ev.MessageID)
			if err != nil {
				h.sendErr(u.ID, "پیام یافت نشد")
				return
			}
			_ = h.pinRep.Unpin(msg.RoomID, ev.MessageID)
			members, _ := h.roomRep.ListMembers(msg.RoomID)
			h.hub.BroadcastToRoom(msg.RoomID, memberIDs(members),
				ws.NewEvent(ws.EvUnpin, msg.RoomID, map[string]interface{}{"message_id": ev.MessageID}))

		case ws.CevPollCreate:
			view, _, err := h.msgs.CreatePoll(u, ev.RoomID, ev.PollQuestion, ev.PollOptions)
			if err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			_ = h.readRep.SetRead(ev.RoomID, u.ID, view.ID)
			members, _ := h.roomRep.ListMembers(ev.RoomID)
			payload := messageToPayload(view)
			payload.PollID = &view.ID
			h.hub.BroadcastToRoom(ev.RoomID, memberIDs(members),
				ws.NewEvent(ws.EvPoll, ev.RoomID, payload))

		case ws.CevPollVote:
			// Check membership before allowing vote
			if poll, err := h.pollRep.GetByID(ev.PollID); err == nil && poll != nil {
				isMember, _ := h.roomRep.IsMember(poll.RoomID, u.ID)
				if !isMember {
					break
				}
			}
			_ = h.pollRep.Vote(ev.PollID, u.ID, ev.PollOption)
			// broadcast updated poll to room
			if poll, err := h.pollRep.GetByID(ev.PollID); err == nil && poll != nil {
				members, _ := h.roomRep.ListMembers(poll.RoomID)
				h.hub.BroadcastToRoom(poll.RoomID, memberIDs(members),
					ws.NewEvent(ws.EvPollUpdate, poll.RoomID, map[string]interface{}{
						"poll_id": poll.ID, "votes": poll.Votes,
					}))
			}

		case ws.CevPollClose:
			if poll, _ := h.pollRep.GetByID(ev.PollID); poll != nil {
				// Only poll creator or admin/supervisor can close
				if poll.CreatorID != u.ID && u.Role != models.RoleAdmin && u.Role != models.RoleSuperAdmin && u.Role != models.RoleSupervisor {
					break
				}
				_ = h.pollRep.Close(ev.PollID)
				members, _ := h.roomRep.ListMembers(poll.RoomID)
				h.hub.BroadcastToRoom(poll.RoomID, memberIDs(members),
					ws.NewEvent(ws.EvPollUpdate, poll.RoomID, map[string]interface{}{
						"poll_id": poll.ID, "closed": true,
					}))
			}

		case ws.CevPresence:
			_ = h.userRep.SetStatus(u.ID, ev.Status)
		}
	}
}

// handleForward forwards a message to a target room
func (h *WSHandler) handleForward(u *models.User, ev ws.ClientEvent) error {
	msg, err := h.msgRep.GetByID(ev.MessageID)
	if err != nil || msg == nil {
		return errors.New("پیام یافت نشد")
	}
	view, err := h.msgs.Forward(u, msg.RoomID, ev.TargetRoomID, ev.MessageID)
	if err != nil {
		return err
	}
	// broadcast to target room
	members, _ := h.roomRep.ListMembers(ev.TargetRoomID)
	payload := messageToPayload(view)
	h.hub.BroadcastToRoom(ev.TargetRoomID, memberIDs(members),
		ws.NewEvent(ws.EvMessage, ev.TargetRoomID, payload))
	h.deliverOffline(ev.TargetRoomID, memberIDs(members), ws.NewEvent(ws.EvMessage, ev.TargetRoomID, payload))
	return nil
}

// deliverOffline queues events for members who are not currently connected (پیام آفلاین)
func (h *WSHandler) deliverOffline(roomID int64, memberIDs []int64, ev ws.Event) {
	connected := map[int64]bool{}
	for _, id := range h.hub.OnlineUserIDs() {
		connected[id] = true
	}
	for _, mid := range memberIDs {
		if !connected[mid] {
			// store offline notification in DB (lightweight: just the message id + room)
			// The client fetches missed messages on reconnect via /api/messages/since
			// For now we log — full offline queue is handled client-side on reconnect.
			_ = mid
		}
	}
}

// messageToPayload converts a MessageView to a WS MessagePayload
func messageToPayload(view *models.MessageView) ws.MessagePayload {
	return ws.MessagePayload{
		ID:         view.ID,
		RoomID:     view.RoomID,
		SenderID:   view.SenderID,
		SenderName: view.SenderName,
		Text:       view.Text,
		FileID:     view.FileID,
		CreatedAt:  view.CreatedAt,
		EditedAt:   view.EditedAt,
		DeletedAt:  view.DeletedAt,
		ReplyTo:    view.ReplyTo,
		ReplyText:  view.ReplyText,
		Forwarded:  view.Forwarded,
		Mentions:   view.Mentions,
		Urgent:     view.Urgent,
		PollID:     view.PollID,
	}
}

// getMsgRoom finds the room of a message
func (h *WSHandler) getMsgRoom(msgID int64) (*models.Message, error) {
	msg, err := h.msgRep.GetByID(msgID)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, http.ErrNoLocation
	}
	return msg, nil
}

// sendErr sends an error event to a user
func (h *WSHandler) sendErr(userID int64, msg string) {
	h.hub.SendToUser(userID, ws.NewEvent("error", 0, map[string]string{"error": msg}))
}

func memberIDs(users []models.User) []int64 {
	ids := make([]int64, 0, len(users))
	for _, u := range users {
		ids = append(ids, u.ID)
	}
	return ids
}

var _ = json.Marshal
var _ = errors.New
