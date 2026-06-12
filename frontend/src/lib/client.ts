import axios, { AxiosError } from "axios";

type ApiErrorResponse = {
    code?: string;
    detail?: string;
    message?: string;
};

export const API_BASE =
    process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") ||
    "http://localhost:8080";

const apiClient = axios.create({
    baseURL: `${API_BASE}`,
    withCredentials: true,
});

apiClient.interceptors.response.use(
    (response) => response,
    (error: AxiosError<ApiErrorResponse>) => {
        const status = error.response?.status;
        const code = error.response?.data?.code;
        const currentPath =
            typeof window !== "undefined" ? window.location.pathname : "";

        if (typeof window === "undefined") {
            return Promise.reject(error);
        }

        if (status === 401) {
            if (currentPath !== "/login") {
                window.location.replace("/login");
            }
            return Promise.reject(error);
        }

        if (status === 403) {
            if (
                code === "ONBOARDING_REQUIRED" &&
                currentPath !== "/onboarding"
            ) {
                window.location.replace("/onboarding");
            } else if (currentPath !== "/403") {
                window.location.replace("/403");
            }
            return Promise.reject(error);
        }

        return Promise.reject(error);
    },
);

export default apiClient;
