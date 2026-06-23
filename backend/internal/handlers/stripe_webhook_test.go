package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStripeWebhookRejectsInvalidSignature(t *testing.T) {
	handler := &StripeWebhookHandler{webhookSecret: "whsec_test"}
	request := httptest.NewRequest(http.MethodPost, "/webhooks/stripe", strings.NewReader(`{"type":"charge.refunded"}`))
	request.Header.Set("Stripe-Signature", "invalid")
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}
