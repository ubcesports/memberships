package handlers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	filters, err := parseAdminUserFilters(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

func (h *AdminUserHandler) ExportUsersCSV(w http.ResponseWriter, r *http.Request) {
	filters, err := parseAdminUserFilters(r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := h.adminUserService.ExportUsers(r.Context(), filters)
	if err != nil {
		http.Error(w, "unable to export users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="users.csv"`)

	writer := csv.NewWriter(w)
	if err := writer.Write([]string{
		"ID",
		"Full Name",
		"Email",
		"Student ID",
		"Role",
		"Is Student",
		"Groups",
		"Created At",
		"Updated At",
		"Email Verified At",
		"Onboarding Completed At",
		"Avatar URL",
	}); err != nil {
		return
	}

	for _, user := range users {
		groups := make([]string, 0, len(user.Groups))
		for _, group := range user.Groups {
			groups = append(groups, string(group))
		}

		if err := writer.Write([]string{
			user.ID,
			safeCSVCell(user.FullName),
			safeCSVCell(user.Email),
			safeCSVCell(optionalString(user.StudentID)),
			string(user.Role),
			strconv.FormatBool(user.IsStudent),
			strings.Join(groups, ";"),
			user.CreatedAt.Format(time.RFC3339),
			user.UpdatedAt.Format(time.RFC3339),
			optionalTime(user.EmailVerifiedAt),
			optionalTime(user.OnboardingCompletedAt),
			safeCSVCell(optionalString(user.AvatarURL)),
		}); err != nil {
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return
	}
}

func parseAdminUserFilters(r *http.Request, includePagination bool) (service.AdminUserFilters, error) {
	query := r.URL.Query()
	filters := service.AdminUserFilters{
		FullName:  query.Get("full_name"),
		StudentID: query.Get("student_id"),
		Email:     query.Get("email"),
		Role:      query.Get("role"),
		Group:     query.Get("group"),
		Limit:     25,
	}

	if value := query.Get("is_student"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return service.AdminUserFilters{}, errors.New("is_student must be true or false")
		}
		filters.IsStudent = &parsed
	}

	if filters.Role != "" {
		switch dto.RoleType(filters.Role) {
		case dto.RoleMember, dto.RoleAdmin:
		default:
			return service.AdminUserFilters{}, errors.New("invalid role")
		}
	}

	if filters.Group != "" {
		switch dto.GroupType(filters.Group) {
		case dto.GroupMember,
			dto.GroupCompetitiveTeam,
			dto.GroupExecutive,
			dto.GroupDirector,
			dto.GroupBoard:
		default:
			return service.AdminUserFilters{}, errors.New("invalid group")
		}
	}

	if !includePagination {
		filters.Limit = 0
		return filters, nil
	}

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 32)
		if err != nil || parsed <= 0 {
			return service.AdminUserFilters{}, errors.New("limit must be a positive integer")
		}
		filters.Limit = int32(parsed)
	}

	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 32)
		if err != nil || parsed < 0 {
			return service.AdminUserFilters{}, errors.New("offset must be a non-negative integer")
		}
		filters.Offset = int32(parsed)
	}

	return filters, nil
}

func optionalString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func optionalTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(time.RFC3339)
}

func safeCSVCell(value string) string {
	if value == "" || !strings.ContainsRune("=+-@", rune(value[0])) {
		return value
	}
	return "'" + value
}
