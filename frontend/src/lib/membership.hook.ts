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
