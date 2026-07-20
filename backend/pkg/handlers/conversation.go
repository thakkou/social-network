package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"01social/pkg/db/sqlite"
	"01social/pkg/utilities"
	"01social/pkg/ws"
)

type SendMessageRequest struct {
	ReceiverID     int    `json:"receiver_id"`
	Text           string `json:"text"`
	ConversationID *int   `json:"conversation_id"`
}

type User struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type Profile struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type ConversationPreview struct {
	ConversationID *int    `json:"conversationId"`
	Date           *string `json:"date"`
	LastMessage    *string `json:"lastMessage"`
	Status         string  `json:"status"`

	LastSeen    *string `json:"lastSeen"`
	UnreadCount int     `json:"unreadCount"`
	LastSender  string  `json:"lastSender"`
}

type UserFeedItem struct {
	Profile      Profile             `json:"profile"`
	Conversation ConversationPreview `json:"conversation"`
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	// -------------------------
	// Get sender from session
	// -------------------------
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// fmt.Println("[AUTH] missing session cookie")
		utilities.WriteJSON(w, 401, "unauthorized", nil)
		return
	}

	senderID, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		// fmt.Println("[AUTH] invalid session:", err)
		utilities.WriteJSON(w, 401, "unauthorized", nil)
		return
	}

	// fmt.Printf("[AUTH] sender=%d\n", senderID)

	// -------------------------
	// Decode request
	// -------------------------
	var req SendMessageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// fmt.Println("[REQUEST] decode error:", err)
		utilities.WriteJSON(w, 400, "invalid request body", nil)
		return
	}
	req.Text = strings.TrimSpace(req.Text)

	// -------------------------
	// Validate
	// -------------------------
	if req.ReceiverID == 0 || req.Text == "" {
		// fmt.Println("[VALIDATION] missing fields")
		utilities.WriteJSON(w, 400, "missing fields", nil)
		return
	}

	if senderID == req.ReceiverID {
		// fmt.Println("[VALIDATION] user tried to message himself")
		utilities.WriteJSON(w, 400, "cannot send message to yourself", nil)
		return
	}

	// -------------------------
	// Normalize pair
	// -------------------------
	user1 := senderID
	user2 := req.ReceiverID

	if user1 > user2 {
		user1, user2 = user2, user1
	}

	// -------------------------
	// Start transaction
	// -------------------------
	tx, err := sqlite.DB().Begin()
	if err != nil {
		// fmt.Println("[DB] begin transaction error:", err)
		utilities.WriteJSON(w, 500, "db error", nil)
		return
	}
	defer tx.Rollback()

	var conversationID int
	isNewConversation := false

	// -------------------------
	// CASE 1: conversation_id provided
	// -------------------------
	if req.ConversationID != nil {
		conversationID = *req.ConversationID

		var exists int
		err := tx.QueryRow(`
			SELECT id
			FROM CONVERSATIONS
			WHERE id = ?
			AND user1_id = ?
			AND user2_id = ?
		`,
			conversationID,
			user1,
			user2,
		).Scan(&exists)
		if err != nil {
			// fmt.Println("[CONVERSATION] invalid conversation:", err)
			utilities.WriteJSON(w, http.StatusNotFound, "conversation not found", nil)
			return
		}

	} else {
		// -------------------------
		// CASE 2: Find or create conversation
		// -------------------------
		err := tx.QueryRow(`
			SELECT id
			FROM CONVERSATIONS
			WHERE user1_id = ?
			AND user2_id = ?
		`,
			user1,
			user2,
		).Scan(&conversationID)

		switch {
		case err == sql.ErrNoRows:
			res, err := tx.Exec(`
				INSERT INTO CONVERSATIONS (user1_id, user2_id)
				VALUES (?, ?)
			`,
				user1,
				user2,
			)
			if err != nil {
				// fmt.Println("[CONVERSATION] create error:", err)
				utilities.WriteJSON(w, 500, "failed to create conversation", nil)
				return
			}

			id, err := res.LastInsertId()
			if err != nil {
				// fmt.Println("[CONVERSATION] last insert id error:", err)
				utilities.WriteJSON(w, 500, "failed to create conversation", nil)
				return
			}

			conversationID = int(id)
			isNewConversation = true
			// fmt.Printf("[CONVERSATION] created id=%d\n", conversationID)

		case err != nil:
			// fmt.Println("[CONVERSATION] lookup error:", err)
			utilities.WriteJSON(w, 500, "db error", nil)
			return

		default:
			// fmt.Printf("[CONVERSATION] found id=%d\n", conversationID)
		}
	}

	// -------------------------
	// Insert message
	// -------------------------
	result, err := tx.Exec(`
		INSERT INTO MESSAGES (conversation_id, sender_id, text)
		VALUES (?, ?, ?)
	`,
		conversationID,
		senderID,
		req.Text,
	)
	if err != nil {
		// fmt.Println("[MESSAGE] insert error:", err)
		utilities.WriteJSON(w, 500, "failed to send message", nil)
		return
	}

	messageID, _ := result.LastInsertId()

	// -------------------------
	// Update conversation preview
	// -------------------------
	_, err = tx.Exec(`
		UPDATE CONVERSATIONS
		SET last_message = ?, last_message_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`,
		req.Text,
		conversationID,
	)
	if err != nil {
		// fmt.Println("[CONVERSATION] update preview error:", err)
		utilities.WriteJSON(w, 500, "failed to update conversation", nil)
		return
	}

	// -------------------------
	// Get sender nickname
	// ⚠️ MUST happen before commit — tx is unusable afterwards
	// -------------------------
	var nickname string
	err = tx.QueryRow(`
		SELECT nickname
		FROM USERS
		WHERE id = ?
	`, senderID).Scan(&nickname)
	if err != nil {
		// fmt.Println("[USER] failed to get nickname:", err)
		if err == sql.ErrNoRows {
			utilities.WriteJSON(w, 404, "user not found", nil)
		} else {
			utilities.WriteJSON(w, 500, "db error", nil)
		}
		return
	}

	// -------------------------
	// Commit
	// -------------------------
	if err := tx.Commit(); err != nil {
		// fmt.Println("[DB] commit error:", err)
		utilities.WriteJSON(w, 500, "commit failed", nil)
		return
	}

	// -------------------------
	// Notify both users over websocket
	// -------------------------
	ws.NotifyUser(
		strconv.Itoa(req.ReceiverID),
		"new_message",
		map[string]interface{}{
			"isMine":            false,
			"isNewConversation": isNewConversation,
			"nickname":          nickname,
			"conversation_id":   conversationID,
			"message_id":        messageID,
			"sender_id":         senderID,
			"text":              req.Text,
		},
	)

	ws.NotifyUser(
		strconv.Itoa(senderID),
		"new_message",
		map[string]interface{}{
			"isMine":            true,
			"isNewConversation": isNewConversation,
			"conversation_id":   conversationID,
			"message_id":        messageID,
			"sender_id":         senderID,
			"text":              req.Text,
		},
	)

	// -------------------------
	// Respond
	// -------------------------
	utilities.WriteJSON(
		w,
		200,
		"message sent success",
		map[string]interface{}{
			"conversation_id": conversationID,
			"message_id":      messageID,
		},
	)
}

