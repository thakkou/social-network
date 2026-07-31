package fixtures

// EventInfo describes a seeded group event.
type EventInfo struct {
	ID        int
	GroupID   int
	CreatorID int
	Title     string
}

// Named event constants.
var (
	GoMeetup       = EventInfo{ID: 1, GroupID: GophersUnited.ID, CreatorID: Alice.ID, Title: "Go Meetup"}
	FootballMatch  = EventInfo{ID: 2, GroupID: SportsClub.ID, CreatorID: Farid.ID, Title: "Football Match"}
	MuseumVisit    = EventInfo{ID: 3, GroupID: CultureArts.ID, CreatorID: Hugo.ID, Title: "Museum Visit"}
	WeekendRoadTrip = EventInfo{ID: 4, GroupID: TravelExplorers.ID, CreatorID: Bob.ID, Title: "Weekend Road Trip"}
)

// AllEvents returns all seeded events.
func AllEvents() []EventInfo {
	return []EventInfo{GoMeetup, FootballMatch, MuseumVisit, WeekendRoadTrip}
}

// EventsByGroup maps group IDs to their events.
func EventsByGroup(groupID int) []EventInfo {
	var events []EventInfo
	for _, e := range AllEvents() {
		if e.GroupID == groupID {
			events = append(events, e)
		}
	}
	return events
}

// EventByID returns the event with the given ID, or nil.
func EventByID(eventID int) *EventInfo {
	for _, e := range AllEvents() {
		if e.ID == eventID {
			return &e
		}
	}
	return nil
}

// EventResponse represents a single user's RSVP to an event.
type EventResponse struct {
	UserID int
	Status string // "going" | "not_going"
}

// EventResponses lists responses: event_id -> []EventResponse.
var EventResponses = map[int][]EventResponse{
	GoMeetup.ID: {
		{Alice.ID, "going"},
		{David.ID, "going"},
		{Farid.ID, "not_going"},
	},
	FootballMatch.ID: {
		{Farid.ID, "going"},
		{Alice.ID, "going"},
		{Hugo.ID, "not_going"},
	},
}
