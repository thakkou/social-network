package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	db "01social/pkg/db/sqlite"
	"01social/pkg/middlewares"
	"01social/pkg/utilities"
	"01social/pkg/ws"
)

// SendMessage handles sending a DM and triggering real-time WS notifications.

type SendMessageRequest struct {
	Type           string `json:"type"`            // "direct" | "group"
	Text           string `json:"text"`            // Message body
	ReceiverID     int    `json:"receiver_id"`     // Required for direct messages
	ConversationID *int   `json:"conversation_id"` // Optional for direct messages
	GroupID        int    `json:"group_id"`        // Required for group messages
}

// SendMessage handles sending both direct and group messages and dispatching WS notifications.
func SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	senderID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	req.Text = strings.TrimSpace(req.Text)
	if req.Text == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "message text cannot be empty", nil)
		return
	}

	// Default to direct if type is omitted for backward compatibility
	if req.Type == "" {
		req.Type = "direct"
	}

	tx, err := db.Database.Begin()
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "db error", nil)
		return
	}
	defer tx.Rollback()

	// Fetch sender profile details before commit for WS metadata
	senderProfile, err := Repos.User.GetByID(senderID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "user profile fetch error", nil)
		return
	}

	switch req.Type {
	case "direct":
		if req.ReceiverID == 0 {
			utilities.WriteJSON(w, http.StatusBadRequest, "missing receiver_id for direct message", nil)
			return
		}

		if senderID == req.ReceiverID {
			utilities.WriteJSON(w, http.StatusBadRequest, "cannot send message to yourself", nil)
			return
		}

		// Permission check: at least one user must be following the other (accepted follow)
		senderFollowsReceiver, err := Repos.Follow.IsFollowing(senderID, req.ReceiverID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to check follow status", nil)
			return
		}
		receiverFollowsSender, err := Repos.Follow.IsFollowing(req.ReceiverID, senderID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to check follow status", nil)
			return
		}

		if !senderFollowsReceiver && !receiverFollowsSender {
			utilities.WriteJSON(w, http.StatusForbidden, "you can only message users who follow you or you follow", nil)
			return
		}

		var conversationID int
		var isNew bool

		if req.ConversationID != nil && *req.ConversationID > 0 {
			conversationID = *req.ConversationID
		} else {
			convID, newlyCreated, err := Repos.Conversation.FindOrCreateDirectConversation(tx, senderID, req.ReceiverID)
			if err != nil {
				utilities.WriteJSON(w, http.StatusInternalServerError, "failed to resolve conversation", nil)
				return
			}
			conversationID = convID
			isNew = newlyCreated
		}

		messageID, err := Repos.Conversation.SaveMessage(tx, conversationID, senderID, req.Text)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to save message", nil)
			return
		}

		if err := tx.Commit(); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "commit failed", nil)
			return
		}

		// Dispatch WebSocket push notifications for direct message
		wsPayload := map[string]interface{}{
			"type":              "direct",
			"isNewConversation": isNew,
			"nickname":          senderProfile.Nickname,
			"conversation_id":   conversationID,
			"message_id":        messageID,
			"sender_id":         senderID,
			"text":              req.Text,
		}

		wsPayload["isMine"] = false
		ws.NotifyUser(strconv.Itoa(req.ReceiverID), "new_message", wsPayload)

		wsPayload["isMine"] = true
		ws.NotifyUser(strconv.Itoa(senderID), "new_message", wsPayload)

		utilities.WriteJSON(
			w,
			http.StatusOK,
			"message sent success",
			map[string]interface{}{
				"type":            "direct",
				"conversation_id": conversationID,
				"message_id":      messageID,
			},
		)

	case "group":
		if req.GroupID == 0 {
			utilities.WriteJSON(w, http.StatusBadRequest, "missing group_id for group message", nil)
			return
		}

		messageID, err := Repos.Conversation.SaveGroupMessage(tx, req.GroupID, senderID, req.Text)
		if err == sql.ErrNoRows {
			utilities.WriteJSON(w, http.StatusForbidden, "you are not a member of this group", nil)
			return
		} else if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "failed to save group message", nil)
			return
		}

		if err := tx.Commit(); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "commit failed", nil)
			return
		}

		// Fetch group members and broadcast via WebSockets
		memberIDs, err := Repos.Conversation.GetGroupMemberIDs(req.GroupID)
		if err == nil {
			for _, memberID := range memberIDs {
				wsPayload := map[string]interface{}{
					"type":       "group",
					"group_id":   req.GroupID,
					"message_id": messageID,
					"sender_id":  senderID,
					"nickname":   senderProfile.Nickname,
					"text":       req.Text,
					"isMine":     memberID == senderID,
				}
				ws.NotifyUser(strconv.Itoa(memberID), "new_group_message", wsPayload)
			}
		}

		utilities.WriteJSON(
			w,
			http.StatusOK,
			"group message sent success",
			map[string]interface{}{
				"type":       "group",
				"group_id":   req.GroupID,
				"message_id": messageID,
			},
		)

	default:
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid message type (must be 'direct' or 'group')", nil)
	}
}