// -------------------------------------------------------------------------
// UNIFIED FEED TYPES
// -------------------------------------------------------------------------

type ConversationFeedItem struct {
	Type          string  `json:"type"` // "direct" | "group"
	ID            int     `json:"id"`   // conversation_id or group_id
	DisplayName   string  `json:"display_name"`
	Avatar        *string `json:"avatar,omitempty"`
	LastMessage   *string `json:"last_message"`
	LastMessageAt *string `json:"last_message_at"`
	UnreadCount   int     `json:"unread_count"`

	// position in the combined (groups + direct) recency ranking —
	// use this to interleave the two arrays client-side if needed,
	// lower = more recent
	Rank int `json:"rank"`

	// direct-only
	OtherUserID *int `json:"other_user_id,omitempty"`

	// group-only
	MemberCount *int `json:"member_count,omitempty"`
}

type feedRef struct {
	id            int
	kind          string // "direct" | "group"
	lastMessageAt sql.NullString
}

// -------------------------------------------------------------------------
// GET /api/conversations  -> merged, paginated by last activity
// -------------------------------------------------------------------------

func GetConversation(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	userId, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not authorized", nil)
		return
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 30
	}
	if limit > 30 {
		limit = 30
	}

	// =========================================================
	// STEP 1: rank DMs + groups together by last activity,
	// return only the page window of (id, type).
	// =========================================================
	rankRows, err := sqlite.DB().Query(`
		SELECT id, type, last_message_at FROM (
			SELECT c.id AS id, 'direct' AS type, c.last_message_at AS last_message_at
			FROM CONVERSATIONS c
			WHERE c.user1_id = ? OR c.user2_id = ?

			UNION ALL

			SELECT g.id AS id, 'group' AS type, gm.last_message_at AS last_message_at
			FROM GROUPS g
			JOIN GROUP_MEMBERS mem
				ON mem.group_id = g.id AND mem.user_id = ?
			LEFT JOIN (
				SELECT group_id, MAX(created_at) AS last_message_at
				FROM GROUP_MESSAGES
				GROUP BY group_id
			) gm ON gm.group_id = g.id
		)
		ORDER BY (last_message_at IS NULL) ASC, last_message_at DESC
		LIMIT ? OFFSET ?;
	`, userId, userId, userId, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var refs []feedRef
	var directIDs []int
	var groupIDs []int

	for rankRows.Next() {
		var ref feedRef
		if err := rankRows.Scan(&ref.id, &ref.kind, &ref.lastMessageAt); err != nil {
			rankRows.Close()
			http.Error(w, err.Error(), 500)
			return
		}
		refs = append(refs, ref)
		if ref.kind == "direct" {
			directIDs = append(directIDs, ref.id)
		} else {
			groupIDs = append(groupIDs, ref.id)
		}
	}
	rankRows.Close()

	// =========================================================
	// STEP 2: fetch full details for the DMs on this page
	// =========================================================
	directDetails := map[int]ConversationFeedItem{}
	if len(directIDs) > 0 {
		query, args := buildInClause(`
			SELECT
				c.id,
				CASE WHEN c.user1_id = ? THEN u2.id ELSE u1.id END,
				CASE WHEN c.user1_id = ? THEN u2.nickname ELSE u1.nickname END,
				CASE WHEN c.user1_id = ? THEN u2.firstname ELSE u1.firstname END,
				CASE WHEN c.user1_id = ? THEN u2.lastname ELSE u1.lastname END,
				CASE WHEN c.user1_id = ? THEN u2.avatar ELSE u1.avatar END,
				c.last_message,
				c.last_message_at
			FROM CONVERSATIONS c
			JOIN USERS u1 ON u1.id = c.user1_id
			JOIN USERS u2 ON u2.id = c.user2_id
			WHERE c.id IN (%s)
		`, []interface{}{userId, userId, userId, userId, userId}, directIDs)

		rows, err := sqlite.DB().Query(query, args...)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		for rows.Next() {
			var (
				convID              int
				otherID             int
				nickname            sql.NullString
				firstname, lastname string
				avatar              sql.NullString
				lastMsg, lastDate   sql.NullString
			)
			if err := rows.Scan(&convID, &otherID, &nickname, &firstname, &lastname, &avatar, &lastMsg, &lastDate); err != nil {
				rows.Close()
				http.Error(w, err.Error(), 500)
				return
			}

			displayName := firstname + " " + lastname
			if nickname.Valid && nickname.String != "" {
				displayName = nickname.String
			}

			var unread int
			_ = sqlite.DB().QueryRow(`
				SELECT COUNT(*) FROM MESSAGES
				WHERE conversation_id = ? AND sender_id != ? AND is_read = 0
			`, convID, userId).Scan(&unread)

			item := ConversationFeedItem{
				Type:        "direct",
				ID:          convID,
				DisplayName: displayName,
				UnreadCount: unread,
				OtherUserID: &otherID,
			}
			if avatar.Valid {
				item.Avatar = &avatar.String
			}
			if lastMsg.Valid {
				item.LastMessage = &lastMsg.String
			}
			if lastDate.Valid {
				item.LastMessageAt = &lastDate.String
			}
			directDetails[convID] = item
		}
		rows.Close()
	}

	// =========================================================
	// STEP 3: fetch full details for the groups on this page
	// =========================================================
	groupDetails := map[int]ConversationFeedItem{}
	if len(groupIDs) > 0 {
		query, args := buildInClause(`
			SELECT
				g.id,
				g.title,
				(SELECT COUNT(*) FROM GROUP_MEMBERS gm WHERE gm.group_id = g.id),
				lm.text,
				lm.created_at
			FROM GROUPS g
			LEFT JOIN (
				SELECT gm1.group_id, gm1.text, gm1.created_at
				FROM GROUP_MESSAGES gm1
				INNER JOIN (
					SELECT group_id, MAX(created_at) AS max_created
					FROM GROUP_MESSAGES
					GROUP BY group_id
				) gm2 ON gm1.group_id = gm2.group_id AND gm1.created_at = gm2.max_created
			) lm ON lm.group_id = g.id
			WHERE g.id IN (%s)
		`, nil, groupIDs)

		rows, err := sqlite.DB().Query(query, args...)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		for rows.Next() {
			var (
				groupID           int
				title             string
				memberCount       int
				lastMsg, lastDate sql.NullString
			)
			if err := rows.Scan(&groupID, &title, &memberCount, &lastMsg, &lastDate); err != nil {
				rows.Close()
				http.Error(w, err.Error(), 500)
				return
			}

			// unread = messages after this user's last_read_message_id
			// (requires GROUP_MESSAGE_READS table)
			var unread int
			_ = sqlite.DB().QueryRow(`
				SELECT COUNT(*) FROM GROUP_MESSAGES gmsg
				WHERE gmsg.group_id = ?
				AND gmsg.sender_id != ?
				AND gmsg.id > COALESCE(
					(SELECT last_read_message_id FROM GROUP_MESSAGE_READS
					 WHERE group_id = ? AND user_id = ?), 0
				)
			`, groupID, userId, groupID, userId).Scan(&unread)

			item := ConversationFeedItem{
				Type:        "group",
				ID:          groupID,
				DisplayName: title,
				UnreadCount: unread,
				MemberCount: &memberCount,
			}
			if lastMsg.Valid {
				item.LastMessage = &lastMsg.String
			}
			if lastDate.Valid {
				item.LastMessageAt = &lastDate.String
			}
			groupDetails[groupID] = item
		}
		rows.Close()
	}

	// =========================================================
	// STEP 4: reassemble in the ranked order from STEP 1,
	// but split back out into their own arrays by type
	// =========================================================
	directItems := make([]ConversationFeedItem, 0, len(directIDs))
	groupItems := make([]ConversationFeedItem, 0, len(groupIDs))

	for i, ref := range refs {
		if ref.kind == "direct" {
			if it, ok := directDetails[ref.id]; ok {
				it.Rank = i
				directItems = append(directItems, it)
			}
		} else {
			if it, ok := groupDetails[ref.id]; ok {
				it.Rank = i
				groupItems = append(groupItems, it)
			}
		}
	}

	utilities.WriteJSON(w, 200, "ok", map[string]interface{}{
		"groups": groupItems,
		"direct": directItems,
	})
}

