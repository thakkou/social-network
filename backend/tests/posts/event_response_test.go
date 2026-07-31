package posts

import (
	"testing"

	"01social/tests/setup"
)

func TestEventResponse_NotGoing(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// Event 1 in Group 1, User 4 responds not_going
	sessionID, err := setup.CreateSession(db, 4)
	if err != nil {
		t.Fatal(err)
	}

	req := setup.JSON(server, "POST", "/api/groups/1/events/1/respond", map[string]string{
		"status": "not_going",
	})
	req.WithAuth(sessionID)
	resp, err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}
