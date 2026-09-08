import type { Metadata } from "next";

export const SITE_URL =
  process.env.NEXT_PUBLIC_SITE_URL?.replace(/\/$/, "") || "https://app.ubcesports.ca";

export const SITE_NAME = "UBC Esports Memberships";

export const SITE_DESCRIPTION =
  "Compare UBC Esports Association membership passes, pricing, and benefits, including Legion Gaming Lounge access and member event perks.";

export const NO_INDEX_ROBOTS: Metadata["robots"] = {
  index: false,
  follow: false,
  googleBot: {
    index: false,
    follow: false,
    noimageindex: true,
  },
};
