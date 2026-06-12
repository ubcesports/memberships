"use client";

import { type ReactNode, useEffect, useMemo, useRef, useState } from "react";
import {
  AlertTriangle,
  ExternalLink,
  Loader2,
  LogOut,
  RefreshCw,
  Trash2,
} from "lucide-react";
import { BasePage } from "@/components/layout/base-page";
import apiClient from "@/lib/client";
import { redirectToJasperLabsSignIn } from "@/lib/auth";

type UserDetails = Record<string, unknown>;

type SessionResponse = {
  user?: UserDetails;
};

const JASPERLABS_ACCOUNT_URL =
  process.env.NEXT_PUBLIC_JASPERLABS_ACCOUNT_URL ||
  "https://auth.jasperlabs.net/dashboard";

function stringValue(value: unknown) {
  return typeof value === "string" && value.trim() ? value : null;
}

function booleanValue(value: unknown) {
  return typeof value === "boolean" ? value : false;
}

function formatDate(value: unknown) {
  const raw = stringValue(value);
  if (!raw) {
    return null;
  }

  const date = new Date(raw);
  if (Number.isNaN(date.getTime())) {
    return raw;
  }

  return new Intl.DateTimeFormat("en", {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "numeric",
    minute: "2-digit",
  }).format(date);
}

function formatDateOnly(value: unknown) {
  const raw = stringValue(value);
  if (!raw) {
    return null;
  }

  const date = new Date(raw);
  if (Number.isNaN(date.getTime())) {
    return raw;
  }

  return new Intl.DateTimeFormat("en", {
    month: "short",
    day: "numeric",
    year: "numeric",
  }).format(date);
}

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

type StatusBadgeProps = {
  children: string;
  tone?: "default" | "success" | "warning" | "muted";
};

