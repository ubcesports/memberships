package handlers

import (
	"encoding/json"
	"net/http"

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

	healthy := h.healthService.IsDatabaseHealthy(r.Context())

	response := map[string]interface{}{
		"status": "healthy",
		"database": map[string]string{
			"status": "connected",
		},
	}

	if !healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
		response["status"] = "unhealthy"
		response["database"].(map[string]string)["status"] = "disconnected"
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
