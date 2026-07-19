import { Loader2 } from "lucide-react";
import Image from "next/image";
import { StatusBadge } from "@/components/status-badge";
import { SurfacePanel } from "@/components/surface-panel";
import type { User } from "@/lib/admin/admin.types";
import { formatTime } from "@/lib/utils/formatting";
import { getGroupBadgeClass, titleCase } from "@/lib/utils/groups";

type UsersTableProps = {
  users: User[];
  isLoading: boolean;
  isFetching: boolean;
};

function EmptyValue() {
  return <span className="text-brand-text-muted">—</span>;
}

function formatOptionalTime(value: string | null) {
  if (!value) {
    return <EmptyValue />;
  }

  return formatTime(value);
}

function UserAvatar({ user }: { user: User }) {
  if (!user.avatar_url) {
    return <EmptyValue />;
  }

  return (
    <Image
      src={user.avatar_url}
      alt=""
      width={32}
      height={32}
      className="size-8 border border-brand-border object-cover"
      unoptimized
    />
  );
}

const TABLE_HEADERS = [
  "Full name",
  "Email",
  "Student ID",
  "Role",
  "Is student",
  "Groups",
  "Created at",
  "Updated at",
  "Email verified at",
  "Onboarding completed at",
  "Avatar",
  "ID",
] as const;

export function UsersTable({ users, isLoading, isFetching }: UsersTableProps) {
  if (isLoading) {
    return (
      <SurfacePanel className="flex min-h-0 flex-1 flex-col">
        <div className="flex flex-1 items-center justify-center gap-3 px-6 py-12 text-brand-text-muted">
          <Loader2 aria-hidden="true" className="size-5 animate-spin" />
          <span>Loading users</span>
        </div>
      </SurfacePanel>
    );
  }

  if (users.length === 0) {
    return (
      <SurfacePanel className="flex min-h-0 flex-1 flex-col">
        <div className="flex flex-1 items-center justify-center px-6 py-12 text-brand-text-muted">
          No users match your search and filters.
        </div>
      </SurfacePanel>
    );
  }

  return (
    <SurfacePanel className="flex min-h-0 flex-1 flex-col overflow-hidden">
      <div
        className={`min-h-0 flex-1 overflow-auto ${isFetching ? "opacity-70 transition-opacity" : ""}`}
      >
        <table className="min-w-full border-collapse text-left text-sm">
          <thead>
            <tr className="border-b border-brand-border bg-white/[0.02]">
              {TABLE_HEADERS.map((header) => (
                <th
                  key={header}
                  scope="col"
                  className="whitespace-nowrap px-4 py-3 text-xs font-semibold uppercase tracking-wide text-brand-text-subtle"
                >
                  {header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id} className="border-b border-brand-border/70 last:border-b-0">
                <td className="whitespace-nowrap px-4 py-3 font-medium text-brand-text">
                  {user.full_name}
                </td>
                <td className="whitespace-nowrap px-4 py-3 text-brand-text-muted">{user.email}</td>
                <td className="whitespace-nowrap px-4 py-3 font-mono text-brand-text-muted">
                  {user.student_id ?? <EmptyValue />}
                </td>
                <td className="whitespace-nowrap px-4 py-3">
                  <StatusBadge tone={user.role === "admin" ? "warning" : "default"}>
                    {titleCase(user.role)}
                  </StatusBadge>
                </td>
                <td className="whitespace-nowrap px-4 py-3">
                  <StatusBadge tone={user.is_student ? "success" : "muted"}>
                    {user.is_student ? "Yes" : "No"}
                  </StatusBadge>
                </td>
                <td className="px-4 py-3">
                  {user.groups.length > 0 ? (
                    <div className="flex flex-wrap gap-1.5">
                      {user.groups.map((group) => (
                        <StatusBadge
                          key={group}
                          tone="default"
                          className={getGroupBadgeClass(group)}
                        >
                          {titleCase(group)}
                        </StatusBadge>
                      ))}
                    </div>
                  ) : (
                    <EmptyValue />
                  )}
                </td>
                <td className="whitespace-nowrap px-4 py-3 text-brand-text-muted">
                  {formatTime(user.created_at)}
                </td>
                <td className="whitespace-nowrap px-4 py-3 text-brand-text-muted">
                  {formatTime(user.updated_at)}
                </td>
                <td className="whitespace-nowrap px-4 py-3 text-brand-text-muted">
                  {formatOptionalTime(user.email_verified_at)}
                </td>
                <td className="whitespace-nowrap px-4 py-3 text-brand-text-muted">
                  {formatOptionalTime(user.onboarding_completed_at)}
                </td>
                <td className="whitespace-nowrap px-4 py-3">
                  <UserAvatar user={user} />
                </td>
                <td className="whitespace-nowrap px-4 py-3 font-mono text-xs text-brand-text-subtle">
                  {user.id}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </SurfacePanel>
  );
}
