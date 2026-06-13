"use client";

import { type ReactNode, useState } from "react";
import {
  ExternalLink,
  Loader2,
  LogOut,
  RefreshCw,
} from "lucide-react";
import { BasePage } from "@/components/layout/base-page";
import apiClient from "@/lib/client";
import { redirectToSignIn } from "@/lib/auth";
import { StatusBadge, StatusBadgeProps } from "@/components/status-badge";
import { formatDate } from "@/lib/utils/formatting";
import { useProfile } from "@/lib/profile.hook";

const JASPERLABS_ACCOUNT_URL =
  process.env.NEXT_PUBLIC_JASPERLABS_ACCOUNT_URL ||
  "https://auth.jasperlabs.net/dashboard";

function titleCase(value: string) {
  return value
    .replace(/_/g, " ")
    .replace(/\b\w/g, (letter) => letter.toUpperCase());
}

function getInitials(name: string, email: string) {
  const source = name !== "Profile" ? name : email;
  const parts = source
    .split(/[\s@.]+/)
    .map((part) => part.trim())
    .filter(Boolean);

  if (parts.length === 0) {
    return "UB";
  }

  return parts
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");
}

type SummaryTileProps = {
  label: string;
  value: string;
  detail: string;
  tone?: StatusBadgeProps["tone"];
};

function SummaryTile({
  label,
  value,
  detail,
  tone = "default",
}: SummaryTileProps) {
  return (
    <div className="min-w-0 border border-brand-border bg-white/[0.03] p-4">
      <div className="grid min-w-0 grid-cols-[minmax(0,1fr)_auto] items-center gap-4">
        <p className="text-sm font-medium leading-5 text-brand-text-subtle">
          {label}
        </p>
        <StatusBadge tone={tone}>{value}</StatusBadge>
      </div>
      <p className="mt-3 text-sm leading-6 text-brand-text-muted">{detail}</p>
    </div>
  );
}

type DetailRowProps = {
  label: string;
  children: ReactNode;
};

function DetailRow({ label, children }: DetailRowProps) {
  return (
    <div className="grid gap-2 border-t border-brand-border px-5 py-3.5 sm:grid-cols-[130px_minmax(0,1fr)] sm:items-center">
      <dt className="text-sm font-medium text-brand-text-subtle">{label}</dt>
      <dd className="min-w-0 text-sm text-brand-text">{children}</dd>
    </div>
  );
}

