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

type CommentPayload struct {
	PostID    int    `json:"post_id"`
	CommentID int    `json:"comment_id"`
	Snippet   string `json:"snippet,omitempty"`
}

type FollowRequestPayload struct {
	FollowStatus string `json:"follow_status"` // "pending"
}

type FollowAcceptedPayload struct {
	FollowStatus string `json:"follow_status"` // "accepted"
}
type PostReactionPayload struct {
	PostID    int    `json:"post_id"`
	PostImage string `json:"post_image,omitempty"`
	Reaction  string `json:"reaction"` // e.g. "like" | "dislike"
}

type GroupInvitePayload struct {
	GroupID     int    `json:"group_id"`
	GroupName   string `json:"group_name"`
	GroupAvatar string `json:"group_avatar,omitempty"`
}

type GroupJoinRequestPayload struct {
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
}

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusBadRequest, "method not allowed", nil)
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

	fmt.Println("get notification for user", notifType, userID)
	notifs, err := Repos.Notification.GetByUserID(userID, notifType)

	// start enrish the notifications
	notifRes := EnrishNotif(notifs)
	fmt.Println("", notifRes)
	if err != nil {
		fmt.Println("errors", err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch notifications", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "notifications fetched", notifRes)
}

func DeletNotif(w http.ResponseWriter, r *http.Request) {
	utilities.WriteJSON(w, http.StatusAccepted, "start delete one notif", nil)
}

func DeletAllNotif(w http.ResponseWriter, r *http.Request) {
	utilities.WriteJSON(w, http.StatusAccepted, "start delete one notif", nil)
}

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

/*
actor:{
avatar,nickname,first and last name}
*/
/*
all types
1-  "type": "post_reaction",
     "object_type": "post",
2- "type": "comment",
     "object_type": "comment",
3-"type": "follow_request",
"object_type": "follow",
4-"type": "follow_accepted",
    "object_type": "follow",
5-"type": "group_invite",
            "object_type": "group_invite",
6-"type": "group_join_request",
    "object_type": "group_request",
*/

func EnrishNotif(notifs []repository.Notification) []NotifResponse {
	enriched := make([]NotifResponse, 0, len(notifs))

	for _, n := range notifs {
		var actor ActorInfo

		// NOTE: assumes repository.Notification has a SenderID field and
		// Repos.User has a GetMinimalByID method — rename to match your
		// actual repo/user model.
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

// buildNotifPayload returns the right typed payload based on n.Type.
// NOTE: assumes repository.Notification has an ObjectID field pointing at
// the related post/comment/group/etc — adjust field names to match your model.
func buildNotifPayload(n repository.Notification) interface{} {
	switch n.Type {
	case "post_reaction":
		return PostReactionPayload{
			PostID: n.ObjectID,
		}
	case "comment":
		return CommentPayload{
			PostID: n.ObjectID,
		}
	case "follow_request":
		return FollowRequestPayload{
			FollowStatus: "pending",
		}
	case "follow_accepted":
		return FollowAcceptedPayload{
			FollowStatus: "accepted",
		}
	case "group_invite":
		return GroupInvitePayload{
			GroupID: n.ObjectID,
		}
	case "group_join_request":
		return GroupJoinRequestPayload{
			GroupID: n.ObjectID,
		}
	default:
		return nil
	}
}
