"use client";

import { Ban, Loader2 } from "lucide-react";
import { Fragment, useMemo, useState } from "react";
import { toast } from "sonner";
import { ActionButton } from "@/components/action-button";
import { DetailRow } from "@/components/detail-row";
import { StatusBadge } from "@/components/status-badge";
import { SurfacePanel } from "@/components/surface-panel";
import type { UpdateUserRequest } from "@/lib/types/admin.types";
import { useMembershipCatalog } from "@/lib/membership.hook";
import { formatDate, formatTime } from "@/lib/utils/formatting";
import { titleCase } from "@/lib/utils/groups";

type UserMembershipsPanelProps = {
  memberships: Membership[];
  onSave: (body: UpdateUserRequest) => Promise<unknown>;
  isSaving: boolean;
};

import type { Membership, MembershipStatus } from "@/lib/types/membership.types";

function getMembershipState(membership: Membership): MembershipStatus {
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

function TransactionDetails({ membership }: { membership: Membership }) {
  const tx = membership.transaction;

  return (
    <div className="px-5 py-4 text-sm text-brand-text-muted">
      <div className="grid grid-cols-2 gap-4">
        <div>
          <div className="font-medium text-brand-text">Transaction ID</div>
          <div>{tx.id}</div>
        </div>
        <div>
          <div className="font-medium text-brand-text">Amount Paid</div>
          <div>{tx.amount_paid}</div>
        </div>
        <div>
          <div className="font-medium text-brand-text">Group at purchase</div>
          <div>{titleCase(tx.group_at_purchase)}</div>
        </div>
        <div>
          <div className="font-medium text-brand-text">Student at purchase</div>
          <div>{titleCase(tx.student_at_purchase.toString())}</div>
        </div>
        <div>
          <div className="font-medium text-brand-text">Purchase type</div>
          <div>{titleCase(tx.purchase_type)}</div>
        </div>
        <div>
          <div className="font-medium text-brand-text">Stripe payment intent ID</div>
          <div>{tx.stripe_payment_intent_id}</div>
        </div>
      </div>
    </div>
  );
}

export function UserMembershipsPanel({ memberships, onSave, isSaving }: UserMembershipsPanelProps) {
  const { data: catalog } = useMembershipCatalog();
  const [pendingCancellationId, setPendingCancellationId] = useState<string | null>(null);

  const activeMemberships = memberships.filter(
    (membership) => getMembershipState(membership) === "active",
  );
  const activeMembershipIds = new Set(activeMemberships.map((membership) => membership.id));
  const past = memberships.filter((membership) => !activeMembershipIds.has(membership.id));

  const [expanded, setExpanded] = useState<Record<string, boolean>>({});

  const toggleExpanded = (id: string) => {
    setExpanded((prev) => ({ ...prev, [id]: !prev[id] }));
  };

  const cancelMembership = async (membership: Membership) => {
    try {
      await onSave({ cancel_membership_id: membership.id });
      toast.success(`${membership.tier_title} membership cancelled`);
    } catch {
      // The API client already surfaces the error message as a toast.
    } finally {
      setPendingCancellationId(null);
    }
  };

  return (
    <div className="flex flex-col gap-4">
      <SurfacePanel className="bg-transparent">
        <div className="border-b border-brand-border px-5 py-4">
          <h2 className="text-base font-semibold text-brand-text">Current memberships</h2>
          <p className="mt-1 text-sm text-brand-text-subtle">
            The active memberships this user holds across all programs.
          </p>
        </div>

        {activeMemberships.length > 0 ? (
          <div className="divide-y divide-brand-border">
            {activeMemberships.map((membership) => (
              <section key={membership.id} aria-labelledby={`active-membership-${membership.id}`}>
                <div className="flex flex-wrap items-start justify-between gap-3 bg-white/2 px-5 py-4">
                  <div>
                    <h3
                      id={`active-membership-${membership.id}`}
                      className="mt-1 text-base font-semibold text-brand-text"
                    >
                      Tier: {membership.tier_title}
                    </h3>
                    <p className="text-sm font-semibold text-brand-text-subtle">
                      Program: {membership.program_name}
                    </p>
                  </div>
                  <div className="flex flex-wrap items-center justify-end gap-2">
                    <StatusBadge tone="success">Active</StatusBadge>
                    {pendingCancellationId === membership.id ? (
                      <>
                        <ActionButton
                          onClick={() => cancelMembership(membership)}
                          loading={isSaving}
                          loadingIcon={
                            <Loader2 aria-hidden="true" className="size-4 animate-spin" />
                          }
                          className="border-amber-300/40 text-amber-100"
                        >
                          Confirm cancel
                        </ActionButton>
                        <ActionButton
                          onClick={() => setPendingCancellationId(null)}
                          disabled={isSaving}
                        >
                          Keep membership
                        </ActionButton>
                      </>
                    ) : (
                      <ActionButton
                        onClick={() => setPendingCancellationId(membership.id)}
                        disabled={isSaving || pendingCancellationId !== null}
                        icon={<Ban aria-hidden="true" className="size-4" />}
                      >
                        Cancel membership
                      </ActionButton>
                    )}
                  </div>
                </div>

                <dl>
                  <DetailRow label="Started">{formatTime(membership.started_at)}</DetailRow>
                  <DetailRow label="Expires">{formatTime(membership.expires_at)}</DetailRow>
                  <DetailRow label="Transaction">
                    <div className="flex items-center justify-between gap-4">
                      <div className="flex flex-wrap items-center gap-2">
                        <StatusBadge tone="muted">
                          {titleCase(membership.transaction.status)}
                        </StatusBadge>
                      </div>
                      <button
                        type="button"
                        onClick={() => toggleExpanded(membership.id)}
                        className="text-sm text-brand-primary underline"
                      >
                        {expanded[membership.id] ? "Hide details" : "Details"}
                      </button>
                    </div>
                  </DetailRow>
                </dl>
                {expanded[membership.id] ? <TransactionDetails membership={membership} /> : null}
              </section>
            ))}
          </div>
        ) : (
          <p className="px-5 py-6 text-sm text-brand-text-muted">
            This user has no active memberships.
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
            <table className="w-full min-w-2xl border-collapse text-left">
              <thead>
                <tr className="border-b border-brand-border bg-white/2">
                  {["Program", "Tier", "Status", "Started", "Expires", "Cancelled"].map(
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
                        <td className="px-4 py-3 text-brand-text-muted">
                          {membership.program_name}
                        </td>
                        <td className="px-4 py-3 text-brand-text">{membership.tier_title}</td>
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
                        <tr key={`${membership.id}-details`} className="bg-white/2">
                          <td colSpan={7} className="px-0">
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