export default function ProfilePage() {
  const [isSigningOut, setIsSigningOut] = useState(false);
  const [isSyncingAccount, setIsSyncingAccount] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const { data: profile, isPending } = useProfile();

  const displayName = profile?.name ?? profile?.email ?? "Profile";

  async function handleSignOut() {
    setIsSigningOut(true);
    setError(null);

    try {
      await apiClient.post("/auth/signout", {});
      window.location.replace("/");
    } catch {
      setError("Sign out failed. Try again.");
      setIsSigningOut(false);
    }
  }

  async function handleSyncAccount() {
    setIsSyncingAccount(true);
    setError(null);

    try {
      await redirectToSignIn(window.location.href);
    } catch {
      setError("Unable to start account sync. Try again.");
      setIsSyncingAccount(false);
    }
  }

  return (
    <BasePage>
      <div className="flex flex-1 items-center py-12">
        <section className="mx-auto w-full max-w-4xl">
          <div className="mt-10 border border-brand-border bg-brand-surface/80 shadow-2xl shadow-black/25">
            <div className="flex flex-col gap-4 border-b border-brand-border px-5 py-5 sm:flex-row sm:items-center sm:justify-between sm:px-6">
              <div>
                <h2 className="text-lg font-semibold text-brand-text">
                  UBCEA Membership
                </h2>
                <p className="mt-1 text-sm text-brand-text-subtle">
                  Your profile and account status.
                </p>
              </div>
              {profile ? (
                <button
                  type="button"
                  onClick={handleSignOut}
                  disabled={isSigningOut}
                  className="inline-flex h-10 cursor-pointer items-center justify-center gap-2 border border-brand-border px-4 text-sm font-semibold text-brand-text transition hover:border-brand-text-muted hover:bg-white/5 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {isSigningOut ? (
                    <Loader2
                      aria-hidden="true"
                      className="size-4 animate-spin"
                    />
                  ) : (
                    <LogOut aria-hidden="true" className="size-4" />
                  )}
                  <span>{isSigningOut ? "Signing out" : "Sign out"}</span>
                </button>
              ) : null}
            </div>

            {isPending ? (
              <div className="flex min-h-56 items-center justify-center gap-3 px-6 py-12 text-brand-text-muted">
                <Loader2 aria-hidden="true" className="size-5 animate-spin" />
                <span>Loading profile</span>
              </div>
            ) : error ? (
              <div className="px-6 py-12 text-brand-text-muted">{error}</div>
            ) : profile ? (
              <div className="p-5 sm:p-6">
                <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_minmax(360px,0.78fr)]">
                  <div className="border border-brand-border bg-white/[0.03] p-5 sm:p-6">
                    <div className="flex flex-col gap-5 sm:flex-row sm:items-center">
                      {profile.avatarUrl ? (
                        <img
                          src={profile.avatarUrl}
                          alt=""
                          className="size-18 shrink-0 border border-brand-primary/40 bg-brand-primary/15 object-cover"
                        />
                      ) : (
                        <div className="flex size-18 shrink-0 items-center justify-center border border-brand-primary/40 bg-brand-primary/15 text-2xl font-semibold text-brand-text">
                          {getInitials(profile.name, profile.email)}
                        </div>
                      )}
                      <div className="min-w-0">
                        <div className="flex flex-wrap items-center gap-3">
                          <h3 className="break-words text-2xl font-semibold text-brand-text">
                            {displayName}
                          </h3>
                          <StatusBadge
                            tone={
                              profile.role === "admin" ? "warning" : "default"
                            }
                          >
                            {titleCase(profile.role)}
                          </StatusBadge>
                        </div>
                        <p className="mt-2 break-words text-sm text-brand-text-muted">
                          {profile.createdAt
                            ? `Member since ${formatDate(profile.createdAt)}`
                            : "Membership start date unavailable"}
                        </p>
                      </div>
                    </div>

                    <div className="mt-6 grid gap-3.5">
                      <SummaryTile
                        label="Account type"
                        value={titleCase(profile.role)}
                        detail="Standard UBCEA membership account."
                        tone="default"
                      />
                      <SummaryTile
                        label="Student status"
                        value={profile.isStudent ? "Student" : "Non-student"}
                        detail={
                          profile.isStudent
                            ? "Student pricing and eligibility can apply."
                            : "Registered as a community member."
                        }
                        tone={profile.isStudent ? "success" : "muted"}
                      />
                    </div>
                  </div>

                  <div className="border border-brand-border bg-white/[0.03]">
                    <div className="px-5 py-4">
                      <h3 className="text-base font-semibold text-brand-text">
                        Account Details
                      </h3>
                    </div>
                    <dl>
                      <DetailRow label="Email">
                        <span className="break-words">{profile.email}</span>
                      </DetailRow>
                      <DetailRow label="Email verified">
                        {profile.emailVerifiedAt ? (
                          <StatusBadge tone="success">Verified</StatusBadge>
                        ) : (
                          <StatusBadge tone="muted">Not verified</StatusBadge>
                        )}
                      </DetailRow>
                      <DetailRow label="Student ID">
                        {profile.studentId ? (
                          <span className="break-words font-mono">
                            {profile.studentId}
                          </span>
                        ) : (
                          <StatusBadge tone="muted">Not provided</StatusBadge>
                        )}
                      </DetailRow>
                      <DetailRow label="Onboarding">
                        <StatusBadge
                          tone={
                            profile.onboardingCompletedAt
                              ? "success"
                              : "warning"
                          }
                        >
                          {profile.onboardingCompletedAt
                            ? `Completed ${formatDate(profile.onboardingCompletedAt)}`
                            : "Pending"}
                        </StatusBadge>
                      </DetailRow>
                    </dl>
                    <div className="grid gap-3 border-t border-brand-border px-5 py-4 sm:grid-cols-2">
                      <a
                        href={JASPERLABS_ACCOUNT_URL}
                        className="inline-flex h-10 items-center justify-center gap-2 border border-brand-border px-4 text-sm font-semibold text-brand-text transition hover:border-brand-text-muted hover:bg-white/5"
                      >
                        <ExternalLink aria-hidden="true" className="size-4" />
                        <span>Manage</span>
                      </a>
                      <button
                        type="button"
                        onClick={handleSyncAccount}
                        disabled={isSyncingAccount}
                        className="inline-flex h-10 cursor-pointer items-center justify-center gap-2 border border-brand-border px-4 text-sm font-semibold text-brand-text transition hover:border-brand-text-muted hover:bg-white/5 disabled:cursor-not-allowed disabled:opacity-60"
                      >
                        {isSyncingAccount ? (
                          <Loader2
                            aria-hidden="true"
                            className="size-4 animate-spin"
                          />
                        ) : (
                          <RefreshCw aria-hidden="true" className="size-4" />
                        )}
                        <span>{isSyncingAccount ? "Syncing" : "Sync"}</span>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ) : (
              <div className="px-6 py-12 text-brand-text-muted">
                No profile details were returned.
              </div>
            )}
          </div>
        </section>
      </div>
    </BasePage>
  );
}
