"use client";

import { useMutation } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { useEffect, useRef } from "react";
import { ActionButton } from "@/components/action-button";
import { BasePage } from "@/components/layout/base-page";
import { redirectToSignIn } from "@/lib/auth";

const POST_AUTH_PATH = "/onboard/check";

export default function LoginPage() {
  const hasStartedSignIn = useRef(false);
  const {
    mutate: startSignInFlow,
    isError,
    isPending,
  } = useMutation({
    mutationFn: () => redirectToSignIn(`${window.location.origin}${POST_AUTH_PATH}`),
  });

  useEffect(() => {
    if (hasStartedSignIn.current) return;

    hasStartedSignIn.current = true;
    startSignInFlow();
  }, [startSignInFlow]);

  return (
    <BasePage>
      <div className="flex flex-1 items-center justify-center py-12">
        <div className="flex flex-col items-center gap-4 text-brand-text-muted" role="status">
          {isError ? (
            <ActionButton
              onClick={() => {
                hasStartedSignIn.current = true;
                startSignInFlow();
              }}
              loading={isPending}
              loadingIcon={<Loader2 aria-hidden="true" className="size-5 animate-spin" />}
            >
              Try again
            </ActionButton>
          ) : (
            <>
              <Loader2 aria-hidden="true" className="size-6 animate-spin" />
              <p>Redirecting to sign in…</p>
            </>
          )}
        </div>
      </div>
    </BasePage>
  );
}
