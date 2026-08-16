package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ubcesports/memberships/internal/database/db"
	"github.com/ubcesports/memberships/internal/dto"
)

func TestParseAdminFiltersForExportIgnoresPagination(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/admin/users/export?full_name=dip&role=member&group=competitive_team&is_student=false&limit=invalid",
		nil,
	)

	filters, err := parseAdminUserFilters(req, false)
	if err != nil {
		t.Fatalf("expected valid filters, got %v", err)
	}
	if filters.FullName != "dip" || filters.Role != "member" || filters.Group != "competitive_team" {
		t.Fatalf("unexpected filters: %#v", filters)
	}
	if filters.IsStudent == nil || *filters.IsStudent {
		t.Fatalf("expected is_student=false, got %#v", filters.IsStudent)
	}
	if filters.Limit != 0 || filters.Offset != 0 {
		t.Fatalf("expected export pagination to be disabled, got limit=%d offset=%d", filters.Limit, filters.Offset)
	}
}

func TestParseAdminFiltersRejectsInvalidGroup(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin/users?group=invalid", nil)

	if _, err := parseAdminUserFilters(req, true); err == nil {
		t.Fatal("expected invalid group error")
	}
}

func TestParseAdminAuditLogFilters(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin/audit-logs?actor_name=dip&limit=10&offset=20", nil)

	filters, err := parseAdminAuditLogFilters(req)
	if err != nil {
		t.Fatalf("expected valid filters, got %v", err)
	}
	if filters.ActorName != "dip" || filters.Limit != 10 || filters.Offset != 20 {
		t.Fatalf("unexpected filters: %#v", filters)
	}
}

func TestParseAdminAuditLogFiltersRejectsInvalidPagination(t *testing.T) {
	tests := []string{
		"/admin/audit-logs?limit=0",
		"/admin/audit-logs?offset=-1",
	}

	for _, target := range tests {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		if _, err := parseAdminAuditLogFilters(req); err == nil {
			t.Fatalf("expected invalid pagination error for %s", target)
		}
	}
}

func TestBuildUpdateUserRequestMapsBodyToServiceTypes(t *testing.T) {
	studentID := "12345678"
	isStudent := true
	role := dto.RoleAdmin

	request := buildUpdateUserRequest(dto.AdminUpdateUserRequest{
		StudentID:        &studentID,
		IsStudent:        &isStudent,
		GroupsAdd:        []dto.GroupType{dto.GroupBoard},
		GroupsRemove:     []dto.GroupType{dto.GroupMember, dto.GroupExecutive},
		Role:             &role,
		CancelMembership: true,
	})

	if request.StudentID == nil || *request.StudentID != "12345678" {
		t.Fatalf("expected student ID to be carried over, got %#v", request.StudentID)
	}
	if request.IsStudent == nil || !*request.IsStudent {
		t.Fatalf("expected is_student=true, got %#v", request.IsStudent)
	}
	if request.Role == nil || *request.Role != db.RoleTypeAdmin {
		t.Fatalf("expected admin role, got %#v", request.Role)
	}
	if len(request.GroupsAdd) != 1 || request.GroupsAdd[0] != db.GroupTypeBoard {
		t.Fatalf("unexpected groups to add: %v", request.GroupsAdd)
	}
	if len(request.GroupsRemove) != 2 || request.GroupsRemove[1] != db.GroupTypeExecutive {
		t.Fatalf("unexpected groups to remove: %v", request.GroupsRemove)
	}
	if !request.CancelMembership {
		t.Fatalf("unexpected membership flags: %#v", request)
	}
}

func TestBuildUpdateUserRequestLeavesAbsentRoleNil(t *testing.T) {
	request := buildUpdateUserRequest(dto.AdminUpdateUserRequest{})

	if request.Role != nil {
		t.Fatalf("expected an absent role to stay nil, got %#v", request.Role)
	}
}

func TestSafeCSVCellEscapesFormulaPrefix(t *testing.T) {
	if got := safeCSVCell("=1+1"); got != "'=1+1" {
		t.Fatalf("expected formula prefix to be escaped, got %q", got)
	}
	if got := safeCSVCell("Sudipto Islam"); got != "Sudipto Islam" {
		t.Fatalf("expected ordinary text to be unchanged, got %q", got)
	}
}
