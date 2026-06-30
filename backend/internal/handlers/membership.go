package handlers

import (
	"fmt"
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
		util.WriteApiResponse(w, 500, "ErrorGetPublicTiersWithPrices", fmt.Sprintf("Error retrieving public tiers and prices. Error message: %v", err))
		return
	}
	util.WriteJson(w, 200, tiers)
}

func (h *MembershipHandler) GetEligibleTiersWithPrices(w http.ResponseWriter, r *http.Request) {
	userId, ok := currentUserID(r)
	if !ok {
		util.WriteApiResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	tiers, err := h.membershipService.GetEligibleTiersWithPrices(r.Context(), userId)
	if err != nil {
		log.Printf("failed to get eligible tiers and prices: %v", err)
		util.WriteApiResponse(w, 500, "ErrorGetEligibleTiersWithPrices", fmt.Sprintf("Error retrieving eligible tiers and prices. Error message: %v", err))
		return
	}

	if tiers != nil {
		util.WriteJson(w, 200, *tiers)
	} else {
		util.WriteJson(w, 200, nil)
	}
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
		util.WriteApiResponse(w, 500, "ErrorGetCurrentMembershipWithTransaction", fmt.Sprintf("Error retrieving current user's membership. Error message: %v", err))
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
