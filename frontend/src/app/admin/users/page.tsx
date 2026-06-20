"use client";

import { Loader2 } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import type { AxiosError } from "axios";
import { UsersPagination } from "@/components/admin/users/users-pagination";
import { UsersTable } from "@/components/admin/users/users-table";
import { UsersToolbar } from "@/components/admin/users/users-toolbar";
import { BasePage } from "@/components/layout/base-page";
import { downloadCSVBlob, exportAdminUsersCSV } from "@/lib/admin-users.api";
import { useAdminUsers } from "@/lib/admin-users.hook";
import type {
  AdminUserFilters,
  AppliedSearch,
  GroupType,
  RoleType,
  SearchMode,
} from "@/lib/admin-users.types";
import { DEFAULT_PAGE_SIZE } from "@/lib/admin-users.types";
import { useProfile } from "@/lib/profile.hook";
import { useDebouncedValue } from "@/lib/use-debounced-value";

function getApiErrorMessage(error: unknown, fallback: string) {
  const axiosError = error as AxiosError<string>;
  const message = axiosError.response?.data;

  if (typeof message === "string" && message.trim()) {
    return message;
  }

  return fallback;
}

export default function AdminUsersPage() {
  const router = useRouter();
  const { data: profile, isPending: isProfilePending } = useProfile();

  const [searchMode, setSearchMode] = useState<SearchMode>("full_name");
  const [searchInput, setSearchInput] = useState("");
  const debouncedSearch = useDebouncedValue(searchInput.trim(), 300);
  const appliedSearch = useMemo<AppliedSearch>(
    () => (debouncedSearch ? { mode: searchMode, value: debouncedSearch } : null),
    [debouncedSearch, searchMode],
  );
  const [filters, setFilters] = useState<AdminUserFilters>({});
  const [limit, setLimit] = useState(DEFAULT_PAGE_SIZE);
  const [pagination, setPagination] = useState({ offset: 0, anchor: "" });

  const searchAnchor = `${searchMode}|${debouncedSearch}`;
  const offset = pagination.anchor === searchAnchor ? pagination.offset : 0;

  const isAdmin = profile?.role === "admin";

  const { data, error, isPending, isFetching, isPlaceholderData } = useAdminUsers(
    appliedSearch,
    filters,
    { limit, offset },
    {
      enabled: isAdmin,
    },
  );

  const { mutate: exportUsers, isPending: isExporting } = useMutation({
    mutationFn: () => exportAdminUsersCSV(appliedSearch, filters),
    onSuccess: (blob) => {
      downloadCSVBlob(blob);
      toast.success("Users exported");
    },
    onError: (exportError) => {
      toast.error(getApiErrorMessage(exportError, "Unable to export users"));
    },
  });

  useEffect(() => {
    if (!isProfilePending && profile && profile.role !== "admin") {
      router.replace("/403");
    }
  }, [isProfilePending, profile, router]);

  useEffect(() => {
    if (error) {
      toast.error(getApiErrorMessage(error, "Unable to load users"));
    }
  }, [error]);

  const handleSearchModeChange = (mode: SearchMode) => {
    setSearchMode(mode);
  };

  const handleResetSearch = () => {
    setSearchInput("");
  };

  const handleOffsetChange = (nextOffset: number) => {
    setPagination({ offset: nextOffset, anchor: searchAnchor });
  };

  const handleRoleChange = (role: RoleType | undefined) => {
    setFilters((current) => ({ ...current, role }));
    setPagination({ offset: 0, anchor: searchAnchor });
  };

  const handleGroupChange = (group: GroupType | undefined) => {
    setFilters((current) => ({ ...current, group }));
    setPagination({ offset: 0, anchor: searchAnchor });
  };

  const handleIsStudentChange = (value: "all" | "yes" | "no") => {
    setFilters((current) => ({
      ...current,
      isStudent: value === "all" ? undefined : value === "yes",
    }));
    setPagination({ offset: 0, anchor: searchAnchor });
  };

  const handleResetFilters = () => {
    setFilters({});
    setPagination({ offset: 0, anchor: searchAnchor });
  };

  const handleLimitChange = (nextLimit: number) => {
    setLimit(nextLimit);
    setPagination({ offset: 0, anchor: searchAnchor });
  };

  const users = data?.users ?? [];
  const total = data?.total ?? 0;
  const showInitialLoading = isProfilePending || (isAdmin && isPending && !isPlaceholderData);

  if (showInitialLoading) {
    return (
      <BasePage>
        <div className="flex flex-1 items-center py-6">
          <section className="mx-auto flex min-h-[85vh] w-full max-h-[calc(100vh-3rem)] items-center justify-center border border-brand-border bg-brand-surface/80 shadow-2xl shadow-black/25">
            <div className="flex items-center gap-3 text-brand-text-muted">
              <Loader2 aria-hidden="true" className="size-5 animate-spin" />
              <span>Loading users</span>
            </div>
          </section>
        </div>
      </BasePage>
    );
  }

  if (!isAdmin) {
    return null;
  }

  return (
    <BasePage>
      <div className="flex flex-1 items-center py-6">
        <section className="mx-auto flex min-h-[85vh] w-full max-h-[calc(100vh-3rem)] flex-col">
          <div className="flex min-h-[85vh] flex-1 flex-col border border-brand-border bg-brand-surface/80 shadow-2xl shadow-black/25">
            <div className="shrink-0 border-b border-brand-border px-5 py-5 sm:px-6">
              <h1 className="text-lg font-semibold text-brand-text">Users</h1>
              <p className="mt-1 text-sm text-brand-text-subtle">
                Search, filter, and export member records.
              </p>
            </div>

            <div className="shrink-0">
              <UsersToolbar
                searchMode={searchMode}
                searchInput={searchInput}
                filters={filters}
                total={total}
                isExporting={isExporting}
                onSearchModeChange={handleSearchModeChange}
                onSearchInputChange={setSearchInput}
                onResetSearch={handleResetSearch}
                onRoleChange={handleRoleChange}
                onGroupChange={handleGroupChange}
                onIsStudentChange={handleIsStudentChange}
                onResetFilters={handleResetFilters}
                onExport={() => exportUsers()}
              />
            </div>

            <div className="flex min-h-0 flex-1 flex-col p-5 sm:p-6">
              <UsersTable
                users={users}
                isLoading={isPending && !isPlaceholderData}
                isFetching={isFetching}
              />
            </div>

            <div className="shrink-0">
              <UsersPagination
                offset={offset}
                limit={limit}
                total={total}
                usersCount={users.length}
                onOffsetChange={handleOffsetChange}
                onLimitChange={handleLimitChange}
              />
            </div>
          </div>
        </section>
      </div>
    </BasePage>
  );
}
