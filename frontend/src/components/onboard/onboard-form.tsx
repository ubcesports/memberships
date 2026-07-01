import type { FormEvent, ReactNode } from "react";
import { GraduationCap, Loader2, UserRound } from "lucide-react";
import { ActionButton } from "@/components/action-button";
import type { StudentStatus } from "@/lib/onboard/onboard.types";

type OnboardFormProps = {
  studentStatus: StudentStatus;
  studentId: string;
  validationError: string | null;
  canSubmit: boolean;
  isPending: boolean;
  onStudentStatusChange: (status: Exclude<StudentStatus, null>) => void;
  onStudentIdChange: (value: string) => void;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
};

export function OnboardForm({
  studentStatus,
  studentId,
  validationError,
  canSubmit,
  isPending,
  onStudentStatusChange,
  onStudentIdChange,
  onSubmit,
}: OnboardFormProps) {
  const isStudent = studentStatus === "student";

  return (
    <form className="px-5 py-5 sm:px-6" onSubmit={onSubmit}>
      <fieldset disabled={isPending} className="space-y-5 disabled:opacity-70">
        <div>
          <label className="text-sm font-medium text-brand-text">Are you a student?</label>
          <div className="mt-3 grid gap-3 sm:grid-cols-2">
            <StudentStatusButton
              selected={studentStatus === "student"}
              icon={
                <GraduationCap
                  aria-hidden="true"
                  className="mt-0.5 size-5 shrink-0 text-brand-primary"
                />
              }
              title="Yes"
              description="I have a student number."
              onClick={() => onStudentStatusChange("student")}
            />
            <StudentStatusButton
              selected={studentStatus === "not_student"}
              icon={
                <UserRound
                  aria-hidden="true"
                  className="mt-0.5 size-5 shrink-0 text-brand-primary"
                />
              }
              title="No"
              description="Continue as a community member."
              onClick={() => onStudentStatusChange("not_student")}
            />
          </div>
        </div>

        {isStudent ? (
          <div>
            <label htmlFor="student-id" className="text-sm font-medium text-brand-text">
              Student number
            </label>
            <input
              id="student-id"
              inputMode="numeric"
              pattern="[0-9]{8}"
              autoComplete="off"
              maxLength={8}
              value={studentId}
              onChange={(event) => onStudentIdChange(event.target.value)}
              className="mt-2 h-12 w-full border border-brand-border bg-white/3 px-4 font-mono text-base text-brand-text outline-none transition placeholder:text-brand-text-subtle focus:border-brand-primary disabled:cursor-not-allowed"
              placeholder="12345678"
              aria-describedby={validationError ? "onboard-validation-error" : undefined}
            />
            <p className="mt-2 text-sm leading-6 text-brand-text-subtle">
              Enter the 8 digit number on your student account.
            </p>
          </div>
        ) : null}

        {validationError ? (
          <p id="onboard-validation-error" className="text-sm leading-6 text-brand-warning">
            {validationError}
          </p>
        ) : null}

        <ActionButton
          type="submit"
          className="h-12 w-full border-brand-primary bg-brand-primary text-base hover:border-brand-primary-hover hover:bg-brand-primary-hover"
          loading={isPending}
          disabled={!canSubmit}
          icon={<GraduationCap aria-hidden="true" className="size-5" />}
          loadingIcon={<Loader2 aria-hidden="true" className="size-5 animate-spin" />}
        >
          {isPending ? "Submitting" : "Continue"}
        </ActionButton>
      </fieldset>
    </form>
  );
}

type StudentStatusButtonProps = {
  selected: boolean;
  icon: ReactNode;
  title: string;
  description: string;
  onClick: () => void;
};

function StudentStatusButton({
  selected,
  icon,
  title,
  description,
  onClick,
}: StudentStatusButtonProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`flex min-h-24 cursor-pointer items-start gap-3 border px-4 py-4 text-left transition disabled:cursor-not-allowed ${
        selected
          ? "border-brand-primary bg-brand-primary/15"
          : "border-brand-border bg-white/3 hover:border-brand-text-muted hover:bg-white/5"
      }`}
    >
      {icon}
      <span>
        <span className="block text-sm font-semibold text-brand-text">{title}</span>
        <span className="mt-1 block text-sm leading-5 text-brand-text-muted">{description}</span>
      </span>
    </button>
  );
}
