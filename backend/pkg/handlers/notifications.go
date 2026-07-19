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
3-"type": "follow_request", req
"object_type": "follow",
4-"type": "follow_accepted",
    "object_type": "follow",
5-"type": "group_invite",   req
            "object_type": "group_invite",
6-"type": "group_join_request",    user a find group b and send a request join
    "object_type": "group_request",
*/

/*payload data on each type
1-follow_request
 ->{follow request ID
    + follow status
 }
2-group_invite
->{
groupid,grouptitle,avatar,invitationID
}
3-group_join_request (later + not required)
 ->{
     user senderid + invitationid
 }
	 4-postreaction or comment reaction={
	only like + post title,createt at}
	5-new comment=>

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
			// skippe or ignore the notification

			return nil
		}
		return CommentPayload{
			PostID:    commentData.PostID,
			CommentID: n.ObjectID, // fake for now
			Snippet:   commentData.Text,
			CreatedAt: commentData.CreatedAt.String(),
		}

	case "follow_request":
		return FollowRequestPayload{
			FollowRequestID: n.ObjectID, // fake
			FollowStatus:    "pending",
		}

	case "follow_accepted":
		return FollowAcceptedPayload{
			FollowStatus: "accepted",
		}

	case "group_invite":
		groupdata, err := Repos.Group.GetPublicGroupDetails(n.ObjectID)
		if err != nil {
			// skippe or ignore the notification

			return nil
		}
		return GroupInvitePayload{
			GroupID:      groupdata.ID,
			GroupName:    groupdata.Title,
			GroupAvatar:  "https://fake-image.com/group.png",
			InvitationID: n.ObjectID, // fake
		}

	case "group_join_request":
		groupdata, err := Repos.Group.GetPublicGroupDetails(n.ObjectID)
		if err != nil {
			// skippe or ignore the notification

			return nil
		}
		return GroupJoinRequestPayload{
			GroupID:      n.ObjectID,
			SenderID:     n.ActorID,
			InvitationID: groupdata.ID, // fake
		}

	default:
		return nil
	}
}
