import type {
  MembershipTier,
  EligibleMembershipTier,
  Membership,
} from "@/lib/types/membership.types";
import { queryOptions, useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/client";

const catalogQuery = queryOptions({
  queryKey: ["membership", "catalog"],
  queryFn: async ({ signal }) => {
    const response = await apiClient.get<MembershipTier[]>("/membership/tiers", { signal });

    return response.data;
  },
});

export const useMembershipCatalog = (initialData?: MembershipTier[]) =>
  useQuery({
    ...catalogQuery,
    initialData,
  });

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

export const useAllMemberships = () =>
  useQuery({
    queryKey: ["membership", "all"],
    queryFn: async ({ signal }) => {
      const response = await apiClient.get<Membership[]>("/membership/me/all", { signal });
      return response.data;
    },
  });
