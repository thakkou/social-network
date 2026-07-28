package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"01social/pkg/utilities"
)

// FeedAuthor holds basic profile info for the creator of a feed item.
type FeedAuthor struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Avatar    string `json:"avatar"`
}

type Group struct {
	ID          int       `json:"id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Logo        string    `json:"logo,omitempty"`
	Background  string    `json:"background,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupFeedItem struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"` // "post" or "event"`
	GroupID     int       `json:"group_id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title,omitempty"`
	Text        string    `json:"text,omitempty"`
	Image       string    `json:"image"`
	Description string    `json:"description,omitempty"`
	EventTime   time.Time `json:"event_time,omitempty"`
	CreatedAt   time.Time `json:"created_at"`

	// Creator profile
	Author FeedAuthor `json:"author"`

	// Populated for Type == "post"
	LikesCount    int                           `json:"likes_count,omitempty"`
	DislikesCount int                           `json:"dislikes_count,omitempty"`
	IsLiked       int                           `json:"is_liked"`
	CommentsCount int                           `json:"comments_count,omitempty"`
	Comments      []GroupPostCommentWithAuthor `json:"comments,omitempty"`

	// Populated for Type == "event"
	EventResponses []EventResponder `json:"event_responses,omitempty"`
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
	query := `INSERT INTO GROUP_REQUESTS (group_id, user_id, status) VALUES (?, ?, 'pending')
		ON CONFLICT(group_id, user_id) DO UPDATE SET status = 'pending'`
	_, err := r.DB.Exec(query, groupID, userID)
	return err
}

func (r *GroupRepository) InviteToGroup(groupID, inviterID, invitedUserID int) error {
	_, err := r.DB.Exec(`INSERT INTO GROUP_INVITES (group_id, inviter_id, invited_user_id, status) VALUES (?, ?, ?, 'pending')`, groupID, inviterID, invitedUserID)
	return err
}

func (r *GroupRepository) AcceptGroupInvite(groupID, userID int) error {
	res, err := r.DB.Exec(`UPDATE GROUP_INVITES SET status = 'accepted' WHERE group_id = ? AND invited_user_id = ? AND status = 'pending'`, groupID, userID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no pending invite found")
	}
	_, err = r.DB.Exec(`INSERT OR IGNORE INTO GROUP_MEMBERS (group_id, user_id, role) VALUES (?, ?, 'member')`, groupID, userID)
	return err
}

func (r *GroupRepository) RejectGroupInvite(groupID, userID int) error {
	res, err := r.DB.Exec(`UPDATE GROUP_INVITES SET status = 'rejected' WHERE group_id = ? AND invited_user_id = ? AND status = 'pending'`, groupID, userID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no pending invite found")
	}
	return nil
}

func (r *GroupRepository) AcceptGroupRequest(groupID, userID int) error {
	res, err := r.DB.Exec(`UPDATE GROUP_REQUESTS SET status = 'accepted' WHERE group_id = ? AND user_id = ? AND status = 'pending'`, groupID, userID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no pending request found")
	}
	_, err = r.DB.Exec(`INSERT OR IGNORE INTO GROUP_MEMBERS (group_id, user_id, role) VALUES (?, ?, 'member')`, groupID, userID)
	return err
}

func (r *GroupRepository) RejectGroupRequest(groupID, userID int) error {
	res, err := r.DB.Exec(`UPDATE GROUP_REQUESTS SET status = 'rejected' WHERE group_id = ? AND user_id = ? AND status = 'pending'`, groupID, userID)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no pending request found")
	}
	return nil
}

