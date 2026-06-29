import type { ReactNode } from "react";
import { StatusBadge, type StatusBadgeProps } from "@/components/status-badge";

export type SummaryTileProps = {
  label: string;
  value?: string;
  detail: ReactNode;
  tone?: StatusBadgeProps["tone"];
};

export function SummaryTile({ label, value, detail, tone = "default" }: SummaryTileProps) {
  return (
    <div className="min-w-0 border border-brand-border bg-white/[0.03] p-4">
      <div className="grid min-w-0 grid-cols-[minmax(0,1fr)_auto] items-center gap-4">
        <p className="text-sm font-medium leading-5 text-brand-text-subtle">{label}</p>
        {value && <StatusBadge tone={tone}>{value}</StatusBadge>}
      </div>
      <div className="mt-3 text-sm leading-6 text-brand-text-muted">{detail}</div>
    </div>
  );
}
