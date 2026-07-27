import { SurfacePanel } from "@/components/surface-panel";
import { AuditLogActor, AuditLogEntry, User } from "@/lib/admin/admin.types";
import { formatTime } from "@/lib/utils/formatting";
import { Loader2 } from "lucide-react";
import Image from "next/image";

type AuditLogsTableProps = {
    logs: AuditLogEntry[];
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

function UserAvatar({ actor }: { actor: AuditLogActor }) {
    if (!actor.actor_avatar_url) {
        return <EmptyValue />;
    }

    return (
        <Image 
            src={actor.actor_avatar_url}
            alt=""
            width={32}
            height={32}
            className="size-8 border border-brand-border object-cover"
            unoptimized
        />
    );
}

const TABLE_HEADERS = [
    "Actor",
    "Actor Avatar",
    "Occurred At",
    "Action",
    "Description",
    "Outcome",
    "Request ID",
    "Target User",
    "Target User Avatar"
] as const;

export function AuditLogsTable({ logs, isLoading, isFetching }: AuditLogsTableProps) {
    if (isLoading) {
        return (
            <SurfacePanel className="flex min-h-0 flex-1 flex-col">
                <div className="flex flex-1 items-center justify-center gap-3">
                    <Loader2 aria-hidden="true" className="size-5 animate-spin" />
                    <span>Loading audit logs</span>
                </div>
            </SurfacePanel>
        )
    }

    if (logs.length === 0) {
        return (
            <SurfacePanel className="flex min-h-0 flex-1 flex-col">
                <div className="flex flex-1 items-center justify-center px-6 py-12 text-brand-text-muted">
                    No logs found.
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
                        {logs.map((log) => (
                            <tr key={log.request_id} className="border-b border-brand-border/70 last:border-b-0">
                                <td className="whitespace-nowrap px-4 py-3 font-medium text-brand-text">
                                    {log.actor.actor_full_name}
                                </td>
                                <td className="whitespace-nowrap px-4 py-3">
                                    <UserAvatar actor={log.actor} />
                                </td>
                                <td className="whitespace-nowrap px-4 py-3 text-brand-text-muted">
                                    {formatOptionalTime(log.occured_at)}
                                </td>
                                <td className="whitespace-nowrap px-4 py-3 text-brand-text-muted">
                                    {log.action}
                                </td>
                                <td className="whitespace-nowrap px-4 py-3 text-brand-text-muted">
                                    {log.description}
                                </td>
                                <td className="whitespace-nowrap px-4 py-3 text-brand-text-muted">
                                    {log.outcome}
                                </td>
                                <td className="whitespace-nowrap px-4 py-3 font-mono text-xs text-brand-text-subtle">
                                    {log.request_id}
                                </td>
                                <td className="whitespace-nowrap px-4 py-3 text-brand-text-muted">
                                    {log.target_user?.actor_full_name ? log.target_user.actor_full_name : <EmptyValue />}
                                </td>
                                <td className="whitespace-nowrap px-4 py-3">
                                    {log.target_user ? <UserAvatar actor={log.target_user} /> : <EmptyValue />}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </SurfacePanel>
    )
}