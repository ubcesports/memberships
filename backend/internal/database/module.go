package database

import (
	"context"
	"log"

	"github.com/ubcesports/memberships/internal/database/db"
	"go.uber.org/fx"
)

var Module = fx.Module("database",
	fx.Provide(provideDatabase),
)

func provideDatabase(lc fx.Lifecycle) (*db.Queries, error) {
	pool, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing database connection pool...")
			pool.Close()
			return nil
		},
	})

	log.Println("Database connection established successfully.")
	return db.New(pool), nil
}
