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

	// 2. Limit request body size to 2 MB (for logo + background uploads combined)
	const maxUploadSize int64 = 2 << 20 // 2 MB
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
		eventID, err := strconv.Atoi(segments[4])
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid event id", nil)
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
		switch action {
		case "accept":
			if err := Repos.Group.AcceptGroupRequest(groupID, targetUserID); err != nil {
				utilities.WriteJSON(w, http.StatusInternalServerError, "could not accept request", nil)
				return
			}
			utilities.WriteJSON(w, http.StatusOK, "join request accepted", nil)
		case "reject":
			if err := Repos.Group.RejectGroupRequest(groupID, targetUserID); err != nil {
				utilities.WriteJSON(w, http.StatusInternalServerError, "could not reject request", nil)
				return
			}
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
				utilities.WriteJSON(w, http.StatusInternalServerError, "could not accept invite", nil)
				return
			}
			utilities.WriteJSON(w, http.StatusOK, "group invite accepted", nil)
		case "reject":
			if err := Repos.Group.RejectGroupInvite(groupID, userID); err != nil {
				utilities.WriteJSON(w, http.StatusInternalServerError, "could not reject invite", nil)
				return
			}
			utilities.WriteJSON(w, http.StatusOK, "group invite rejected", nil)
		default:
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid action", nil)
		}
		return
	}

	switch endpoint {
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
					"text":       payload.Text,
				})
			}
		}

		utilities.WriteJSON(w, http.StatusCreated, "group message sent", map[string]any{"message_id": msg.ID})
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
		}
		utilities.WriteJSON(w, http.StatusOK, "join request sent", nil)
	case "invite":
		if r.Method != http.MethodPost {
			utilities.WriteJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
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
			CreatedAt: time.Now(),
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
		var payload struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			EventTime   string `json:"event_time"`
		}
		if err := utilities.ReadJSONRequestIntoStruct(r, &payload); err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		parsedTime, err := time.Parse("2006-01-02 15:04:05", payload.EventTime)
		if err != nil {
			utilities.WriteJSON(w, http.StatusBadRequest, "invalid event_time", nil)
			return
		}
		event := &repository.GroupEvent{GroupID: groupID, CreatorID: userID, Title: payload.Title, Description: payload.Description, EventTime: parsedTime}
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
		"group":    group,
		"is_member": isMember,
	})
}

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
	fmt.Println("get content groups", groupID)

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
