package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"01social/pkg/middlewares"
	"01social/pkg/repository"
	"01social/pkg/utilities"
)

type ActorInfo struct {
	UserID    int    `json:"user_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type NotifResponse struct {
	ID         int         `json:"id"`
	Type       string      `json:"type"`        // post_reaction, comment, follow_request, follow_accepted, group_invite, group_join_request
	ObjectType string      `json:"object_type"` // post, comment, follow, group_invite, group_request
	IsRead     bool        `json:"is_read"`
	CreatedAt  time.Time   `json:"created_at"`
	Actor      ActorInfo   `json:"actor"`
	Payload    interface{} `json:"payload"`
}
type FollowRequestPayload struct {
	FollowRequestID int    `json:"follow_request_id"`
	FollowStatus    string `json:"follow_status"` // pending
}

type FollowAcceptedPayload struct {
	FollowStatus string `json:"follow_status"` // accepted
}

type PostReactionPayload struct {
	PostID    int    `json:"post_id"`
	PostTitle string `json:"post_title"`
	CreatedAt string `json:"created_at"`
	Reaction  string `json:"reaction"` // like/dislike
}

type CommentPayload struct {
	PostID    int    `json:"post_id"`
	CommentID int    `json:"comment_id"`
	Snippet   string `json:"snippet,omitempty"`
	CreatedAt string `json:"created_at"`
}

type GroupInvitePayload struct {
	GroupID      int    `json:"group_id"`
	GroupName    string `json:"group_name"`
	GroupAvatar  string `json:"group_avatar,omitempty"`
	InvitationID int    `json:"invitation_id"`
}

type GroupJoinRequestPayload struct {
	GroupID      int `json:"group_id"`
	SenderID     int `json:"sender_id"`
	InvitationID int `json:"invitation_id"`
}

type GroupEventPayload struct {
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
	EventID   int    `json:"event_id"`
}

// GetNotifications fetches the current user's notifications.
// @Summary Get notifications
// @Description Returns notifications for the logged-in user. Can filter by type (all or unread).
// @Tags Notifications
// @Produce json
// @Param type query string false "Notification type filter: all or unread" Enums(all, unread)
// @Success 200 {array} NotifResponse "Notifications fetched"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/notifications [get]
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusBadRequest, "method not allowed", nil)
		return
	}
	notifType := r.URL.Query().Get("type")
	if notifType == "" {
		notifType = "all"
	}
	if !(notifType == "unread" || notifType == "all") {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)

	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	notifs, err := Repos.Notification.GetByUserID(userID, notifType)

	// start enrish the notifications
	notifRes := EnrishNotif(notifs)
	if err != nil {
		fmt.Println("errors", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch notifications", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "notifications fetched", notifRes)
}

// DeletNotif deletes a single notification by ID.
// @Summary Delete a notification
// @Description Deletes a single notification belonging to the current user.
// @Tags Notifications
// @Produce json
// @Param id query int true "Notification ID to delete"
// @Success 200 {object} map[string]string "Notification deleted"
// @Failure 400 {object} map[string]string "Invalid notification ID"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/notifications/delete [delete]
func DeletNotif(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid notification id", nil)
		return
	}

	if err := Repos.Notification.DeleteByIDAndUserID(id, userID); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not delete notification", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "notification deleted", nil)
}

// DeletAllNotif deletes all notifications for the current user.
// @Summary Delete all notifications
// @Description Deletes every notification belonging to the current user.
// @Tags Notifications
// @Produce json
// @Success 200 {object} map[string]string "All notifications deleted"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/notifications/deletAll [delete]
func DeletAllNotif(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	if err := Repos.Notification.DeleteAllByUserID(userID); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not delete notifications", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "all notifications deleted", nil)
}

// MarkNotificationRead marks a single notification as read.
// @Summary Mark notification as read
// @Description Marks a single notification as read by its ID.
// @Tags Notifications
// @Produce json
// @Param id query int true "Notification ID to mark as read"
// @Success 200 {object} map[string]string "Notification marked as read"
// @Failure 400 {object} map[string]string "Invalid notification ID"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/notifications/read [post]
func MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid notification id", nil)
		return
	}

	if err := Repos.Notification.MarkAsRead(id); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not mark notification as read", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "notification marked as read", nil)
}

// MarkAllNotificationsRead marks all notifications as read for the current user.
// @Summary Mark all notifications as read
// @Description Marks every notification belonging to the current user as read.
// @Tags Notifications
// @Produce json
// @Success 200 {object} map[string]string "All notifications marked as read"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/notifications/readAll [post]
func MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	if err := Repos.Notification.MarkAllAsReadByUserID(userID); err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not mark all notifications as read", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "all notifications marked as read", nil)
}

func EnrishNotif(notifs []repository.Notification) []NotifResponse {
	enriched := make([]NotifResponse, 0, len(notifs))

	for _, n := range notifs {
		var actor ActorInfo

		if user, err := Repos.Profile.GetPublicProfile(n.ActorID); err == nil && user != nil {
			actor = ActorInfo{
				UserID:    user.ID,
				Nickname:  user.Nickname,
				Avatar:    user.Avatar,
				Firstname: user.Firstname,
				Lastname:  user.Lastname,
			}
		}

		enriched = append(enriched, NotifResponse{
			ID:         n.ID,
			Type:       n.Type,
			ObjectType: n.ObjectType,
			IsRead:     n.IsRead,
			CreatedAt:  n.CreatedAt,
			Actor:      actor,
			Payload:    buildNotifPayload(n),
		})
	}

	return enriched
}

func buildNotifPayload(n repository.Notification) interface{} {
	switch n.Type {

	case "post_reaction":
		postData, err := Repos.Post.GetPostByID(n.ObjectID)
		if err != nil {
			return nil
		}
		return PostReactionPayload{
			PostID:    n.ObjectID,
			PostTitle: "My first post title",
			CreatedAt: postData.CreatedAt.String(),
			Reaction:  "like",
		}

	case "comment":
		commentData, err := Repos.Comment.GetCommentByID(n.ObjectID)
		if err != nil {
			return nil
		}
		return CommentPayload{
			PostID:    commentData.PostID,
			CommentID: n.ObjectID,
			Snippet:   commentData.Text,
			CreatedAt: commentData.CreatedAt.String(),
		}

	case "follow_request":
		return FollowRequestPayload{
			FollowRequestID: n.ObjectID,
			FollowStatus:    "pending",
		}

	case "follow_accepted":
		return FollowAcceptedPayload{
			FollowStatus: "accepted",
		}

	case "group_invite":
		groupdata, err := Repos.Group.GetPublicGroupDetails(n.ObjectID)
		if err != nil {
			return nil
		}
		return GroupInvitePayload{
			GroupID:      groupdata.ID,
			GroupName:    groupdata.Title,
			GroupAvatar:  "https://fake-image.com/group.png",
			InvitationID: n.ObjectID,
		}

	case "group_join_request":
		groupdata, err := Repos.Group.GetPublicGroupDetails(n.ObjectID)
		if err != nil {
			return nil
		}
		return GroupJoinRequestPayload{
			GroupID:      n.ObjectID,
			SenderID:     n.ActorID,
			InvitationID: groupdata.ID,
		}

	case "group_event":
		event, err := Repos.Group.GetGroupEventByID(n.ObjectID)
		if err != nil {
			return nil
		}
		groupdata, err := Repos.Group.GetPublicGroupDetails(event.GroupID)
		if err != nil {
			return nil
		}
		return GroupEventPayload{
			GroupID:   event.GroupID,
			GroupName: groupdata.Title,
			EventID:   event.ID,
		}

	default:
		return nil
	}
}
