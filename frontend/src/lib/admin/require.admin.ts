import { useRouter } from "next/dist/client/components/navigation";
import { useProfile } from "../profile.hook";
import { useEffect } from "react";

export function useRequireAdmin() {
  const router = useRouter();
  const { data: profile, isPending: isProfilePending } = useProfile();
  const isAdmin = profile?.role === "admin";

  useEffect(() => {
    if (!isProfilePending && profile && profile.role !== "admin") {
      router.replace("/403");
    }
  }, [isProfilePending, profile, router]);

  return { isAdmin, isProfilePending };
}
