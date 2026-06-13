export type StatusBadgeProps = {
  children: string;
  tone?: "default" | "success" | "warning" | "muted";
};

const TONE_STYLES = {
  default: "border-brand-primary/40 bg-brand-primary/15 text-brand-text",
  success: "border-green-400/35 bg-green-400/10 text-green-100",
  warning: "border-amber-300/35 bg-amber-300/10 text-amber-100",
  muted: "border-brand-border bg-white/5 text-brand-text-muted",
} as const;

export function StatusBadge({ children, tone = "default" }: StatusBadgeProps) {
  const toneClass = TONE_STYLES[tone];

  return (
    <span
      className={`inline-flex min-h-6 items-center border px-2 text-xs font-semibold ${toneClass}`}
    >
      {children}
    </span>
  );
}
