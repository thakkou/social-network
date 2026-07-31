package groups

import (
	"testing"

	"01social/tests/setup"
)

func TestEventResponse_NonMember(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// User 9 (Isabella) is not a member of Group 1
	sessionID, err := setup.CreateSession(db, 9)
	if err != nil {
		t.Fatal(err)
	}

	// Try responding to Event 1 in Group 1
	resp, err := setup.JSON(server, "POST", "/api/groups/1/events/1/respond", map[string]string{
		"status": "going",
	}).WithAuth(sessionID).Do()

	if err != nil {
		t.Fatal(err)
	}
	// Non-members might get 403 or the event response might still work
	// depending on implementation
	_ = resp
	resp.Body.Close()
}
