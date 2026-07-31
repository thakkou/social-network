package auth

import (
	"testing"

	"01social/tests/setup"
)

func TestRegister_Success(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.POST(server, "/api/register").
		WithMultipart(map[string]string{
			"firstname": "New",
			"lastname":  "User",
			"email":     "newuser@example.com",
			"password":  "password123",
			"birthDate": "2000-01-01",
			"nickname":  "newuser",
		}, nil).
		DoStatus(200)

	if err != nil {
		t.Fatalf("expected 200, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, err := setup.ParseResponseWrapper(body)
	if err != nil {
		t.Fatal(err)
	}
	setup.AssertResponseOK(t, wrapper)
	setup.AssertMessage(t, wrapper, "registration success")
}

func TestRegister_MissingFields(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.POST(server, "/api/register").
		WithMultipart(map[string]string{
			"firstname": "",
			"lastname":  "User",
			"email":     "test@example.com",
			"password":  "password123",
			"birthDate": "2000-01-01",
		}, nil).
		DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseError(t, wrapper, 400)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	// Try registering with Alice's email (already seeded)
	body, err := setup.POST(server, "/api/register").
		WithMultipart(map[string]string{
			"firstname": "Fake",
			"lastname":  "Alice",
			"email":     "alice@example.com",
			"password":  "password123",
			"birthDate": "1990-01-01",
		}, nil).
		DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}

	wrapper, _ := setup.ParseResponseWrapper(body)
	setup.AssertResponseError(t, wrapper, 400)
}

func TestRegister_InvalidEmail(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.POST(server, "/api/register").
		WithMultipart(map[string]string{
			"firstname": "Bad",
			"lastname":  "Email",
			"email":     "not-an-email",
			"password":  "password123",
			"birthDate": "2000-01-01",
		}, nil).
		DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	server, db, err := setup.NewTestServerWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	defer setup.CleanupTestServer(server, db)

	body, err := setup.POST(server, "/api/register").
		WithMultipart(map[string]string{
			"firstname": "Weak",
			"lastname":  "Password",
			"email":     "weak@example.com",
			"password":  "123",
			"birthDate": "2000-01-01",
		}, nil).
		DoStatus(400)

	if err != nil {
		t.Fatalf("expected 400, got error: %v\nBody: %s", err, string(body))
	}
}
