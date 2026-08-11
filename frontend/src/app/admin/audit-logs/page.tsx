"use client";

import { useAuditLogs } from "@/lib/admin/admin.hook";
import { AuditLogEntry, DEFAULT_PAGE_SIZE } from "@/lib/admin/admin.types";
import { useState } from "react";
import { AuditLogsToolbar } from "@/components/admin/audit-logs/audit-logs-toolbar";
import { AuditLogsTable } from "@/components/admin/audit-logs/audit-logs-table";
import { useDebouncedValue } from "@/lib/use-debounced-value.hook";
import { useResettablePagination } from "@/lib/utils/pagination.hook";
import { useRequireAdmin } from "@/lib/admin/require.admin";
import { AdminTablePage } from "../admin-table-page";
import { AdminTablePagination } from "@/components/admin/admin-table-pagination";

export default function AuditLogsPage() {
  const [searchInput, setSearchInput] = useState<string>("");
  const debouncedSearch = useDebouncedValue(searchInput.trim(), 300);

  const { isAdmin, isProfilePending } = useRequireAdmin();

  const { offset, limit, setOffset, changeLimit } = useResettablePagination(
    debouncedSearch, 
    DEFAULT_PAGE_SIZE
  );

  const { data, isPending, isFetching, isPlaceholderData } = useAuditLogs(
    debouncedSearch,
    { limit, offset },
    {
      enabled: isAdmin,
    },
  );

  const auditLogs: AuditLogEntry[] = data?.logs ?? [];
  const total: number = data?.total ?? 0;

  return (
    <AdminTablePage 
      title="Audit Logs"
      description="View the audit logs of user actions."
      isLoading={isProfilePending || (isAdmin && isPending && !isPlaceholderData )}
      loadingLabel="Loading audit logs"
      toolbar={
        <AuditLogsToolbar 
          searchInput={searchInput}
          onSearchInputChange={setSearchInput}
          onResetSearch={() => setSearchInput("")}
          total={total}
        />
      }
      table={
        <AuditLogsTable 
          logs={auditLogs}
          isLoading={isPending && !isPlaceholderData}
          isFetching={isFetching}
        />
      }
      pagination={
        <AdminTablePagination 
          offset={offset}
          limit={limit}
          total={total}
          itemsCount={auditLogs.length}
          itemsName="Logs"
          onOffsetChange={setOffset}
          onLimitChange={changeLimit}
        />
      }
    />
  )
}
