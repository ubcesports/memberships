import { Check, Clock3 } from "lucide-react";
import { PurchaseButton } from "@/components/membership/purchase-button";
import {
  formatMembershipPrice,
  getPriceByStudentStatus,
  isMembershipTierPrice,
  membershipPriceLabel,
} from "@/components/membership/pricing";
import type { EligibleMembershipTier, MembershipTier } from "@/lib/membership.hook";

type DayPassCardProps = {
  tier: MembershipTier;
  eligibleTier?: EligibleMembershipTier;
  checkoutPending: boolean;
  isSignedIn: boolean;
  onCheckout: (tier: EligibleMembershipTier) => void;
  onSignIn: () => void;
  signInPending: boolean;
};

export function DayPassCard({
  tier,
  eligibleTier,
  checkoutPending,
  isSignedIn,
  onCheckout,
  onSignIn,
  signInPending,
}: DayPassCardProps) {
  return (
    <article className="grid border border-brand-border bg-brand-surface/75 lg:grid-cols-[minmax(0,1fr)_minmax(19rem,0.65fr)]">
      <div className="flex gap-5 p-6 sm:p-7">
        <div className="flex size-11 shrink-0 items-center justify-center border border-brand-border bg-white/4">
          <Clock3 aria-hidden="true" className="size-5 text-brand-text-muted" />
        </div>
        <div>
          <h3 className="text-xl font-semibold text-brand-text">{tier.title}</h3>
          <p className="mt-2 max-w-2xl text-sm leading-6 text-brand-text-muted">
            {tier.description || "Access for a short visit or single event."}
          </p>
          <BenefitList benefits={tier.benefits} />
          {!eligibleTier ? <InlinePublicPrices tier={tier} /> : null}
        </div>
      </div>

      <div className="border-t border-brand-border p-6 lg:border-l lg:border-t-0">
        {eligibleTier ? (
          <div className="mb-4 flex items-end justify-between gap-4">
            <div>
              <p className="text-xs font-medium uppercase tracking-wider text-brand-text-subtle">
                Your eligible price
              </p>
              <p className="mt-1 text-3xl font-semibold text-brand-text">
                {formatMembershipPrice(eligibleTier.prices.price)}
              </p>
            </div>
            <Check aria-hidden="true" className="mb-2 size-4 text-blue-200" />
          </div>
        ) : null}
        <PurchaseButton
          tier={eligibleTier}
          isSignedIn={isSignedIn}
          checkoutPending={checkoutPending}
          onCheckout={onCheckout}
          onSignIn={onSignIn}
          signInPending={signInPending}
        />
      </div>
    </article>
  );
}

function BenefitList({ benefits }: { benefits: string[] }) {
  if (benefits.length === 0) {
    return null;
  }

  return (
    <ul className="mt-4 grid gap-2 text-sm text-brand-text-muted">
      {benefits.map((benefit) => (
        <li key={benefit} className="flex gap-3">
          <Check aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-blue-200" />
          <span>{benefit}</span>
        </li>
      ))}
    </ul>
  );
}

function InlinePublicPrices({ tier }: { tier: MembershipTier }) {
  const prices = [getPriceByStudentStatus(tier, true), getPriceByStudentStatus(tier, false)].filter(
    isMembershipTierPrice,
  );

  return (
    <div className="mt-4 flex flex-wrap gap-x-6 gap-y-2 text-sm">
      {prices.map((price) => (
        <span key={price.price_id} className="text-brand-text-muted">
          {membershipPriceLabel(price)}{" "}
          <strong className="font-semibold text-brand-text">
            {formatMembershipPrice(price.price)}
          </strong>
        </span>
      ))}
    </div>
  );
}
