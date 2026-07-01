import { useQuery } from "@tanstack/react-query";
import { checkOnboardingStatus } from "./onboard.api";

export function useOnboardCheck() {
  return useQuery({
    queryKey: ["onboard", "check"],
    retry: false,
    queryFn: async () => {
      const result = await checkOnboardingStatus();
      window.location.replace(result.destination);
      return result;
    },
  });
}
