import { queryOptions, useQuery } from "@tanstack/react-query";
import { redirectToSignIn } from "./auth";
import apiClient from "./client";

type SessionResponse = {
  user?: {
    avatar_url: string | null;
    created_at: string;
    email: string;
    email_verified_at: string | null;
    full_name: string;
    groups: string[];
    id: string;
    is_student: boolean;
    onboarding_completed_at: string | null;
    role: string;
    student_id: string | null;
    updated_at: string;
  };
};

const query = queryOptions({
  queryKey: ["auth", "profile"],
  queryFn: async ({ signal }) => {
    const response = await apiClient.get<SessionResponse>("/profile", {
      signal,
      validateStatus: (status) => status === 200 || status === 401,
    });

    if (response.status === 401 || !response.data.user) {
      await redirectToSignIn(window.location.href);
      throw new Error("Authentication required");
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
  },
});

export const useProfile = (options?: Partial<typeof query>) => useQuery({ ...query, ...options });
