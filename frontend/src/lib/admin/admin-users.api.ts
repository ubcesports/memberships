import apiClient from "../client";
import type {
  AdminUserFilters,
  AdminUserPagination,
  AdminUsersResponse,
  AppliedSearch,
} from "./admin-users.types";

export function buildAdminUserParams(
  appliedSearch: AppliedSearch,
  filters: AdminUserFilters,
  pagination?: AdminUserPagination,
): Record<string, string | number | boolean> {
  const params: Record<string, string | number | boolean> = {};

  if (appliedSearch?.value.trim()) {
    params[appliedSearch.mode] = appliedSearch.value.trim();
  }

  if (filters.role) {
    params.role = filters.role;
  }

  if (filters.group) {
    params.group = filters.group;
  }

  if (filters.isStudent !== undefined) {
    params.is_student = filters.isStudent;
  }

  if (pagination) {
    params.limit = pagination.limit;
    params.offset = pagination.offset;
  }

  return params;
}

export async function fetchAdminUsers(
  appliedSearch: AppliedSearch,
  filters: AdminUserFilters,
  pagination: AdminUserPagination,
  signal?: AbortSignal,
): Promise<AdminUsersResponse> {
  const response = await apiClient.get<AdminUsersResponse>("/admin/users", {
    params: buildAdminUserParams(appliedSearch, filters, pagination),
    signal,
  });

  return response.data;
}

export async function exportAdminUsersCSV(
  appliedSearch: AppliedSearch,
  filters: AdminUserFilters,
  signal?: AbortSignal,
): Promise<Blob> {
  const response = await apiClient.get<Blob>("/admin/users/export", {
    params: buildAdminUserParams(appliedSearch, filters),
    responseType: "blob",
    signal,
  });

  return response.data;
}

export function downloadCSVBlob(blob: Blob, filename = "users.csv") {
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = filename;
  anchor.click();
  URL.revokeObjectURL(url);
}
