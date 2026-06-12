import axios from "axios";
import { API_BASE } from "@/lib/client";

type OAuthAuthorizeResponse = {
  url: string;
};

export async function redirectToJasperLabsSignIn(returnTo: string) {
  const response = await axios.get<OAuthAuthorizeResponse>(
    `${API_BASE}/auth/oauth/jasperlabs/authorize`,
    {
      params: {
        redirect_uri: returnTo,
        error_redirect_uri: returnTo,
      },
      withCredentials: true,
    },
  );

  window.location.assign(response.data.url);
}
