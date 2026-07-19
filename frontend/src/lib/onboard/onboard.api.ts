import apiClient from "@/lib/client";
import type { ApiErrorResponse } from "@/lib/client";
import type { CompleteOnboardingPayload } from "./onboard.types";

export const ONBOARDING_REQUIRED = "ONBOARDING_REQUIRED";

export async function checkOnboardingStatus() {
  const response = await apiClient.get<ApiErrorResponse>("/onboard/check", {
    validateStatus: (status) => status === 200 || status === 401 || status === 403,
  });

  if (response.status === 200) {
    return { destination: "/" };
  }

  if (response.status === 401) {
    return { destination: "/login" };
  }

  if (response.status === 403 && response.data.code === ONBOARDING_REQUIRED) {
    return { destination: "/onboard" };
  }

  return { destination: "/403" };
}

export async function completeOnboarding(payload: CompleteOnboardingPayload) {
  await apiClient.post("/onboard", payload);
  return { destination: "/" };
}
