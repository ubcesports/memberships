import { queryOptions, useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/client";

export type MembershipExpirationType = "day" | "semester" | "year";

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
  program_id: string;
  program_name: string;
  expiration_type: MembershipExpirationType;
};

export type EligibleMembershipTier = Omit<MembershipTier, "prices"> & {
  purchase_type: "new" | "upgrade";
  prices: MembershipTierPrice;
};

const catalogQuery = queryOptions({
  queryKey: ["membership", "catalog"],
  queryFn: async ({ signal }) => {
    const response = await apiClient.get<MembershipTier[]>("/membership/tiers", { signal });

    return response.data;
  },
});

export const useMembershipCatalog = () => useQuery(catalogQuery);

export const useEligibleMembershipTiers = (enabled: boolean) =>
  useQuery({
    queryKey: ["membership", "eligible"],
    queryFn: async ({ signal }) => {
      const response = await apiClient.get<EligibleMembershipTier[] | null>(
        "/membership/tiers/eligible",
        {
          signal,
          validateStatus: (status) => status === 200 || status === 401 || status === 403,
        },
      );

      if (response.status !== 200) {
        return [];
      }

      return response.data ?? [];
    },
    enabled,
  });

export type Transaction = {
  id: string;
  amount_paid: string;
  status: "pending" | "completed" | "failed" | "refunded" | "expired";
  group_at_purchase: string;
};

export type Membership = {
  id: string;
  tier_id: string;
  tier_title: string;
  started_at: string;
  expires_at: string;
  cancelled_at: string | null;
  transaction: Transaction;
  slug: string;
  program_id: string;
  program_name: string;
};

export const useAllMemberships = () =>
  useQuery({
    queryKey: ["membership", "all"],
    queryFn: async ({ signal }) => {
      const response = await apiClient.get<Membership[]>("/membership/me/all", { signal });
      return response.data;
    },
  });
