"use client";

import { Loader2 } from "lucide-react";
import { BasePage } from "@/components/layout/base-page";
import { useOnboardCheck } from "@/lib/onboard/onboard.hook";

export default function OnboardCheckPage() {
  useOnboardCheck();

  return (
    <BasePage>
      <div className="flex flex-1 items-center justify-center py-12">
        <section className="w-full max-w-md border border-brand-border bg-brand-surface/85 px-5 py-6 shadow-2xl shadow-black/25 sm:px-6">
          <div className="flex items-center gap-3 text-brand-text">
            <Loader2 aria-hidden="true" className="size-5 animate-spin text-brand-primary" />
            <h1 className="text-lg font-semibold">Logging you in...</h1>
          </div>
          <p className="mt-3 text-sm leading-6 text-brand-text-muted">
            You will be redirected once your account status is confirmed.
          </p>
        </section>
      </div>
    </BasePage>
  );
}
