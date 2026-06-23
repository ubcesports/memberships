package dto

import "time"

type MembershipTierPriceDTO struct {
	AmountMinor int64     `json:"amount_minor"`
	Currency    string    `json:"currency"`
	Group       GroupType `json:"group"`
}

type PublicMembershipTierDTO struct {
	ID          string                   `json:"id"`
	Slug        string                   `json:"slug"`
	Title       string                   `json:"title"`
	Description *string                  `json:"description"`
	Prices      []MembershipTierPriceDTO `json:"prices"`
	ExpiresAt   time.Time                `json:"expires_at"`
}

type MembershipTierDTO struct {
	ID                string                 `json:"id"`
	Slug              string                 `json:"slug"`
	Title             string                 `json:"title"`
	Description       *string                `json:"description"`
	Price             MembershipTierPriceDTO `json:"price"`
	IsUpgrade         bool                   `json:"is_upgrade"`
	CreditAmountMinor int64                  `json:"credit_amount_minor"`
	AmountDueMinor    int64                  `json:"amount_due_minor"`
	ExpiresAt         time.Time              `json:"expires_at"`
}

type MembershipDTO struct {
	ID              string     `json:"id"`
	TierID          string     `json:"tier_id"`
	TierSlug        string     `json:"tier_slug"`
	TierTitle       string     `json:"tier_title"`
	TierDescription *string    `json:"tier_description"`
	GroupAtPurchase GroupType  `json:"group_at_purchase"`
	StartedAt       time.Time  `json:"started_at"`
	ExpiresAt       time.Time  `json:"expires_at"`
	CancelledAt     *time.Time `json:"cancelled_at"`
	Status          string     `json:"status"`
}

type CreateCheckoutSessionDTO struct {
	TierID string `json:"tier_id"`
}

type CheckoutSessionDTO struct {
	URL string `json:"url"`
}
