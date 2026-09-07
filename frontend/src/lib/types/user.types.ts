import type { RedirectUrlResponse } from "./api.types";

export type RoleType = "member" | "admin";

export type GroupType = "member" | "competitive_team" | "executive" | "director" | "board";

export type User = {
  id: string;
  email: string;
  student_id: string | null;
  role: RoleType;
  created_at: string;
  updated_at: string;
  full_name: string;
  email_verified_at: string | null;
  is_student: boolean;
  onboarding_completed_at: string | null;
  avatar_url: string | null;
  groups: GroupType[];
};

export type SessionResponse = { user?: User };

export type OAuthAuthorizeResponse = RedirectUrlResponse;

export type CompleteOnboardingPayload =
  { is_student: true; student_id: string } | { is_student: false };

export type StudentStatus = "student" | "not_student";
