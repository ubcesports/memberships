import { DataTable, type Column } from "../admin-data-table";
import { AvatarCell, EmptyValue, formatOptionalTime } from "../admin-table-cells";
import { StatusBadge } from "@/components/status-badge";
import type { User } from "@/lib/admin/admin.types";
import { formatTime } from "@/lib/utils/formatting";
import { getGroupBadgeClass, titleCase } from "@/lib/utils/groups";

type UsersTableProps = {
  users: User[];
  isLoading: boolean;
  isFetching: boolean;
};

const columns: Column<User>[] = [
  {
    header: "Full name",
    cellClassName: "whitespace-nowrap px-4 py-3 font-medium text-brand-text",
    cell: (user) => user.full_name,
  },
  { header: "Email", cell: (user) => user.email },
  {
    header: "Student ID",
    cellClassName: "whitespace-nowrap px-4 py-3 font-mono text-brand-text-muted",
    cell: (user) => user.student_id ?? <EmptyValue />,
  },
  {
    header: "Role",
    cell: (user) => (
      <StatusBadge tone={user.role === "admin" ? "warning" : "default"}>
        {titleCase(user.role)}
      </StatusBadge>
    ),
  },
  {
    header: "Is student",
    cell: (user) => (
      <StatusBadge tone={user.is_student ? "success" : "muted"}>
        {user.is_student ? "Yes" : "No"}
      </StatusBadge>
    ),
  },
  {
    header: "Groups",
    cellClassName: "px-4 py-3",
    cell: (user) =>
      user.groups.length > 0 ? (
        <div className="flex flex-wrap gap-1.5">
          {user.groups.map((group) => (
            <StatusBadge key={group} tone="default" className={getGroupBadgeClass(group)}>
              {titleCase(group)}
            </StatusBadge>
          ))}
        </div>
      ) : (
        <EmptyValue />
      ),
  },
  { header: "Created at", cell: (user) => formatTime(user.created_at) },
  { header: "Updated at", cell: (user) => formatTime(user.updated_at) },
  {
    header: "Email verified at",
    cell: (user) => formatOptionalTime(user.email_verified_at),
  },
  {
    header: "Onboarding completed at",
    cell: (user) => formatOptionalTime(user.onboarding_completed_at),
  },
  { header: "Avatar", cell: (user) => <AvatarCell src={user.avatar_url} /> },
  {
    header: "ID",
    cellClassName: "whitespace-nowrap px-4 py-3 font-mono text-xs text-brand-text-subtle",
    cell: (user) => user.id,
  },
];

export function UsersTable({ users, isLoading, isFetching }: UsersTableProps) {
  return (
    <DataTable
      data={users}
      columns={columns}
      getRowKey={(user) => user.id}
      isLoading={isLoading}
      isFetching={isFetching}
      loadingLabel="Loading users"
      emptyLabel="No users match your search and filters."
    />
  );
}
