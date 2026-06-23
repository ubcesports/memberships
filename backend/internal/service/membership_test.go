package service

import (
	"testing"
	"time"

	"github.com/stripe/stripe-go/v85"
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
		{name: "April 30", purchase: time.Date(2027, time.April, 30, 23, 59, 0, 0, location), year: 2027},
		{name: "May 1", purchase: time.Date(2027, time.May, 1, 0, 0, 0, 0, location), year: 2028},
		{name: "September", purchase: time.Date(2026, time.September, 15, 12, 0, 0, 0, location), year: 2027},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expiresAt, err := membershipExpiry(test.purchase)
			if err != nil {
				t.Fatal(err)
			}
			expected := time.Date(test.year, time.May, 1, 0, 0, 0, 0, location)
			if !expiresAt.Equal(expected) {
				t.Fatalf("expected %s, got %s", expected, expiresAt)
			}
		})
	}
}

func TestIsFullRefund(t *testing.T) {
	tests := []struct {
		name   string
		charge *stripe.Charge
		want   bool
	}{
		{name: "nil", charge: nil},
		{name: "partial", charge: &stripe.Charge{Amount: 1000, AmountRefunded: 500, Refunded: false}},
		{name: "amount without Stripe full flag", charge: &stripe.Charge{Amount: 1000, AmountRefunded: 1000, Refunded: false}},
		{name: "full", charge: &stripe.Charge{Amount: 1000, AmountRefunded: 1000, Refunded: true}, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isFullRefund(test.charge); got != test.want {
				t.Fatalf("expected %v, got %v", test.want, got)
			}
		})
	}
}
