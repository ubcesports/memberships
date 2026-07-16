import { Download, Loader2, RotateCcw } from "lucide-react";
import { ActionButton } from "@/components/action-button";
import type { AdminFilters, GroupType, RoleType, SearchMode } from "@/lib/admin/admin.types";
import { GROUP_OPTIONS, ROLE_OPTIONS, SEARCH_MODE_OPTIONS } from "@/lib/admin/admin.types";

const FIELD_CLASS_NAME =
  "h-10 border border-brand-border bg-brand-surface px-3 text-sm text-brand-text";

type IsStudentFilter = "all" | "yes" | "no";

type UsersToolbarProps = {
  searchMode: SearchMode;
  searchInput: string;
  filters: AdminFilters;
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

function getIsStudentFilterValue(filters: AdminFilters): IsStudentFilter {
  if (filters.isStudent === true) {
    return "yes";
  }

  if (filters.isStudent === false) {
    return "no";
  }

  return "all";
}

function hasActiveFilters(filters: AdminFilters) {
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
    <div className="flex flex-col gap-5 border-b border-brand-border px-5 py-5 sm:px-6">
      <div className="flex flex-col gap-3 lg:flex-row lg:items-end">
        <label className="flex min-w-40 flex-col gap-1.5 text-sm text-brand-text-subtle">
          <span>Search by</span>
          <select
            value={searchMode}
            onChange={(event) => onSearchModeChange(event.target.value as SearchMode)}
            className={FIELD_CLASS_NAME}
            aria-label="Search method"
          >
            {SEARCH_MODE_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>

        <label className="flex min-w-0 flex-1 flex-col gap-1.5 text-sm text-brand-text-subtle">
          <span>Search</span>
          <input
            type="search"
            value={searchInput}
            onChange={(event) => onSearchInputChange(event.target.value)}
            placeholder={`Search by ${SEARCH_MODE_OPTIONS.find((option) => option.value === searchMode)?.label.toLowerCase()}`}
            className={`${FIELD_CLASS_NAME} w-full`}
            aria-label="Search users"
          />
        </label>

        <div className="flex flex-wrap gap-2">
          <ActionButton
            onClick={onResetSearch}
            disabled={searchInput.trim().length === 0}
            icon={<RotateCcw aria-hidden="true" className="size-4" />}
          >
            Reset search
          </ActionButton>
        </div>
      </div>

      <div className="flex flex-col gap-3 xl:flex-row xl:items-end xl:justify-between">
        <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
          <label className="flex flex-col gap-1.5 text-sm text-brand-text-subtle">
            <span>Group</span>
            <select
              value={filters.group ?? ""}
              onChange={(event) =>
                onGroupChange(event.target.value ? (event.target.value as GroupType) : undefined)
              }
              className={FIELD_CLASS_NAME}
              aria-label="Filter by group"
            >
              <option value="">All</option>
              {GROUP_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </label>

          <label className="flex flex-col gap-1.5 text-sm text-brand-text-subtle">
            <span>Role</span>
            <select
              value={filters.role ?? ""}
              onChange={(event) =>
                onRoleChange(event.target.value ? (event.target.value as RoleType) : undefined)
              }
              className={FIELD_CLASS_NAME}
              aria-label="Filter by role"
            >
              <option value="">All</option>
              {ROLE_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </label>

          <label className="flex flex-col gap-1.5 text-sm text-brand-text-subtle">
            <span>Is student</span>
            <select
              value={getIsStudentFilterValue(filters)}
              onChange={(event) => onIsStudentChange(event.target.value as IsStudentFilter)}
              className={FIELD_CLASS_NAME}
              aria-label="Filter by student status"
            >
              <option value="all">All</option>
              <option value="yes">Yes</option>
              <option value="no">No</option>
            </select>
          </label>
        </div>

        <div className="flex flex-wrap gap-2">
          <ActionButton
            onClick={onResetFilters}
            disabled={!hasActiveFilters(filters)}
            icon={<RotateCcw aria-hidden="true" className="size-4" />}
          >
            Reset filters
          </ActionButton>
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
    </div>
  );
}
