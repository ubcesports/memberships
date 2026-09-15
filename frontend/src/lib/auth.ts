import apiClient from "@/lib/client";

import type { OAuthAuthorizeResponse } from "@/lib/types/user.types";

export async function redirectToSignIn(returnTo: string) {
  const response = await apiClient.get<OAuthAuthorizeResponse>("/auth/oauth/zetrova/authorize", {
    params: {
      redirect_uri: returnTo,
      error_redirect_uri: returnTo,
    },
  });

  window.location.assign(response.data.url);
}