type PendingRequest struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Avatar    string `json:"avatar"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (r *GroupRepository) ListGroupPendingRequests(groupID int) ([]PendingRequest, error) {
	query := `
		SELECT
			gr.id,
			gr.user_id,
			COALESCE(u.nickname, '') AS nickname,
			u.firstname,
			u.lastname,
			COALESCE(u.avatar, '') AS avatar,
			gr.status,
			gr.created_at
		FROM GROUP_REQUESTS gr
		JOIN USERS u ON u.id = gr.user_id
		WHERE gr.group_id = ? AND gr.status = 'pending'
		ORDER BY gr.created_at DESC
	`
	rows, err := r.DB.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []PendingRequest
	for rows.Next() {
		var req PendingRequest
		if err := rows.Scan(&req.ID, &req.UserID, &req.Nickname, &req.Firstname, &req.Lastname, &req.Avatar, &req.Status, &req.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

func (r *GroupRepository) GetUserGroups(userID int) ([]Group, error) {
	query := `
		SELECT g.id, g.creator_id, g.title, COALESCE(g.description, ''), COALESCE(g.logo, ''), COALESCE(g.background, ''), g.created_at
		FROM GROUPS g
		JOIN GROUP_MEMBERS gm ON gm.group_id = g.id
		WHERE gm.user_id = ?
		ORDER BY g.created_at DESC
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		var createdAt string
		if err := rows.Scan(&g.ID, &g.CreatorID, &g.Title, &g.Description, &g.Logo, &g.Background, &createdAt); err != nil {
			return nil, err
		}
		g.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *GroupRepository) GetGroupCreatorID(groupID int) (int, error) {
	var creatorID int
	err := r.DB.QueryRow(`SELECT creator_id FROM GROUPS WHERE id = ?`, groupID).Scan(&creatorID)
	return creatorID, err
}

