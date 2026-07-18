package repository

import "database/sql"

type Repositories struct {
	Post         *PostRepository
	Category     *CategoryRepository
	Chat         *ChatRepository
	User         *UserRepository
	Group        *GroupRepository
	Follow       *FollowRepository
	Profile      *ProfileRepository
	Reaction     *ReactionRepository
	Notification *NotificationRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Post:         NewPostRepository(db),
		Category:     NewCategoryRepository(db),
		Chat:         NewChatRepository(db),
		User:         NewUserRepository(db),
		Group:        NewGroupRepository(db),
		Follow:       NewFollowRepository(db),
		Profile:      NewProfileRepository(db),
		Reaction:     NewReactionRepository(db),
		Notification: NewNotificationRepository(db),
	}
}
