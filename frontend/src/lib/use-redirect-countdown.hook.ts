"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export function useRedirectCountdown(destination: string, delaySeconds: number) {
  const router = useRouter();
  const [secondsRemaining, setSecondsRemaining] = useState(delaySeconds);

  useEffect(() => {
    const countdown = window.setInterval(() => {
      setSecondsRemaining((current) => Math.max(0, current - 1));
    }, 1_000);
    const redirect = window.setTimeout(() => {
      router.replace(destination);
    }, delaySeconds * 1_000);

    return () => {
      window.clearInterval(countdown);
      window.clearTimeout(redirect);
    };
  }, [delaySeconds, destination, router]);

  return secondsRemaining;
}
