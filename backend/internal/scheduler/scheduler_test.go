package scheduler

import (
	"testing"
	"time"
)

func vancouver(t *testing.T) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	return loc
}

func TestDurationUntilNextRunLaterToday(t *testing.T) {
	loc := vancouver(t)
	now := time.Date(2026, time.April, 20, 6, 0, 0, 0, loc)

	got := durationUntilNextRun(now, 9)

	want := 3 * time.Hour
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestDurationUntilNextRunRollsOverToTomorrow(t *testing.T) {
	loc := vancouver(t)
	now := time.Date(2026, time.April, 20, 9, 0, 0, 0, loc)

	got := durationUntilNextRun(now, 9)

	want := 24 * time.Hour
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
