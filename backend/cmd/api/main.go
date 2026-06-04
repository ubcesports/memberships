package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/ubcesports/memberships/internal/database"
	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/handlers"
	"github.com/ubcesports/memberships/internal/mailer"
	"github.com/ubcesports/memberships/internal/repository"
	"github.com/ubcesports/memberships/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found!")
		return
	}

	mailer.Init()

	pool, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Database connection established successfully")

	store := db.New(pool)

	healthRepo := repository.NewHealthRepository(store)
	healthService := service.NewHealthService(healthRepo)
	healthHandler := handlers.NewHealthHandler(healthService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler.IsDatabaseHealthy)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
