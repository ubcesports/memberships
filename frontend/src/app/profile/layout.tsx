import type { Metadata } from "next";
import type { ReactNode } from "react";
import { NO_INDEX_ROBOTS } from "@/lib/site";

export const metadata: Metadata = {
  title: "Profile",
  robots: NO_INDEX_ROBOTS,
};

export default function ProfileLayout({ children }: { children: ReactNode }) {
  return children;
}
