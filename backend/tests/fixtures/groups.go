package fixtures

// GroupInfo holds key details about a seeded group.
type GroupInfo struct {
	ID          int
	Title       string
	CreatorID   int
	Description string
	MemberIDs   []int
}

// Named group constants matching the seeder_test.sql data.
var (
	GophersUnited = GroupInfo{
		ID: 1, Title: "Gophers United", CreatorID: Alice.ID,
		Description: "Everything about Go programming, backend development, concurrency and open source.",
		MemberIDs:   []int{Alice.ID, Bob.ID, David.ID, Farid.ID},
	}
	SportsClub = GroupInfo{
		ID: 2, Title: "Sports Club", CreatorID: Farid.ID,
		Description: "Football, basketball, running, fitness and every kind of sport.",
		MemberIDs:   []int{Farid.ID, Alice.ID, Grace.ID, Hugo.ID},
	}
	CultureArts = GroupInfo{
		ID: 3, Title: "Culture & Arts", CreatorID: Hugo.ID,
		Description: "Books, music, cinema, painting and cultural events.",
		MemberIDs:   []int{Hugo.ID, Chloe.ID, Emma.ID, Bob.ID},
	}
	TravelExplorers = GroupInfo{
		ID: 4, Title: "Travel Explorers", CreatorID: Bob.ID,
		Description: "Share destinations, travel tips, hiking adventures and unforgettable experiences.",
		MemberIDs:   []int{Bob.ID, Alice.ID, Emma.ID, Grace.ID},
	}
)

// AllGroups returns all seeded groups.
func AllGroups() []GroupInfo {
	return []GroupInfo{GophersUnited, SportsClub, CultureArts, TravelExplorers}
}

// GroupByID returns the GroupInfo for the given ID, or nil if not found.
func GroupByID(id int) *GroupInfo {
	for _, g := range AllGroups() {
		if g.ID == id {
			return &g
		}
	}
	return nil
}

// IsGroupMember returns true if userID is a member of the group.
func (g GroupInfo) IsGroupMember(userID int) bool {
	for _, id := range g.MemberIDs {
		if id == userID {
			return true
		}
	}
	return false
}

// GroupInvites lists group_id -> invited_user_id from the seeder.
var GroupInvites = map[int][]int{
	GophersUnited.ID: {Grace.ID},
	SportsClub.ID:    {Emma.ID},
	CultureArts.ID:   {David.ID},
	TravelExplorers.ID: {Chloe.ID},
}

// GroupJoinRequests lists group_id -> requester_user_id from the seeder.
var GroupJoinRequests = map[int][]int{
	GophersUnited.ID: {Hugo.ID},
	SportsClub.ID:    {David.ID},
}
