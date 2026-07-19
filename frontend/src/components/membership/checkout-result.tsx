"use client";

import { Check, CircleCheckBig, RotateCcw, XCircle } from "lucide-react";
import Link from "next/link";
import { useRedirectCountdown } from "@/lib/use-redirect-countdown.hook";

const REDIRECT_DELAY_SECONDS = 6;

type CheckoutResultProps = {
  successful: boolean;
};

export function CheckoutResult({ successful }: CheckoutResultProps) {
  const secondsRemaining = useRedirectCountdown("/pricing", REDIRECT_DELAY_SECONDS);

  return (
    <section className="mx-auto flex w-full max-w-xl flex-col items-center border border-brand-border bg-brand-surface/90 px-6 py-10 text-center shadow-2xl shadow-black/25 sm:px-10 sm:py-12">
      <div
        className={`flex size-14 items-center justify-center border ${
          successful
            ? "border-emerald-300/40 bg-emerald-400/10 text-emerald-200"
            : "border-amber-300/40 bg-amber-400/10 text-amber-200"
        }`}
      >
        {successful ? (
          <CircleCheckBig aria-hidden="true" className="size-7" />
        ) : (
          <XCircle aria-hidden="true" className="size-7" />
        )}
      </div>

      <p className="mt-6 font-mono text-xs font-semibold uppercase tracking-[0.22em] text-brand-text-subtle">
        {successful ? "Checkout complete" : "Checkout cancelled"}
      </p>
      <h1 className="mt-3 text-3xl font-semibold text-brand-text sm:text-4xl">
        {successful ? "Your checkout is complete" : "Nothing was changed"}
      </h1>
      <p className="mt-4 max-w-md text-sm leading-6 text-brand-text-muted sm:text-base">
        {successful
          ? "We're confirming your payment and membership now. It may take a moment to appear on your account."
          : "You left checkout before completing payment. Your previous membership, if any, remains active."}
      </p>

      <div className="mt-8 grid w-full gap-3 sm:grid-cols-2">
        <Link
          href="/pricing"
          className="inline-flex h-11 items-center justify-center gap-2 bg-brand-primary px-5 text-sm font-semibold text-white transition hover:bg-brand-primary-hover"
        >
          {successful ? (
            <Check aria-hidden="true" className="size-4" />
          ) : (
            <RotateCcw aria-hidden="true" className="size-4" />
          )}
          {successful ? "View passes" : "Return to passes"}
        </Link>
        <Link
          href="/profile"
          className="inline-flex h-11 items-center justify-center border border-brand-border px-5 text-sm font-semibold text-brand-text transition hover:border-brand-text-muted hover:bg-white/5"
        >
          View profile
        </Link>
      </div>

      <p className="mt-6 text-xs text-brand-text-subtle" aria-live="polite">
        Returning to passes in {secondsRemaining} second
        {secondsRemaining === 1 ? "" : "s"}.
      </p>
    </section>
  );
}
