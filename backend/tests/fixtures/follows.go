package fixtures

// FollowRel describes a follow relationship from the seeder data.
type FollowRel struct {
	FollowerID int
	FollowingID int
	Status     string // "accepted" | "pending"
}

// FollowRelationships lists all follow relationships from the seed.
// These are the seeds: follower_id, following_id, status.
var FollowRelationships = []FollowRel{
	{Alice.ID, Bob.ID, "accepted"},
	{Bob.ID, Alice.ID, "accepted"},
	{Alice.ID, David.ID, "accepted"},
	{David.ID, Farid.ID, "accepted"},
	{Farid.ID, Alice.ID, "accepted"},
	{Grace.ID, Chloe.ID, "pending"},
	{Hugo.ID, Emma.ID, "pending"},
	{Bob.ID, Emma.ID, "accepted"},
	{Chloe.ID, Alice.ID, "accepted"},
	{Alice.ID, Grace.ID, "accepted"},
	{Bob.ID, David.ID, "accepted"},
	{Hugo.ID, Alice.ID, "accepted"},
	{David.ID, Alice.ID, "accepted"},
	{Farid.ID, Hugo.ID, "accepted"},
}

// FollowersOf returns all user IDs that follow (accepted) the given user.
func FollowersOf(userID int) []int {
	var ids []int
	for _, rel := range FollowRelationships {
		if rel.FollowingID == userID && rel.Status == "accepted" {
			ids = append(ids, rel.FollowerID)
		}
	}
	return ids
}

// FollowingOf returns all user IDs that the given user follows (accepted).
func FollowingOf(userID int) []int {
	var ids []int
	for _, rel := range FollowRelationships {
		if rel.FollowerID == userID && rel.Status == "accepted" {
			ids = append(ids, rel.FollowingID)
		}
	}
	return ids
}

// PendingRequestsFor returns all pending follow request follower IDs for a user.
func PendingRequestsFor(userID int) []int {
	var ids []int
	for _, rel := range FollowRelationships {
		if rel.FollowingID == userID && rel.Status == "pending" {
			ids = append(ids, rel.FollowerID)
		}
	}
	return ids
}

// AreFollowingEachOther returns true if both users follow each other (mutual).
func AreFollowingEachOther(a, b int) bool {
	return IsFollowing(a, b) && IsFollowing(b, a)
}

// IsFollowing returns true if follower follows following with accepted status.
func IsFollowing(follower, following int) bool {
	for _, rel := range FollowRelationships {
		if rel.FollowerID == follower && rel.FollowingID == following && rel.Status == "accepted" {
			return true
		}
	}
	return false
}
