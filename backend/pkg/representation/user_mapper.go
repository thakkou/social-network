package representation

import (
	repModal "01social/pkg/models/representation"
	"01social/pkg/repository"
)

func UserToProfileResponse(
	user *repository.User,
	followers []repository.User,
	following []repository.User,
	followingStatus string,
) repModal.ProfileResponse {
	response := repModal.ProfileResponse{
		ID:              user.ID,
		Firstname:       user.Firstname,
		Lastname:        user.Lastname,
		Nickname:        user.Nickname,
		Avatar:          user.Avatar,
		AboutMe:         user.AboutMe,
		IsPrivate:       user.IsPrivate,
		FollowingStatus: followingStatus,
		Posts:           0,
		PostsList:       []repModal.Post{},
		Followers:       mapFollowUsers(followers),
		Following:       mapFollowUsers(following),
	}

	return response
}

func mapFollowUsers(users []repository.User) []repModal.UserFollow {
	result := make([]repModal.UserFollow, 0)

	for _, user := range users {
		result = append(result, repModal.UserFollow{
			ID:        user.ID,
			Firstname: user.Firstname,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
		})
	}

	return result
}
