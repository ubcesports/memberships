"use client";

import { useMutation } from "@tanstack/react-query";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { UsersTable } from "@/components/admin/users/users-table";
import { UsersToolbar } from "@/components/admin/users/users-toolbar";
import { downloadCSVBlob, exportUsersCSV } from "@/lib/admin/admin.api";
import { useUsers } from "@/lib/admin/admin.hook";
import type {
  AdminUserFilters,
  AppliedSearch,
  GroupType,
  RoleType,
  SearchMode,
} from "@/lib/admin/admin.types";
import { DEFAULT_PAGE_SIZE } from "@/lib/admin/admin.types";
import { useDebouncedValue } from "@/lib/use-debounced-value.hook";
import { AdminTablePagination } from "@/components/admin/admin-table-pagination";
import { AdminTablePage } from "../admin-table-page";
import { useRequireAdmin } from "@/lib/admin/require.admin";
import { useResettablePagination } from "@/lib/utils/pagination.hook";

export default function UsersPage() {
  const { isAdmin, isProfilePending } = useRequireAdmin();

  const [searchMode, setSearchMode] = useState<SearchMode>("full_name");
  const [searchInput, setSearchInput] = useState("");
  const debouncedSearch = useDebouncedValue(searchInput.trim(), 300);
  const appliedSearch = useMemo<AppliedSearch>(
    () => (debouncedSearch ? { mode: searchMode, value: debouncedSearch } : null),
    [debouncedSearch, searchMode],
  );
  const [filters, setFilters] = useState<AdminUserFilters>({});

  const searchAnchor = `${searchMode}|${debouncedSearch}`;
  const { offset, limit, setOffset, changeLimit, resetOffset } = useResettablePagination(
    searchAnchor,
    DEFAULT_PAGE_SIZE,
  );

  const { data, isPending, isFetching, isPlaceholderData } = useUsers(
    appliedSearch,
    filters,
    { limit, offset },
    {
      enabled: isAdmin,
    },
  );

  const { mutate: exportUsers, isPending: isExporting } = useMutation({
    mutationFn: () => exportUsersCSV(appliedSearch, filters),
    onSuccess: (blob) => {
      downloadCSVBlob(blob);
      toast.success("Users exported");
    },
  });

  const handleRoleChange = (role: RoleType | undefined) => {
    setFilters((current) => ({ ...current, role }));
    resetOffset();
  };

  const handleGroupChange = (group: GroupType | undefined) => {
    setFilters((current) => ({ ...current, group }));
    resetOffset();
  };

  const handleIsStudentChange = (value: "all" | "yes" | "no") => {
    setFilters((current) => ({
      ...current,
      isStudent: value === "all" ? undefined : value === "yes",
    }));
    resetOffset();
  };

  const handleResetFilters = () => {
    setFilters({});
    resetOffset();
  };

  const users = data?.users ?? [];
  const total = data?.total ?? 0;

  return(
    <AdminTablePage 
      title="Users"
      description="Search, filter, and export member records."
      isLoading={isProfilePending || (isAdmin && isPending && !isPlaceholderData)}
      loadingLabel="Loading users"
      toolbar={
        <UsersToolbar
          searchMode={searchMode}
          searchInput={searchInput}
          filters={filters}
          total={total}
          isExporting={isExporting}
          onSearchModeChange={setSearchMode}
          onSearchInputChange={setSearchInput}
          onResetSearch={() => setSearchInput("")}
          onRoleChange={handleRoleChange}
          onGroupChange={handleGroupChange}
          onIsStudentChange={handleIsStudentChange}
          onResetFilters={handleResetFilters}
          onExport={() => exportUsers()}
        />
      }
      table={
        <UsersTable 
          users={users}
          isLoading={isPending && !isPlaceholderData}
          isFetching={isFetching}
        />
      }
      pagination={
        <AdminTablePagination
          offset={offset}
          limit={limit}
          total={total}
          itemsCount={users.length}
          itemsName="Users"
          onOffsetChange={setOffset}
          onLimitChange={changeLimit}
        />
      }
    />
  );
}
