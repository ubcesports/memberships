import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { fetchAdminUsers } from "./admin-users.api";
import type { AdminUserFilters, AdminUserPagination, AppliedSearch } from "./admin-users.types";

export function useAdminUsers(
  appliedSearch: AppliedSearch,
  filters: AdminUserFilters,
  pagination: AdminUserPagination,
  options?: { enabled?: boolean },
) {
  return useQuery({
    queryKey: ["admin", "users", { appliedSearch, filters, pagination }],
    queryFn: ({ signal }) => fetchAdminUsers(appliedSearch, filters, pagination, signal),
    placeholderData: keepPreviousData,
    enabled: options?.enabled ?? true,
  });
}
