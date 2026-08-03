import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "cdn.zetrova.com",
        pathname: "/envoy/avatars/**",
      },
    ],
  },
};

export default nextConfig;
