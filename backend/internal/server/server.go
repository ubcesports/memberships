package server

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thecodearcher/limen"
	"github.com/ubcesports/memberships/internal/auth"
	"github.com/ubcesports/memberships/internal/handlers"
	"go.uber.org/fx"
)

var Module = fx.Module("server",
	fx.Provide(provideRouter),
	fx.Invoke(startServer),
)

// Add all new routes here
func provideRouter(healthHandler *handlers.HealthHandler, profileHandler *handlers.ProfileHandler, limen *limen.Limen) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Mount("/", limen.Handler())

	// All public routes
	r.Get("/health", healthHandler.IsDatabaseHealthy)

	// All protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth(limen))
		r.Get("/profile", profileHandler.GetCurrentProfile)
	})

	// All onboarded routes
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth(limen))
		r.Use(auth.RequireOnboarded)
	})

	// All admin routes
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth(limen))
		r.Use(auth.RequireRole("admin"))
	})

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := make(map[string]struct{})
	for _, origin := range auth.TrustedFrontendOrigins() {
		allowedOrigins[origin] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, ok := allowedOrigins[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func startServer(lc fx.Lifecycle, r *chi.Mux) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatal("BASE_URL environment variable is required")
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Invalid BASE_URL: %v", err)
	}

	srv := &http.Server{Addr: parsed.Host, Handler: r}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Server running on %s", baseURL)
			go func() {
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					log.Fatalf("Server failed to start: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server...")
			return srv.Shutdown(ctx)
		},
	})
}
