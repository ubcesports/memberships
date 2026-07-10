import { queryOptions, useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/client";

export type MembershipTierPrice = {
  price: number;
  price_id: string;
  is_student_required: boolean | null;
};

export type MembershipTier = {
  id: string;
  title: string;
  description: string;
  benefits: string[];
  slug: string;
  product_id: string;
  prices: MembershipTierPrice[];
};

export type EligibleMembershipTier = Omit<MembershipTier, "prices"> & {
  purchase_type: "new" | "upgrade";
  prices: MembershipTierPrice;
};

export type OptionalProfile = {
  name: string;
  email: string;
  onboardingCompletedAt?: string;
};

type ProfileResponse = {
  user?: {
    email: string;
    full_name: string;
    onboarding_completed_at: string | null;
  };
};

const catalogQuery = queryOptions({
  queryKey: ["membership", "catalog"],
  queryFn: async ({ signal }) => {
    const response = await apiClient.get<MembershipTier[]>(
      "/membership/tiers",
      { signal },
    );

    return response.data;
  },
});

const optionalProfileQuery = queryOptions({
  queryKey: ["auth", "optional-profile"],
  queryFn: async ({ signal }) => {
    const response = await apiClient.get<ProfileResponse>("/profile", {
      signal,
      validateStatus: (status) => status === 200 || status === 401,
    });

    if (response.status === 401 || !response.data.user) {
      return null;
    }

    return {
      name: response.data.user.full_name,
      email: response.data.user.email,
      onboardingCompletedAt:
        response.data.user.onboarding_completed_at ?? undefined,
    };
  },
  retry: false,
});

export const useMembershipCatalog = () => useQuery(catalogQuery);

export const useOptionalProfile = () => useQuery(optionalProfileQuery);

export const useEligibleMembershipTiers = (enabled: boolean) =>
  useQuery({
    queryKey: ["membership", "eligible"],
    queryFn: async ({ signal }) => {
      const response = await apiClient.get<EligibleMembershipTier[] | null>(
        "/membership/tiers/eligible",
        {
          signal,
          validateStatus: (status) =>
            status === 200 || status === 401 || status === 403,
        },
      );

      if (response.status !== 200) {
        return [];
      }

      return response.data ?? [];
    },
    enabled,
  });
