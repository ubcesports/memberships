package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseAdminUserFiltersForExportIgnoresPagination(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/admin/users/export?full_name=dip&role=member&group=competitive_team&limit=invalid",
		nil,
	)

	filters, err := parseAdminUserFilters(req, false)
	if err != nil {
		t.Fatalf("expected valid filters, got %v", err)
	}
	if filters.FullName != "dip" || filters.Role != "member" || filters.Group != "competitive_team" {
		t.Fatalf("unexpected filters: %#v", filters)
	}
	if filters.Limit != 0 || filters.Offset != 0 {
		t.Fatalf("expected export pagination to be disabled, got limit=%d offset=%d", filters.Limit, filters.Offset)
	}
}

func TestParseAdminUserFiltersRejectsInvalidGroup(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin/users?group=invalid", nil)

	if _, err := parseAdminUserFilters(req, true); err == nil {
		t.Fatal("expected invalid group error")
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
