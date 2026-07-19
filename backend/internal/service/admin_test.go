package service

import "testing"

func TestBuildAdminQueryParamsLeavesPaginationNullForExport(t *testing.T) {
	isStudent := false
	params := buildAdminQueryParams(AdminUserFilters{
		FullName:  "dip",
		IsStudent: &isStudent,
	})

	if !params.FullName.Valid || params.FullName.String != "dip" {
		t.Fatalf("expected full_name filter, got %#v", params.FullName)
	}
	if !params.IsStudent.Valid || params.IsStudent.Bool {
		t.Fatalf("expected is_student=false filter, got %#v", params.IsStudent)
	}
	if params.Limit.Valid || params.Offset.Valid {
		t.Fatalf("expected null pagination, got limit=%#v offset=%#v", params.Limit, params.Offset)
	}
}
