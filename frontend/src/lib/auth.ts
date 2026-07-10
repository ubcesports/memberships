import { API_BASE } from "@/lib/client";

export function redirectToSignIn(returnTo: string) {
  const url = new URL(`${API_BASE}/auth/oauth/jasperlabs/authorize`);
  url.searchParams.set("redirect_uri", returnTo);
  url.searchParams.set("error_redirect_uri", returnTo);

  window.location.assign(url.toString());
}
