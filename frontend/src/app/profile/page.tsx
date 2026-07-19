"use client";

import { ExternalLink, Loader2, LogOut, RefreshCw } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import { ActionButton } from "@/components/action-button";
import { ActionLink } from "@/components/action-link";
import { DetailRow } from "@/components/detail-row";
import { BasePage } from "@/components/layout/base-page";
import { StatusBadge } from "@/components/status-badge";
import { SummaryTile } from "@/components/summary-tile";
import { SurfacePanel } from "@/components/surface-panel";
import { redirectToSignIn } from "@/lib/auth";
import { useSignOut } from "@/lib/use-sign-out.hook";
import { formatDate, getInitials } from "@/lib/utils/formatting";
import { getGroupBadgeClass, titleCase } from "@/lib/utils/groups";
import { useProfile } from "@/lib/profile.hook";
import Image from "next/image";

const JASPERLABS_ACCOUNT_URL =
  process.env.NEXT_PUBLIC_JASPERLABS_ACCOUNT_URL || "https://auth.jasperlabs.net/dashboard";

export default function ProfilePage() {
  const { data: profile, isPending } = useProfile();

  const { mutate: signOut, error: signOutError, isPending: signOutPending } = useSignOut();

  const {
    mutate: syncAccount,
    error: syncAccountError,
    isPending: syncAccountPending,
  } = useMutation({
    mutationFn: async () => await redirectToSignIn(window.location.href),
  });

  const error = signOutError
    ? "Sign out failed. Try again."
    : syncAccountError
      ? "Unable to sync account. Try again."
      : null;

  const displayName = profile?.name ?? profile?.email ?? "Profile";

  return (
    <BasePage>
      <div className="flex flex-1 items-center py-12">
        <section className="mx-auto w-full max-w-4xl">
          <div className="mt-10 border border-brand-border bg-brand-surface/80 shadow-2xl shadow-black/25">
            <div className="flex flex-col gap-4 border-b border-brand-border px-5 py-5 sm:flex-row sm:items-center sm:justify-between sm:px-6">
              <div>
                <h2 className="text-lg font-semibold text-brand-text">UBCEA Membership</h2>
                <p className="mt-1 text-sm text-brand-text-subtle">
                  Your profile and account status.
                </p>
              </div>
              {profile ? (
                <ActionButton
                  onClick={() => signOut()}
                  loading={signOutPending}
                  icon={<LogOut aria-hidden="true" className="size-4" />}
                  loadingIcon={<Loader2 aria-hidden="true" className="size-4 animate-spin" />}
                >
                  {signOutPending ? "Signing out" : "Sign out"}
                </ActionButton>
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
                  <SurfacePanel className="p-5 sm:p-6">
                    <div className="flex flex-col gap-5 sm:flex-row sm:items-center">
                      {profile.avatarUrl ? (
                        <Image
                          src={profile.avatarUrl}
                          alt=""
                          width={72}
                          height={72}
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
                          <StatusBadge tone={profile.role === "admin" ? "warning" : "default"}>
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
                        label="Assigned groups"
                        detail={
                          profile.groups.length > 0 ? (
                            <div className="flex flex-wrap gap-2">
                              {profile.groups.map((group) => (
                                <StatusBadge
                                  key={group}
                                  tone="default"
                                  className={getGroupBadgeClass(group)}
                                >
                                  {titleCase(group)}
                                </StatusBadge>
                              ))}
                            </div>
                          ) : (
                            "No groups assigned"
                          )
                        }
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
                  </SurfacePanel>

                  <SurfacePanel>
                    <div className="px-5 py-4">
                      <h3 className="text-base font-semibold text-brand-text">Account Details</h3>
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
                          <span className="break-words font-mono">{profile.studentId}</span>
                        ) : (
                          <StatusBadge tone="muted">Not provided</StatusBadge>
                        )}
                      </DetailRow>
                      <DetailRow label="Onboarding">
                        <StatusBadge tone={profile.onboardingCompletedAt ? "success" : "warning"}>
                          {profile.onboardingCompletedAt
                            ? `Completed ${formatDate(profile.onboardingCompletedAt)}`
                            : "Pending"}
                        </StatusBadge>
                      </DetailRow>
                    </dl>
                    <div className="grid gap-3 border-t border-brand-border px-5 py-4 sm:grid-cols-2">
                      <ActionLink
                        href={JASPERLABS_ACCOUNT_URL}
                        icon={<ExternalLink aria-hidden="true" className="size-4" />}
                      >
                        Manage
                      </ActionLink>
                      <ActionButton
                        onClick={() => syncAccount()}
                        loading={syncAccountPending}
                        icon={<RefreshCw aria-hidden="true" className="size-4" />}
                        loadingIcon={<Loader2 aria-hidden="true" className="size-4 animate-spin" />}
                      >
                        {syncAccountPending ? "Syncing" : "Sync"}
                      </ActionButton>
                    </div>
                  </SurfacePanel>
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
