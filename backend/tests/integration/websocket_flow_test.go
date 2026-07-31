package integration

import (
	"testing"

	"01social/tests/setup"
)

func TestWebSocketFlow_TicketCreation(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// Create a session for websocket ticket creation
	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	// WS ticket endpoint requires POST
	body, err := setup.POST(server, "/api/ws-ticket").
		WithAuth(sessionID).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseOK(t, wrapper)

	// Verify the ticket is returned in the response
	var data map[string]interface{}
	if err := wrapper.GetData(&data); err == nil {
		if ticket, ok := data["ticket"]; ok {
			if ticketStr, ok := ticket.(string); ok && ticketStr == "" {
				t.Error("expected non-empty ticket string")
			}
		}
	}
}

func TestWebSocketFlow_Unauthenticated(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	resp, err := setup.POST(server, "/api/ws-ticket").Do()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
