package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ubcesports/memberships/internal/auth"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/service"
)

type MembershipHandler struct {
	membershipService *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService}
}

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

	writeJSON(w, http.StatusOK, response)
}

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

	writeJSON(w, http.StatusOK, response)
}

func (h *MembershipHandler) StartCheckout(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(w, r)
	if !ok {
		return
	}

	var request dto.StartCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request body")
		return
	}

	response, err := h.membershipService.StartCheckout(r.Context(), userID, request.TierCode)
	if err != nil {
		writeMembershipError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *MembershipHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_BODY", "Invalid request body")
		return
	}

	if err := h.membershipService.HandleStripeWebhook(r.Context(), body, r.Header.Get("Stripe-Signature")); err != nil {
		writeMembershipError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func currentUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	session := auth.SessionFromContext(r.Context())
	if session == nil || session.User == nil || session.User.ID == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return "", false
	}
	return fmt.Sprint(session.User.ID), true
}

func writeMembershipError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidUserID):
		writeError(w, http.StatusUnauthorized, "INVALID_USER", "Invalid authenticated user")
	case errors.Is(err, service.ErrInvalidTierCode):
		writeError(w, http.StatusBadRequest, "INVALID_TIER_CODE", "Invalid membership tier code")
	case errors.Is(err, service.ErrMembershipAlreadyExists):
		writeError(w, http.StatusConflict, "MEMBERSHIP_EXISTS", "User already has an active membership")
	case errors.Is(err, service.ErrTierPriceNotFound):
		writeError(w, http.StatusNotFound, "TIER_PRICE_NOT_FOUND", "Membership tier price was not found")
	case errors.Is(err, service.ErrStripePriceMissing):
		writeError(w, http.StatusConflict, "STRIPE_PRICE_MISSING", "Stripe price is not configured for this membership tier")
	case errors.Is(err, service.ErrStripeNotConfigured):
		writeError(w, http.StatusInternalServerError, "STRIPE_NOT_CONFIGURED", "Stripe is not configured")
	case errors.Is(err, service.ErrInvalidStripeSignature):
		writeError(w, http.StatusBadRequest, "INVALID_STRIPE_SIGNATURE", "Invalid Stripe signature")
	case errors.Is(err, service.ErrInvalidStripeEvent):
		writeError(w, http.StatusBadRequest, "INVALID_STRIPE_EVENT", "Invalid Stripe event")
	default:
		log.Printf("membership handler error: %v", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{
		"code":    code,
		"message": message,
	})
}
