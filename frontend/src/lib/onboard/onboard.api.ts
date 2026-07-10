import apiClient from "@/lib/client";
import type { ApiErrorResponse, CompleteOnboardingPayload } from "./onboard.types";

export const ONBOARDING_REQUIRED = "ONBOARDING_REQUIRED";

export function getApiErrorMessage(data?: ApiErrorResponse) {
  return data?.detail || data?.message || "Unable to complete onboarding. Try again.";
}

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
  const response = await apiClient.post<ApiErrorResponse>("/onboard", payload, {
    validateStatus: () => true,
  });

  if (response.status === 200) {
    return { destination: "/" };
  }

  if (response.status === 401) {
    return { destination: "/login" };
  }

  throw new Error(getApiErrorMessage(response.data));
}
