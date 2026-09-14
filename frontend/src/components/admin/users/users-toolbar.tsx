import { Download, Loader2 } from "lucide-react";
import { ActionButton } from "@/components/action-button";
import type {
  AdminMembershipTierOption,
  AdminUserFilters,
  SearchMode,
} from "@/lib/types/admin.types";
import type { GroupType, RoleType } from "@/lib/types/user.types";
import {
  GROUP_OPTIONS,
  ROLE_OPTIONS,
  SEARCH_MODE_OPTIONS,
  IS_STUDENT_OPTIONS,
} from "@/lib/types/admin.types";
import { ToolbarContainer } from "@/components/toolbar/toolbar-container";
import { ToolbarRow } from "@/components/toolbar/toolbar-row";
import { SelectField } from "@/components/toolbar/select-option";
import { SearchField } from "@/components/toolbar/search-field";
import { ResetButton } from "@/components/toolbar/reset-button";
import { CheckboxFilter } from "@/components/toolbar/checkbox-filter";

import type { IsStudentFilter } from "@/lib/types/admin.types";

type UsersToolbarProps = {
  searchMode: SearchMode;
  searchInput: string;
  filters: AdminUserFilters;
  total: number;
  membershipTierOptions: AdminMembershipTierOption[];
  isExporting: boolean;
  onSearchModeChange: (mode: SearchMode) => void;
  onSearchInputChange: (value: string) => void;
  onResetSearch: () => void;
  onRoleChange: (role: RoleType | undefined) => void;
  onGroupChange: (groups: GroupType[]) => void;
  onMembershipTierChange: (tierIds: string[]) => void;
  onIsStudentChange: (value: IsStudentFilter) => void;
  onResetFilters: () => void;
  onExport: () => void;
};

function getIsStudentFilterValue(filters: AdminUserFilters): IsStudentFilter {
  if (filters.isStudent === true) {
    return "yes";
  }

  if (filters.isStudent === false) {
    return "no";
  }

  return "all";
}

function hasActiveFilters(filters: AdminUserFilters) {
  return (
    filters.role !== undefined ||
    filters.groups !== undefined ||
    filters.membershipTierIds !== undefined ||
    filters.isStudent !== undefined
  );
}

export function UsersToolbar({
  searchMode,
  searchInput,
  filters,
  total,
  membershipTierOptions,
  isExporting,
  onSearchModeChange,
  onSearchInputChange,
  onResetSearch,
  onRoleChange,
  onGroupChange,
  onMembershipTierChange,
  onIsStudentChange,
  onResetFilters,
  onExport,
}: UsersToolbarProps) {
  return (
    <ToolbarContainer>
      <ToolbarRow>
        <SelectField
          label="Search by"
          value={searchMode}
          onChange={(event) => onSearchModeChange(event as SearchMode)}
          options={SEARCH_MODE_OPTIONS}
          ariaLabel="Search method"
        />
        <SearchField
          label="Search"
          value={searchInput}
          onChange={(event) => onSearchInputChange(event)}
          placeholder={`Search by ${SEARCH_MODE_OPTIONS.find((option) => option.value === searchMode)?.label.toLowerCase()}`}
          ariaLabel="Search users"
          type="search"
          className="min-w-0 flex-1"
        />
        <div className="flex flex-wrap gap-2">
          <ResetButton
            label="Reset Search"
            onClick={onResetSearch}
            disabled={searchInput.trim().length === 0}
          />
        </div>
      </ToolbarRow>

      <div className="flex flex-col gap-4">
        <div className="flex min-w-0 flex-1 flex-col gap-4">
          <div className="grid gap-3 sm:grid-cols-2">
            <SelectField
              label="Role"
              value={filters.role ?? ""}
              onChange={(event) => onRoleChange(event ? (event as RoleType) : undefined)}
              options={ROLE_OPTIONS}
              allLabel="All"
              ariaLabel="Filter by role"
            />
            <SelectField
              label="Is student"
              value={getIsStudentFilterValue(filters)}
              onChange={(event) => onIsStudentChange(event as IsStudentFilter)}
              options={IS_STUDENT_OPTIONS}
              allLabel="All"
              allValue="all"
              ariaLabel="Filter by student status"
            />
          </div>

          <div className="grid min-w-0 gap-4 xl:grid-cols-2">
            <CheckboxFilter
              label="Groups"
              options={GROUP_OPTIONS}
              selectedValues={filters.groups ?? []}
              onChange={onGroupChange}
            />
            <CheckboxFilter
              label="Active membership tiers"
              options={membershipTierOptions.map((tier) => ({
                value: tier.id,
                label: tier.title,
                description: tier.program_name,
              }))}
              selectedValues={filters.membershipTierIds ?? []}
              onChange={onMembershipTierChange}
            />
          </div>
        </div>

        <div className="flex flex-wrap justify-end gap-2 border-t border-brand-border/70 pt-4">
          <ResetButton
            label="Reset Filters"
            onClick={onResetFilters}
            disabled={!hasActiveFilters(filters)}
          />
          <ActionButton
            onClick={onExport}
            disabled={total === 0 || isExporting}
            loading={isExporting}
            icon={<Download aria-hidden="true" className="size-4" />}
            loadingIcon={<Loader2 aria-hidden="true" className="size-4 animate-spin" />}
          >
            {isExporting ? "Exporting" : "Export CSV"}
          </ActionButton>
        </div>
      </div>
    </ToolbarContainer>
  );
}
