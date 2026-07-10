"use client";

import { FormEvent, useMemo, useState } from "react";
import { toast } from "sonner";
import { useMutation } from "@tanstack/react-query";
import { BasePage } from "@/components/layout/base-page";
import { OnboardForm } from "@/components/onboard/onboard-form";
import { completeOnboarding } from "@/lib/onboard/onboard.api";
import type { CompleteOnboardingPayload, StudentStatus } from "@/lib/onboard/onboard.types";

const STUDENT_ID_PATTERN = /^\d{8}$/;

export default function OnboardPage() {
  const [studentStatus, setStudentStatus] = useState<StudentStatus | null>(null);
  const [studentId, setStudentId] = useState("");
  const [validationError, setValidationError] = useState<string | null>(null);

  const isStudent = studentStatus === "student";
  const normalizedStudentId = studentId.trim();

  const canSubmit = useMemo(() => {
    if (!studentStatus) {
      return false;
    }

    if (isStudent) {
      return STUDENT_ID_PATTERN.test(normalizedStudentId);
    }

    return true;
  }, [isStudent, normalizedStudentId, studentStatus]);

  const { mutate: submitOnboarding, isPending } = useMutation({
    mutationFn: async () => {
      const payload: CompleteOnboardingPayload = isStudent
        ? { is_student: true, student_id: normalizedStudentId }
        : { is_student: false };
      const result = await completeOnboarding(payload);
      window.location.replace(result.destination);
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : "Unable to complete onboarding.");
    },
  });

  function handleStudentStatusChange(status: StudentStatus) {
    setStudentStatus(status);
    setValidationError(null);

    if (status === "not_student") {
      setStudentId("");
    }
  }

  function handleStudentIdChange(value: string) {
    setStudentId(value.replace(/\D/g, "").slice(0, 8));
    setValidationError(null);
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!studentStatus) {
      setValidationError("Select whether you are a student.");
      return;
    }

    if (isStudent && !STUDENT_ID_PATTERN.test(normalizedStudentId)) {
      setValidationError("Enter an 8 digit student number.");
      return;
    }

    setValidationError(null);
    submitOnboarding();
  }

  return (
    <BasePage>
      <div className="flex flex-1 items-center justify-center py-12">
        <section className="w-full max-w-xl border border-brand-border bg-brand-surface/85 shadow-2xl shadow-black/25">
          <div className="border-b border-brand-border px-5 py-5 sm:px-6">
            <p className="text-sm font-semibold text-brand-primary">UBCEA Memberships</p>
            <h1 className="mt-3 text-2xl font-semibold text-brand-text">Complete your profile</h1>
            <p className="mt-2 text-sm leading-6 text-brand-text-muted">
              Confirm your student status before continuing.
            </p>
          </div>

          <OnboardForm
            studentStatus={studentStatus}
            studentId={studentId}
            validationError={validationError}
            canSubmit={canSubmit}
            isPending={isPending}
            onStudentStatusChange={handleStudentStatusChange}
            onStudentIdChange={handleStudentIdChange}
            onSubmit={handleSubmit}
          />
        </section>
      </div>
    </BasePage>
  );
}
