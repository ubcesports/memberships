import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
  title: "Complete Profile",
};

export default function OnboardLayout({ children }: { children: ReactNode }) {
  return children;
}
