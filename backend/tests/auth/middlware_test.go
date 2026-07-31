package auth

import (
	"testing"

	"01social/tests/setup"
)

func TestAuthMiddleware_BlocksUnauthenticated(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// These endpoints require auth and should return 401 without a session cookie
	protectedEndpoints := []string{
		"/api/me",
		"/api/profile/1",
		"/api/posts",
		"/api/groups",
		"/api/follow/follow/2",
		"/api/conversations",
		"/api/notifications",
		"/api/search?q=test",
		"/api/users/search?q=test",
		"/api/users/groups",
		"/api/online-users",
	}

	for _, endpoint := range protectedEndpoints {
		resp, err := setup.GET(server, endpoint).Do()
		if err != nil {
			t.Fatalf("request to %s failed: %v", endpoint, err)
		}
		if resp.StatusCode != 401 {
			t.Errorf("expected 401 for %s, got %d", endpoint, resp.StatusCode)
		}
		resp.Body.Close()
	}
}

func TestAuthMiddleware_AllowsAuthenticated(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	sessionID, err := setup.CreateSession(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	// These should work with auth
	endpoints := []string{
		"/api/me",
		"/api/profile/1",
		"/api/groups",
		"/api/users/groups",
	}

	for _, endpoint := range endpoints {
		resp, err := setup.GET(server, endpoint).WithAuth(sessionID).Do()
		if err != nil {
			t.Fatalf("request to %s failed: %v", endpoint, err)
		}
		if resp.StatusCode == 401 {
			t.Errorf("unexpected 401 for authenticated %s", endpoint)
		}
		resp.Body.Close()
	}
}
