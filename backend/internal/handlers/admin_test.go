package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestSafeCSVCellEscapesFormulaPrefix(t *testing.T) {
	if got := safeCSVCell("=1+1"); got != "'=1+1" {
		t.Fatalf("expected formula prefix to be escaped, got %q", got)
	}
	if got := safeCSVCell("Sudipto Islam"); got != "Sudipto Islam" {
		t.Fatalf("expected ordinary text to be unchanged, got %q", got)
	}
}
