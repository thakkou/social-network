package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

type ConversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

type ConversationFeedItem struct {
	Type          string  `json:"type"` // "direct" | "group"
	ID            int     `json:"id"`   // conversation_id or group_id
	DisplayName   string  `json:"display_name"`
	Avatar        *string `json:"avatar,omitempty"`
	LastMessage   *string `json:"last_message"`
	LastMessageAt *string `json:"last_message_at"`
	Rank          int     `json:"rank"`

	OtherUserID *int `json:"other_user_id,omitempty"`
	MemberCount *int `json:"member_count,omitempty"`
}

type Message struct {
	ID        int    `json:"id"`
	SenderID  int    `json:"sender_id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	Nickname  string `json:"nickname"`
}

type feedRef struct {
	id            int
	kind          string
	lastMessageAt sql.NullString
}

// FindOrCreateDirectConversation gets an existing DM ID or creates a new row.
func (r *ConversationRepository) FindOrCreateDirectConversation(tx *sql.Tx, senderID, receiverID int) (int, bool, error) {
	u1, u2 := senderID, receiverID
	if u1 > u2 {
		u1, u2 = u2, u1
	}

	var conversationID int
	err := tx.QueryRow(`
		SELECT id FROM CONVERSATIONS
		WHERE user1_id = ? AND user2_id = ?
	`, u1, u2).Scan(&conversationID)

	if err == sql.ErrNoRows {
		res, err := tx.Exec(`
			INSERT INTO CONVERSATIONS (user1_id, user2_id)
			VALUES (?, ?)
		`, u1, u2)
		if err != nil {
			return 0, false, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return 0, false, err
		}
		return int(id), true, nil
	} else if err != nil {
		return 0, false, err
	}

	return conversationID, false, nil
}

// SaveMessage inserts the message and updates conversation last_message stats.
func (r *ConversationRepository) SaveMessage(tx *sql.Tx, conversationID, senderID int, text string) (int64, error) {
	res, err := tx.Exec(`
		INSERT INTO MESSAGES (conversation_id, sender_id, text)
		VALUES (?, ?, ?)
	`, conversationID, senderID, text)
	if err != nil {
		return 0, err
	}

	msgID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(`
		UPDATE CONVERSATIONS
		SET last_message = ?, last_message_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, text, conversationID)
	return msgID, err
}

// SaveGroupMessage verifies group membership and saves a message to GROUP_MESSAGES.
func (r *ConversationRepository) SaveGroupMessage(tx *sql.Tx, groupID, senderID int, text string) (int64, error) {
	// Verify sender is an active member of the group
	var isMember bool
	err := tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM GROUP_MEMBERS 
			WHERE group_id = ? AND user_id = ?
		)
	`, groupID, senderID).Scan(&isMember)
	if err != nil {
		return 0, err
	}
	if !isMember {
		return 0, sql.ErrNoRows // Used to signify unauthorized access
	}

	// Insert into GROUP_MESSAGES
	res, err := tx.Exec(`
		INSERT INTO GROUP_MESSAGES (group_id, sender_id, text)
		VALUES (?, ?, ?)
	`, groupID, senderID, text)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// GetGroupMemberIDs fetches all member user IDs in a group (for WS broadcasting).
func (r *ConversationRepository) GetGroupMemberIDs(groupID int) ([]int, error) {
	rows, err := r.db.Query(`
		SELECT user_id FROM GROUP_MEMBERS WHERE group_id = ?
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberIDs []int
	for rows.Next() {
		var uid int
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, uid)
	}
	return memberIDs, nil
}

