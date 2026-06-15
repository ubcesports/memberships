package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/thecodearcher/limen"
	"github.com/ubcesports/memberships/internal/utils"
)

type contextKey string

const sessionKey contextKey = "session"

// Requires user to be signed in to access endpoints under this
func RequireAuth(auth *limen.Limen) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := auth.GetSession(r)
			if err != nil || session == nil {
				utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
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
				utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
				return
			}
			role, _ := session.User.Raw()["role"].(string)
			if !allowed[role] {
				utils.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Requires user to have completed onboarding to access endpoints under this.
// Use after RequireAuth so the current session is already in the request context.
func RequireOnboarded(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := SessionFromContext(r.Context())
		if session == nil {
			utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
			return
		}
		if !isUserOnboarded(session.User) {
			utils.WriteError(w, http.StatusForbidden, "ONBOARDING_REQUIRED", "Onboarding required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isUserOnboarded(user *limen.User) bool {
	if user == nil {
		return false
	}
	completedAt := user.Raw()["onboarding_completed_at"]
	switch v := completedAt.(type) {
	case time.Time:
		return !v.IsZero()
	case *time.Time:
		return v != nil && !v.IsZero()
	default:
		return completedAt != nil
	}
}

// Extract session from context
func SessionFromContext(ctx context.Context) *limen.ValidatedSession {
	session, _ := ctx.Value(sessionKey).(*limen.ValidatedSession)
	return session
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	session := SessionFromContext(ctx)
	if session == nil || session.User == nil || session.User.ID == nil {
		return "", false
	}
	return fmt.Sprint(session.User.ID), true
}
