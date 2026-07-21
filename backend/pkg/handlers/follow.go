package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"01social/pkg/middlewares"
	"01social/pkg/repository"
	"01social/pkg/utilities"
)

// this func handle the follow logic
/*
- if user is public it follow directly
- if user is private it make it pending
- user can cancel the follow by removing it
*/
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

		if status == "pending" {
			if err := Repos.Notification.Create(&repository.Notification{
				UserID:     targetID, // receiver
				ActorID:    userID,   // requester
				Type:       "follow_request",
				ObjectType: "follow",
				ObjectID:   userID, // or the FOLLOWS row ID if you add one
			}); err != nil {
				fmt.Printf("failed to create notification: %v\n", err)
			}
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
			map[string]any{"status": "non"},
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
