package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ubcesports/memberships/internal/auth"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/service"
	"github.com/ubcesports/memberships/internal/util"
)

type MembershipHandler struct {
	membershipService *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService}
}

/*
Returns all public membership tiers and prices.

-> Day pass with student/non-student prices
-> Regular pass with student/non-student prices
-> Premium pass with student/non-student prices

API URL: GET /membership/tiers

Args:

	None

Returns:

	[]dto.MembershipTierDTO (HTTP 200)

Raises:

	500: public membership tiers and prices could not be retrieved
*/
func (h *MembershipHandler) GetPublicTiersWithPrices(w http.ResponseWriter, r *http.Request) {
	tiers, err := h.membershipService.GetPublicTiersAndPrices(r.Context())
	if err != nil {
		log.Printf("failed to get public tiers and prices: %v", err)
		util.WriteApiResponse(w, 500, "ErrorGetPublicTiersWithPrices", fmt.Sprintf("Error retrieving public tiers and prices. Error message: %v", err))
		return
	}
	util.WriteJson(w, 200, tiers)
}

/*
Returns membership tiers and prices the current user is eligible to buy.

For exec/director/board members:

	Executive Pass

For competitive team players:

	Competitive Team Pass

For regular UBC students:

	Day pass with student price
	Regular pass with student price
	Premium pass with student price

For regular non-students:

	Day pass with non-student price
	Regular pass with non-student price
	Premium pass with non-student price

API URL: GET /membership/tiers/eligible

Args:

	auth.Session user id

Returns:

	[]dto.EligibleMembershipTierDTO (HTTP 200)

Raises:

	401: user is not authenticated
	500: eligible membership tiers and prices could not be retrieved
*/
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

/*
Returns the current user's active membership and transaction.

API URL: GET /membership/me/current

Args:

	auth.Session user id

Returns:

	dto.MembershipDTO (HTTP 200)

Raises:

	401: user is not authenticated
	500: current membership could not be retrieved
*/
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

/*
Returns all past and present memberships and transactions for the current user.

API URL: GET /membership/me/all

Args:

	auth.Session user id

Returns:

	[]dto.MembershipDTO (HTTP 200)

Raises:

	401: user is not authenticated
	500: memberships could not be retrieved
*/
func (h *MembershipHandler) GetAllMembershipsWithTransactions(w http.ResponseWriter, r *http.Request) {
	userId, ok := currentUserID(r)
	if !ok {
		util.WriteApiResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	memberships, err := h.membershipService.GetAllMembershipsWithTransactions(r.Context(), userId)
	if err != nil {
		log.Printf("failed to get all memberships: %v", err)
		util.WriteApiResponse(w, 500, "ErrorGetAllMembershipsWithTransactions", fmt.Sprintf("Error retrieving current user's memberships. Error message: %v", err))
		return
	}
	util.WriteJson(w, 200, memberships)
}

/*
Creates a signed Stripe Checkout Session for a membership purchase.

API URL: POST /membership/checkout

Args:

	auth.Session user id
	dto.CheckoutSessionRequest

Returns:

	dto.CheckoutSessionResponse (HTTP 200)

Raises:

	400: request body is unreadable or malformed
	401: user is not authenticated
	403: requested tier is not eligible for the user
	404: requested tier does not exist
	409: user already has a membership that blocks this purchase
	500: checkout session could not be created
*/
func (h *MembershipHandler) CreateMembershipCheckoutSession(w http.ResponseWriter, r *http.Request) {
	userId, ok := currentUserID(r)
	if !ok {
		util.WriteApiResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	var checkoutSessionRequest dto.CheckoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&checkoutSessionRequest); err != nil {
		util.WriteApiResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body. Please try again.")
		return
	}

	checkoutSessionResponse, err := h.membershipService.CreateCheckoutSession(r.Context(), userId, checkoutSessionRequest)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMembershipAlreadyExists):
			util.WriteApiResponse(w, http.StatusConflict, "CONFLICT", err.Error())
			return

		case errors.Is(err, service.ErrTierNotEligible):
			util.WriteApiResponse(w, http.StatusForbidden, "TIER_NOT_AVAILABLE", err.Error())
			return

		case errors.Is(err, service.ErrTierNotFound):
			util.WriteApiResponse(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return

		case errors.Is(err, service.ErrMembershipPurchaseClosed):
			util.WriteApiResponse(w, http.StatusForbidden, "MEMBERSHIP_PURCHASE_CLOSED", err.Error())
			return

		case errors.Is(err, service.ErrPendingCheckoutAlreadyPaid):
			util.WriteApiResponse(w, http.StatusConflict, "CHECKOUT_ALREADY_PAID", err.Error())
			return

		default:
			util.WriteApiResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}

		return
	}

	util.WriteJson(w, http.StatusOK, *checkoutSessionResponse)
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