// FetchConversationFeed returns ranked direct and group feeds for a user.
func (r *ConversationRepository) FetchConversationFeed(userID, limit, offset int) ([]ConversationFeedItem, []ConversationFeedItem, error) {
	rows, err := r.db.Query(`
		SELECT id, type, last_message_at FROM (
			SELECT c.id AS id, 'direct' AS type, c.last_message_at AS last_message_at
			FROM CONVERSATIONS c
			WHERE c.user1_id = ? OR c.user2_id = ?

			UNION ALL

			SELECT g.id AS id, 'group' AS type, gm.last_message_at AS last_message_at
			FROM GROUPS g
			JOIN GROUP_MEMBERS mem ON mem.group_id = g.id AND mem.user_id = ?
			LEFT JOIN (
				SELECT group_id, MAX(created_at) AS last_message_at
				FROM GROUP_MESSAGES
				GROUP BY group_id
			) gm ON gm.group_id = g.id
		)
		ORDER BY (last_message_at IS NULL) ASC, last_message_at DESC
		LIMIT ? OFFSET ?;
	`, userID, userID, userID, limit, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var refs []feedRef
	var directIDs, groupIDs []int

	for rows.Next() {
		var ref feedRef
		if err := rows.Scan(&ref.id, &ref.kind, &ref.lastMessageAt); err != nil {
			return nil, nil, err
		}
		refs = append(refs, ref)
		if ref.kind == "direct" {
			directIDs = append(directIDs, ref.id)
		} else {
			groupIDs = append(groupIDs, ref.id)
		}
	}

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
		`, []interface{}{userID, userID, userID, userID, userID}, directIDs)

		dRows, err := r.db.Query(query, args...)
		if err != nil {
			return nil, nil, err
		}
		defer dRows.Close()

		for dRows.Next() {
			var (
				convID, otherID     int
				nickname            sql.NullString
				firstname, lastname string
				avatar              sql.NullString
				lastMsg, lastDate   sql.NullString
			)
			if err := dRows.Scan(&convID, &otherID, &nickname, &firstname, &lastname, &avatar, &lastMsg, &lastDate); err != nil {
				return nil, nil, err
			}

			displayName := firstname + " " + lastname
			if nickname.Valid && nickname.String != "" {
				displayName = nickname.String
			}

			item := ConversationFeedItem{
				Type:        "direct",
				ID:          convID,
				DisplayName: displayName,
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
	}

	groupDetails := map[int]ConversationFeedItem{}
	if len(groupIDs) > 0 {
		query, args := buildInClause(`
			SELECT
				g.id,
				g.title,
				(SELECT COUNT(*) FROM GROUP_MEMBERS gm WHERE gm.group_id = g.id),
				lm.text,
				lm.created_at,
				g.logo
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

		gRows, err := r.db.Query(query, args...)
		if err != nil {
			return nil, nil, err
		}
		defer gRows.Close()

		for gRows.Next() {
			var (
				groupID           int
				title             string
				memberCount       int
				lastMsg, lastDate sql.NullString
				groupLogo         sql.NullString
			)
			if err := gRows.Scan(&groupID, &title, &memberCount, &lastMsg, &lastDate, &groupLogo); err != nil {
				return nil, nil, err
			}

			item := ConversationFeedItem{
				Type:        "group",
				ID:          groupID,
				DisplayName: title,
				MemberCount: &memberCount,
			}
			if lastMsg.Valid {
				item.LastMessage = &lastMsg.String
			}
			if lastDate.Valid {
				item.LastMessageAt = &lastDate.String
			}
			if groupLogo.Valid && groupLogo.String != "" {
				item.Avatar = &groupLogo.String
			}
			groupDetails[groupID] = item
		}
	}

	var directItems []ConversationFeedItem
	var groupItems []ConversationFeedItem

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

	return directItems, groupItems, nil
}

// FetchDirectMessages gets paginated messages for DMs.
func (r *ConversationRepository) FetchDirectMessages(convID, userID, limit, offset int) ([]Message, error) {
	var validID int
	err := r.db.QueryRow(`
		SELECT id FROM CONVERSATIONS
		WHERE id = ? AND (user1_id = ? OR user2_id = ?)
	`, convID, userID, userID).Scan(&validID)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT m.id, m.sender_id, m.text, m.created_at, COALESCE(u.nickname, u.firstname) 
		FROM MESSAGES m
		JOIN USERS u ON u.id = m.sender_id
		WHERE m.conversation_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`, convID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.SenderID, &m.Text, &m.CreatedAt, &m.Nickname); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// FetchGroupMessages gets paginated messages for Groups.
func (r *ConversationRepository) FetchGroupMessages(groupID, userID, limit, offset int) ([]Message, error) {
	var validID int
	err := r.db.QueryRow(`
		SELECT g.id FROM GROUPS g
		JOIN GROUP_MEMBERS m ON m.group_id = g.id
		WHERE g.id = ? AND m.user_id = ?
	`, groupID, userID).Scan(&validID)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT gm.id, gm.sender_id, gm.text, gm.created_at, COALESCE(u.nickname, u.firstname)
		FROM GROUP_MESSAGES gm
		JOIN USERS u ON u.id = gm.sender_id
		WHERE gm.group_id = ?
		ORDER BY gm.created_at DESC
		LIMIT ? OFFSET ?
	`, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.SenderID, &m.Text, &m.CreatedAt, &m.Nickname); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func buildInClause(queryTemplate string, leadingArgs []interface{}, ids []int) (string, []interface{}) {
	placeholders := make([]string, len(ids))
	args := append([]interface{}{}, leadingArgs...)
	for i, id := range ids {
		placeholders[i] = "?"
		args = append(args, id)
	}
	return fmt.Sprintf(queryTemplate, strings.Join(placeholders, ",")), args
}
