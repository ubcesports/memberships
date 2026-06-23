package service

import (
	"testing"
	"time"
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
