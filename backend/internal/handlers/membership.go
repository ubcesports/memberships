package handlers

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
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
	}
	util.WriteJson(w, 200, tiers)
}

func (h *MembershipHandler) GetEligibleTiersWithPrices(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *MembershipHandler) GetCurrentMembershipWithTransaction(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *MembershipHandler) GetAllMembershipsWithTransactions(w http.ResponseWriter, r *http.Request) {
	return
}

func (h *MembershipHandler) CreateMembershipCheckoutSession(w http.ResponseWriter, r *http.Request) {
	return
}

func currentUserID(w http.ResponseWriter, r *http.Request) (pgtype.UUID, bool) {
	session := auth.SessionFromContext(r.Context())
	if session == nil || session.User == nil {
		util.WriteApiResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return pgtype.UUID{}, false
	}
	value, ok := session.User.ID.(string)
	if !ok || value == "" {
		util.WriteApiResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return pgtype.UUID{}, false
	}
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		util.WriteApiResponse(w, http.StatusInternalServerError, "INVALID_USER", "Invalid session user")
		return pgtype.UUID{}, false
	}
	return id, true
}
