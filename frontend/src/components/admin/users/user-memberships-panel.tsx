"use client";

import { Loader2, Ban } from "lucide-react";
import { Fragment, useMemo, useState } from "react";
import { toast } from "sonner";
import { ActionButton } from "@/components/action-button";
import { DetailRow } from "@/components/detail-row";
import { StatusBadge } from "@/components/status-badge";
import { SurfacePanel } from "@/components/surface-panel";
import type { Membership, UpdateUserRequest } from "@/lib/admin/admin.types";
import { useMembershipCatalog } from "@/lib/membership.hook";
import { formatDate, formatTime } from "@/lib/utils/formatting";
import { getGroupBadgeClass, titleCase } from "@/lib/utils/groups";

type UserMembershipsPanelProps = {
  memberships: Membership[];
  onSave: (body: UpdateUserRequest) => Promise<unknown>;
  isSaving: boolean;
};

type MembershipState = "active" | "cancelled" | "expired";

function getMembershipState(membership: Membership): MembershipState {
  if (membership.cancelled_at) {
    return "cancelled";
  }

  return new Date(membership.expires_at) > new Date() ? "active" : "expired";
}

const STATE_TONE = {
  active: "success",
  cancelled: "warning",
  expired: "muted",
} as const;

function TransactionSummary({ membership }: { membership: Membership }) {
  return (
    <span className="text-brand-text-muted">
      ${membership.transaction.amount_paid} · {titleCase(membership.transaction.status)}
    </span>
  );
}

function TransactionDetails({ membership }: { membership: Membership }) {
  const tx = membership.transaction;

  return (
    <div className="px-5 py-4 text-sm text-brand-text-muted">
      <div className="grid grid-cols-2 gap-4">
        <div>
          <div className="font-medium text-brand-text">Transaction ID</div>
          <div>{tx.id ?? "—"}</div>
        </div>
        <div>
          <div className="font-medium text-brand-text">Status</div>
          <div>{titleCase(tx.status ?? "unknown")}</div>
        </div>
        <div>
          <div className="font-medium text-brand-text">Amount</div>
          <div>
            {tx.currency ? `${tx.currency.toUpperCase()} ` : ""}
            {tx.amount_paid ?? "—"}
          </div>
        </div>
        <div>
          <div className="font-medium text-brand-text">Customer</div>
          <div>{tx.customer_id ?? "—"}</div>
        </div>
        <div>
          <div className="font-medium text-brand-text">Payment intent</div>
          <div>{tx.payment_intent ?? "—"}</div>
        </div>
        <div>
          <div className="font-medium text-brand-text">Charge</div>
          <div>{tx.charge_id ?? "—"}</div>
        </div>
        <div className="col-span-2">
          <div className="font-medium text-brand-text">Created</div>
          <div>{tx.created_at ? formatDate(tx.created_at) : "—"}</div>
        </div>
        {tx.metadata ? (
          <div className="col-span-2">
            <div className="font-medium text-brand-text">Metadata</div>
            <pre className="whitespace-pre-wrap break-words text-xs">
              {JSON.stringify(tx.metadata)}
            </pre>
          </div>
        ) : null}
      </div>
    </div>
  );
}

