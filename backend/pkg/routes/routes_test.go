package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterRoutesMatchesDynamicProfilePath(t *testing.T) {
	oldMux := http.DefaultServeMux
	http.DefaultServeMux = http.NewServeMux()
	defer func() { http.DefaultServeMux = oldMux }()

	RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/profile/42", nil)
	rr := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized for protected profile route, got %d", rr.Code)
	}
}

func TestRegisterRoutesMatchesDynamicUsersPath(t *testing.T) {
	oldMux := http.DefaultServeMux
	http.DefaultServeMux = http.NewServeMux()
	defer func() { http.DefaultServeMux = oldMux }()

	RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/users/42", nil)
	rr := httptest.NewRecorder()

	http.DefaultServeMux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized for protected users route, got %d", rr.Code)
	}
}
