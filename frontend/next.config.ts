import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
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
