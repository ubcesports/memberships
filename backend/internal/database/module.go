package database

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ubcesports/memberships/internal/database/db"
	"go.uber.org/fx"
)

var Module = fx.Module("database",
	fx.Provide(provideDatabase),
	fx.Provide(provideStdlibDB),
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

func provideStdlibDB(lc fx.Lifecycle) (*sql.DB, error) {
	sqlDB, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing Limen database connection...")
			return sqlDB.Close()
		},
	})
	return sqlDB, nil
}
