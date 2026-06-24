import { Check, Clock3 } from "lucide-react";
import { PurchaseButton } from "@/components/membership/purchase-button";
import {
    formatMembershipPrice,
    membershipGroupLabel,
} from "@/components/membership/pricing";
import type {
    EligibleMembershipTier,
    MembershipTier,
} from "@/lib/membership.hook";

type DayPassCardProps = {
    tier: MembershipTier;
    eligibleTier?: EligibleMembershipTier;
    currentTierSlug?: string;
    checkoutPending: boolean;
    isSignedIn: boolean;
    onCheckout: (tier: EligibleMembershipTier) => void;
    onSignIn: () => void;
    signInPending: boolean;
};

export function DayPassCard({
    tier,
    eligibleTier,
    currentTierSlug,
    checkoutPending,
    isSignedIn,
    onCheckout,
    onSignIn,
    signInPending,
}: DayPassCardProps) {
    const isCurrent = currentTierSlug === tier.slug;

    return (
        <article className="grid border border-brand-border bg-brand-surface/75 lg:grid-cols-[minmax(0,1fr)_minmax(19rem,0.65fr)]">
            <div className="flex gap-5 p-6 sm:p-7">
                <div className="flex size-11 shrink-0 items-center justify-center border border-brand-border bg-white/4">
                    <Clock3
                        aria-hidden="true"
                        className="size-5 text-brand-text-muted"
                    />
                </div>
                <div>
                    <div className="flex flex-wrap items-center gap-3">
                        <h3 className="text-xl font-semibold text-brand-text">
                            {tier.title}
                        </h3>
                        {isCurrent ? (
                            <span className="border border-emerald-400/40 bg-emerald-400/10 px-2.5 py-1 text-xs font-semibold text-emerald-100">
                                Active
                            </span>
                        ) : null}
                    </div>
                    <p className="mt-2 max-w-2xl text-sm leading-6 text-brand-text-muted">
                        {tier.description ||
                            "Access for the rest of the purchase day."}{" "}
                        Expires at midnight in Vancouver.
                    </p>
                    {!eligibleTier ? (
                        <div className="mt-4 flex flex-wrap gap-x-6 gap-y-2 text-sm">
                            {["student", "member"].map((group) => {
                                const price = tier.prices.find(
                                    (item) => item.group === group,
                                );
                                return (
                                    <span
                                        key={group}
                                        className="text-brand-text-muted"
                                    >
                                        {membershipGroupLabel(group)}{" "}
                                        <strong className="font-semibold text-brand-text">
                                            {price
                                                ? formatMembershipPrice(
                                                      price.amount_minor,
                                                      price.currency,
                                                  )
                                                : "—"}
                                        </strong>
                                    </span>
                                );
                            })}
                        </div>
                    ) : null}
                </div>
            </div>

            <div className="border-t border-brand-border p-6 lg:border-l lg:border-t-0">
                {eligibleTier ? (
                    <div className="mb-4 flex items-end justify-between gap-4">
                        <div>
                            <p className="text-xs font-medium uppercase tracking-wider text-brand-text-subtle">
                                Your{" "}
                                {membershipGroupLabel(eligibleTier.price.group)}{" "}
                                price
                            </p>
                            <p className="mt-1 text-3xl font-semibold text-brand-text">
                                {formatMembershipPrice(
                                    eligibleTier.amount_due_minor,
                                    eligibleTier.price.currency,
                                )}
                            </p>
                        </div>
                        <Check
                            aria-hidden="true"
                            className="mb-2 size-4 text-blue-200"
                        />
                    </div>
                ) : null}
                <PurchaseButton
                    tier={eligibleTier}
                    isSignedIn={isSignedIn}
                    isCurrent={isCurrent}
                    checkoutPending={checkoutPending}
                    onCheckout={onCheckout}
                    onSignIn={onSignIn}
                    signInPending={signInPending}
                />
            </div>
        </article>
    );
}
