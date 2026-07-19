import { useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient from "./client";

export function useSignOut() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => apiClient.post("/auth/signout", {}),
    onSuccess: () => {
      queryClient.clear();
      window.location.replace("/");
    },
  });
}
