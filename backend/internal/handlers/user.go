package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ubcesports/memberships/internal/auth"
	"github.com/ubcesports/memberships/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) DeleteSelf(w http.ResponseWriter, r *http.Request) {
	session := auth.SessionFromContext(r.Context())
	if session == nil || session.User == nil {
		writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	if err := h.userService.DeleteUser(r.Context(), session.User.ID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "DELETE_USER_FAILED", "Unable to delete account")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code":    code,
		"message": message,
	})
}
