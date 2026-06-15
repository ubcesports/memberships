package dto

import "time"

/*
Request
*/

type StartCheckoutRequest struct {
	TierCode TierCodeType `json:"tier_code"`
}

/*
Response
*/

type MembershipTierDTO struct {
	Code             TierCodeType      `json:"code"`
	Title            string            `json:"title"`
	Description      *string           `json:"description,omitempty"`
	Price            string            `json:"price"`
	Currency         string            `json:"currency"`
	Group            GroupType         `json:"group"`
	StudentStatus    StudentStatusType `json:"student_status"`
	RequiresCheckout bool              `json:"requires_checkout"`
}

type ListMembershipTiersResponse struct {
	Tiers []MembershipTierDTO `json:"tiers"`
}

type CurrentMembershipDTO struct {
	ID                      string                 `json:"id"`
	TierCode                TierCodeType           `json:"tier_code"`
	TierTitle               string                 `json:"tier_title"`
	TierDescription         *string                `json:"tier_description,omitempty"`
	GroupAtPurchase         GroupType              `json:"group_at_purchase"`
	StudentStatusAtPurchase StudentStatusType      `json:"student_status_at_purchase"`
	StartedAt               time.Time              `json:"started_at"`
	ExpiresAt               time.Time              `json:"expires_at"`
	CancelledAt             *time.Time             `json:"cancelled_at,omitempty"`
	Price                   *string                `json:"price,omitempty"`
	Currency                string                 `json:"currency"`
	TransactionStatus       *TransactionStatusType `json:"transaction_status,omitempty"`
	StripePaymentIntentID   *string                `json:"stripe_payment_intent_id,omitempty"`
}

type CurrentMembershipResponse struct {
	Membership *CurrentMembershipDTO `json:"membership"`
}

type StartCheckoutResponse struct {
	Status                  string  `json:"status"`
	CheckoutURL             *string `json:"checkout_url,omitempty"`
	StripeCheckoutSessionID *string `json:"stripe_checkout_session_id,omitempty"`
	MembershipID            *string `json:"membership_id,omitempty"`
}
