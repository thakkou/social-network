package representation

import (
	repModal "01social/pkg/models/representation"
	"01social/pkg/repository"
)

func UserToProfileResponse(
	user *repository.User,
	followers []repository.User,
	following []repository.User,
	posts []repository.Post,
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
		CreatedAt:       user.CreatedAt,
		Followers:       mapFollowUsers(followers),
		Following:       mapFollowUsers(following),
		Posts:           mapPosts(posts, user.Nickname),
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

func mapPosts(posts []repository.Post, Nickname string) []repModal.PostResponse {
	result := make([]repModal.PostResponse, 0)

	for _, post := range posts {
		result = append(result, repModal.PostResponse{
			ID:           post.ID,
			UserID:       post.UserID,
			Nickname:     Nickname,
			CreatedAt:    post.CreatedAt,
			Title:        post.Title,
			Text:         post.Text,
			Image:        post.Image,
			Privacy:      post.Privacy,
			LikeCount:    post.LikeCount,
			DislikeCount: post.DislikeCount,
			IsLiked:      post.IsLiked,
			// Comments: mapComments(post.Comments),
		})
	}

	return result
}
