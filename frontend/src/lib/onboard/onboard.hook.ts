import { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { checkOnboardingStatus } from "./onboard.api";

export function useOnboardCheck() {
  const { data, error } = useQuery({
    queryKey: ["onboard", "check"],
    retry: false,
    queryFn: checkOnboardingStatus,
  });

  useEffect(() => {
    if (data) {
      window.location.replace(data.destination);
    }
  }, [data]);

  useEffect(() => {
    if (error) {
      console.error(error);
    }
  }, [error]);
}
