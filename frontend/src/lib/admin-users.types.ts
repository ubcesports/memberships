export type RoleType = "member" | "admin";

export type GroupType = "member" | "competitive_team" | "executive" | "director" | "board";

export type SearchMode = "full_name" | "email" | "student_id";

export type AppliedSearch = {
  mode: SearchMode;
  value: string;
} | null;

export type AdminUserFilters = {
  role?: RoleType;
  group?: GroupType;
  isStudent?: boolean;
};

export type AdminUser = {
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

export type AdminUsersResponse = {
  users: AdminUser[];
  total: number;
};

export type AdminUserPagination = {
  limit: number;
  offset: number;
};

export const PAGE_SIZE_OPTIONS = [10, 25, 50, 100] as const;

export const DEFAULT_PAGE_SIZE = 25;

export const GROUP_OPTIONS: { value: GroupType; label: string }[] = [
  { value: "member", label: "Member" },
  { value: "competitive_team", label: "Competitive Team" },
  { value: "executive", label: "Executive" },
  { value: "director", label: "Director" },
  { value: "board", label: "Board" },
];

export const ROLE_OPTIONS: { value: RoleType; label: string }[] = [
  { value: "member", label: "Member" },
  { value: "admin", label: "Admin" },
];

export const SEARCH_MODE_OPTIONS: { value: SearchMode; label: string }[] = [
  { value: "full_name", label: "Full name" },
  { value: "email", label: "Email" },
  { value: "student_id", label: "Student ID" },
];
