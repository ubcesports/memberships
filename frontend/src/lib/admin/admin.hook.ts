import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  fetchUser,
  fetchUserMemberships,
  fetchUsers,
  updateUser,
  fetchAuditLogs,
} from "./admin.api";
import type {
  AdminUserFilters,
  AdminPagination,
  AppliedSearch,
  UpdateUserRequest,
  User,
} from "./admin.types";

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

export function useUser(userId: string, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["admin", "user", userId],
    queryFn: ({ signal }) => fetchUser(userId, signal),
    enabled: (options?.enabled ?? true) && Boolean(userId),
    // A 404 for an unknown user is a final answer, not a transient failure.
    retry: false,
  });
}

export function useUserMemberships(userId: string, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["admin", "user", userId, "memberships"],
    queryFn: ({ signal }) => fetchUserMemberships(userId, signal),
    enabled: (options?.enabled ?? true) && Boolean(userId),
    retry: false,
  });
}

export function useUpdateUser(userId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: UpdateUserRequest) => updateUser(userId, body),
    onSuccess: (user: User, body: UpdateUserRequest) => {
      // The response carries the full updated profile, so the detail view can
      // refresh without a second round trip.
      queryClient.setQueryData(["admin", "user", userId], user);
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });

      if (body.cancel_membership) {
        queryClient.invalidateQueries({ queryKey: ["admin", "user", userId, "memberships"] });
      }
    },
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
