import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
  title: "Signing In",
};

export default function OnboardCheckLayout({ children }: { children: ReactNode }) {
  return children;
}