type SendGroupMessageRequest struct {
	GroupID int    `json:"group_id"`
	Text    string `json:"text"`
}

// SendGroupMessage handles sending a group chat message and broadcasting via WebSockets.
func SendGroupMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	senderID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req SendGroupMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	req.Text = strings.TrimSpace(req.Text)
	if req.GroupID == 0 || req.Text == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "missing group_id or text", nil)
		return
	}

	tx, err := db.Database.Begin()
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "db error", nil)
		return
	}
	defer tx.Rollback()

	// Save group message
	messageID, err := Repos.Conversation.SaveGroupMessage(tx, req.GroupID, senderID, req.Text)
	if err == sql.ErrNoRows {
		utilities.WriteJSON(w, http.StatusForbidden, "you are not a member of this group", nil)
		return
	} else if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "failed to save message", nil)
		return
	}

	// Fetch sender details for real-time client metadata
	senderProfile, err := Repos.User.GetByID(senderID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "user profile fetch error", nil)
		return
	}

	if err := tx.Commit(); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "commit failed", nil)
		return
	}

	// Fetch group members to broadcast WS notification
	memberIDs, err := Repos.Conversation.GetGroupMemberIDs(req.GroupID)
	if err == nil {
		for _, memberID := range memberIDs {
			wsPayload := map[string]interface{}{
				"group_id":   req.GroupID,
				"message_id": messageID,
				"sender_id":  senderID,
				"nickname":   senderProfile.Nickname,
				"text":       req.Text,
				"isMine":     memberID == senderID,
				"type":       "group",
			}
			ws.NotifyUser(strconv.Itoa(memberID), "new_group_message", wsPayload)
		}
	}

	utilities.WriteJSON(
		w,
		http.StatusOK,
		"group message sent success",
		map[string]interface{}{
			"group_id":   req.GroupID,
			"message_id": messageID,
		},
	)
}

// GetConversation fetches the active user's conversation feed.
func GetConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 || limit > 30 {
		limit = 30
	}

	directItems, groupItems, err := Repos.Conversation.FetchConversationFeed(userID, limit, offset)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "ok", map[string]interface{}{
		"direct": directItems,
		"groups": groupItems,
	})
}

// GetConversationByID handles message retrieval for direct chats or group chats.
func GetConversationByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	convType := r.PathValue("type") // "direct" or "group"
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	switch convType {
	case "direct":
		messages, err := Repos.Conversation.FetchDirectMessages(id, userID, limit, offset)
		if err == sql.ErrNoRows {
			utilities.WriteJSON(w, http.StatusForbidden, "not allowed", nil)
			return
		} else if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "db error", nil)
			return
		}

		utilities.WriteJSON(w, http.StatusOK, "ok", map[string]interface{}{
			"type":     convType,
			"id":       id,
			"messages": messages,
		})

	case "group":
		messages, err := Repos.Conversation.FetchGroupMessages(id, userID, limit, offset)
		if err == sql.ErrNoRows {
			utilities.WriteJSON(w, http.StatusForbidden, "not allowed", nil)
			return
		} else if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "db error", nil)
			return
		}

		utilities.WriteJSON(w, http.StatusOK, "ok", map[string]interface{}{
			"type":     convType,
			"id":       id,
			"messages": messages,
		})

	default:
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid conversation type", nil)
	}
}
