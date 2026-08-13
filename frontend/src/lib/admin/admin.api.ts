import apiClient from "../client";
import type {
  AdminUserFilters,
  AdminPagination,
  AppliedSearch,
  UsersResponse,
  AuditLogEntry,
  AuditLogResponse,
} from "./admin.types";

export function buildAdminUserParams(
  appliedSearch: AppliedSearch,
  filters: AdminUserFilters,
  pagination?: AdminPagination,
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

export async function fetchUsers(
  appliedSearch: AppliedSearch,
  filters: AdminUserFilters,
  pagination: AdminPagination,
  signal?: AbortSignal,
): Promise<UsersResponse> {
  const response = await apiClient.get<UsersResponse>("/admin/users", {
    params: buildAdminUserParams(appliedSearch, filters, pagination),
    signal,
  });

  return response.data;
}

export async function exportUsersCSV(
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

export function buildAdminAuditLogParams(
  actorName?: string,
  pagination?: AdminPagination,
): Record<string, string | number> {
  const params: Record<string, string | number> = {};

  if (actorName?.trim()) {
    params.actor_name = actorName.trim();
  }

  if (pagination) {
    params.limit = pagination.limit;
    params.offset = pagination.offset;
  }

  return params;
}

export async function fetchAuditLogs(
  actorName?: string,
  pagination?: AdminPagination,
  signal?: AbortSignal,
): Promise<AuditLogResponse> {
  const response = await apiClient.get<AuditLogResponse>("/admin/audit-logs", {
    params: buildAdminAuditLogParams(actorName, pagination),
    signal,
  });
  return response.data;
}

export async function exportAuditLogsCSV(
  actorName?: string,
  signal?: AbortSignal,
): Promise<Blob> {
  const response = await apiClient.get<Blob>("/admin/audit-logs/export", {
    params: buildAdminAuditLogParams(actorName),
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
