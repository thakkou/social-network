package repository

import (
	"database/sql"
	"time"
)

type Group struct {
	ID          int       `json:"id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Logo        string    `json:"logo,omitempty"`
	Background  string    `json:"background,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupPost struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
}

type GroupPostComment struct {
	ID          int       `json:"id"`
	GroupPostID int       `json:"group_post_id"`
	UserID      int       `json:"user_id"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupEvent struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupMessage struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	SenderID  int       `json:"sender_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupRepository struct {
	DB *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{DB: db}
}

func (r *GroupRepository) CreateGroup(g *Group) error {
	description := sql.NullString{String: g.Description, Valid: g.Description != ""}
	logo := sql.NullString{String: g.Logo, Valid: g.Logo != ""}
	background := sql.NullString{String: g.Background, Valid: g.Background != ""}

	query := `INSERT INTO GROUPS (creator_id, title, description, logo, background) VALUES (?, ?, ?, ?, ?)`
	res, err := r.DB.Exec(query, g.CreatorID, g.Title, description, logo, background)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		g.ID = int(id)
	}

	// Creator automatically becomes the group admin
	_, err = r.DB.Exec(`INSERT INTO GROUP_MEMBERS (group_id, user_id, role) VALUES (?, ?, 'admin')`, g.ID, g.CreatorID)
	return err
}

