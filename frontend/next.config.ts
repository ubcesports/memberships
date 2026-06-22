import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "cdn.jasperlabs.net",
        pathname: "/envoy/avatars/**",
      },
    ],
  },
};

export default nextConfig;
