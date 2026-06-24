package service

import (
	"testing"
	"time"

	"github.com/stripe/stripe-go/v85"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
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

func TestCalculateUpgradeAmounts(t *testing.T) {
	tests := []struct {
		name       string
		target     int64
		available  int64
		wantCredit int64
		wantAmount int64
	}{
		{name: "student upgrade", target: 2500, available: 1500, wantCredit: 1500, wantAmount: 1000},
		{name: "community upgrade", target: 3000, available: 2000, wantCredit: 2000, wantAmount: 1000},
		{name: "exact credit", target: 2000, available: 2000, wantCredit: 2000, wantAmount: 0},
		{name: "credit capped at target", target: 1500, available: 2000, wantCredit: 1500, wantAmount: 0},
		{name: "free group price", target: 0, available: 1500, wantCredit: 0, wantAmount: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			credit, amount := calculateUpgradeAmounts(test.target, test.available)
			if credit != test.wantCredit || amount != test.wantAmount {
				t.Fatalf("expected credit %v and amount %v, got credit %v and amount %v", test.wantCredit, test.wantAmount, credit, amount)
			}
		})
	}
}

func TestUpgradeWindowOpen(t *testing.T) {
	now := time.Date(2027, time.August, 31, 20, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{name: "more than one hour remains", expiresAt: now.Add(time.Hour + time.Second), want: true},
		{name: "exactly one hour remains", expiresAt: now.Add(time.Hour), want: false},
		{name: "less than one hour remains", expiresAt: now.Add(30 * time.Minute), want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := upgradeWindowOpen(test.expiresAt, now); got != test.want {
				t.Fatalf("expected %v, got %v", test.want, got)
			}
		})
	}
}

func TestPreferSelectedTier(t *testing.T) {
	tests := []struct {
		name      string
		candidate selectedTier
		current   selectedTier
		want      bool
	}{
		{
			name:      "lower eligible price wins",
			candidate: selectedTier{dto: dto.MembershipTierDTO{Price: dto.MembershipTierPriceDTO{AmountMinor: 0}}, group: db.GroupTypeCompetitiveTeam},
			current:   selectedTier{dto: dto.MembershipTierDTO{Price: dto.MembershipTierPriceDTO{AmountMinor: 2000}}, group: db.GroupTypeMember},
			want:      true,
		},
		{
			name:      "higher eligible price loses",
			candidate: selectedTier{dto: dto.MembershipTierDTO{Price: dto.MembershipTierPriceDTO{AmountMinor: 2000}}, group: db.GroupTypeMember},
			current:   selectedTier{dto: dto.MembershipTierDTO{Price: dto.MembershipTierPriceDTO{AmountMinor: 1000}}, group: db.GroupTypeExecutive},
			want:      false,
		},
		{
			name:      "enum order breaks ties",
			candidate: selectedTier{dto: dto.MembershipTierDTO{Price: dto.MembershipTierPriceDTO{AmountMinor: 1500}}, group: db.GroupTypeExecutive},
			current:   selectedTier{dto: dto.MembershipTierDTO{Price: dto.MembershipTierPriceDTO{AmountMinor: 1500}}, group: db.GroupTypeStudent},
			want:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := preferSelectedTier(test.candidate, test.current); got != test.want {
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
