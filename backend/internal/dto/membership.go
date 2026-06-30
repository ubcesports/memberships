package dto

import "time"

type MembershipTierDTO struct {
	ID          string                   `json:"id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Slug        string                   `json:"slug"`
	Prices      []MembershipTierPriceDTO `json:"prices"`
}

type EligibleMembershipTierDTO struct {
	ID           string                   `json:"id"`
	Title        string                   `json:"title"`
	Description  string                   `json:"description"`
	Slug         string                   `json:"slug"`
	PurchaseType PurchaseType             `json:"purchase_type"`
	Prices       []MembershipTierPriceDTO `json:"prices"`
}

type MembershipTierPriceDTO struct {
	Price             string `json:"price"`
	IsStudentRequired *bool  `json:"is_student_required"`
}

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
