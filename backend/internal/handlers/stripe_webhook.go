package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/stripe/stripe-go/v86"
	"github.com/stripe/stripe-go/v86/webhook"
	"github.com/ubcesports/memberships/internal/service"
	"github.com/ubcesports/memberships/internal/util"
)

type StripeWebhookHandler struct {
	service       *service.MembershipService
	webhookSecret string
}

func NewStripeWebhookHandler(membershipService *service.MembershipService) (*StripeWebhookHandler, error) {
	secret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("STRIPE_WEBHOOK_SECRET is required")
	}
	return &StripeWebhookHandler{service: membershipService, webhookSecret: secret}, nil
}

/*
Processes signed Stripe Checkout, payment failure, and expiration events.

API URL: POST /webhooks/stripe

Args:

	Stripe-Signature header: required Stripe webhook signature
	body: raw Stripe Event JSON, limited to 64 KiB

Handled events:

	checkout.session.completed: fulfills a paid membership
	checkout.session.async_payment_succeeded: fulfills a delayed payment
	checkout.session.async_payment_failed: marks the transaction as failed
	checkout.session.expired: marks the transaction as expired

Returns:

	empty response acknowledging the event (HTTP 200)

Raises:

	400: body is unreadable, signature is invalid, or event data is malformed
	500: the verified event could not be processed
*/
func (h *StripeWebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<10))
	if err != nil {
		log.Printf("Stripe webhook request body failed: %v", err)
		util.WriteApiResponse(w, http.StatusBadRequest, "INVALID_WEBHOOK", "Unable to read webhook")
		return
	}
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), h.webhookSecret)
	if err != nil {
		log.Printf("Stripe webhook signature verification failed: %v", err)
		util.WriteApiResponse(w, http.StatusBadRequest, "INVALID_WEBHOOK_SIGNATURE", "Invalid webhook signature")
		return
	}
	occurredAt := time.Unix(event.Created, 0).UTC()

	// Handle every event type of the webhook properly
	switch event.Type {

	// The checkout is completed successfully and the payment has gone through
	case "checkout.session.completed", "checkout.session.async_payment_succeeded":
		session, ok := decodeCheckoutSession(w, event)
		if !ok {
			return
		}
		err = h.service.HandleCheckoutPaid(r.Context(), session, occurredAt)

	// The checkout session was expired
	case "checkout.session.expired":
		session, ok := decodeCheckoutSession(w, event)
		if !ok {
			return
		}
		err = h.service.HandleCheckoutExpired(r.Context(), session.ID)

	// The checkout session's payment failed to go through
	case "checkout.session.async_payment_failed":
		session, ok := decodeCheckoutSession(w, event)
		if !ok {
			return
		}
		err = h.service.HandleCheckoutFailed(r.Context(), session.ID)
	}
	if err != nil {
		log.Printf("Stripe webhook event %s (%s) processing failed: %v", event.ID, event.Type, err)
		util.WriteApiResponse(w, http.StatusInternalServerError, "WEBHOOK_PROCESSING_FAILED", "Webhook processing failed")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func decodeCheckoutSession(w http.ResponseWriter, event stripe.Event) (*stripe.CheckoutSession, bool) {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		log.Printf("Stripe webhook event %s (%s) decoding failed: %v", event.ID, event.Type, err)
		util.WriteApiResponse(w, http.StatusBadRequest, "INVALID_WEBHOOK", "Invalid Checkout Session")
		return nil, false
	}
	return &session, true
}
