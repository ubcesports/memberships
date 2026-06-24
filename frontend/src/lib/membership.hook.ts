import { queryOptions, useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/client";

export type MembershipPrice = {
    amount_minor: number;
    currency: string;
    group: string;
};

export type MembershipTier = {
    id: string;
    slug: string;
    title: string;
    description: string | null;
    prices: MembershipPrice[];
    expires_at: string;
};

export type EligibleMembershipTier = Omit<MembershipTier, "prices"> & {
    price: MembershipPrice;
    is_upgrade: boolean;
    credit_amount_minor: number;
    amount_due_minor: number;
};

export type ActiveMembership = {
    id: string;
    tier_id: string;
    tier_slug: string;
    tier_title: string;
    tier_description: string | null;
    group_at_purchase: string;
    started_at: string;
    expires_at: string;
    cancelled_at: string | null;
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
        const response = await apiClient.get<{ tiers: MembershipTier[] }>(
            "/membership/tiers",
            { signal },
        );
        return response.data.tiers;
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
            const response = await apiClient.get<{
                tiers: EligibleMembershipTier[];
            }>("/membership/tiers/eligible", { signal });
            return response.data.tiers;
        },
        enabled,
    });

export const useCurrentMembership = (enabled: boolean) =>
    useQuery({
        queryKey: ["membership", "current"],
        queryFn: async ({ signal }) => {
            const response = await apiClient.get<{
                membership: ActiveMembership | null;
            }>("/membership/me", { signal });
            return response.data.membership;
        },
        enabled,
    });
