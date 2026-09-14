import { ChevronDown } from "lucide-react";

type SummaryDropdownProps = {
  summary: string;
  children: React.ReactNode;
};

export function SummaryDropdown({ summary, children }: SummaryDropdownProps) {
  return (
    <details className="group relative inline-block" onClick={(event) => event.stopPropagation()}>
      <summary className="flex min-h-8 cursor-pointer list-none items-center gap-2 whitespace-nowrap border border-brand-border bg-white/3 px-2.5 text-xs font-semibold text-brand-text transition-colors hover:bg-white/[0.07] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand-primary [&::-webkit-details-marker]:hidden">
        <span className="whitespace-nowrap">{summary}</span>
        <ChevronDown
          aria-hidden="true"
          className="size-3.5 transition-transform group-open:rotate-180"
        />
      </summary>

      <div className="absolute left-0 top-full z-30 mt-1 min-w-max border border-brand-border bg-brand-surface p-2 shadow-xl shadow-black/30">
        <div className="flex max-w-72 flex-col items-start gap-1.5">{children}</div>
      </div>
    </details>
  );
}
