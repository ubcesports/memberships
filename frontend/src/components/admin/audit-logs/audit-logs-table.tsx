import { AuditLogEntry } from "@/lib/admin/admin.types";
import { Column, DataTable } from "../admin-data-table";
import { AvatarCell, EmptyValue, formatOptionalTime, LinkCell } from "../admin-table-cells";

type AuditLogsTableProps = {
  logs: AuditLogEntry[];
  isLoading: boolean;
  isFetching: boolean;
};

const columns: Column<AuditLogEntry>[] = [
  {
    header: "Actor",
    cellClassName: "whitespace-nowrap px-4 py-3 font-medium text-brand-text",
    cell: (log) => (
      <LinkCell href={`/admin/users/${log.actor.actor_user_id}`}>
        {log.actor.actor_full_name}
      </LinkCell>
    ),
  },
  {
    header: "Actor Avatar",
    cell: (log) => (
      <LinkCell href={`/admin/users/${log.actor.actor_avatar_url}`}>
        <AvatarCell src={log.actor.actor_avatar_url} />
      </LinkCell>
    ),
  },
  {
    header: "Occurred At",
    cell: (log) => formatOptionalTime(log.occured_at),
  },
  { header: "Action", cell: (log) => log.action },
  { header: "Description", cell: (log) => log.description },
  { header: "Outcome", cell: (log) => log.outcome },
  {
    header: "Request ID",
    cellClassName: "whitespace-nowrap px-4 py-3 font-mono text-xs text-brand-text-subtle",
    cell: (log) => log.request_id,
  },
  {
    header: "Target User",
    cell: (log) =>
      log.target_user?.actor_full_name ? (
        <LinkCell href={`/admin/users/${log.target_user.actor_user_id}`}>
          {log.target_user.actor_full_name}
        </LinkCell>
      ) : (
        <EmptyValue />
      ),
  },
  {
    header: "Target User Avatar",
    cell: (log) =>
      log.target_user?.actor_avatar_url ? (
        <LinkCell href={`/admin/users/${log.target_user.actor_user_id}`}>
          <AvatarCell src={log.target_user.actor_avatar_url} />
        </LinkCell>
      ) : (
        <EmptyValue />
      ),
  },
];

export function AuditLogsTable({ logs, isLoading, isFetching }: AuditLogsTableProps) {
  return (
    <DataTable
      data={logs}
      columns={columns}
      getRowKey={(log) => log.request_id}
      isLoading={isLoading}
      isFetching={isFetching}
      loadingLabel="Loading audit logs"
      emptyLabel="No logs found."
    />
  );
}
