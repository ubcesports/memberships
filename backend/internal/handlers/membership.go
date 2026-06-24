package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/auth"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/service"
)

type MembershipHandler struct {
	service *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{service: membershipService}
}

/*
Returns the public Regular and Premium catalog with the ordinary Member price.

API URL: GET /membership/tiers

Returns:

	tiers: public tiers with the Member price and membership expiry date (HTTP 200)

Raises:

	500: tier mappings or Stripe prices could not be loaded
*/
func (h *MembershipHandler) GetPublicTiers(w http.ResponseWriter, r *http.Request) {
	tiers, err := h.service.ListPublicTiers(r.Context())
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "MEMBERSHIP_TIERS_UNAVAILABLE", "Unable to load membership tiers")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tiers": tiers})
}

/*
Returns the membership products that the current user may purchase or upgrade to.

API URL: GET /membership/tiers/eligible

Returns:

	tiers: personalized tiers with amount due and upgrade credit (HTTP 200)

Raises:

	401: user is not authenticated
	403: user has not completed onboarding
	500: tier mappings, membership state, or Stripe prices could not be loaded
*/
func (h *MembershipHandler) GetEligibleTiers(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	tiers, err := h.service.ListEligibleTiers(r.Context(), userID)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "MEMBERSHIP_TIERS_UNAVAILABLE", "Unable to load membership tiers")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tiers": tiers})
}

/*
Returns the current user's active membership.

API URL: GET /membership/me

Returns:

	membership: active membership details, or null when no active membership exists (HTTP 200)

Raises:

	401: user is not authenticated
	403: user has not completed onboarding
	500: membership could not be loaded
*/
func (h *MembershipHandler) GetCurrentMembership(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	membership, err := h.service.GetActiveMembership(r.Context(), userID)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "MEMBERSHIP_UNAVAILABLE", "Unable to load membership")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"membership": membership})
}

/*
Expires any previous open Checkout Session and creates a new 60-minute purchase or Premium-upgrade Checkout Session.

API URL: POST /membership/checkout

Args (JSON body):

	tier_id: required UUID of an eligible membership tier

Returns:

	url: Stripe-hosted Checkout URL

Raises:

	400: request body or tier_id is invalid
	401: user is not authenticated
	403: user has not completed onboarding
	409: active membership rules block the purchase or a payment is already processing
	422: tier is unavailable or the upgrade window is closed
	500: previous Checkout could not be expired or a new Checkout could not be created
*/
func (h *MembershipHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}
	var request dto.CreateCheckoutSessionDTO
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil || request.TierID == "" {
		writeAPIError(w, http.StatusBadRequest, "INVALID_REQUEST", "A valid tier_id is required")
		return
	}

	checkout, err := h.service.CreateCheckoutSession(r.Context(), userID, request.TierID)
	switch {
	case errors.Is(err, service.ErrMembershipActive):
		writeAPIError(w, http.StatusConflict, "ACTIVE_MEMBERSHIP", "An active membership already exists")
	case errors.Is(err, service.ErrTierUnavailable):
		writeAPIError(w, http.StatusUnprocessableEntity, "TIER_UNAVAILABLE", "The selected membership tier is unavailable")
	case errors.Is(err, service.ErrUpgradeWindowClosed):
		writeAPIError(w, http.StatusUnprocessableEntity, "UPGRADE_WINDOW_CLOSED", "Memberships cannot be upgraded during their final hour")
	case errors.Is(err, service.ErrPaymentProcessing):
		writeAPIError(w, http.StatusConflict, "PAYMENT_PROCESSING", "A membership payment is already processing")
	case err != nil:
		writeAPIError(w, http.StatusInternalServerError, "CHECKOUT_UNAVAILABLE", "Unable to create checkout")
	default:
		status := http.StatusCreated
		writeJSON(w, status, checkout)
	}
}

func currentUserID(w http.ResponseWriter, r *http.Request) (pgtype.UUID, bool) {
	session := auth.SessionFromContext(r.Context())
	if session == nil || session.User == nil {
		writeAPIError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return pgtype.UUID{}, false
	}
	value, ok := session.User.ID.(string)
	if !ok || value == "" {
		writeAPIError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return pgtype.UUID{}, false
	}
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "INVALID_USER", "Invalid session user")
		return pgtype.UUID{}, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"code": code, "message": message})
}
