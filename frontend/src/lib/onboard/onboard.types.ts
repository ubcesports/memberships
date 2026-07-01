export type ApiErrorResponse = {
  code?: string;
  detail?: string;
  message?: string;
};

export type CompleteOnboardingPayload = {
  is_student: boolean;
  student_id: string | null;
};

export type StudentStatus = "student" | "not_student" | null;
