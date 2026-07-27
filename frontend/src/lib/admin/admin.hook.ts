import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { fetchUsers, fetchAuditLogs } from "./admin.api";
import type { AdminUserFilters, AdminPagination, AppliedSearch } from "./admin.types";

export function useUsers(
  appliedSearch: AppliedSearch,
  filters: AdminUserFilters,
  pagination: AdminPagination,
  options?: { enabled?: boolean },
) {
  return useQuery({
    queryKey: ["admin", "users", { appliedSearch, filters, pagination }],
    queryFn: ({ signal }) => fetchUsers(appliedSearch, filters, pagination, signal),
    placeholderData: keepPreviousData,
    enabled: options?.enabled ?? true,
  });
}

export function useAuditLogs(
  actorName?: string,
  pagination?: AdminPagination,
  options?: { enabled?: boolean },
) {
  return useQuery({
    queryKey: ["admin", "audit-logs", { actorName, pagination }],
    queryFn: ({ signal }) => fetchAuditLogs(actorName, pagination, signal),
    placeholderData: keepPreviousData,
    enabled: options?.enabled ?? true,
  });
}
