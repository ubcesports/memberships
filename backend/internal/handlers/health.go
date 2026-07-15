package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/ubcesports/memberships/internal/service"
)

type HealthHandler struct {
	healthService *service.HealthService
}

func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

func (h *HealthHandler) IsDatabaseHealthy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	healthy, err := h.healthService.IsDatabaseHealthy(r.Context())

	response := map[string]interface{}{
		"status": "healthy",
		"database": map[string]string{
			"status": "connected",
		},
	}

	if err != nil {
		slog.ErrorContext(r.Context(), "database health check failed",
			"error", err,
			"request_id", middleware.GetReqID(r.Context()),
		)
	}

	if err != nil || !healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
		response["status"] = "unhealthy"
		response["database"].(map[string]string)["status"] = "disconnected"
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