// -------------------------------------------------------------------------
// GET /api/conversations/{type}/{id}/messages  (type = "direct" | "group")
// -------------------------------------------------------------------------

func GetConversationByID(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utilities.WriteJSON(w, 401, "unauthorized", nil)
		return
	}
	userID, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		utilities.WriteJSON(w, 401, "unauthorized", nil)
		return
	}

	convType := r.PathValue("type") // "direct" or "group"
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utilities.WriteJSON(w, 400, "invalid id", nil)
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	type Message struct {
		ID        int    `json:"id"`
		SenderID  int    `json:"sender_id"`
		Text      string `json:"text"`
		CreatedAt string `json:"created_at"`
	}
	messages := []Message{}

	switch convType {

	case "direct":
		var convID int
		err = sqlite.DB().QueryRow(`
			SELECT id FROM CONVERSATIONS
			WHERE id = ? AND (user1_id = ? OR user2_id = ?)
		`, id, userID, userID).Scan(&convID)
		if err == sql.ErrNoRows {
			utilities.WriteJSON(w, 403, "not allowed", nil)
			return
		}
		if err != nil {
			utilities.WriteJSON(w, 500, "db error", nil)
			return
		}

		rows, err := sqlite.DB().Query(`
			SELECT id, sender_id, text, created_at
			FROM MESSAGES
			WHERE conversation_id = ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`, convID, limit, offset)
		if err != nil {
			utilities.WriteJSON(w, 500, "db error", nil)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var m Message
			if err := rows.Scan(&m.ID, &m.SenderID, &m.Text, &m.CreatedAt); err != nil {
				utilities.WriteJSON(w, 500, "scan error", nil)
				return
			}
			messages = append(messages, m)
		}

	case "group":
		var groupID int
		err = sqlite.DB().QueryRow(`
			SELECT g.id FROM GROUPS g
			JOIN GROUP_MEMBERS m ON m.group_id = g.id
			WHERE g.id = ? AND m.user_id = ?
		`, id, userID).Scan(&groupID)
		if err == sql.ErrNoRows {
			utilities.WriteJSON(w, 403, "not allowed", nil)
			return
		}
		if err != nil {
			utilities.WriteJSON(w, 500, "db error", nil)
			return
		}

		rows, err := sqlite.DB().Query(`
			SELECT id, sender_id, text, created_at
			FROM GROUP_MESSAGES
			WHERE group_id = ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`, groupID, limit, offset)
		if err != nil {
			utilities.WriteJSON(w, 500, "db error", nil)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var m Message
			if err := rows.Scan(&m.ID, &m.SenderID, &m.Text, &m.CreatedAt); err != nil {
				utilities.WriteJSON(w, 500, "scan error", nil)
				return
			}
			messages = append(messages, m)
		}

	default:
		utilities.WriteJSON(w, 400, "invalid conversation type", nil)
		return
	}

	utilities.WriteJSON(w, 200, "ok", map[string]interface{}{
		"type":     convType,
		"id":       id,
		"messages": messages,
	})
}

// -------------------------------------------------------------------------
// helper: build a `col IN (?,?,?)` clause, prepending any leading args
// -------------------------------------------------------------------------

func buildInClause(queryTemplate string, leadingArgs []interface{}, ids []int) (string, []interface{}) {
	placeholders := make([]string, len(ids))
	args := append([]interface{}{}, leadingArgs...)
	for i, id := range ids {
		placeholders[i] = "?"
		args = append(args, id)
	}
	query := fmt.Sprintf(queryTemplate, strings.Join(placeholders, ","))
	return query, args
}
