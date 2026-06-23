package service

import "testing"

func TestBuildAdminUserQueryParamsLeavesPaginationNullForExport(t *testing.T) {
	params := buildAdminUserQueryParams(AdminUserFilters{
		FullName: "dip",
	})

	if !params.FullName.Valid || params.FullName.String != "dip" {
		t.Fatalf("expected full_name filter, got %#v", params.FullName)
	}
	if params.Limit.Valid || params.Offset.Valid {
		t.Fatalf("expected null pagination, got limit=%#v offset=%#v", params.Limit, params.Offset)
	}
}
