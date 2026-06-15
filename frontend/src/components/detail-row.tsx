import type { ReactNode } from "react";

export type DetailRowProps = {
  label: string;
  children: ReactNode;
};

export function DetailRow({ label, children }: DetailRowProps) {
  return (
    <div className="grid gap-2 border-t border-brand-border px-5 py-3.5 sm:grid-cols-[130px_minmax(0,1fr)] sm:items-center">
      <dt className="text-sm font-medium text-brand-text-subtle">{label}</dt>
      <dd className="min-w-0 text-sm text-brand-text">{children}</dd>
    </div>
  );
}
