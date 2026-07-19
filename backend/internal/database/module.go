package database

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ubcesports/memberships/internal/database/db"
	"go.uber.org/fx"
)

var Module = fx.Module("database",
	fx.Provide(providePool),
	fx.Provide(provideDatabase),
	fx.Provide(provideStdlibDB),
)

func providePool(lc fx.Lifecycle) (*pgxpool.Pool, error) {
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
	return pool, nil
}

func provideDatabase(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
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
