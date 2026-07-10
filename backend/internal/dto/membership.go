package dto

import "time"

// Membership tiers

type MembershipTierDTO struct {
	ID          string                   `json:"id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Slug        string                   `json:"slug"`
	ProductId   string                   `json:"product_id"`
	Prices      []MembershipTierPriceDTO `json:"prices"`
}

type EligibleMembershipTierDTO struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Slug         string                 `json:"slug"`
	PurchaseType PurchaseType           `json:"purchase_type"`
	ProductId    string                 `json:"product_id"`
	Price        MembershipTierPriceDTO `json:"prices"`
}

type MembershipTierPriceDTO struct {
	Price             float64 `json:"price"`
	PriceId           string  `json:"price_id"`
	IsStudentRequired *bool   `json:"is_student_required"`
}

// Memberships

type MembershipDTO struct {
	ID          string         `json:"id"`
	TierId      string         `json:"tier_id"`
	StartedAt   time.Time      `json:"started_at"`
	ExpiresAt   time.Time      `json:"expires_at"`
	CancelledAt *time.Time     `json:"cancelled_at"`
	Transaction TransactionDTO `json:"transaction"`
}

type TransactionDTO struct {
	ID              string                `json:"id"`
	AmountPaid      string                `json:"amount_paid"`
	Status          TransactionStatusType `json:"status"`
	GroupAtPurchase GroupType             `json:"group_at_purchase"`
}

// Request

type CheckoutSessionRequest struct {
	TierId string `json:"tier_id"`
}

// Response

type CheckoutSessionResponse struct {
	Url string `json:"url"`
}
