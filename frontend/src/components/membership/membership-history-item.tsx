import { ChevronRight } from "lucide-react";
import { StatusBadge } from "@/components/status-badge";
import type { Membership } from "@/lib/membership.hook";
import { formatDate } from "@/lib/utils/formatting";
import { titleCase } from "@/lib/utils/groups";

type MembershipStatus = "active" | "expired" | "cancelled";

const MEMBERSHIP_STATUS_TONE = {
  active: "success",
  expired: "muted",
  cancelled: "muted",
} as const;

const TRANSACTION_STATUS_TONE = {
  completed: "success",
  pending: "warning",
  failed: "muted",
  refunded: "muted",
  expired: "muted",
} as const;

function getMembershipStatus(membership: Membership): MembershipStatus {
  if (membership.cancelled_at) return "cancelled";
  if (new Date(membership.expires_at) < new Date()) return "expired";
  return "active";
}

type MembershipHistoryItemProps = {
  membership: Membership;
};

export function MembershipHistoryItem({ membership }: MembershipHistoryItemProps) {
  const status = getMembershipStatus(membership);
  const endDateLabel = membership.cancelled_at ? "Cancelled on" : "Valid until";
  const endDate = membership.cancelled_at ?? membership.expires_at;

  return (
    <li>
      <details className="group/transaction">
        <summary className="flex cursor-pointer list-none items-center gap-3 px-5 py-4 transition hover:bg-white/[0.03] focus-visible:outline-2 focus-visible:outline-offset-[-2px] focus-visible:outline-brand-primary [&::-webkit-details-marker]:hidden">
          <ChevronRight
            aria-hidden="true"
            className="size-4 shrink-0 text-brand-text-subtle transition-transform group-open/transaction:rotate-90"
          />
          <span className="flex min-w-0 flex-1 flex-wrap items-center justify-between gap-x-8 gap-y-2">
            <span className="flex items-center gap-3">
              <span className="text-sm font-medium text-brand-text">{membership.tier_title}</span>
              <StatusBadge tone={MEMBERSHIP_STATUS_TONE[status]}>{titleCase(status)}</StatusBadge>
            </span>

            <span className="flex flex-wrap gap-x-6 gap-y-1 text-sm text-brand-text-muted">
              <span>
                Started <span className="text-brand-text">{formatDate(membership.started_at)}</span>
              </span>
              <span>
                {endDateLabel} <span className="text-brand-text">{formatDate(endDate)}</span>
              </span>
            </span>
          </span>
        </summary>

        <dl className="grid gap-x-8 gap-y-5 border-t border-brand-border/70 bg-white/[0.02] px-5 py-5 pl-12 sm:grid-cols-2">
          <TransactionDetail label="Payment status">
            <StatusBadge tone={TRANSACTION_STATUS_TONE[membership.transaction.status]}>
              {titleCase(membership.transaction.status)}
            </StatusBadge>
          </TransactionDetail>
          <TransactionDetail label="Amount paid">
            ${membership.transaction.amount_paid} CAD
          </TransactionDetail>
          <TransactionDetail label="Purchase group">
            {titleCase(membership.transaction.group_at_purchase)}
          </TransactionDetail>
          <TransactionDetail label="Transaction ID">
            <span className="break-all font-mono text-xs">{membership.transaction.id}</span>
          </TransactionDetail>
        </dl>
      </details>
    </li>
  );
}

function TransactionDetail({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div>
      <dt className="text-xs font-medium tracking-wide text-brand-text-subtle uppercase">
        {label}
      </dt>
      <dd className="mt-1.5 text-sm font-medium text-brand-text">{children}</dd>
    </div>
  );
}
