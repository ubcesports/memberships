"use client";

import { useMutation } from "@tanstack/react-query";
import { Loader2, LogIn } from "lucide-react";
import { ActionButton } from "@/components/action-button";
import { BasePage } from "@/components/layout/base-page";
import { redirectToSignIn } from "@/lib/auth";

const POST_AUTH_PATH = "/onboard/check";

export default function LoginPage() {
  const {
    mutate: continueWithJasperLabs,
    error,
    isPending,
  } = useMutation({
    mutationFn: async () => {
      await redirectToSignIn(`${window.location.origin}${POST_AUTH_PATH}`);
    },
  });

  return (
    <BasePage>
      <div className="flex flex-1 items-center justify-center py-12">
        <section className="w-full max-w-md border border-brand-border bg-brand-surface/85 shadow-2xl shadow-black/25">
          <div className="border-b border-brand-border px-5 py-5 sm:px-6">
            <p className="text-sm font-semibold text-brand-primary">UBCEA Memberships</p>
            <h1 className="mt-3 text-2xl font-semibold text-brand-text">Log in or sign up</h1>
            <p className="mt-2 text-sm leading-6 text-brand-text-muted">
              Continue with your JasperLabs account to access your membership profile.
            </p>
          </div>

          <div className="px-5 py-5 sm:px-6">
            <ActionButton
              className="h-12 w-full border-brand-primary bg-brand-primary text-base hover:border-brand-primary-hover hover:bg-brand-primary-hover"
              onClick={() => continueWithJasperLabs()}
              loading={isPending}
              icon={<LogIn aria-hidden="true" className="size-5" />}
              loadingIcon={<Loader2 aria-hidden="true" className="size-5 animate-spin" />}
            >
              {isPending ? "Redirecting" : "Continue with JasperLabs"}
            </ActionButton>

            {error && (
              <p className="mt-4 text-sm leading-6 text-brand-text-muted">
                Unable to start JasperLabs sign in. Try again.
              </p>
            )}
          </div>
        </section>
      </div>
    </BasePage>
  );
}
