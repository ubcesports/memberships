export const GROUP_BADGE_STYLES = {
  member: "border-sky-300/35 bg-sky-400/12 text-sky-100",
  competitive_team: "border-fuchsia-300/35 bg-fuchsia-400/12 text-fuchsia-100",
  executive: "border-emerald-300/35 bg-emerald-400/12 text-emerald-100",
  director: "border-amber-300/35 bg-amber-400/12 text-amber-100",
  board: "border-rose-300/35 bg-rose-400/12 text-rose-100",
} as const;

export function titleCase(value: string) {
  return value
    .replace(/_/g, " ")
    .replace(/\b\w/g, (letter) => letter.toUpperCase());
}

export function getGroupBadgeClass(group: string) {
  return (
    GROUP_BADGE_STYLES[group as keyof typeof GROUP_BADGE_STYLES] ||
    "border-brand-border bg-white/5 text-brand-text-muted"
  );
}
