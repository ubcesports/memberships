import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { ThemedToaster } from "@/components/notifications/themed-toaster";
import { AppProviders } from "./providers";

import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "UBCEA Memberships",
  description: "Compare and purchase UBCEA membership passes.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col">
        <AppProviders>{children}</AppProviders>
        <ThemedToaster />
      </body>
    </html>
  );
}
