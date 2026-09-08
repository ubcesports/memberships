import { queryOptions, useQuery } from "@tanstack/react-query";
import { redirectToSignIn } from "./auth";
import apiClient from "./client";

import type { SessionResponse } from "@/lib/types/user.types";

async function fetchProfile(signal?: AbortSignal) {
  const response = await apiClient.get<SessionResponse>("/profile", {
    signal,
    validateStatus: (status) => status === 200 || status === 401,
  });

  if (response.status === 401 || !response.data.user) {
    return null;
  }

  const { user } = response.data;

  return {
    name: user.full_name,
    email: user.email,
    emailVerifiedAt: user.email_verified_at ? new Date(user.email_verified_at) : undefined,
    groups: user.groups ?? [],
    avatarUrl: user.avatar_url ?? undefined,
    onboardingCompletedAt: user.onboarding_completed_at
      ? new Date(user.onboarding_completed_at)
      : undefined,
    isStudent: user.is_student,
    studentId: user.student_id ?? undefined,
    role: user.role,
    createdAt: new Date(user.created_at),
    updatedAt: new Date(user.updated_at),
  };
}

const query = queryOptions({
  queryKey: ["auth", "profile"],
  queryFn: async ({ signal }) => {
    const profile = await fetchProfile(signal);

    if (!profile) {
      await redirectToSignIn(window.location.href);
      throw new Error("Authentication required");
    }

    return profile;
  },
});

export const useProfile = (options?: Partial<typeof query>) => useQuery({ ...query, ...options });

const optionalProfileQuery = queryOptions({
  queryKey: ["auth", "optional-profile"],
  queryFn: ({ signal }) => fetchProfile(signal),
});

export const useOptionalProfile = () => useQuery(optionalProfileQuery);
