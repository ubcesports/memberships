"use client";

import { useAuditLogs } from "@/lib/admin/admin.hook";
import { AuditLogEntry, AuditLogResponse, DEFAULT_PAGE_SIZE } from "@/lib/admin/admin.types";
import { useRouter } from "next/dist/client/components/navigation";
import { useProfile } from "@/lib/profile.hook";
import { useEffect, useState } from "react";
import { BasePage } from "@/components/layout/base-page";
import { Loader2 } from "lucide-react";
import { AuditLogsToolbar } from "@/components/admin/audit-logs/audit-logs-toolbar";
import { AuditLogsTable } from "@/components/admin/audit-logs/audit-logs-table";
import { useDebouncedValue } from "@/lib/use-debounced-value.hook";
import { AuditLogsPagination } from "@/components/admin/audit-logs/audit-logs-pagination";

export default function AuditLogsPage() {
  const router = useRouter();
  const { data: profile, isPending: isProfilePending } = useProfile();

  const [searchInput, setSearchInput] = useState<string>("");
  const debouncedSearch = useDebouncedValue(searchInput.trim(), 300);

  const [limit, setLimit] = useState(DEFAULT_PAGE_SIZE);
  const [pagination, setPagination] = useState({ offset: 0, anchor: "" });

  const offset = pagination.anchor === debouncedSearch ? pagination.offset : 0;
  const isAdmin = profile?.role === "admin";

  const { data, isPending, isFetching, isPlaceholderData } = useAuditLogs(
    debouncedSearch,
    { limit, offset },
    {
      enabled: isAdmin,
    },
  );

  // TODO: leaving out exporting for now unless needed

  useEffect(() => {
    if (!isProfilePending && profile && profile.role !== "admin") {
      router.replace("/403");
    }
  }, [isProfilePending, profile, router]);

  const handleOffsetChange = (newOffset: number) => {
    setPagination({ offset: newOffset, anchor: debouncedSearch });
  };

  const handleLimitChange = (nextLimit: number) => {
    setLimit(nextLimit);
    setPagination({ offset: 0, anchor: debouncedSearch });
  };

  const auditLogs: AuditLogEntry[] = data?.logs ?? [];
  const total: number = data?.total ?? 0;
  const showInitialLoading = isProfilePending || (isAdmin && isPending && !isPlaceholderData);

  if (showInitialLoading) {
    return (
      <BasePage>
        <div className="flex flex-1 items-center py-6">
          <section className="mx-auto flex min-h-[85vh] w-full max-h-[calc(100vh-3rem)] items-center justify-center border border-brand-border bg-brand-surface/80 shadow-2xl shadow-black/25">
            <div className="flex items-center gap-3 text-brand-text-muted">
              <Loader2 aria-hidden="true" className="size-5 animate-spin" />
              <span>Loading audit logs</span>
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
              <h1 className="text-lg font-semibold text-brand-text">Audit Logs</h1>
              <p className="mt-1 text-sm text-brand-text-subtle">
                View the audit logs of user actions.
              </p>
            </div>

            <div className="shrink-0">
              <AuditLogsToolbar
                searchInput={searchInput}
                onSearchInputChange={setSearchInput}
                onResetSearch={() => setSearchInput("")}
                total={total}
              />
            </div>
            <div className="flex min-h-0 flex-1 flex-col p-5 sm:p-6">
              <AuditLogsTable
                logs={auditLogs}
                isLoading={isPending && !isPlaceholderData}
                isFetching={isFetching}
              />
            </div>
            <div className="shrink-0">
              {/* _TODO: implement pagination for audit logs specifically */}
              <AuditLogsPagination
                offset={offset}
                limit={limit}
                total={total}
                logsCount={auditLogs.length}
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