function StatusBadge({ children, tone = "default" }: StatusBadgeProps) {
  const toneClass = {
    default: "border-brand-primary/40 bg-brand-primary/15 text-brand-text",
    success: "border-green-400/35 bg-green-400/10 text-green-100",
    warning: "border-amber-300/35 bg-amber-300/10 text-amber-100",
    muted: "border-brand-border bg-white/5 text-brand-text-muted",
  }[tone];

  return (
    <span
      className={`inline-flex min-h-7 items-center border px-2.5 text-xs font-semibold ${toneClass}`}
    >
      {children}
    </span>
  );
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
  const [user, setUser] = useState<UserDetails | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSigningOut, setIsSigningOut] = useState(false);
  const [isSyncingAccount, setIsSyncingAccount] = useState(false);
  const [isDeletingAccount, setIsDeletingAccount] = useState(false);
  const [isConfirmingDelete, setIsConfirmingDelete] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const didStartAuth = useRef(false);

  useEffect(() => {
    let isMounted = true;

    async function loadProfile() {
      try {
        const response = await apiClient.get<SessionResponse>("/auth/me", {
          validateStatus: (status) => status === 200 || status === 401,
        });

        if (!isMounted) {
          return;
        }

        if (response.status === 401) {
          if (!didStartAuth.current) {
            didStartAuth.current = true;
            await redirectToJasperLabsSignIn(window.location.href);
          }
          return;
        }

        setUser(response.data.user ?? null);
      } catch {
        if (isMounted) {
          setError("Profile details could not be loaded.");
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    loadProfile();

    return () => {
      isMounted = false;
    };
  }, []);

  const displayName = useMemo(() => {
    if (!user) {
      return "Profile";
    }

    return (
      stringValue(user.full_name) ||
      stringValue(user.name) ||
      stringValue(user.email) ||
      "Profile"
    );
  }, [user]);

  const profile = useMemo(() => {
    if (!user) {
      return null;
    }

    const email = stringValue(user.email) ?? "No email provided";
    const role = stringValue(user.role) ?? "member";
    const isAdmin = role === "admin";
    const isStudent = booleanValue(user.is_student);
    const studentId = stringValue(user.student_id);
    const avatarURL = stringValue(user.avatar_url);
    const emailVerifiedAt = formatDate(user.email_verified_at);
    const onboardingCompletedAt = formatDate(user.onboarding_completed_at);
    const createdAt = formatDateOnly(user.created_at);

    return {
      email,
      role,
      roleLabel: titleCase(role),
      initials: getInitials(displayName, email),
      avatarURL,
      accountType: isAdmin ? "Admin" : "Member",
      accountDetail: isAdmin
        ? "Can manage membership administration."
        : "Standard UBCEA membership account.",
      accountTone: isAdmin ? "warning" : "default",
      studentStatus: isStudent ? "Student" : "Non-student",
      studentDetail: isStudent
        ? "Student pricing and eligibility can apply."
        : "Registered as a community member.",
      studentTone: isStudent ? "success" : "muted",
      studentId,
      emailVerifiedAt,
      onboardingCompletedAt,
      onboardingStatus: onboardingCompletedAt ? "Complete" : "Needs setup",
      onboardingDetail: onboardingCompletedAt
        ? `Completed ${onboardingCompletedAt}.`
        : "Finish setup to unlock the full account flow.",
      onboardingTone: onboardingCompletedAt ? "success" : "warning",
      createdAt,
    } satisfies {
      email: string;
      role: string;
      roleLabel: string;
      initials: string;
      avatarURL: string | null;
      accountType: string;
      accountDetail: string;
      accountTone: StatusBadgeProps["tone"];
      studentStatus: string;
      studentDetail: string;
      studentTone: StatusBadgeProps["tone"];
      studentId: string | null;
      emailVerifiedAt: string | null;
      onboardingCompletedAt: string | null;
      onboardingStatus: string;
      onboardingDetail: string;
      onboardingTone: StatusBadgeProps["tone"];
      createdAt: string | null;
    };
  }, [displayName, user]);

  async function handleSignOut() {
    setIsSigningOut(true);
    setError(null);

    try {
      await apiClient.post(
        "/auth/signout",
        {},
        {
          validateStatus: (status) => status === 204 || status === 401,
        },
      );

      window.location.replace("/");
    } catch {
      setError("Sign out failed. Try again.");
      setIsSigningOut(false);
    }
  }

  async function handleDeleteAccount() {
    setIsDeletingAccount(true);
    setError(null);

    try {
      await apiClient.delete("/auth/users/me", {
        data: {},
        headers: {
          "Content-Type": "application/json",
        },
        validateStatus: (status) => status === 204 || status === 401,
      });

      window.location.replace("/");
    } catch {
      setError("Account deletion failed. Try again.");
      setIsDeletingAccount(false);
    }
  }

  async function handleSyncAccount() {
    setIsSyncingAccount(true);
    setError(null);

    try {
      await redirectToJasperLabsSignIn(window.location.href);
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
              {user ? (
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

            {isLoading ? (
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
                      {profile.avatarURL ? (
                        // eslint-disable-next-line @next/next/no-img-element
                        <img
                          src={profile.avatarURL}
                          alt=""
                          className="size-18 shrink-0 border border-brand-primary/40 bg-brand-primary/15 object-cover"
                        />
                      ) : (
                        <div className="flex size-18 shrink-0 items-center justify-center border border-brand-primary/40 bg-brand-primary/15 text-2xl font-semibold text-brand-text">
                          {profile.initials}
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
                            {profile.roleLabel}
                          </StatusBadge>
                        </div>
                        <p className="mt-2 break-words text-sm text-brand-text-muted">
                          {profile.createdAt
                            ? `Member since ${profile.createdAt}`
                            : "Membership start date unavailable"}
                        </p>
                      </div>
                    </div>

                    <div className="mt-6 grid gap-3.5">
                      <SummaryTile
                        label="Account type"
                        value={profile.accountType}
                        detail={profile.accountDetail}
                        tone={profile.accountTone}
                      />
                      <SummaryTile
                        label="Student status"
                        value={profile.studentStatus}
                        detail={profile.studentDetail}
                        tone={profile.studentTone}
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
                        <StatusBadge tone={profile.onboardingTone}>
                          {profile.onboardingStatus}
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

                <div className="mt-6 border border-brand-border bg-white/[0.03]">
                  <div className="flex flex-col gap-4 px-5 py-5 sm:flex-row sm:items-center sm:justify-between">
                    <div className="flex min-w-0 gap-4">
                      <div className="flex size-10 shrink-0 items-center justify-center border border-red-400/35 bg-red-400/10 text-red-100">
                        <AlertTriangle aria-hidden="true" className="size-5" />
                      </div>
                      <div className="min-w-0">
                        <h3 className="text-base font-semibold text-brand-text">
                          Delete account
                        </h3>
                        <p className="mt-1 max-w-xl text-sm leading-6 text-brand-text-muted">
                          Permanently remove this user account and membership
                          records from this system.
                        </p>
                      </div>
                    </div>

                    <button
                      type="button"
                      onClick={() => setIsConfirmingDelete(true)}
                      disabled={isDeletingAccount}
                      className="inline-flex h-10 cursor-pointer items-center justify-center gap-2 border border-red-400/40 px-4 text-sm font-semibold text-red-100 transition hover:bg-red-400/10 disabled:cursor-not-allowed disabled:opacity-60"
                    >
                      <Trash2 aria-hidden="true" className="size-4" />
                      <span>Delete account</span>
                    </button>
                  </div>

                  {isConfirmingDelete ? (
                    <div className="border-t border-brand-border px-5 py-5">
                      <p className="text-sm leading-6 text-brand-text-muted">
                        This action cannot be undone. Confirm that you want to
                        delete this account.
                      </p>
                      <div className="mt-4 flex flex-col gap-3 sm:flex-row">
                        <button
                          type="button"
                          onClick={handleDeleteAccount}
                          disabled={isDeletingAccount}
                          className="inline-flex h-10 cursor-pointer items-center justify-center gap-2 bg-red-500 px-4 text-sm font-semibold text-white transition hover:bg-red-600 disabled:cursor-not-allowed disabled:opacity-60"
                        >
                          {isDeletingAccount ? (
                            <Loader2
                              aria-hidden="true"
                              className="size-4 animate-spin"
                            />
                          ) : (
                            <Trash2 aria-hidden="true" className="size-4" />
                          )}
                          <span>
                            {isDeletingAccount
                              ? "Deleting account"
                              : "Confirm deletion"}
                          </span>
                        </button>
                        <button
                          type="button"
                          onClick={() => setIsConfirmingDelete(false)}
                          disabled={isDeletingAccount}
                          className="inline-flex h-10 cursor-pointer items-center justify-center border border-brand-border px-4 text-sm font-semibold text-brand-text transition hover:border-brand-text-muted hover:bg-white/5 disabled:cursor-not-allowed disabled:opacity-60"
                        >
                          Cancel
                        </button>
                      </div>
                    </div>
                  ) : null}
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
