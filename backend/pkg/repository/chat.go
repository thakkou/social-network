package repository

import (
	"database/sql"
	"time"
)

type Conversation struct {
	ID            int       `json:"id"`
	User1ID       int       `json:"user1_id"`
	User2ID       int       `json:"user2_id"`
	LastMessage   string    `json:"last_message"`
	LastMessageAt time.Time `json:"last_message_at"`
	Created_At    time.Time `json:"created_at"`
}

type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	SenderID       int       `json:"sender_id"`
	Text           string    `json:"text"`
	CreatedAt      time.Time `json:"created_at"`
	IsRead         int       `json:"is_read"`
}

type ChatRepository struct {
	DB *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{DB: db}
}

func (r *ChatRepository) GetOrCreateConversation(u1, u2 int) (int, error) {
	if u1 > u2 {
		u1, u2 = u2, u1
	}

	var id int
	query := `SELECT id FROM CONVERSATIONS WHERE user1_id = ? AND user2_id = ?`
	err := r.DB.QueryRow(query, u1, u2).Scan(&id)
	if err == nil {
		return id, nil
	}

	insertQuery := `INSERT INTO CONVERSATIONS (user1_id, user2_id) VALUES (?, ?)`
	res, err := r.DB.Exec(insertQuery, u1, u2)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	return int(lastID), err
}

func (r *ChatRepository) SaveMessage(m *Message) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tStr := m.CreatedAt.Format("2006-01-02 15:04:05.000")
	query := `INSERT INTO MESSAGES (conversation_id, sender_id, text, created_at) VALUES (?, ?, ?, ?)`
	res, err := tx.Exec(query, m.ConversationID, m.SenderID, m.Text, tStr)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	m.ID = int(id)

	updateConv := `UPDATE CONVERSATIONS SET last_message = ?, last_message_at = ? WHERE id = ?`
	_, err = tx.Exec(updateConv, m.Text, tStr, m.ConversationID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ChatRepository) GetChatHistory(convID, limit, offset int) ([]Message, error) {
	query := `SELECT id, conversation_id, sender_id, text, created_at, is_read 
	          FROM MESSAGES WHERE conversation_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.DB.Query(query, convID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		var tStr string
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Text, &tStr, &m.IsRead); err != nil {
			return nil, err
		}
		m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", tStr[:19])
		messages = append(messages, m)
	}
	return messages, nil
}
