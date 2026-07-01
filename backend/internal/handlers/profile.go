package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ubcesports/memberships/internal/auth"
	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/service"
	"github.com/ubcesports/memberships/internal/util"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetCurrentProfile(w http.ResponseWriter, r *http.Request) {
	// Get current user id
	userId, ok := currentUserID(r)
	if !ok {
		util.WriteApiResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	profile, err := h.profileService.GetProfileByUserID(r.Context(), userId)
	if err != nil {
		http.Error(w, "Unable to load profile", http.StatusInternalServerError)
		return
	}

	util.WriteJson(w, http.StatusOK, map[string]dto.ProfileDTO{"user": *profile})
}

func (h *ProfileHandler) GetIsUserOnboarded(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *ProfileHandler) OnboardUser(w http.ResponseWriter, r *http.Request) {
	// Get current user id
	userId, ok := currentUserID(r)
	if !ok {
		log.Printf("unauthorized onboard attempt: missing current user id")
		util.WriteApiResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	log.Printf("onboard user request started: user_id=%s", userId)

	// Ensure request body is valid
	var onboardUserRequest dto.OnboardUserRequest
	if err := json.NewDecoder(r.Body).Decode(&onboardUserRequest); err != nil {
		log.Printf("invalid onboard request body: user_id=%s error=%v", userId, err)
		util.WriteApiResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body. Please try again.")
		return
	}

	// Onboard user
	if err := h.profileService.OnboardUser(r.Context(), userId, onboardUserRequest); err != nil {
		switch {
		case errors.Is(err, service.ErrValidation):
			log.Printf("onboard validation failed: user_id=%s error=%v", userId, err)
			util.WriteApiResponse(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())

		case errors.Is(err, service.ErrConflict):
			log.Printf("onboard conflict: user_id=%s error=%v", userId, err)
			util.WriteApiResponse(w, http.StatusConflict, "CONFLICT", err.Error())

		default:
			log.Printf("onboard failed: user_id=%s error=%v", userId, err)
			util.WriteApiResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Error onboarding user. Please try again.")
		}

		return
	}

	log.Printf("onboard user request completed: user_id=%s", userId)
	util.WriteApiResponse(w, http.StatusOK, "OK", "User onboarded successfully!")
}

func currentUserID(r *http.Request) (string, bool) {
	session := auth.SessionFromContext(r.Context())
	if session == nil || session.User == nil {
		return "", false
	}
	value, ok := session.User.ID.(string)
	if !ok || value == "" {
		return "", false
	}
	return value, true
}
