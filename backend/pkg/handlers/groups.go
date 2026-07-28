package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"01social/pkg/middlewares"
	"01social/pkg/repository"
	"01social/pkg/utilities"
	"01social/pkg/ws"
)

func notifyGroupUsers(
	groupID int,
	eventType string,
	objectType string,
	objectID int,
	actorID int,
	excludedUserIDs ...int,
) {
	memberIDs, err := Repos.Group.GetGroupMemberIDs(groupID)
	if err != nil {
		return
	}

	excluded := make(map[int]struct{}, len(excludedUserIDs))

	for _, userID := range excludedUserIDs {
		excluded[userID] = struct{}{}
	}

	// Fetch group details once (for WS event payload)
	groupData, _ := Repos.Group.GetPublicGroupDetails(groupID)

	for _, userID := range memberIDs {

		if _, skip := excluded[userID]; skip {
			continue
		}

		err := Repos.Notification.Create(&repository.Notification{
			UserID:     userID,
			ActorID:    actorID,
			Type:       eventType,
			ObjectType: objectType,
			ObjectID:   objectID,
		})
		if err != nil {
			fmt.Printf("notification error: %v\n", err)
		}

		// Send live WS event for group events
		if eventType == "group_event" {
			ws.NotifyUser(strconv.Itoa(userID), eventType, map[string]any{
				"event_id":   objectID,
				"group_id":   groupID,
				"group_name": groupData.Title,
				"user_id":    actorID,
			})
		}
	}
}

func parseGroupPath(path string) (groupID int, endpoint string, targetID int, action string, ok bool) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) < 3 || segments[0] != "api" || segments[1] != "groups" {
		return 0, "", 0, "", false
	}

	groupID, err := strconv.Atoi(segments[2])
	if err != nil {
		return 0, "", 0, "", false
	}

	if len(segments) >= 4 {
		endpoint = segments[3]
	}

	if len(segments) >= 6 && endpoint == "requests" {
		targetID, err = strconv.Atoi(segments[4])
		if err == nil {
			action = strings.ToLower(segments[5])
		}
		return groupID, endpoint, targetID, action, true
	}

	if len(segments) >= 5 && endpoint == "invites" {
		action = strings.ToLower(segments[4])
		return groupID, endpoint, 0, action, true
	}

	return groupID, endpoint, 0, "", true
}

