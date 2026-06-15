package service

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/database/db"
)

func TestProfileServiceGetProfileByUserIDAggregatesGroups(t *testing.T) {
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	userID := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	studentID := pgtype.Text{String: "12345678", Valid: true}
	emailVerifiedAt := pgtype.Timestamptz{Time: now.Add(-24 * time.Hour), Valid: true}
	onboardingCompletedAt := pgtype.Timestamptz{Time: now.Add(-48 * time.Hour), Valid: true}
	avatarURL := pgtype.Text{String: "https://example.com/avatar.png", Valid: true}

	row := db.GetProfileByUserIDRow{
		ID:                    userID,
		Email:                 "sudi@example.com",
		StudentID:             studentID,
		Role:                  db.RoleTypeMember,
		CreatedAt:             pgtype.Timestamptz{Time: now.Add(-72 * time.Hour), Valid: true},
		UpdatedAt:             pgtype.Timestamptz{Time: now, Valid: true},
		FullName:              "Sudi Mango",
		EmailVerifiedAt:       emailVerifiedAt,
		IsStudent:             true,
		OnboardingCompletedAt: onboardingCompletedAt,
		AvatarUrl:             avatarURL,
		Groups:                []string{"member", "board"},
	}

	profile := buildProfile(row)

	if len(profile.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(profile.Groups))
	}
	if profile.Groups[0] != "member" || profile.Groups[1] != "board" {
		t.Fatalf("expected ordered groups [member board], got %v", profile.Groups)
	}
	if profile.StudentID == nil || *profile.StudentID != "12345678" {
		t.Fatalf("expected student ID to be preserved, got %#v", profile.StudentID)
	}
}

func TestProfileServiceGetProfileByUserIDReturnsEmptyGroupsSlice(t *testing.T) {
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	userID := pgtype.UUID{Bytes: [16]byte{2}, Valid: true}

	profile := buildProfile(db.GetProfileByUserIDRow{
		ID:        userID,
		Email:     "sudi@example.com",
		Role:      db.RoleTypeMember,
		CreatedAt: pgtype.Timestamptz{Time: now.Add(-72 * time.Hour), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		FullName:  "Sudi Mango",
		IsStudent: false,
		Groups:    []string{},
	})

	if profile.Groups == nil {
		t.Fatal("expected groups to be an empty slice, got nil")
	}
	if len(profile.Groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(profile.Groups))
	}
}
