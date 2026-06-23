package repository

import (
	"context"

	"github.com/ubcesports/memberships/internal/database/db"
)

type HealthRepository struct{
	store *db.Queries
}

func NewHealthRepository(store *db.Queries) *HealthRepository {
	return &HealthRepository{
		store: store,
	}
}

func (r *HealthRepository) IsDatabaseHealthy(context context.Context) (bool, error) {
	res, err := r.store.PingDatabase(context)
	if err != nil {
		return false, err
	}

	return res == 1, nil
}