func (r *GroupRepository) ListGroups() ([]Group, error) {
	rows, err := r.DB.Query(`
SELECT 
    id,
    creator_id,
    title,
    COALESCE(description, ''),
    created_at
FROM GROUPS
ORDER BY created_at DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		var createdAt string
		if err := rows.Scan(&g.ID, &g.CreatorID, &g.Title, &g.Description, &createdAt); err != nil {
			return nil, err
		}
		g.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *GroupRepository) RequestToJoin(groupID, userID int) error {
	query := `INSERT INTO GROUP_REQUESTS (group_id, user_id, status) VALUES (?, ?, 'pending')`
	_, err := r.DB.Exec(query, groupID, userID)
	return err
}

func (r *GroupRepository) InviteToGroup(groupID, inviterID, invitedUserID int) error {
	_, err := r.DB.Exec(`INSERT INTO GROUP_INVITES (group_id, inviter_id, invited_user_id, status) VALUES (?, ?, ?, 'pending')`, groupID, inviterID, invitedUserID)
	return err
}

func (r *GroupRepository) AcceptGroupInvite(groupID, userID int) error {
	_, err := r.DB.Exec(`UPDATE GROUP_INVITES SET status = 'accepted' WHERE group_id = ? AND invited_user_id = ? AND status = 'pending'`, groupID, userID)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(`INSERT OR IGNORE INTO GROUP_MEMBERS (group_id, user_id, role) VALUES (?, ?, 'member')`, groupID, userID)
	return err
}

func (r *GroupRepository) RejectGroupInvite(groupID, userID int) error {
	_, err := r.DB.Exec(`UPDATE GROUP_INVITES SET status = 'rejected' WHERE group_id = ? AND invited_user_id = ? AND status = 'pending'`, groupID, userID)
	return err
}

func (r *GroupRepository) AcceptGroupRequest(groupID, userID int) error {
	_, err := r.DB.Exec(`UPDATE GROUP_REQUESTS SET status = 'accepted' WHERE group_id = ? AND user_id = ? AND status = 'pending'`, groupID, userID)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(`INSERT OR IGNORE INTO GROUP_MEMBERS (group_id, user_id, role) VALUES (?, ?, 'member')`, groupID, userID)
	return err
}

func (r *GroupRepository) RejectGroupRequest(groupID, userID int) error {
	_, err := r.DB.Exec(`UPDATE GROUP_REQUESTS SET status = 'rejected' WHERE group_id = ? AND user_id = ? AND status = 'pending'`, groupID, userID)
	return err
}

func (r *GroupRepository) GetGroupCreatorID(groupID int) (int, error) {
	var creatorID int
	err := r.DB.QueryRow(`SELECT creator_id FROM GROUPS WHERE id = ?`, groupID).Scan(&creatorID)
	return creatorID, err
}

func (r *GroupRepository) GetGroupMemberIDs(groupID int) ([]int, error) {
	rows, err := r.DB.Query(`SELECT user_id FROM GROUP_MEMBERS WHERE group_id = ?`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]int, 0)
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		members = append(members, userID)
	}
	return members, rows.Err()
}

func (r *GroupRepository) IsGroupMember(groupID, userID int) (bool, error) {
	var exists int
	err := r.DB.QueryRow(`SELECT 1 FROM GROUP_MEMBERS WHERE group_id = ? AND user_id = ?`, groupID, userID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *GroupRepository) CreateGroupMessage(groupID, userID int, text string) (*GroupMessage, error) {
	res, err := r.DB.Exec(`INSERT INTO GROUP_MESSAGES (group_id, sender_id, text) VALUES (?, ?, ?)`, groupID, userID, text)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &GroupMessage{ID: int(id), GroupID: groupID, SenderID: userID, Text: text}, nil
}

func (r *GroupRepository) ListGroupMessages(groupID, limit, offset int) ([]GroupMessage, error) {
	rows, err := r.DB.Query(`SELECT id, group_id, sender_id, text, created_at FROM GROUP_MESSAGES WHERE group_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]GroupMessage, 0)
	for rows.Next() {
		var m GroupMessage
		var createdAt string
		if err := rows.Scan(&m.ID, &m.GroupID, &m.SenderID, &m.Text, &createdAt); err != nil {
			return nil, err
		}
		m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (r *GroupRepository) ListGroupPosts(groupID int) ([]GroupPost, error) {
	rows, err := r.DB.Query(`SELECT id, group_id, user_id, created_at, title, text, image FROM GROUP_POSTS WHERE group_id = ? ORDER BY created_at DESC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]GroupPost, 0)
	for rows.Next() {
		var p GroupPost
		var createdAt string
		if err := rows.Scan(&p.ID, &p.GroupID, &p.UserID, &createdAt, &p.Title, &p.Text, &p.Image); err != nil {
			return nil, err
		}
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *GroupRepository) ListGroupPostComments(groupPostID int) ([]GroupPostComment, error) {
	rows, err := r.DB.Query(`SELECT id, group_post_id, user_id, text, created_at FROM GROUP_POST_COMMENTS WHERE group_post_id = ? ORDER BY created_at ASC`, groupPostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]GroupPostComment, 0)
	for rows.Next() {
		var c GroupPostComment
		var createdAt string
		if err := rows.Scan(&c.ID, &c.GroupPostID, &c.UserID, &c.Text, &createdAt); err != nil {
			return nil, err
		}
		c.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *GroupRepository) CreateGroupPost(gp *GroupPost) error {
	query := `INSERT INTO GROUP_POSTS (group_id, user_id, title, text, image) VALUES (?, ?, ?, ?, ?)`
	res, err := r.DB.Exec(query, gp.GroupID, gp.UserID, gp.Title, gp.Text, gp.Image)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	gp.ID = int(id)
	return nil
}

func (r *GroupRepository) CreateGroupPostComment(groupPostID, userID int, text string) error {
	_, err := r.DB.Exec(`INSERT INTO GROUP_POST_COMMENTS (group_post_id, user_id, text) VALUES (?, ?, ?)`, groupPostID, userID, text)
	return err
}

func (r *GroupRepository) ListGroupEvents(groupID int) ([]GroupEvent, error) {
	rows, err := r.DB.Query(`SELECT id, group_id, creator_id, title, description, event_time, created_at FROM GROUP_EVENTS WHERE group_id = ? ORDER BY event_time ASC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]GroupEvent, 0)
	for rows.Next() {
		var e GroupEvent
		var eventTime string
		var createdAt string
		if err := rows.Scan(&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &eventTime, &createdAt); err != nil {
			return nil, err
		}
		e.EventTime, _ = time.Parse("2006-01-02 15:04:05", eventTime)
		e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *GroupRepository) CreateEvent(e *GroupEvent) error {
	query := `INSERT INTO GROUP_EVENTS (group_id, creator_id, title, description, event_time) VALUES (?, ?, ?, ?, ?)`
	res, err := r.DB.Exec(query, e.GroupID, e.CreatorID, e.Title, e.Description, e.EventTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	e.ID = int(id)
	return nil
}

func (r *GroupRepository) RespondToEvent(eventID, userID int, status string) error {
	query := `INSERT INTO EVENT_RESPONSES (event_id, user_id, status) VALUES (?, ?, ?)
	          ON CONFLICT(event_id, user_id) DO UPDATE SET status = ?`
	_, err := r.DB.Exec(query, eventID, userID, status, status)
	return err
}

func (r *GroupRepository) SearchGroups(text string) ([]Group, error) {
	query := `
	SELECT
		id,
		creator_id,
		title,
		COALESCE(description, ''),	
		logo,
		background,
	created_at
	FROM GROUPS
	WHERE
		title LIKE ?
		OR description LIKE ?
	ORDER BY created_at DESC
	LIMIT 20
	`

	search := "%" + text + "%"

	rows, err := r.DB.Query(query, search, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group

	for rows.Next() {
		var g Group
		var createdAt string

		if err := rows.Scan(
			&g.ID,
			&g.CreatorID,
			&g.Title,
			&g.Description,
			&g.Logo,
			&g.Background,
			&createdAt,
		); err != nil {
			return nil, err
		}

		g.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

		groups = append(groups, g)
	}

	return groups, rows.Err()
}

func (r *GroupRepository) GetPublicGroupDetails(groupID int) (*Group, error) {
	var group Group

	var createdAt string

	err := r.DB.QueryRow(`
		SELECT
    id,
    title,
    COALESCE(description, ''),
    COALESCE(logo, ''),
    COALESCE(backgroun, ''),
    created_at
FROM GROUPS
		WHERE id = ?
	`, groupID).Scan(
		&group.ID,
		&group.Title,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	group.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

	return &group, nil
}
