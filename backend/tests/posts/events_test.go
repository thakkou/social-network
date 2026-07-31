package posts

import (
	"testing"

	"01social/tests/setup"
)

func TestEventResponse_Going(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// Event 1 is a Go Meetup in Group 1 (Gophers United)
	// User 4 (David) is a member of Group 1
	sessionID, err := setup.CreateSession(db, 4)
	if err != nil {
		t.Fatal(err)
	}

	// POST /api/groups/1/events/1/respond
	req := setup.JSON(server, "POST", "/api/groups/1/events/1/respond", map[string]string{
		"status": "going",
	})
	req.WithAuth(sessionID)
	body, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	body.Body.Close()
}

// event_response_test.go placeholder
func TestEventsList(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Just verify the test server works
	resp, err := setup.GET(server, "/api/groups/1/events").WithAuth(sessionID).Do()
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}
