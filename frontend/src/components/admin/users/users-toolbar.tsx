import { Download, Loader2 } from "lucide-react";
import { ActionButton } from "@/components/action-button";
import type { AdminUserFilters, GroupType, RoleType, SearchMode } from "@/lib/admin/admin.types";
import {
  GROUP_OPTIONS,
  ROLE_OPTIONS,
  SEARCH_MODE_OPTIONS,
  IS_STUDENT_OPTIONS,
} from "@/lib/admin/admin.types";
import { ToolbarContainer } from "@/components/toolbar/toolbar-container";
import { ToolbarRow } from "@/components/toolbar/toolbar-row";
import { SelectField } from "@/components/toolbar/select-option";
import { SearchField } from "@/components/toolbar/search-field";
import { ResetButton } from "@/components/toolbar/reset-button";

type IsStudentFilter = "all" | "yes" | "no";

type UsersToolbarProps = {
  searchMode: SearchMode;
  searchInput: string;
  filters: AdminUserFilters;
  total: number;
  isExporting: boolean;
  onSearchModeChange: (mode: SearchMode) => void;
  onSearchInputChange: (value: string) => void;
  onResetSearch: () => void;
  onRoleChange: (role: RoleType | undefined) => void;
  onGroupChange: (group: GroupType | undefined) => void;
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
    filters.role !== undefined || filters.group !== undefined || filters.isStudent !== undefined
  );
}

export function UsersToolbar({
  searchMode,
  searchInput,
  filters,
  total,
  isExporting,
  onSearchModeChange,
  onSearchInputChange,
  onResetSearch,
  onRoleChange,
  onGroupChange,
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

      <ToolbarRow justify>
        <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
          <SelectField
            label="Group"
            value={filters.group ?? ""}
            onChange={(event) => onGroupChange(event ? (event as GroupType) : undefined)}
            options={GROUP_OPTIONS}
            allLabel="All"
            ariaLabel="Filter by group"
          />
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

        <div className="flex flex-wrap gap-2">
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
      </ToolbarRow>
    </ToolbarContainer>
  );
}
