package auth

import (
	"context"
	"net/http"

	"github.com/thecodearcher/limen"
)

type contextKey string

const sessionKey contextKey = "session"

// Requires user to be signed in to access endpoints under this
func RequireAuth(auth *limen.Limen) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := auth.GetSession(r)
			if err != nil || session == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), sessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Requires user to have these roles to access endpoints under this
// eg. r.Use(auth.RequireRole("admin")) or r.Use(auth.RequireRole("admin", "member"))
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := SessionFromContext(r.Context())
			if session == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			role, _ := session.User.Raw()["role"].(string)
			if !allowed[role] {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Extract session from context
func SessionFromContext(ctx context.Context) *limen.ValidatedSession {
	session, _ := ctx.Value(sessionKey).(*limen.ValidatedSession)
	return session
}
