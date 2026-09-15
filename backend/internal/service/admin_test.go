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

func TestBuildAdminQueryParamsNormalizesExactSetFilters(t *testing.T) {
	params := buildAdminQueryParams(AdminUserFilters{
		Groups:            []string{"member", "executive", "member"},
		MembershipTierIDs: []string{"b", "a", "b"},
	})

	if len(params.Groups) != 2 || params.Groups[0] != "executive" || params.Groups[1] != "member" {
		t.Fatalf("expected sorted unique groups, got %#v", params.Groups)
	}
	if len(params.MembershipTierIds) != 2 || params.MembershipTierIds[0] != "a" || params.MembershipTierIds[1] != "b" {
		t.Fatalf("expected sorted unique membership tiers, got %#v", params.MembershipTierIds)
	}
}

func TestBuildAdminQueryParamsLeavesEmptySetsUnfiltered(t *testing.T) {
	params := buildAdminQueryParams(AdminUserFilters{})

	if params.Groups != nil || params.MembershipTierIds != nil {
		t.Fatalf("expected omitted sets to remain nil, got groups=%#v tiers=%#v", params.Groups, params.MembershipTierIds)
	}
}
