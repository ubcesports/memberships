package main

import (
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/fx"

	"github.com/ubcesports/memberships/internal/database"
	"github.com/ubcesports/memberships/internal/handlers"
	"github.com/ubcesports/memberships/internal/mailer"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/server"
	"github.com/ubcesports/memberships/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found!")
		return
	}

	fx.New(
		database.Module,
		repository.Module,
		service.Module,
		handlers.Module,
		server.Module,
		mailer.Module,
	).Run()
}
