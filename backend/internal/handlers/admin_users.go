package handlers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
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

/*
Returns a filtered and paginated list of users.

API URL: GET /admin/users

Args (query params):

	full_name: optional case-insensitive name substring
	student_id: optional case-insensitive student ID substring
	email: optional case-insensitive email substring
	role: optional role (member or admin)
	group: optional group membership
	limit: optional page size (default 25, maximum 100)
	offset: optional number of users to skip (default 0)

Returns:

	users: paginated user profiles (HTTP 200)
	total: number of users matching the filters

Raises:

	400: invalid filter or pagination value
	401: user is not authenticated
	403: user is not an admin
	500: users could not be loaded
*/
func (h *AdminUserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	filters, err := parseAdminUserFilters(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, total, err := h.adminUserService.GetUsers(r.Context(), filters)
	if err != nil {
		log.Printf("unable to load users: %v", err)
		http.Error(w, "unable to load users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"users": users,
		"total": total,
	})
}

/*
Exports every user matching the supplied filters as CSV.

API URL: GET /admin/users/export

Args (query params):

	full_name: optional case-insensitive name substring
	student_id: optional case-insensitive student ID substring
	email: optional case-insensitive email substring
	role: optional role (member or admin)
	group: optional group membership

Returns:

	users.csv: CSV file containing all matching users (HTTP 200)

Raises:

	400: invalid filter value
	401: user is not authenticated
	403: user is not an admin
	500: users could not be exported
*/
func (h *AdminUserHandler) ExportUsersCSV(w http.ResponseWriter, r *http.Request) {
	filters, err := parseAdminUserFilters(r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := h.adminUserService.ExportUsers(r.Context(), filters)
	if err != nil {
		log.Printf("unable to export users: %v", err)
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
			dto.GroupStudent,
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
