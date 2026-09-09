import type { GroupType, RoleType, User } from "./user.types";
import type { Membership, Transaction } from "./membership.types";

export type SearchMode = "full_name" | "email" | "student_id";

export type IsStudent = "yes" | "no";

export type AppliedSearch = {
  mode: SearchMode;
  value: string;
} | null;

export type AdminUserFilters = {
  role?: RoleType;
  group?: GroupType;
  isStudent?: boolean;
};

export type AuditLogActor = {
  actor_user_id: string;
  actor_full_name: string;
  actor_avatar_url: string;
};

export type AuditLogOutcome = "success" | "failed" | "denied";

export type AuditLogResponse = {
  logs: AuditLogEntry[];
  total: number;
};

export type AuditLogEntry = {
  actor: AuditLogActor;
  occured_at: string;
  action: string;
  description: string | null;
  outcome: AuditLogOutcome;
  request_id: string;
  target_user: AuditLogActor | null;
};

export type UsersResponse = {
  users: User[];
  total: number;
};

export type UserResponse = {
  user: User;
};

/*
  Every field is optional. An absent field is left untouched, so only send the
  ones the admin actually changed.
*/
export type UpdateUserRequest = {
  full_name?: string;
  student_id?: string;
  is_student?: boolean;
  groups_add?: GroupType[];
  groups_remove?: GroupType[];
  role?: RoleType;
  cancel_membership_id?: string;
};

export type AdminPagination = {
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

export const IS_STUDENT_OPTIONS: { value: IsStudent; label: string }[] = [
  { value: "yes", label: "Yes" },
  { value: "no", label: "No" },
];

export type IsStudentFilter = "all" | IsStudent;
