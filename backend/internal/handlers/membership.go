package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/ubcesports/memberships/internal/auth"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/service"
	"github.com/ubcesports/memberships/internal/utils"
)

type MembershipHandler struct {
	membershipService *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService}
}

// ListTiers returns the active membership tiers with prices resolved for the
// current authenticated user.
//
// API URL: GET /membership/tiers
//
// Args:
//   - w: HTTP response writer
//   - r: HTTP request with authenticated, onboarded user context
//
// Returns:
//   - 200: dto.ListMembershipTiersResponse
//
// Raises:
//   - 401: user is not authenticated
//   - 404: tier price is missing for the user's group/student status
//   - 500: unexpected backend error
func (h *MembershipHandler) ListTiers(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}

	response, err := h.membershipService.ListAvailableTiers(r.Context(), userID)
	if err != nil {
		writeMembershipError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetCurrent returns the current uncancelled membership for the authenticated
// user, if one exists.
//
// API URL: GET /membership/current
//
// Args:
//   - w: HTTP response writer
//   - r: HTTP request with authenticated, onboarded user context
//
// Returns:
//   - 200: dto.CurrentMembershipResponse. membership is null when the user has
//     no current membership.
//
// Raises:
//   - 401: user is not authenticated
//   - 500: unexpected backend error
func (h *MembershipHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}

	response, err := h.membershipService.GetCurrentMembership(r.Context(), userID)
	if err != nil {
		writeMembershipError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// StartCheckout starts a membership purchase for the authenticated user.
//
// API URL: POST /membership/checkout
//
// Args:
//   - w: HTTP response writer
//   - r: HTTP request with JSON body dto.StartCheckoutRequest
//
// Returns:
//   - 200: dto.StartCheckoutResponse. For paid tiers, the response includes a
//     Stripe Checkout URL. For free tiers, the membership is completed
//     immediately.
//
// Raises:
//   - 400: invalid JSON or invalid tier code
//   - 401: user is not authenticated
//   - 404: tier price is missing for the user's group/student status
//   - 409: user already has an active membership, or Stripe price is not set
//   - 500: Stripe is not configured or an unexpected backend error occurred
func (h *MembershipHandler) StartCheckout(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}

	var request dto.StartCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request body")
		return
	}

	response, err := h.membershipService.StartCheckout(r.Context(), userID, request.TierCode)
	if err != nil {
		writeMembershipError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// HandleStripeWebhook receives Stripe webhook events and completes paid
// membership purchases after Stripe confirms payment.
//
// API URL: POST /stripe/webhook
//
// Args:
//   - w: HTTP response writer
//   - r: HTTP request containing the raw Stripe webhook body and
//     Stripe-Signature header
//
// Returns:
//   - 200: {"status":"ok"} for handled or ignored Stripe events
//
// Raises:
//   - 400: invalid body, invalid Stripe signature, or invalid Stripe event
//   - 500: Stripe webhook secret is not configured or an unexpected backend
//     error occurred
func (h *MembershipHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_BODY", "Invalid request body")
		return
	}

	if err := h.membershipService.HandleStripeWebhook(r.Context(), body, r.Header.Get("Stripe-Signature")); err != nil {
		writeMembershipError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

/*
Private functions
*/

func currentUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return "", false
	}
	return userID, true
}

func writeMembershipError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidUserID):
		utils.WriteError(w, http.StatusUnauthorized, "INVALID_USER", "Invalid authenticated user")
	case errors.Is(err, service.ErrInvalidTierCode):
		utils.WriteError(w, http.StatusBadRequest, "INVALID_TIER_CODE", "Invalid membership tier code")
	case errors.Is(err, service.ErrMembershipAlreadyExists):
		utils.WriteError(w, http.StatusConflict, "MEMBERSHIP_EXISTS", "User already has an active membership")
	case errors.Is(err, service.ErrTierPriceNotFound):
		utils.WriteError(w, http.StatusNotFound, "TIER_PRICE_NOT_FOUND", "Membership tier price was not found")
	case errors.Is(err, service.ErrStripePriceMissing):
		utils.WriteError(w, http.StatusConflict, "STRIPE_PRICE_MISSING", "Stripe price is not configured for this membership tier")
	case errors.Is(err, service.ErrStripeNotConfigured):
		utils.WriteError(w, http.StatusInternalServerError, "STRIPE_NOT_CONFIGURED", "Stripe is not configured")
	case errors.Is(err, service.ErrInvalidStripeSignature):
		utils.WriteError(w, http.StatusBadRequest, "INVALID_STRIPE_SIGNATURE", "Invalid Stripe signature")
	case errors.Is(err, service.ErrInvalidStripeEvent):
		utils.WriteError(w, http.StatusBadRequest, "INVALID_STRIPE_EVENT", "Invalid Stripe event")
	default:
		log.Printf("membership handler error: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
