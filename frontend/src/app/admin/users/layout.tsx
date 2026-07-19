import type { Metadata } from "next";
import type { ReactNode } from "react";

export const metadata: Metadata = {
  title: "Users",
};

export default function AdminUsersLayout({ children }: { children: ReactNode }) {
  return children;
}
