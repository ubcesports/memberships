package service

import (
	"testing"
	"time"

	"github.com/stripe/stripe-go/v85"
	"github.com/ubcesports/memberships/internal/database/db"
)

func TestMembershipExpiry(t *testing.T) {
	location, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name     string
		purchase time.Time
		year     int
	}{
		{name: "February", purchase: time.Date(2027, time.February, 10, 12, 0, 0, 0, location), year: 2027},
		{name: "August 31", purchase: time.Date(2027, time.August, 31, 23, 59, 0, 0, location), year: 2027},
		{name: "September 1", purchase: time.Date(2027, time.September, 1, 0, 0, 0, 0, location), year: 2028},
		{name: "September", purchase: time.Date(2026, time.September, 15, 12, 0, 0, 0, location), year: 2027},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expiresAt, err := membershipExpiry(test.purchase)
			if err != nil {
				t.Fatal(err)
			}
			expected := time.Date(test.year, time.September, 1, 0, 0, 0, 0, location)
			if !expiresAt.Equal(expected) {
				t.Fatalf("expected %s, got %s", expected, expiresAt)
			}
		})
	}
}

func TestCalculateUpgradeAmount(t *testing.T) {
	tests := []struct {
		name   string
		target int64
		credit int64
		want   int64
	}{
		{name: "student upgrade", target: 2500, credit: 1500, want: 1000},
		{name: "community upgrade", target: 3000, credit: 2000, want: 1000},
		{name: "zero difference", target: 2000, credit: 2000, want: 0},
		{name: "negative difference", target: 1500, credit: 2000, want: -500},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := calculateUpgradeAmount(test.target, test.credit); got != test.want {
				t.Fatalf("expected %v, got %v", test.want, got)
			}
		})
	}
}

func TestIsNonChargeableUpgrade(t *testing.T) {
	tests := []struct {
		name   string
		kind   db.TransactionKindType
		amount int64
		want   bool
	}{
		{name: "free purchase", kind: db.TransactionKindTypePurchase, amount: 0, want: false},
		{name: "paid purchase", kind: db.TransactionKindTypePurchase, amount: 1000, want: false},
		{name: "free upgrade", kind: db.TransactionKindTypeUpgrade, amount: 0, want: true},
		{name: "negative upgrade", kind: db.TransactionKindTypeUpgrade, amount: -500, want: true},
		{name: "paid upgrade", kind: db.TransactionKindTypeUpgrade, amount: 1000, want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isNonChargeableUpgrade(test.kind, test.amount); got != test.want {
				t.Fatalf("expected %v, got %v", test.want, got)
			}
		})
	}
}

func TestCheckoutSessionReadyForFulfillment(t *testing.T) {
	tests := []struct {
		status stripe.CheckoutSessionPaymentStatus
		want   bool
	}{
		{status: stripe.CheckoutSessionPaymentStatusPaid, want: true},
		{status: stripe.CheckoutSessionPaymentStatusNoPaymentRequired, want: true},
		{status: stripe.CheckoutSessionPaymentStatusUnpaid, want: false},
	}
	for _, test := range tests {
		if got := checkoutSessionReadyForFulfillment(test.status); got != test.want {
			t.Errorf("status %q: expected %v, got %v", test.status, test.want, got)
		}
	}
}

func TestPaidCheckoutMissingPaymentIntent(t *testing.T) {
	tests := []struct {
		name    string
		session stripe.CheckoutSession
		want    bool
	}{
		{
			name:    "free paid checkout has no PaymentIntent",
			session: stripe.CheckoutSession{PaymentStatus: stripe.CheckoutSessionPaymentStatusPaid, AmountTotal: 0},
			want:    false,
		},
		{
			name:    "chargeable paid checkout requires PaymentIntent",
			session: stripe.CheckoutSession{PaymentStatus: stripe.CheckoutSessionPaymentStatusPaid, AmountTotal: 1000},
			want:    true,
		},
		{
			name: "chargeable paid checkout has PaymentIntent",
			session: stripe.CheckoutSession{
				PaymentStatus: stripe.CheckoutSessionPaymentStatusPaid,
				AmountTotal:   1000,
				PaymentIntent: &stripe.PaymentIntent{ID: "pi_test"},
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := paidCheckoutMissingPaymentIntent(&test.session); got != test.want {
				t.Fatalf("expected %v, got %v", test.want, got)
			}
		})
	}
}
