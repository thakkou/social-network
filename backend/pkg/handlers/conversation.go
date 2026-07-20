package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "01social/pkg/db/sqlite"
	"01social/pkg/middlewares"
	"01social/pkg/utilities"
	"01social/pkg/ws"
)

type SendMessageRequest struct {
	ReceiverID     int    `json:"receiver_id"`
	Text           string `json:"text"`
	ConversationID *int   `json:"conversation_id"`
}

// SendMessage handles sending a DM and triggering real-time WS notifications.
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
	if req.ReceiverID == 0 || req.Text == "" {
		utilities.WriteJSON(w, http.StatusBadRequest, "missing fields", nil)
		return
	}

	if senderID == req.ReceiverID {
		utilities.WriteJSON(w, http.StatusBadRequest, "cannot send message to yourself", nil)
		return
	}

	tx, err := db.Database.Begin()
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "db error", nil)
		return
	}
	defer tx.Rollback()

	var conversationID int
	var isNew bool

	if req.ConversationID != nil {
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

	// Fetch sender profile details before commit
	senderProfile, err := Repos.User.GetByID(senderID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "user profile fetch error", nil)
		return
	}

	if err := tx.Commit(); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "commit failed", nil)
		return
	}

	// WebSocket push notifications
	wsPayload := map[string]interface{}{
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
			"conversation_id": conversationID,
			"message_id":      messageID,
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
	fmt.Println("get the conversation")
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
	fmt.Println("convType", convType, id)

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
