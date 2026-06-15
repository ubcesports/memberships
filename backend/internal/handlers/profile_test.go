package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfileHandlerRejectsUnauthorizedRequests(t *testing.T) {
	handler := &ProfileHandler{}
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rec := httptest.NewRecorder()

	handler.GetCurrentProfile(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}
