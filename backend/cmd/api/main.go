package main

import (
	"log"
	"log/slog"

	"github.com/joho/godotenv"
	"go.uber.org/fx"

	_ "time/tzdata"

	"github.com/ubcesports/memberships/internal/auth"
	"github.com/ubcesports/memberships/internal/database"
	"github.com/ubcesports/memberships/internal/handlers"
	"github.com/ubcesports/memberships/internal/logging"
	"github.com/ubcesports/memberships/internal/mailer"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/server"
	"github.com/ubcesports/memberships/internal/service"
	"github.com/ubcesports/memberships/internal/stripeclient"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	logging.Configure()
	slog.Info("starting backend...")

	fx.New(
		database.Module,
		auth.Module,
		repository.Module,
		service.Module,
		handlers.Module,
		server.Module,
		mailer.Module,
		stripeclient.Module,
	).Run()
}
