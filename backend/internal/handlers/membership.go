package handlers

import (
	"log"
	"net/http"

	"github.com/ubcesports/memberships/internal/auth"
	"github.com/ubcesports/memberships/internal/service"
	"github.com/ubcesports/memberships/internal/util"
)

type MembershipHandler struct {
	membershipService *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService}
}

func (h *MembershipHandler) GetPublicTiersWithPrices(w http.ResponseWriter, r *http.Request) {
	tiers, err := h.membershipService.GetPublicTiersAndPrices(r.Context())
	if err != nil {
		log.Printf("failed to get public tiers and prices: %v", err)
		util.WriteApiResponse(w, 500, "ErrorGetPublicTiersWithPrices", "Error retrieving public tiers and prices. Please try again.")
		return
	}
	util.WriteJson(w, 200, tiers)
}

func (h *MembershipHandler) GetEligibleTiersWithPrices(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *MembershipHandler) GetCurrentMembershipWithTransaction(w http.ResponseWriter, r *http.Request) {
	userId, ok := currentUserID(r)
	if !ok {
		util.WriteApiResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	membership, err := h.membershipService.GetCurrentMembershipWithTransaction(r.Context(), userId)
	if err != nil {
		log.Printf("failed to get current membership: %v", err)
		util.WriteApiResponse(w, 500, "ErrorGetCurrentMembershipWithTransaction", "Error retrieving current user's membership. Please try again.")
		return
	}
	util.WriteJson(w, 200, membership)
}

func (h *MembershipHandler) GetTierByTierId(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *MembershipHandler) GetAllMembershipsWithTransactions(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *MembershipHandler) CreateMembershipCheckoutSession(w http.ResponseWriter, r *http.Request) {
	return
}

func currentUserID(r *http.Request) (string, bool) {
	session := auth.SessionFromContext(r.Context())
	if session == nil || session.User == nil {
		return "", false
	}
	value, ok := session.User.ID.(string)
	if !ok || value == "" {
		return "", false
	}
	return value, true
}
