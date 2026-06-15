package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ubcesports/memberships/internal/auth"
	"github.com/ubcesports/memberships/internal/service"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetCurrentProfile(w http.ResponseWriter, r *http.Request) {
	session := auth.SessionFromContext(r.Context())
	if session == nil || session.User == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.User.ID.(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var id pgtype.UUID
	if err := id.Scan(userID); err != nil {
		http.Error(w, "Invalid session user ID", http.StatusInternalServerError)
		return
	}

	profile, err := h.profileService.GetProfileByUserID(r.Context(), id)
	if err != nil {
		http.Error(w, "Unable to load profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user": profile,
	})
}