// CreateGroup creates a new group.
// @Summary Create a new group
// @Description Creates a new group with title, description, and optional logo/background images.
// @Tags Groups
// @Accept mpfd
// @Produce json
// @Param title formData string true "Group title"
// @Param description formData string false "Group description"
// @Param logo formData file false "Group logo image"
// @Param background formData file false "Group background image"
// @Success 201 {object} map[string]any "Group created successfully"
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/groups/create [post]
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	log.Println("[CREATE_GROUP] Start group creation process")

	if r.Method != http.MethodPost {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	// 1. Check User Session / Auth
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		log.Printf("[CREATE_GROUP] Unauthorized attempt to create group")
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	// 2. Limit request body size to 10 MB (for logo + background uploads combined)
	const maxUploadSize int64 = 10 << 20 // 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Printf("[CREATE_GROUP] Failed parsing multipart form (file limit exceeded): %v", err)
		utilities.WriteJSON(w, http.StatusBadRequest, "Max combined file size is 2MB", nil)
		return
	}

	// 3. Extract & Validate Required Fields
	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))

	if title == "" {
		log.Printf("[CREATE_GROUP] Validation failed: empty title field by User %d", userID)
		utilities.WriteJSON(w, http.StatusBadRequest, "title is required", nil)
		return
	}

	if !utilities.IsValidName(title) { // Or your custom string check
		log.Printf("[CREATE_GROUP] Validation failed: invalid title format %q", title)
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid group title", nil)
		return
	}

	if description != "" && !utilities.IsValidDescription(description) {
		log.Printf("[CREATE_GROUP] Validation failed: description exceeds length limit")
		utilities.WriteJSON(w, http.StatusBadRequest, "description is too long", nil)
		return
	}

	// 4. Handle Optional Logo Upload
	var logoPath string
	logoFile, logoHeader, err := r.FormFile("logo")
	if err == nil {
		defer logoFile.Close()
		path, err := utilities.SaveImage(logoFile, logoHeader, "uploads/groups/logos")
		if err != nil {
			log.Printf("[CREATE_GROUP] Failed to process logo upload: %v", err)
			utilities.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		logoPath = path
	} else if err != http.ErrMissingFile {
		log.Printf("[CREATE_GROUP] Non-standard error reading logo file: %v", err)
	}

	// 5. Handle Optional Background Upload
	var backgroundPath string
	bgFile, bgHeader, err := r.FormFile("background")
	if err == nil {
		defer bgFile.Close()
		path, err := utilities.SaveImage(bgFile, bgHeader, "uploads/groups/backgrounds")
		if err != nil {
			log.Printf("[CREATE_GROUP] Failed to process background upload: %v", err)
			utilities.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		backgroundPath = path
	} else if err != http.ErrMissingFile {
		log.Printf("[CREATE_GROUP] Non-standard error reading background file: %v", err)
	}

	// 6. Construct Model & Save via Repository
	group := &repository.Group{
		CreatorID:   userID,
		Title:       title,
		Description: description,
		Logo:        logoPath,
		Background:  backgroundPath,
	}

	if err := Repos.Group.CreateGroup(group); err != nil {
		log.Printf("[CREATE_GROUP] DB insertion failed for group %q by user %d: %v", title, userID, err)
		utilities.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	log.Printf("[CREATE_GROUP] Successfully created group %q (ID: %d) by User %d", group.Title, group.ID, userID)

	// 7. Success Response
	utilities.WriteJSON(w, http.StatusCreated, "group created success", map[string]any{
		"group_id": group.ID,
		"title":    group.Title,
	})
}

// ListGroups fetches all groups.
// @Summary List all groups
// @Description Returns a list of all groups.
// @Tags Groups
// @Produce json
// @Success 200 {array} repository.Group "Groups fetched"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/groups [get]
func ListGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	groups, err := Repos.Group.ListGroups()
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch groups", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "groups fetched", groups)
}

// GroupResolver handles group operations including posts, events, messages, invites, and membership management.
// @Summary Group operations resolver
// @Description Handles various group operations: posts CRUD, events CRUD, messages, invites, join/leave, member management.
// @Tags Groups
// @Accept json
// @Produce json
// @Param groupId path int true "Group ID"
// @Success 200 {object} map[string]any "Operation successful"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Not logged in"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /api/groups/{groupId}/posts [get]
// @Router /api/groups/{groupId}/posts [post]
// @Router /api/groups/{groupId}/events [get]
// @Router /api/groups/{groupId}/events [post]
// @Router /api/groups/{groupId}/messages [get]
// @Router /api/groups/{groupId}/messages [post]
// @Router /api/groups/{groupId}/join [post]
// @Router /api/groups/{groupId}/leave [post]
// @Router /api/groups/{groupId}/invite [post]
// @Router /api/groups/{groupId}/update [put]
// @Router /api/groups/{groupId}/requests [get]
// @Router /api/groups/{groupId}/requests/{userId}/accept [post]
// @Router /api/groups/{groupId}/requests/{userId}/reject [post]
// @Router /api/groups/{groupId}/invites/accept [post]
// @Router /api/groups/{groupId}/invites/reject [post]
// @Router /api/groups/{groupId}/members/{userId}/kick [post]
func GroupResolver(w http.ResponseWriter, r *http.Request) {
	groupID, endpoint, _, _, ok := parseGroupPath(r.URL.Path)
	if !ok {
		utilities.WriteJSON(w, http.StatusNotFound, "not found", nil)
		return
	}

	segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(segments) < 3 || segments[0] != "api" || segments[1] != "groups" {
		utilities.WriteJSON(w, http.StatusNotFound, "not found", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	if endpoint == "" && len(segments) >= 4 {
		endpoint = segments[3]
	}

	if len(segments) >= 6 && segments[3] == "posts" && segments[5] == "reaction" {
		member, err := Repos.Group.IsGroupMember(groupID, userID)
		if err != nil || !member {
			utilities.WriteJSON(w, http.StatusForbidden, "not a group member", nil)
			return
		}

		postID, err := strconv.Atoi(segments[4])
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid post id", nil)
			return
		}

		belongs, err := Repos.Group.DoesGroupPostBelongToGroup(groupID, postID)
		if err != nil || !belongs {
			utilities.WriteJSON(w, http.StatusNotFound, "post not found in this group", nil)
			return
		}

		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		var payload struct {
			IsLike int `json:"is_like"`
		}
		if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}

		code, err := ReactToGroupPost(userID, postID, payload.IsLike)
		if err != nil {
			if code == http.StatusBadRequest || code == http.StatusNotFound {
				utilities.WriteJSON(w, code, err.Error(), nil)
			} else {
				utilities.WriteJSON(w, http.StatusInternalServerError, "could not process reaction", nil)
			}
			return
		}

		utilities.WriteJSON(w, http.StatusOK, "reaction saved", nil)
		return
	}

	if len(segments) >= 6 && segments[3] == "posts" && segments[5] == "comments" {
		member, err := Repos.Group.IsGroupMember(groupID, userID)
		if err != nil || !member {
			utilities.WriteJSON(w, http.StatusForbidden, "not a group member", nil)
			return
		}
		postID, err := strconv.Atoi(segments[4])
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid post id", nil)
			return
		}

		belongs, err := Repos.Group.DoesGroupPostBelongToGroup(groupID, postID)
		if err != nil || !belongs {
			utilities.WriteJSON(w, http.StatusNotFound, "post not found in this group", nil)
			return
		}

		// Handle comment reactions: /api/groups/{id}/posts/{postId}/comments/{commentId}/reaction
		if len(segments) >= 8 && segments[7] == "reaction" {
			if r.Method != http.MethodPost {
				utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
				return
			}
			commentID, err := strconv.Atoi(segments[6])
			if err != nil {
				utilities.WriteJSON(w, http.StatusBadRequest, "invalid comment id", nil)
				return
			}
			var payload struct {
				IsLike int `json:"is_like"`
			}
			if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
				utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
				return
			}
			if err := Repos.Group.ReactToGroupPostComment(commentID, userID, payload.IsLike); err != nil {
				utilities.WriteJSON(w, http.StatusInternalServerError, "could not react to comment", nil)
				return
			}
			utilities.WriteJSON(w, http.StatusOK, "reaction saved", nil)
			return
		}

		// Handle comment deletion: /api/groups/{id}/posts/{postId}/comments/{commentId}/delete
		if len(segments) >= 8 && segments[7] == "delete" {
			if r.Method != http.MethodPost {
				utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
				return
			}
			commentID, err := strconv.Atoi(segments[6])
			if err != nil {
				utilities.WriteJSON(w, http.StatusBadRequest, "invalid comment id", nil)
				return
			}
			creatorID, err := Repos.Group.GetGroupCreatorID(groupID)
			if err != nil {
				utilities.WriteJSON(w, http.StatusNotFound, "group not found", nil)
				return
			}
			if err := Repos.Group.DeleteGroupPostComment(commentID, userID, creatorID); err != nil {
				utilities.WriteJSON(w, http.StatusForbidden, err.Error(), nil)
				return
			}
			utilities.WriteJSON(w, http.StatusOK, "comment deleted", nil)
			return
		}

		if r.Method == http.MethodGet {
			comments, err := Repos.Group.ListGroupPostComments(postID)
			if err != nil {
				utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch comments", nil)
				return
			}
			utilities.WriteJSON(w, http.StatusOK, "comments fetched", comments)
			return
		}
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}
		var payload struct {
			Text string `json:"text"`
		}
		if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		payload.Text = strings.TrimSpace(payload.Text)
		if payload.Text == "" {
			utilities.WriteJSON(w, http.StatusBadRequest, "comment text is required", nil)
			return
		}
		if err := Repos.Group.CreateGroupPostComment(postID, userID, payload.Text); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not create comment", nil)
			return
		}
		utilities.WriteJSON(w, http.StatusCreated, "comment created", nil)
		return
	}

	if len(segments) >= 6 && segments[3] == "events" && segments[5] == "respond" {
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		member, err := Repos.Group.IsGroupMember(groupID, userID)
		if err != nil || !member {
			utilities.WriteJSON(w, http.StatusForbidden, "not a group member", nil)
			return
		}

		eventID, err := strconv.Atoi(segments[4])
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid event id", nil)
			return
		}

		belongs, err := Repos.Group.DoesGroupEventBelongToGroup(groupID, eventID)
		if err != nil || !belongs {
			utilities.WriteJSON(w, http.StatusNotFound, "event not found in this group", nil)
			return
		}

		var payload struct {
			Status string `json:"status"`
		}
		if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		payload.Status = strings.ToLower(strings.TrimSpace(payload.Status))
		if payload.Status != "going" && payload.Status != "not_going" {
			utilities.WriteJSON(w, http.StatusBadRequest, "status must be going or not_going", nil)
			return
		}
		if err := Repos.Group.RespondToEvent(eventID, userID, payload.Status); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not respond to event", nil)
			return
		}
		utilities.WriteJSON(w, http.StatusOK, "event response saved", nil)
		return
	}

	// Delete a group post: POST /api/groups/{groupId}/posts/{postId}/delete
	if len(segments) >= 6 && segments[3] == "posts" && segments[5] == "delete" {
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}
		postID, err := strconv.Atoi(segments[4])
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid post id", nil)
			return
		}
		belongs, err := Repos.Group.DoesGroupPostBelongToGroup(groupID, postID)
		if err != nil || !belongs {
			utilities.WriteJSON(w, http.StatusNotFound, "post not found in this group", nil)
			return
		}
		creatorID, err := Repos.Group.GetGroupCreatorID(groupID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusNotFound, "group not found", nil)
			return
		}
		if err := Repos.Group.DeleteGroupPost(postID, userID, creatorID); err != nil {
			utilities.WriteJSON(w, http.StatusForbidden, err.Error(), nil)
			return
		}
		utilities.WriteJSON(w, http.StatusOK, "post deleted", nil)
		return
	}

	// Delete a group event: POST /api/groups/{groupId}/events/{eventId}/delete
	if len(segments) >= 6 && segments[3] == "events" && segments[5] == "delete" {
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}
		eventID, err := strconv.Atoi(segments[4])
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid event id", nil)
			return
		}
		belongs, err := Repos.Group.DoesGroupEventBelongToGroup(groupID, eventID)
		if err != nil || !belongs {
			utilities.WriteJSON(w, http.StatusNotFound, "event not found in this group", nil)
			return
		}
		creatorID, err := Repos.Group.GetGroupCreatorID(groupID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusNotFound, "group not found", nil)
			return
		}
		if err := Repos.Group.DeleteGroupEvent(eventID, userID, creatorID); err != nil {
			utilities.WriteJSON(w, http.StatusForbidden, err.Error(), nil)
			return
		}
		utilities.WriteJSON(w, http.StatusOK, "event deleted", nil)
		return
	}

	// Kick a member: POST /api/groups/{groupId}/members/{userId}/kick
	if len(segments) >= 6 && segments[3] == "members" && segments[5] == "kick" {
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}
		memberID, err := strconv.Atoi(segments[4])
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid user id", nil)
			return
		}
		if err := Repos.Group.KickMember(groupID, memberID, userID); err != nil {
			utilities.WriteJSON(w, http.StatusForbidden, err.Error(), nil)
			return
		}
		// Remove notifications about this group for the kicked user
		_ = Repos.Notification.DeleteNotificationsByTypeAndObject(memberID, "group_invite", groupID)
		_ = Repos.Notification.DeleteNotificationsByTypeAndObject(memberID, "group_join_request", groupID)
		utilities.WriteJSON(w, http.StatusOK, "member kicked", nil)
		return
	}

	// GET /api/groups/{groupId}/requests — list pending join requests
	if len(segments) == 4 && segments[3] == "requests" {
		if r.Method != http.MethodGet {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		creatorID, err := Repos.Group.GetGroupCreatorID(groupID)
		if err != nil || creatorID != userID {
			utilities.WriteJSON(w, http.StatusForbidden, "only the group creator can view requests", nil)
			return
		}

		requests, err := Repos.Group.ListGroupPendingRequests(groupID)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch requests", nil)
			return
		}

		utilities.WriteJSON(w, http.StatusOK, "pending requests fetched", requests)
		return
	}

	if len(segments) >= 6 && segments[3] == "requests" {
		targetUserID, err := strconv.Atoi(segments[4])
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid user id", nil)
			return
		}
		action := strings.ToLower(segments[5])
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		// Only the group creator can accept/reject join requests
		creatorID, err := Repos.Group.GetGroupCreatorID(groupID)
		if err != nil || creatorID != userID {
			utilities.WriteJSON(w, http.StatusForbidden, "only the group creator can manage requests", nil)
			return
		}

		switch action {
		case "accept":
			if err := Repos.Group.AcceptGroupRequest(groupID, targetUserID); err != nil {
				if err.Error() == "no pending request found" {
					utilities.WriteJSON(w, http.StatusNotFound, "no pending request from this user", nil)
				} else {
					utilities.WriteJSON(w, http.StatusInternalServerError, "could not accept request", nil)
				}
				return
			}
			// Remove the join request notification for the creator
			_ = Repos.Notification.DeleteNotificationsByTypeAndObject(userID, "group_join_request", groupID)
			utilities.WriteJSON(w, http.StatusOK, "join request accepted", nil)
		case "reject":
			if err := Repos.Group.RejectGroupRequest(groupID, targetUserID); err != nil {
				if err.Error() == "no pending request found" {
					utilities.WriteJSON(w, http.StatusNotFound, "no pending request from this user", nil)
				} else {
					utilities.WriteJSON(w, http.StatusInternalServerError, "could not reject request", nil)
				}
				return
			}
			// Remove the join request notification for the creator
			_ = Repos.Notification.DeleteNotificationsByTypeAndObject(userID, "group_join_request", groupID)
			utilities.WriteJSON(w, http.StatusOK, "join request rejected", nil)
		default:
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid action", nil)
		}
		return
	}

	if len(segments) >= 5 && segments[3] == "invites" {
		action := strings.ToLower(segments[4])
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}
		switch action {
		case "accept":
			if err := Repos.Group.AcceptGroupInvite(groupID, userID); err != nil {
				if err.Error() == "no pending invite found" {
					utilities.WriteJSON(w, http.StatusNotFound, "no pending invite for this group", nil)
				} else {
					utilities.WriteJSON(w, http.StatusInternalServerError, "could not accept invite", nil)
				}
				return
			}
			// Remove the group invite notification
			_ = Repos.Notification.DeleteNotificationsByTypeAndObject(userID, "group_invite", groupID)
			utilities.WriteJSON(w, http.StatusOK, "group invite accepted", nil)
		case "reject":
			if err := Repos.Group.RejectGroupInvite(groupID, userID); err != nil {
				if err.Error() == "no pending invite found" {
					utilities.WriteJSON(w, http.StatusNotFound, "no pending invite for this group", nil)
				} else {
					utilities.WriteJSON(w, http.StatusInternalServerError, "could not reject invite", nil)
				}
				return
			}
			// Remove the group invite notification
			_ = Repos.Notification.DeleteNotificationsByTypeAndObject(userID, "group_invite", groupID)
			utilities.WriteJSON(w, http.StatusOK, "group invite rejected", nil)
		default:
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid action", nil)
		}
		return
	}

	switch endpoint {
	case "update":
		if r.Method != http.MethodPut {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		creatorID, err := Repos.Group.GetGroupCreatorID(groupID)
		if err != nil || creatorID != userID {
			utilities.WriteJSON(w, http.StatusForbidden, "only the group creator can update", nil)
			return
		}

		var title, description, logoPath, backgroundPath string

		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			const maxUploadSize int64 = 10 << 20
			r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
			if err := r.ParseMultipartForm(maxUploadSize); err != nil {
				utilities.WriteJSON(w, http.StatusBadRequest, "file too large or invalid form", nil)
				return
			}
			title = strings.TrimSpace(r.FormValue("title"))
			description = strings.TrimSpace(r.FormValue("description"))

			if logoFile, logoHeader, err := r.FormFile("logo"); err == nil {
				defer logoFile.Close()
				path, err := utilities.SaveImage(logoFile, logoHeader, "uploads/groups/logos")
				if err != nil {
					utilities.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
					return
				}
				logoPath = path
			}

			if bgFile, bgHeader, err := r.FormFile("background"); err == nil {
				defer bgFile.Close()
				path, err := utilities.SaveImage(bgFile, bgHeader, "uploads/groups/backgrounds")
				if err != nil {
					utilities.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
					return
				}
				backgroundPath = path
			}
		} else {
			var payload struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Logo        string `json:"logo"`
				Background  string `json:"background"`
			}
			if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
				utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
				return
			}
			title = strings.TrimSpace(payload.Title)
			description = strings.TrimSpace(payload.Description)
		}

		if err := Repos.Group.UpdateGroup(groupID, userID, title, description, logoPath, backgroundPath); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		utilities.WriteJSON(w, http.StatusOK, "group updated", nil)

	case "messages":
		if r.Method == http.MethodGet {
			limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
			if err != nil || limit <= 0 {
				limit = 20
			}
			if limit > 50 {
				limit = 50
			}
			offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
			if err != nil || offset < 0 {
				offset = 0
			}

			member, err := Repos.Group.IsGroupMember(groupID, userID)
			if err != nil || !member {
				utilities.WriteJSON(w, http.StatusForbidden, "not a group member", nil)
				return
			}

			messages, err := Repos.Group.ListGroupMessages(groupID, limit, offset)
			if err != nil {
				utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch group messages", nil)
				return
			}
			utilities.WriteJSON(w, http.StatusOK, "group messages fetched", messages)
			return
		}

		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		member, err := Repos.Group.IsGroupMember(groupID, userID)
		if err != nil || !member {
			utilities.WriteJSON(w, http.StatusForbidden, "not a group member", nil)
			return
		}

		var payload struct {
			Text string `json:"text"`
		}
		if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		payload.Text = strings.TrimSpace(payload.Text)
		if payload.Text == "" {
			utilities.WriteJSON(w, http.StatusBadRequest, "message text is required", nil)
			return
		}

		msg, err := Repos.Group.CreateGroupMessage(groupID, userID, payload.Text)
		if err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not send group message", nil)
			return
		}

		// Fetch sender profile for the WS nickname
		senderProfile, _ := Repos.User.GetByID(userID)
		senderNickname := ""
		if senderProfile != nil {
			senderNickname = senderProfile.Nickname
			if senderNickname == "" {
				senderNickname = senderProfile.Firstname
				if senderProfile.Lastname != "" {
					senderNickname += " " + senderProfile.Lastname
				}
			}
		}

		notifyGroupUsers(
			groupID,
			"group_message",
			"group_message",
			msg.ID,
			userID,
			userID,
		)
		for _, memberID := range []int{userID} {
			_ = memberID
		}
		memberIDs, err := Repos.Group.GetGroupMemberIDs(groupID)
		if err == nil {
			for _, memberID := range memberIDs {
				if memberID == userID {
					continue
				}
				ws.NotifyUser(strconv.Itoa(memberID), "group_message", map[string]any{
					"group_id":   groupID,
					"message_id": msg.ID,
					"sender_id":  userID,
					"nickname":   senderNickname,
					"text":       payload.Text,
				})
			}
		}

		utilities.WriteJSON(w, http.StatusCreated, "group message sent", map[string]any{"message_id": msg.ID})
	case "leave":
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}
		if err := Repos.Group.LeaveGroup(groupID, userID); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		utilities.WriteJSON(w, http.StatusOK, "left group", nil)

	case "join":
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}
		if err := Repos.Group.RequestToJoin(groupID, userID); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not request to join", nil)
			return
		}
		creatorID, err := Repos.Group.GetGroupCreatorID(groupID)
		if err == nil && creatorID > 0 {
			_ = Repos.Notification.Create(&repository.Notification{
				UserID:     creatorID,
				ActorID:    userID,
				Type:       "group_join_request",
				ObjectType: "group_request",
				ObjectID:   groupID,
			})

			// Send live WS notification
			if requester, err := Repos.User.GetByID(userID); err == nil {
				ws.NotifyUser(strconv.Itoa(creatorID), "group_join_request", map[string]any{
					"user_id":  userID,
					"nickname": requester.Nickname,
					"group_id": groupID,
				})
			}
		}
		utilities.WriteJSON(w, http.StatusOK, "join request sent", nil)
	case "invite":
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		// Only group members can invite others
		member, err := Repos.Group.IsGroupMember(groupID, userID)
		if err != nil || !member {
			utilities.WriteJSON(w, http.StatusForbidden, "only group members can invite users", nil)
			return
		}

		var payload struct {
			UserID int `json:"user_id"`
		}
		if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		if payload.UserID <= 0 {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid user id", nil)
			return
		}
		if err := Repos.Group.InviteToGroup(groupID, userID, payload.UserID); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not invite user", nil)
			return
		}
		_ = Repos.Notification.Create(&repository.Notification{
			UserID:     payload.UserID,
			ActorID:    userID,
			Type:       "group_invite",
			ObjectType: "group_invite",
			ObjectID:   groupID,
		})

		// Send live WS invite notification
		groupData, _ := Repos.Group.GetPublicGroupDetails(groupID)
		ws.NotifyUser(strconv.Itoa(payload.UserID), "group_invite", map[string]any{
			"user_id":    userID,
			"group_id":   groupID,
			"group_name": groupData.Title,
		})

		utilities.WriteJSON(w, http.StatusOK, "invite sent", nil)
	case "posts":
		member, err := Repos.Group.IsGroupMember(groupID, userID)
		if err != nil || !member {
			utilities.WriteJSON(w, http.StatusForbidden, "not a group member", nil)
			return
		}
		if r.Method == http.MethodGet {
			posts, err := Repos.Group.ListGroupPosts(groupID)
			if err != nil {
				utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch group posts", nil)
				return
			}
			utilities.WriteJSON(w, http.StatusOK, "group posts fetched", posts)
			return
		}
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		var title, text, imagePath string

		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			// Handle multipart form (with optional image)
			const maxUploadSize int64 = 10 << 20 // 10 MB
			r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
			if err := r.ParseMultipartForm(maxUploadSize); err != nil {
				utilities.WriteJSON(w, http.StatusBadRequest, "file too large or invalid form", nil)
				return
			}
			title = strings.TrimSpace(r.FormValue("title"))
			text = strings.TrimSpace(r.FormValue("text"))

			if imageFile, imageHeader, err := r.FormFile("image"); err == nil {
				defer imageFile.Close()
				path, err := utilities.SaveImage(imageFile, imageHeader, "uploads/groups/posts")
				if err != nil {
					utilities.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
					return
				}
				imagePath = path
			} else if err != http.ErrMissingFile {
				log.Printf("[CREATE_GROUP_POST] Error reading image: %v", err)
			}
		} else {
			// Handle regular JSON body
			var payload struct {
				Title string `json:"title"`
				Text  string `json:"text"`
			}
			if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
				utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
				return
			}
			title = strings.TrimSpace(payload.Title)
			text = strings.TrimSpace(payload.Text)
		}

		post := &repository.GroupPost{
			GroupID:   groupID,
			UserID:    userID,
			Title:     title,
			Text:      text,
			Image:     imagePath,
			CreatedAt: time.Now().UTC(),
		}
		if err := Repos.Group.CreateGroupPost(post); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not create group post", nil)
			return
		}
		utilities.WriteJSON(w, http.StatusCreated, "group post created", map[string]any{"post_id": post.ID})
	case "events":
		member, err := Repos.Group.IsGroupMember(groupID, userID)
		if err != nil || !member {
			utilities.WriteJSON(w, http.StatusForbidden, "not a group member", nil)
			return
		}
		if r.Method == http.MethodGet {
			events, err := Repos.Group.ListGroupEvents(groupID)
			if err != nil {
				utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch group events", nil)
				return
			}
			utilities.WriteJSON(w, http.StatusOK, "group events fetched", events)
			return
		}
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}
		// Parse multipart form if Content-Type is multipart (for image upload), otherwise use JSON
		var eventTitle, eventDesc, eventImg, eventTimeStr string

		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				utilities.WriteJSON(w, http.StatusBadRequest, "failed to parse form", nil)
				return
			}
			eventTitle = strings.TrimSpace(r.FormValue("title"))
			eventDesc = strings.TrimSpace(r.FormValue("description"))
			eventTimeStr = strings.TrimSpace(r.FormValue("event_time"))

			if file, header, err := r.FormFile("image"); err == nil {
				defer file.Close()
				path, err := utilities.SaveImage(file, header, "uploads/events")
				if err != nil {
					utilities.WriteJSON(w, http.StatusBadRequest, err.Error(), nil)
					return
				}
				eventImg = path
			}
		} else {
			var payload struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Image       string `json:"image"`
				EventTime   string `json:"event_time"`
			}
			if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
				utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
				return
			}
			eventTitle = strings.TrimSpace(payload.Title)
			eventDesc = strings.TrimSpace(payload.Description)
			eventImg = strings.TrimSpace(payload.Image)
			eventTimeStr = strings.TrimSpace(payload.EventTime)
		}

		if eventTitle == "" {
			utilities.WriteJSON(w, http.StatusBadRequest, "title is required", nil)
			return
		}

		parsedTime, err := time.Parse("2006-01-02 15:04:05", eventTimeStr)
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid event_time, use format: YYYY-MM-DD HH:MM:SS", nil)
			return
		}

		// Event time must be at least 2 hours from now
		if parsedTime.Before(time.Now().Add(2 * time.Hour)) {
			utilities.WriteJSON(w, http.StatusBadRequest, "event_time must be at least 2 hours from now", nil)
			return
		}

		event := &repository.GroupEvent{
			GroupID:     groupID,
			CreatorID:   userID,
			Title:       eventTitle,
			Description: eventDesc,
			Image:       eventImg,
			EventTime:   parsedTime,
		}
		if err := Repos.Group.CreateEvent(event); err != nil {
			utilities.WriteJSON(w, http.StatusInternalServerError, "could not create event", nil)
			return
		}
		notifyGroupUsers(
			groupID,
			"group_event",
			"event",
			event.ID,
			userID,
			userID,
		)
		utilities.WriteJSON(w, http.StatusCreated, "event created", map[string]any{"event_id": event.ID})
	default:
		utilities.WriteJSON(w, http.StatusNotFound, "unknown endpoint", nil)
	}
}

