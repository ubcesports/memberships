package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

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

/*
Get the currently authenticated user's profile.

API URL: GET /profile

Args:

	None

Returns:

	response body containing the current user's profile under the "user" key (HTTP 200).

Raises:

	401: unauthorized user
	500: unable to load profile
*/
func (h *ProfileHandler) GetCurrentProfile(w http.ResponseWriter, r *http.Request) {
	// Get current user id
	userId, ok := util.CurrentUserID(r)
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

/*
Check whether the currently authenticated user is onboarded.
Only returns HTTP 200, as middleware auto-handles non-onboarded users.

API URL: GET /onboard/check

Args:

	None

Returns:

	None (HTTP 200)

Raises:

	None
*/
func (h *ProfileHandler) GetIsUserOnboarded(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

/*
Onboard the currently authenticated user.

API URL: POST /onboard

Args:

	request body: dto.OnboardUserRequest

Returns:

	response body with success code and message (HTTP 200)

Raises:

	400: invalid request body or validation error
	401: unauthorized user
	409: conflict, when the user is already onboarded
	500: internal error, when onboarding fails unexpectedly
*/
func (h *ProfileHandler) OnboardUser(w http.ResponseWriter, r *http.Request) {
	// Get current user id
	userId, ok := util.CurrentUserID(r)
	if !ok {
		log.Printf("unauthorized onboard attempt: missing current user id")
		util.WriteApiResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

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

	util.WriteApiResponse(w, http.StatusOK, "OK", "User onboarded successfully!")
}