// GetGroupInvitedUserIDs returns the IDs of users who have been invited to the group
// (with pending or accepted status), excluding rejected invites.
func (r *GroupRepository) GetGroupInvitedUserIDs(groupID int) ([]int, error) {
	rows, err := r.DB.Query(`SELECT invited_user_id FROM GROUP_INVITES WHERE group_id = ? AND status != 'rejected'`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
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

func (r *GroupRepository) DoesGroupPostBelongToGroup(groupID, postID int) (bool, error) {
	var exists int
	err := r.DB.QueryRow(`SELECT 1 FROM GROUP_POSTS WHERE id = ? AND group_id = ?`, postID, groupID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *GroupRepository) DoesGroupEventBelongToGroup(groupID, eventID int) (bool, error) {
	var exists int
	err := r.DB.QueryRow(`SELECT 1 FROM GROUP_EVENTS WHERE id = ? AND group_id = ?`, eventID, groupID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
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

func (r *GroupRepository) UpdateGroup(groupID, creatorID int, title, description, logo, background string) error {
	// Verify the caller is the creator
	var actualCreator int
	err := r.DB.QueryRow(`SELECT creator_id FROM GROUPS WHERE id = ?`, groupID).Scan(&actualCreator)
	if err != nil {
		return fmt.Errorf("group not found")
	}
	if actualCreator != creatorID {
		return fmt.Errorf("only the group creator can update the group")
	}

	query := `UPDATE GROUPS SET
		title = COALESCE(NULLIF(?, ''), title),
		description = COALESCE(NULLIF(?, ''), description),
		logo = COALESCE(NULLIF(?, ''), logo),
		background = COALESCE(NULLIF(?, ''), background)
	WHERE id = ?`

	_, err = r.DB.Exec(query, title, description, logo, background, groupID)
	return err
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
    creator_id,
    title,
    COALESCE(description, ''),
    COALESCE(logo, ''),
    COALESCE(background, ''),
    created_at
FROM GROUPS
		WHERE id = ?
	`, groupID).Scan(
		&group.ID,
		&group.CreatorID,
		&group.Title,
		&group.Description,
		&group.Logo,
		&group.Background,
		&group.CreatedAt,
	)
	if err != nil {
		fmt.Println("errors", err)
		return nil, err
	}

	group.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

	return &group, nil
}

// GetUsersByIDs batch-fetches basic profile info for a set of user IDs.
// Returns a map keyed by user ID; every requested ID is present (zero-valued
// if the user is not found).
func (r *GroupRepository) GetUsersByIDs(userIDs []int) (map[int]FeedAuthor, error) {
	authors := make(map[int]FeedAuthor, len(userIDs))
	if len(userIDs) == 0 {
		return authors, nil
	}

	placeholders, args := utilities.PlaceholdersForInts(userIDs)

	query := `
SELECT
    id,
    COALESCE(nickname, '') AS nickname,
    firstname,
    lastname,
    COALESCE(avatar, '') AS avatar
FROM USERS
WHERE id IN (` + placeholders + `)
`
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a FeedAuthor
		if err := rows.Scan(&a.ID, &a.Nickname, &a.Firstname, &a.Lastname, &a.Avatar); err != nil {
			return nil, err
		}
		authors[a.ID] = a
	}

	return authors, rows.Err()
}

// KickMember removes a member from a group. Only the group creator can kick.
func (r *GroupRepository) KickMember(groupID, memberID, creatorID int) error {
	// Verify the caller is the group creator
	var actualCreator int
	err := r.DB.QueryRow(`SELECT creator_id FROM GROUPS WHERE id = ?`, groupID).Scan(&actualCreator)
	if err != nil {
		return fmt.Errorf("group not found")
	}
	if actualCreator != creatorID {
		return fmt.Errorf("only the group creator can kick members")
	}
	if memberID == creatorID {
		return fmt.Errorf("cannot kick the group creator")
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete member's reactions on group posts
	_, _ = tx.Exec(`DELETE FROM GROUP_POST_REACTIONS WHERE user_id = ? AND group_post_id IN (SELECT id FROM GROUP_POSTS WHERE group_id = ?)`, memberID, groupID)
	// Delete member's comments on group posts
	_, _ = tx.Exec(`DELETE FROM GROUP_POST_COMMENTS WHERE user_id = ? AND group_post_id IN (SELECT id FROM GROUP_POSTS WHERE group_id = ?)`, memberID, groupID)
	// Delete member's group posts
	_, _ = tx.Exec(`DELETE FROM GROUP_POSTS WHERE user_id = ? AND group_id = ?`, memberID, groupID)
	// Delete member's event responses
	_, _ = tx.Exec(`DELETE FROM EVENT_RESPONSES WHERE user_id = ? AND event_id IN (SELECT id FROM GROUP_EVENTS WHERE group_id = ?)`, memberID, groupID)
	// Remove member from group members
	_, err = tx.Exec(`DELETE FROM GROUP_MEMBERS WHERE group_id = ? AND user_id = ?`, groupID, memberID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *GroupRepository) LeaveGroup(groupID, userID int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete user's reactions on group posts
	_, _ = tx.Exec(`DELETE FROM GROUP_POST_REACTIONS WHERE user_id = ? AND group_post_id IN (SELECT id FROM GROUP_POSTS WHERE group_id = ?)`, userID, groupID)

	// Delete user's comments on group posts
	_, _ = tx.Exec(`DELETE FROM GROUP_POST_COMMENTS WHERE user_id = ? AND group_post_id IN (SELECT id FROM GROUP_POSTS WHERE group_id = ?)`, userID, groupID)

	// Delete user's group posts
	_, _ = tx.Exec(`DELETE FROM GROUP_POSTS WHERE user_id = ? AND group_id = ?`, userID, groupID)

	// Delete user's group messages
	_, _ = tx.Exec(`DELETE FROM GROUP_MESSAGES WHERE sender_id = ? AND group_id = ?`, userID, groupID)

	// Delete user's event responses
	_, _ = tx.Exec(`DELETE FROM EVENT_RESPONSES WHERE user_id = ? AND event_id IN (SELECT id FROM GROUP_EVENTS WHERE group_id = ?)`, userID, groupID)

	// Remove user from group members
	result, err := tx.Exec(`DELETE FROM GROUP_MEMBERS WHERE group_id = ? AND user_id = ?`, groupID, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("not a member of this group")
	}

	return tx.Commit()
}

func (r *GroupRepository) GetGroupContent(groupID, userID, limit, lastID int) ([]GroupFeedItem, error) {
	if limit <= 0 {
		limit = 20
	}

	log.Printf("[GetGroupContent] start groupID=%d userID=%d limit=%d lastID=%d",
		groupID, userID, limit, lastID)

	query := `
SELECT *
FROM (
	SELECT
		id,
		'post' AS type,
		group_id,
		user_id,
		COALESCE(title, '') AS title,
		COALESCE(text, '') AS text,
		COALESCE(image, '') AS image,
		'' AS description,
		NULL AS event_time,
		created_at
	FROM GROUP_POSTS
	WHERE group_id = ?

	UNION ALL

	SELECT
		id,
		'event' AS type,
		group_id,
		creator_id AS user_id,
		COALESCE(title, '') AS title,
		'' AS text,
		COALESCE(image, '') AS image,
		COALESCE(description, '') AS description,
		event_time,
		created_at
	FROM GROUP_EVENTS
	WHERE group_id = ?
) feed
WHERE (? = 0 OR id < ?)
ORDER BY created_at DESC
LIMIT ?
`

	log.Printf("[GetGroupContent] executing feed query")

	rows, err := r.DB.Query(query, groupID, groupID, lastID, lastID, limit)
	if err != nil {
		log.Printf("[GetGroupContent] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	feed := make([]GroupFeedItem, 0)

	for rows.Next() {
		var item GroupFeedItem

		var createdAt string
		var eventTime sql.NullString

		err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.GroupID,
			&item.UserID,
			&item.Title,
			&item.Text,
			&item.Image,
			&item.Description,
			&eventTime,
			&createdAt,
		)
		if err != nil {
			log.Printf("[GetGroupContent] row scan failed: %v", err)
			return nil, err
		}

		item.CreatedAt = utilities.ParseSQLiteTime(createdAt)
		if err != nil {
			log.Printf("[GetGroupContent] failed parsing createdAt id=%d value=%s error=%v",
				item.ID, createdAt, err)
		}

		if eventTime.Valid {
			item.EventTime = utilities.ParseSQLiteTime(eventTime.String)
			if err != nil {
				log.Printf("[GetGroupContent] failed parsing eventTime id=%d value=%s error=%v",
					item.ID, eventTime.String, err)
			}
		}

		feed = append(feed, item)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[GetGroupContent] rows iteration failed: %v", err)
		return nil, err
	}

	log.Printf("[GetGroupContent] fetched %d feed items", len(feed))

	var postIDs, eventIDs []int

	for _, item := range feed {
		if item.Type == "post" {
			postIDs = append(postIDs, item.ID)
		} else {
			eventIDs = append(eventIDs, item.ID)
		}
	}

	log.Printf("[GetGroupContent] enrichment IDs posts=%v events=%v",
		postIDs, eventIDs)

	postEngagement, err := r.GetPostsEngagement(postIDs, userID)
	if err != nil {
		log.Printf("[GetGroupContent] failed getting post engagement: %v", err)
		return nil, err
	}

	eventResponses, err := r.GetEventsResponses(eventIDs)
	if err != nil {
		log.Printf("[GetGroupContent] failed getting event responses: %v", err)
		return nil, err
	}

	log.Printf("[GetGroupContent] enrichment completed posts=%d events=%d",
		len(postEngagement), len(eventResponses))

	// Batch-fetch comments for all posts.
	postComments, err := r.GetGroupPostCommentsWithAuthors(postIDs, userID)
	if err != nil {
		log.Printf("[GetGroupContent] failed getting comments: %v", err)
		return nil, err
	}

	// Batch-fetch author profiles for all feed items.
	authorIDs := make([]int, 0, len(feed))
	seen := make(map[int]bool, len(feed))
	for _, item := range feed {
		if !seen[item.UserID] {
			seen[item.UserID] = true
			authorIDs = append(authorIDs, item.UserID)
		}
	}

	authorMap, err := r.GetUsersByIDs(authorIDs)
	if err != nil {
		log.Printf("[GetGroupContent] failed getting authors: %v", err)
		return nil, err
	}

	for i := range feed {
		item := &feed[i]

		if author, ok := authorMap[item.UserID]; ok {
			item.Author = author
		}

		if item.Type == "post" {
			if e, ok := postEngagement[item.ID]; ok {
				item.LikesCount = e.LikesCount
				item.DislikesCount = e.DislikesCount
				item.IsLiked = e.IsLiked
				item.CommentsCount = e.CommentsCount
			}
			if comments, ok := postComments[item.ID]; ok {
				item.Comments = comments
			}
		} else {
			item.EventResponses = eventResponses[item.ID]
		}
	}

	log.Printf("[GetGroupContent] completed successfully returned=%d items", len(feed))

	return feed, nil
}
