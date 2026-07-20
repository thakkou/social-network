package seeder

import (
	"log"
	"time"

	"database/sql"
)

func seedConversations(db *sql.DB, u []int) ([]int, error) {
	type cv struct {
		User1, User2  int
		LastMessage   string
		LastMessageAt time.Time
	}
	convs := []cv{
		{u[0], u[1], "See you tomorrow!", time.Now().Add(-2 * time.Hour)},
		{u[0], u[5], "Sounds good, thanks!", time.Now().Add(-45 * time.Minute)},
	}
	query := `INSERT INTO CONVERSATIONS (user1_id, user2_id, last_message, last_message_at) VALUES (?, ?, ?, ?)`
	ids := make([]int, 0, len(convs))
	for _, c := range convs {
		res, err := db.Exec(query, c.User1, c.User2, c.LastMessage, c.LastMessageAt.Format("2006-01-02 15:04:05"))
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	log.Printf("[SEED] conversations: %d\n", len(ids))
	return ids, nil
}

func seedMessages(db *sql.DB, convIDs, u []int) error {
	type m struct {
		ConversationID, SenderID int
		Text                     string
		IsRead                   int
	}
	messages := []m{
		{convIDs[0], u[0], "Hey Bob, are we still on for tomorrow?", 1},
		{convIDs[0], u[1], "Yep, see you tomorrow!", 0},
		{convIDs[1], u[0], "Can you send me the trail map?", 1},
		{convIDs[1], u[5], "Sounds good, thanks!", 0},
	}
	query := `INSERT INTO MESSAGES (conversation_id, sender_id, text, is_read) VALUES (?, ?, ?, ?)`
	for _, msg := range messages {
		if _, err := db.Exec(query, msg.ConversationID, msg.SenderID, msg.Text, msg.IsRead); err != nil {
			return err
		}
	}
	log.Printf("[SEED] messages: %d\n", len(messages))
	return nil
}
