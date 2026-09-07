import axios, { AxiosError } from "axios";
import { toast } from "sonner";

import type { ApiErrorResponse } from "@/lib/types/api.types";

export const API_BASE =
  process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") || "http://localhost:8080";

const apiClient = axios.create({
  baseURL: `${API_BASE}`,
  withCredentials: true,
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiErrorResponse | Blob>) => {
    if (axios.isCancel(error)) {
      return Promise.reject(error);
    }

    const status = error.response?.status;
    const currentPath = typeof window !== "undefined" ? window.location.pathname : "";

    if (typeof window === "undefined") {
      return Promise.reject(error);
    }

    let data = error.response?.data;
    if (data instanceof Blob && data.type.includes("application/json")) {
      try {
        data = JSON.parse(await data.text()) as ApiErrorResponse;
      } catch {
        data = undefined;
      }
    }

    const apiError = data instanceof Blob ? undefined : data;
    const code = apiError?.code;
    const requestId = apiError?.request_id;
    const toastId = `api-error:${error.config?.method ?? "request"}:${error.config?.url ?? "unknown"}`;

    if (status === 401) {
      if (currentPath !== "/login") {
        window.location.replace("/login");
      }
      return Promise.reject(error);
    }

    if (status === 403 && code === "ONBOARDING_REQUIRED") {
      if (currentPath !== "/onboard") {
        window.location.replace("/onboard");
      }
      return Promise.reject(error);
    }

    if (status === 403 && code === "FORBIDDEN") {
      if (currentPath !== "/403") {
        window.location.replace("/403");
      }
      return Promise.reject(error);
    }

    const fallbackMessage = error.response
      ? "Something went wrong. Please try again."
      : "Unable to reach the server. Check your connection and try again.";

    toast.error(apiError?.message || fallbackMessage, {
      id: toastId,
      description: requestId ? `Request ID: ${requestId}` : undefined,
    });

    return Promise.reject(error);
  },
);

export default apiClient;