export function UserMembershipsPanel({ memberships, onSave, isSaving }: UserMembershipsPanelProps) {
  const [pendingAction, setPendingAction] = useState<"cancel" | null>(null);
  const { data: catalog } = useMembershipCatalog();

  const tierTitleById = useMemo(() => {
    const map = new Map<string, string>();
    for (const tier of catalog ?? []) {
      map.set(tier.id, tier.title);
    }
    return map;
  }, [catalog]);

  const tierTitle = (tierId: string) => tierTitleById.get(tierId) ?? "Unknown tier";

  const tierTitleForMembership = (membership: Membership) => {
    // Prefer server-provided title when available (joined on membership_tiers).
    if (membership.tier_title) return membership.tier_title;

    return tierTitle(membership.tier_id);
  };

  const current =
    memberships.find(
      (membership) => !membership.cancelled_at && new Date(membership.expires_at) > new Date(),
    ) ?? null;
  const past = memberships.filter((membership) => membership.id !== current?.id);

  const runAction = async (body: UpdateUserRequest, message: string) => {
    try {
      await onSave(body);
      toast.success(message);
    } catch {
      // The API client already surfaces the error message as a toast.
    } finally {
      setPendingAction(null);
    }
  };

  const [expanded, setExpanded] = useState<Record<string, boolean>>({});

  const toggleExpanded = (id: string) => {
    setExpanded((prev) => ({ ...prev, [id]: !prev[id] }));
  };

  return (
    <div className="flex flex-col gap-4">
      <SurfacePanel className="bg-transparent">
        <div className="flex flex-col gap-3 border-b border-brand-border px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 className="text-base font-semibold text-brand-text">Current membership</h2>
            <p className="mt-1 text-sm text-brand-text-subtle">
              The membership this user holds right now.
            </p>
          </div>

          <div className="flex shrink-0 gap-2">
            {current ? (
              pendingAction === "cancel" ? (
                <>
                  <ActionButton
                    onClick={() => runAction({ cancel_membership: true }, "Membership cancelled")}
                    loading={isSaving}
                    loadingIcon={<Loader2 aria-hidden="true" className="size-4 animate-spin" />}
                    className="border-amber-300/40 text-amber-100"
                  >
                    Confirm cancel
                  </ActionButton>
                  <ActionButton onClick={() => setPendingAction(null)} disabled={isSaving}>
                    Keep
                  </ActionButton>
                </>
              ) : (
                <ActionButton
                  onClick={() => setPendingAction("cancel")}
                  icon={<Ban aria-hidden="true" className="size-4" />}
                >
                  Cancel membership
                </ActionButton>
              )
            ) : null}
          </div>
        </div>

        {current ? (
          <>
            <dl>
              <DetailRow label="Tier">{tierTitleForMembership(current)}</DetailRow>
              <DetailRow label="Started">{formatTime(current.started_at)}</DetailRow>
              <DetailRow label="Expires">{formatTime(current.expires_at)}</DetailRow>
              <DetailRow label="Amount paid">${current.transaction.amount_paid}</DetailRow>
              <DetailRow label="Transaction">
                <div className="flex items-center justify-between gap-4">
                  <div className="flex flex-wrap items-center gap-2">
                    <StatusBadge tone="muted">{titleCase(current.transaction.status)}</StatusBadge>
                    {current.transaction.group_at_purchase ? (
                      <StatusBadge
                        className={getGroupBadgeClass(current.transaction.group_at_purchase)}
                      >
                        {titleCase(current.transaction.group_at_purchase)}
                      </StatusBadge>
                    ) : null}
                  </div>
                  <button
                    type="button"
                    onClick={() => toggleExpanded(current.id)}
                    className="text-sm text-brand-primary underline"
                  >
                    {expanded[current.id] ? "Hide details" : "Details"}
                  </button>
                </div>
              </DetailRow>
            </dl>
            {expanded[current.id] ? <TransactionDetails membership={current} /> : null}
          </>
        ) : (
          <p className="px-5 py-6 text-sm text-brand-text-muted">
            This user has no active membership.
          </p>
        )}
      </SurfacePanel>

      <SurfacePanel className="bg-transparent">
        <div className="border-b border-brand-border px-5 py-4">
          <h2 className="text-base font-semibold text-brand-text">Membership history</h2>
          <p className="mt-1 text-sm text-brand-text-subtle">
            Every previous membership and the transaction that paid for it.
          </p>
        </div>

        {past.length === 0 ? (
          <p className="px-5 py-6 text-sm text-brand-text-muted">No previous memberships.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[42rem] border-collapse text-left">
              <thead>
                <tr className="border-b border-brand-border bg-white/[0.02]">
                  {["Tier", "Status", "Started", "Expires", "Cancelled", "Transaction"].map(
                    (header) => (
                      <th
                        key={header}
                        className="px-4 py-2.5 text-xs font-semibold uppercase tracking-wide text-brand-text-subtle"
                      >
                        {header}
                      </th>
                    ),
                  )}
                </tr>
              </thead>
              <tbody>
                {past.map((membership) => {
                  const state = getMembershipState(membership);

                  return (
                    <Fragment key={membership.id}>
                      <tr className="border-b border-brand-border/70 text-sm last:border-b-0">
                        <td className="px-4 py-3 text-brand-text">
                          {tierTitleForMembership(membership)}
                        </td>
                        <td className="px-4 py-3">
                          <StatusBadge tone={STATE_TONE[state]}>{titleCase(state)}</StatusBadge>
                        </td>
                        <td className="px-4 py-3 text-brand-text-muted">
                          {formatDate(membership.started_at)}
                        </td>
                        <td className="px-4 py-3 text-brand-text-muted">
                          {formatDate(membership.expires_at)}
                        </td>
                        <td className="px-4 py-3 text-brand-text-muted">
                          {membership.cancelled_at ? formatDate(membership.cancelled_at) : "—"}
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex items-center justify-between">
                            <TransactionSummary membership={membership} />
                            <button
                              type="button"
                              onClick={() => toggleExpanded(membership.id)}
                              className="text-sm text-brand-primary underline"
                            >
                              {expanded[membership.id] ? "Hide details" : "Details"}
                            </button>
                          </div>
                        </td>
                      </tr>
                      {expanded[membership.id] ? (
                        <tr key={`${membership.id}-details`} className="bg-white/[0.02]">
                          <td colSpan={6} className="px-0">
                            <TransactionDetails membership={membership} />
                          </td>
                        </tr>
                      ) : null}
                    </Fragment>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </SurfacePanel>
    </div>
  );
}
