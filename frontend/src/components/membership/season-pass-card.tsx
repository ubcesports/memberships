import { CalendarDays, Check, Sparkles } from "lucide-react";
import { PurchaseButton } from "@/components/membership/purchase-button";
import {
    formatMembershipPrice,
    membershipGroupLabel,
} from "@/components/membership/pricing";
import type {
    EligibleMembershipTier,
    MembershipTier,
} from "@/lib/membership.hook";
import { formatDate } from "@/lib/utils/formatting";

type SeasonPassCardProps = {
    tier: MembershipTier;
    eligibleTier?: EligibleMembershipTier;
    currentTierSlug?: string;
    checkoutPending: boolean;
    isSignedIn: boolean;
    onCheckout: (tier: EligibleMembershipTier) => void;
    onSignIn: () => void;
    signInPending: boolean;
};

export function SeasonPassCard({
    tier,
    eligibleTier,
    currentTierSlug,
    checkoutPending,
    isSignedIn,
    onCheckout,
    onSignIn,
    signInPending,
}: SeasonPassCardProps) {
    const featured = tier.slug === "premium";
    const isCurrent = currentTierSlug === tier.slug;

    return (
        <article
            className={`relative flex min-h-124 flex-col border bg-brand-surface/88 shadow-2xl shadow-black/20 ${
                featured
                    ? "border-brand-primary shadow-brand-primary/10"
                    : "border-brand-border"
            }`}
        >
            {featured ? (
                <div className="absolute inset-x-0 top-0 h-1 bg-brand-primary" />
            ) : null}
            <div className="flex items-start justify-between gap-4 border-b border-brand-border px-6 py-6 sm:px-7">
                <div>
                    <p className="mb-2 font-mono text-xs font-semibold uppercase tracking-[0.2em] text-brand-text-subtle">
                        Season pass
                    </p>
                    <h2 className="text-2xl font-semibold text-brand-text">
                        {tier.title}
                    </h2>
                </div>
                {featured ? (
                    <span className="inline-flex items-center gap-1.5 border border-brand-primary/70 bg-brand-primary/15 px-3 py-1.5 text-xs font-semibold text-blue-100">
                        <Sparkles aria-hidden="true" className="size-3.5" />
                        Premium
                    </span>
                ) : isCurrent ? (
                    <span className="border border-emerald-400/40 bg-emerald-400/10 px-3 py-1.5 text-xs font-semibold text-emerald-100">
                        Active
                    </span>
                ) : null}
            </div>

            <div className="flex flex-1 flex-col p-6 sm:p-7">
                {eligibleTier ? (
                    <div className="border border-brand-primary/35 bg-brand-primary/10 p-5">
                        <div className="flex items-center justify-between gap-3">
                            <span className="text-sm font-medium text-blue-100">
                                Your{" "}
                                {membershipGroupLabel(eligibleTier.price.group)}{" "}
                                price
                            </span>
                            <Check
                                aria-hidden="true"
                                className="size-4 text-blue-200"
                            />
                        </div>
                        <div className="mt-3 flex items-end gap-2">
                            <span className="text-4xl font-semibold tracking-tight text-brand-text">
                                {formatMembershipPrice(
                                    eligibleTier.amount_due_minor,
                                    eligibleTier.price.currency,
                                )}
                            </span>
                            <span className="pb-1 text-sm text-brand-text-subtle">
                                one time
                            </span>
                        </div>
                        {eligibleTier.is_upgrade ? (
                            <div className="mt-3 flex justify-between border-t border-brand-primary/25 pt-3 text-xs text-brand-text-muted">
                                <span>
                                    Pass price{" "}
                                    {formatMembershipPrice(
                                        eligibleTier.price.amount_minor,
                                        eligibleTier.price.currency,
                                    )}
                                </span>
                                <span>
                                    −
                                    {formatMembershipPrice(
                                        eligibleTier.credit_amount_minor,
                                        eligibleTier.price.currency,
                                    )}{" "}
                                    credit
                                </span>
                            </div>
                        ) : null}
                    </div>
                ) : (
                    <div className="grid grid-cols-2 border border-brand-border">
                        {["student", "member"].map((group, index) => {
                            const price = tier.prices.find(
                                (item) => item.group === group,
                            );
                            return (
                                <div
                                    key={group}
                                    className={`p-4 sm:p-5 ${index === 0 ? "border-r border-brand-border" : ""}`}
                                >
                                    <p className="text-xs font-medium uppercase tracking-wider text-brand-text-subtle">
                                        {membershipGroupLabel(group)}
                                    </p>
                                    <p className="mt-2 text-2xl font-semibold text-brand-text">
                                        {price
                                            ? formatMembershipPrice(
                                                  price.amount_minor,
                                                  price.currency,
                                              )
                                            : "—"}
                                    </p>
                                </div>
                            );
                        })}
                    </div>
                )}

                <div className="mt-6 flex-1">
                    <p className="text-sm leading-6 text-brand-text-muted">
                        {tier.description ||
                            `${tier.title} UBCEA membership pass.`}
                    </p>
                    <div className="mt-5 flex items-start gap-3 border-t border-brand-border/70 pt-5">
                        <CalendarDays
                            aria-hidden="true"
                            className="mt-0.5 size-4 shrink-0 text-brand-text-subtle"
                        />
                        <div>
                            <p className="text-sm font-medium text-brand-text">
                                Valid through {formatDate(tier.expires_at)}
                            </p>
                            <p className="mt-1 text-xs leading-5 text-brand-text-subtle">
                                Regular and Premium passes expire at the end of
                                the membership year.
                            </p>
                        </div>
                    </div>
                </div>

                <div className="mt-7">
                    <PurchaseButton
                        tier={eligibleTier}
                        isSignedIn={isSignedIn}
                        isCurrent={isCurrent}
                        checkoutPending={checkoutPending}
                        onCheckout={onCheckout}
                        onSignIn={onSignIn}
                        signInPending={signInPending}
                        featured={featured}
                    />
                </div>
            </div>
        </article>
    );
}
