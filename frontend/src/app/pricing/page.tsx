"use client";

import { Check } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { BasePage } from "@/components/layout/base-page";
import { AssignedPassCard } from "@/components/membership/assigned-pass-card";
import { CatalogLoading } from "@/components/membership/catalog-loading";
import { DayPassCard } from "@/components/membership/day-pass-card";
import { SeasonPassCard } from "@/components/membership/season-pass-card";
import { redirectToSignIn } from "@/lib/auth";
import apiClient from "@/lib/client";
import {
  type EligibleMembershipTier,
  type MembershipTier,
  useEligibleMembershipTiers,
  useMembershipCatalog,
  useOptionalProfile,
} from "@/lib/membership.hook";

type CheckoutResponse = {
  url: string;
};

type ApiErrorResponse = {
  code?: string;
  message?: string;
};

const RESTRICTED_TIER_SLUGS = ["competitive_team", "executive"];

function getApiErrorMessage(data: CheckoutResponse | ApiErrorResponse) {
  return "message" in data && data.message
    ? data.message
    : "Unable to open checkout. Refresh and try again.";
}

export default function PricingPage() {
  const {
    data: catalog,
    isPending: catalogPending,
    isError: catalogError,
  } = useMembershipCatalog();
  const { data: profile } = useOptionalProfile();
  const isSignedIn = Boolean(profile);
  const canLoadEligibility = Boolean(profile?.onboardingCompletedAt);
  const {
    data: eligibleTiers = [],
    isPending: eligibilityPending,
    isError: eligibilityError,
  } = useEligibleMembershipTiers(canLoadEligibility);

  const { mutate: signIn, isPending: signInPending } = useMutation({
    mutationFn: async () => await redirectToSignIn(window.location.href),
    onError: () => toast.error("Unable to start sign in. Try again."),
  });

  const {
    mutate: checkout,
    variables: checkoutTier,
    isPending: checkoutPending,
  } = useMutation({
    mutationFn: async (tier: EligibleMembershipTier) => {
      const response = await apiClient.post<
        CheckoutResponse | ApiErrorResponse
      >(
        "/membership/checkout",
        { tier_id: tier.id },
        {
          validateStatus: (status) => status >= 200 && status < 500,
        },
      );

      if (response.status !== 200 || !("url" in response.data)) {
        throw new Error(getApiErrorMessage(response.data));
      }

      window.location.assign(response.data.url);
    },
    onError: (error) =>
      toast.error(
        error instanceof Error
          ? error.message
          : "Unable to open checkout. Refresh and try again.",
      ),
  });

  const tierBySlug = (slug: string) =>
    catalog?.find((tier) => tier.slug === slug);
  const eligibleBySlug = (slug: string) =>
    eligibleTiers.find((tier) => tier.slug === slug);

  const mainTiers = [tierBySlug("regular"), tierBySlug("premium")].filter(
    (tier): tier is MembershipTier => Boolean(tier),
  );
  const dayTier = tierBySlug("day");
  const assignedTiers = eligibleTiers.filter((tier) =>
    RESTRICTED_TIER_SLUGS.includes(tier.slug),
  );

  return (
    <BasePage>
      <section className="pb-12 pt-16 text-center sm:pb-14 sm:pt-20">
        <p className="font-mono text-xs font-semibold uppercase tracking-[0.24em] text-blue-200">
          UBCEA membership passes
        </p>
        <h1 className="mx-auto mt-5 max-w-3xl text-4xl font-semibold text-brand-text sm:text-5xl lg:text-6xl">
          Choose your UBCEA pass
        </h1>
        <p className="mx-auto mt-5 max-w-2xl text-base leading-7 text-brand-text-muted sm:text-lg">
          Compare Regular and Premium membership pricing, then sign in to see
          the lowest price available to your account.
        </p>
      </section>

      {profile && !profile.onboardingCompletedAt ? (
        <div className="mb-6 border border-amber-300/35 bg-amber-300/10 px-5 py-4 text-sm text-amber-100">
          Finish your account setup to see your eligible membership prices and
          purchase a pass.
        </div>
      ) : profile && eligibilityError ? (
        <div className="mb-6 border border-red-400/35 bg-red-400/10 px-5 py-4 text-sm text-red-100">
          Your personalized prices could not be loaded. Refresh the page to try
          again.
        </div>
      ) : profile ? (
        <div className="mb-6 flex flex-col gap-2 border border-brand-primary/35 bg-brand-primary/10 px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-3">
            <Check
              aria-hidden="true"
              className="size-4 shrink-0 text-blue-200"
            />
            <p className="text-sm text-blue-50">
              Showing eligible checkout options for{" "}
              {profile.name || profile.email}.
            </p>
          </div>
          <p className="text-xs font-medium text-blue-100">
            Restricted passes stay hidden unless assigned to your account.
          </p>
        </div>
      ) : null}

      {catalogPending || (canLoadEligibility && eligibilityPending) ? (
        <CatalogLoading />
      ) : catalogError || !catalog ? (
        <div className="border border-red-400/35 bg-red-400/10 px-6 py-10 text-center text-sm text-red-100">
          Membership prices are unavailable right now. Refresh the page to try
          again.
        </div>
      ) : (
        <>
          {assignedTiers.length > 0 ? (
            <section
              className="pb-10"
              aria-labelledby="assigned-passes-heading"
            >
              <div className="mb-5">
                <p className="font-mono text-xs font-semibold uppercase tracking-[0.2em] text-brand-text-subtle">
                  Assigned access
                </p>
                <h2
                  id="assigned-passes-heading"
                  className="mt-2 text-xl font-semibold text-brand-text"
                >
                  Your restricted passes
                </h2>
              </div>
              <div className="grid gap-5">
                {assignedTiers.map((tier) => (
                  <AssignedPassCard
                    key={tier.id}
                    tier={tier}
                    checkoutPending={
                      checkoutPending && checkoutTier?.id === tier.id
                    }
                    onCheckout={checkout}
                  />
                ))}
              </div>
            </section>
          ) : null}

          <section aria-labelledby="season-passes-heading">
            <div className="mb-5 flex items-end justify-between gap-4">
              <div>
                <p className="font-mono text-xs font-semibold uppercase tracking-[0.2em] text-brand-text-subtle">
                  Main passes
                </p>
                <h2
                  id="season-passes-heading"
                  className="mt-2 text-xl font-semibold text-brand-text"
                >
                  Season memberships
                </h2>
              </div>
              <p className="hidden text-sm text-brand-text-subtle sm:block">
                One-time payment
              </p>
            </div>
            <div className="grid gap-5 lg:grid-cols-2">
              {mainTiers.map((tier) => (
                <SeasonPassCard
                  key={tier.id}
                  tier={tier}
                  eligibleTier={eligibleBySlug(tier.slug)}
                  checkoutPending={
                    checkoutPending && checkoutTier?.id === tier.id
                  }
                  isSignedIn={isSignedIn}
                  onCheckout={checkout}
                  onSignIn={() => signIn()}
                  signInPending={signInPending}
                />
              ))}
            </div>
          </section>

          {dayTier ? (
            <section
              className="pb-20 pt-12"
              aria-labelledby="additional-passes-heading"
            >
              <div className="mb-5">
                <p className="font-mono text-xs font-semibold uppercase tracking-[0.2em] text-brand-text-subtle">
                  Additional option
                </p>
                <h2
                  id="additional-passes-heading"
                  className="mt-2 text-xl font-semibold text-brand-text"
                >
                  Just here for the day?
                </h2>
              </div>
              <DayPassCard
                tier={dayTier}
                eligibleTier={eligibleBySlug(dayTier.slug)}
                checkoutPending={
                  checkoutPending && checkoutTier?.id === dayTier.id
                }
                isSignedIn={isSignedIn}
                onCheckout={checkout}
                onSignIn={() => signIn()}
                signInPending={signInPending}
              />
            </section>
          ) : null}
        </>
      )}
    </BasePage>
  );
}
