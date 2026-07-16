export type CompleteOnboardingPayload =
  | { is_student: true; student_id: string }
  | { is_student: false };

export type StudentStatus = "student" | "not_student";
