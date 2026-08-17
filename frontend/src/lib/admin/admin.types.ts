export type RoleType = "member" | "admin";

export type GroupType = "member" | "competitive_team" | "executive" | "director" | "board";

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

export type TransactionStatusType = "pending" | "completed" | "failed" | "refunded" | "expired";

export type MembershipTransaction = {
  id: string;
  amount_paid: string;
  status: TransactionStatusType;
  group_at_purchase: GroupType;
  currency?: string;
  customer_id?: string;
  payment_intent?: string;
  charge_id?: string;
  created_at?: string;
  metadata?: Record<string, unknown> | null;
};

export type Membership = {
  id: string;
  tier_id: string;
  tier_title?: string;
  started_at: string;
  expires_at: string;
  cancelled_at: string | null;
  transaction: MembershipTransaction;
};

/*
  Every field is optional. An absent field is left untouched, so only send the
  ones the admin actually changed.
*/
export type UpdateUserRequest = {
  student_id?: string;
  is_student?: boolean;
  groups_add?: GroupType[];
  groups_remove?: GroupType[];
  role?: RoleType;
  cancel_membership?: boolean;
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