// GetGroupPublic fetches public group details.
// @Summary Get public group details
// @Description Returns public information about a group including membership status.
// @Tags Groups
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} map[string]any "Group public details"
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 404 {object} map[string]string "Group not found"
// @Router /api/groups/public/{id} [get]
func GetGroupPublic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	groupID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || groupID <= 0 {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid group id", nil)
		return
	}

	group, err := Repos.Group.GetPublicGroupDetails(groupID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusNotFound, "group not found", nil)
		return
	}

	// Include membership status
	isMember := false
	if userID, ok := middlewares.GetUserID(r); ok {
		member, err := Repos.Group.IsGroupMember(groupID, userID)
		if err == nil {
			isMember = member
		}
	}

	utilities.WriteJSON(w, http.StatusOK, "group fetched", map[string]any{
		"group":     group,
		"is_member": isMember,
	})
}

// GetMyGroups fetches groups the current user is a member of.
// @Summary Get my groups
// @Description Returns all groups where the current user is a member.
// @Tags Groups
// @Produce json
// @Success 200 {object} map[string]any "My groups fetched"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/users/groups [get]
func GetMyGroups(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	groups, err := Repos.Group.GetUserGroups(userID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch groups", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "groups fetched", groups)
}

// GetInviteCandidates returns users the current user follows who are not yet members of the group.
// GetInviteCandidates fetches users who can be invited to a group.
// @Summary Get invite candidates
// @Description Returns users the current user follows who aren't already members of the group.
// @Tags Groups
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} map[string]any "Invite candidates fetched"
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/groups/invite-candidates/{id} [get]
func GetInviteCandidates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	userID, ok := middlewares.GetUserID(r)
	if !ok {
		utilities.WriteJSON(w, http.StatusUnauthorized, "not logged in", nil)
		return
	}

	groupID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || groupID <= 0 {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid group id", nil)
		return
	}

	// Get users the current user follows (accepted follows)
	following, err := Repos.Follow.GetFollowing(userID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch following", nil)
		return
	}

	// Also get users who follow the current user
	followers, err := Repos.Follow.GetFollowers(userID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch followers", nil)
		return
	}

	// Merge both lists, deduplicating by ID
	merged := make(map[int]repository.User)
	for _, u := range following {
		merged[u.ID] = u
	}
	for _, u := range followers {
		merged[u.ID] = u
	}

	// Get existing member IDs to exclude them
	memberIDs, err := Repos.Group.GetGroupMemberIDs(groupID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch members", nil)
		return
	}

	memberSet := make(map[int]struct{}, len(memberIDs))
	for _, id := range memberIDs {
		memberSet[id] = struct{}{}
	}

	// Also add the current user to excluded set (can't invite yourself)
	memberSet[userID] = struct{}{}

	// Exclude already invited users
	invitedIDs, err := Repos.Group.GetGroupInvitedUserIDs(groupID)
	if err == nil {
		for _, id := range invitedIDs {
			memberSet[id] = struct{}{}
		}
	}

	type candidate struct {
		ID        int    `json:"id"`
		Nickname  string `json:"nickname"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Avatar    string `json:"avatar"`
	}

	candidates := make([]candidate, 0, len(merged))
	for _, u := range merged {
		if _, excluded := memberSet[u.ID]; excluded {
			continue
		}
		candidates = append(candidates, candidate{
			ID:        u.ID,
			Nickname:  u.Nickname,
			Firstname: u.Firstname,
			Lastname:  u.Lastname,
			Avatar:    u.Avatar,
		})
	}

	utilities.WriteJSON(w, http.StatusOK, "candidates fetched", candidates)
}

// GetGroupMembers fetches group members.
// @Summary Get group members
// @Description Returns the list of members in a group.
// @Tags Groups
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} map[string]any "Group members fetched"
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/groups/members/{id} [get]
func GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	groupID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || groupID <= 0 {
		utilities.WriteJSON(w, http.StatusBadRequest, "invalid group id", nil)
		return
	}

	memberIDs, err := Repos.Group.GetGroupMemberIDs(groupID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch members", nil)
		return
	}

	members, err := Repos.Group.GetUsersByIDs(memberIDs)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch member profiles", nil)
		return
	}

	result := make([]repository.FeedAuthor, 0, len(memberIDs))
	for _, id := range memberIDs {
		if m, ok := members[id]; ok {
			result = append(result, m)
		}
	}

	utilities.WriteJSON(w, http.StatusOK, "members fetched", result)
}

// GetGroupContent fetches the group content feed (posts and events).
// @Summary Get group content feed
// @Description Returns the content feed for a group (posts and events combined).
// @Tags Groups
// @Produce json
// @Param id path int true "Group ID"
// @Param limit query int false "Number of items per page" minimum(1) maximum(50)
// @Param last_id query int false "Last item ID for pagination"
// @Success 200 {object} map[string]any "Group content fetched"
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 401 {object} map[string]string "Not logged in"
// @Router /api/groups/content/{id} [get]
func GetGroupContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	userID, ok := middlewares.GetUserID(r)
	if !ok {

		utilities.WriteJSON(w, http.StatusNotFound, "not user found", nil)
		return
	}
	groupID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utilities.WriteJSON(w, http.StatusNotFound, "not found", nil)
		return
	}

	member, err := Repos.Group.IsGroupMember(groupID, userID)
	if err != nil || !member {
		utilities.WriteJSON(w, http.StatusForbidden, "not a group member", nil)
		return
	}

	limit := 20
	if v, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && v > 0 {
		limit = v
	}

	lastID := 0
	if v, err := strconv.Atoi(r.URL.Query().Get("last_id")); err == nil {
		lastID = v
	}

	feed, err := Repos.Group.GetGroupContent(groupID, userID, limit, lastID)
	if err != nil {
		utilities.WriteJSON(w, http.StatusInternalServerError, "could not fetch content", nil)
		return
	}

	utilities.WriteJSON(w, http.StatusOK, "content fetched", feed)
}
