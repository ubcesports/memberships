package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/database/db"
)

type OnboardingRepository struct {
	store *db.Queries
}

func NewOnboardingRepository(store *db.Queries) *OnboardingRepository {
	return &OnboardingRepository{
		store: store,
	}
}

func (r *OnboardingRepository) CompleteOnboarding(ctx context.Context, userID string, isStudent bool, studentID string) error {
	uuid := pgtype.UUID{}
	uuid.Scan(userID)

	studentIDPtr := pgtype.Text{}
	if studentID != "" {
		studentIDPtr.Scan(studentID)
	}

	return r.store.CompleteUserOnboarding(ctx, db.CompleteUserOnboardingParams{
		ID:        uuid,
		IsStudent: isStudent,
		StudentID: studentIDPtr,
	})
}
