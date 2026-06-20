package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ubcesports/memberships/internal/dto"
	"github.com/ubcesports/memberships/internal/service"
)

type AdminUserHandler struct {
	adminUserService *service.AdminUserService
}

func NewAdminUserHandler(adminUserService *service.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{adminUserService: adminUserService}
}

func (h *AdminUserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit := int32(25)
	offset := int32(0)

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 32)
		if err != nil || parsed <= 0 {
			http.Error(w, "limit must be a positive integer", http.StatusBadRequest)
			return
		}
		limit = int32(parsed)
	}

	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 32)
		if err != nil || parsed < 0 {
			http.Error(w, "offset must be a non-negative integer", http.StatusBadRequest)
			return
		}
		offset = int32(parsed)
	}

	var isStudent *bool
	if value := query.Get("is_student"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			http.Error(w, "is_student must be true or false", http.StatusBadRequest)
			return
		}
		isStudent = &parsed
	}

	role := query.Get("role")
	if role != "" {
		switch dto.RoleType(role) {
		case dto.RoleMember, dto.RoleAdmin:
		default:
			http.Error(w, "invalid role", http.StatusBadRequest)
			return
		}
	}

	group := query.Get("group")
	if group != "" {
		switch dto.GroupType(group) {
		case dto.GroupMember,
			dto.GroupCompetitiveTeam,
			dto.GroupExecutive,
			dto.GroupDirector,
			dto.GroupBoard:
		default:
			http.Error(w, "invalid group", http.StatusBadRequest)
			return
		}
	}

	filters := service.AdminUserFilters{
		FullName:  query.Get("full_name"),
		StudentID: query.Get("student_id"),
		Email:     query.Get("email"),
		Role:      role,
		IsStudent: isStudent,
		Group:     group,
		Limit:     limit,
		Offset:    offset,
	}

	users, err := h.adminUserService.GetUsers(r.Context(), filters)
	if err != nil {
		http.Error(w, "unable to load users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"users": users,
	})
}
