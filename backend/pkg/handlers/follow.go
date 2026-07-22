package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"01social/pkg/middlewares"
	"01social/pkg/repository"
	"01social/pkg/utilities"
	"01social/pkg/ws"
)

// this func handle the follow logic
// FollowResolver handles follow/unfollow/accept/reject actions.
// @Summary Follow operations
// @Description Follow, unfollow, accept follow request, or reject follow request.
// @Tags Follow
// @Produce json
// @Param resolver path string true "Action: follow, unfollow, accept, or reject" Enums(follow, unfollow, accept, reject)
// @Param id path int true "Target user ID"
// @Success 200 {object} map[string]any "Action successful"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/follow/{resolver}/{id} [put]
func FollowResolver(w http.ResponseWriter, r *http.Request) {
	// Only allow PUT requests
	if r.Method != http.MethodPut {
		utilities.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
			nil,
		)
		return
	}

	// Expected format:
	// /api/follow/{resolver}/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 4 ||
		parts[0] != "api" ||
		parts[1] != "follow" {

		utilities.WriteJSON(
			w,
			http.StatusNotFound,
			"invalid endpoint",
			nil,
		)
		return
	}

	resolver := parts[2]

	// Validate resolver
	switch resolver {
	case "follow", "unfollow", "accept", "reject":
		break
	default:
		utilities.WriteJSON(
			w,
			http.StatusBadRequest,
			"invalid follow action",
			nil,
		)
		return
	}

	// Convert target id
	targetID, err := strconv.Atoi(parts[3])
	if err != nil {
		utilities.WriteJSON(
			w,
			http.StatusBadRequest,
			"invalid user id",
			nil,
		)
		return
	}

	// Get logged user id
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(
			w,
			http.StatusUnauthorized,
			"not logged in",
			nil,
		)
		return
	}

	// Prevent following yourself
	if targetID == userID {
		utilities.WriteJSON(
			w,
			http.StatusBadRequest,
			"you cannot follow yourself",
			nil,
		)
		return
	}

	switch resolver {

	case "follow":

		// Check if target user is private
		isPrivate, err := Repos.User.IsPrivateUser(targetID)
		if err != nil {
			utilities.WriteJSON(
				w,
				http.StatusNotFound,
				"user not found",
				nil,
			)
			return
		}

		status := "accepted"

		// Private account requires approval
		if isPrivate {
			status = "pending"
		}

		err = Repos.Follow.Follow(
			userID,
			targetID,
			status,
		)
		if err != nil {
			utilities.WriteJSON(
				w,
				http.StatusInternalServerError,
				"could not follow user",
				nil,
			)
			return
		}

		// Notify based on follow type:
		// - Public follow: "new_follower" (someone started following you)
		// - Private follow (pending): "follow_request" (someone wants to follow you)
		notifType := "new_follower"
		if status == "pending" {
			notifType = "follow_request"
		}

		// Remove any existing notification of this type to prevent duplicates
		if err := Repos.Notification.DeleteNotificationsByTypeAndObject(targetID, notifType, userID); err != nil {
			fmt.Printf("failed to clean up old follow notification: %v\n", err)
		}

		// Create DB notification for the target user
		if err := Repos.Notification.Create(&repository.Notification{
			UserID:     targetID,
			ActorID:    userID,
			Type:       notifType,
			ObjectType: "follow",
			ObjectID:   userID,
		}); err != nil {
			fmt.Printf("failed to create follow notification: %v\n", err)
		}

		// Notify target user via WS
		if requester, err := Repos.User.GetByID(userID); err == nil {
			ws.NotifyUser(strconv.Itoa(targetID), notifType, map[string]any{
				"user_id":   userID,
				"nickname":  requester.Nickname,
				"avatar":    requester.Avatar,
				"firstname": requester.Firstname,
				"lastname":  requester.Lastname,
			})
		}
		utilities.WriteJSON(
			w,
			http.StatusOK,
			"follow request sent",
			map[string]any{"status": status},
		)
		return

	case "unfollow":

		err := Repos.Follow.Unfollow(
			userID,
			targetID,
		)
		if err != nil {
			utilities.WriteJSON(
				w,
				http.StatusInternalServerError,
				"could not unfollow user",
				nil,
			)
			return
		}
		// Clean up any existing follow notifications between these two users
		if err := Repos.Notification.DeleteFollowNotification(userID, targetID); err != nil {
			fmt.Printf("failed to clean up follow notification: %v\n", err)
		}

		utilities.WriteJSON(
			w,
			http.StatusOK,
			"user unfollowed",
			map[string]any{"status": "none"},
		)
		return

	case "accept":

		// targetID is the requester
		// userID is the owner accepting
		err := Repos.Follow.AcceptFollow(
			targetID,
			userID,
		)
		if err != nil {
			utilities.WriteJSON(
				w,
				http.StatusNotFound,
				err.Error(),
				nil,
			)
			return
		}

		// Remove the follow_request notification for the owner (userID) since it's been handled
		if err := Repos.Notification.DeleteNotificationsByTypeAndObject(userID, "follow_request", targetID); err != nil {
			fmt.Printf("failed to clean up follow_request notification: %v\n", err)
		}

		// Notify the requester that their follow was accepted
		if err := Repos.Notification.Create(&repository.Notification{
			UserID:     targetID, // the requester gets notified
			ActorID:    userID,   // the owner who accepted
			Type:       "follow_accepted",
			ObjectType: "follow",
			ObjectID:   targetID,
		}); err != nil {
			fmt.Printf("failed to create follow_accepted notification: %v\n", err)
		}

		// Notify via WS
		if acceptor, err := Repos.User.GetByID(userID); err == nil {
			ws.NotifyUser(strconv.Itoa(targetID), "follow_accepted", map[string]any{
				"user_id":  userID,
				"nickname": acceptor.Nickname,
			})
		}

		utilities.WriteJSON(
			w,
			http.StatusOK,
			"follow accepted",
			nil,
		)
		return

	case "reject":

		// targetID is the requester
		// userID is the owner rejecting
		err := Repos.Follow.RejectFollow(
			targetID,
			userID,
		)
		if err != nil {
			utilities.WriteJSON(
				w,
				http.StatusInternalServerError,
				"could not reject follow",
				nil,
			)
			return
		}

		utilities.WriteJSON(
			w,
			http.StatusOK,
			"follow rejected",
			nil,
		)
		return
	}
}
