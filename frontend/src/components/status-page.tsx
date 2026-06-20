import Link from "next/link";

type StatusPageProps = {
  code: string;
  eyebrow: string;
  title: string;
  description: string;
  primaryAction: {
    href: string;
    label: string;
  };
  secondaryAction?: {
    href: string;
    label: string;
  };
};

export function StatusPage({
  code,
  eyebrow,
  title,
  description,
  primaryAction,
  secondaryAction,
}: StatusPageProps) {
  return (
    <section className="mx-auto grid w-full max-w-6xl gap-12 lg:grid-cols-[minmax(0,0.9fr)_minmax(320px,0.55fr)] lg:items-center">
      <div className="max-w-2xl">
        <p className="mb-5 text-sm font-semibold uppercase tracking-[0.22em] text-brand-text-muted">
          {eyebrow}
        </p>
        <h1 className="text-5xl font-semibold leading-[1.02] text-brand-text sm:text-7xl">
          {title}
        </h1>
        <p className="mt-6 max-w-xl text-lg leading-8 text-brand-text-muted">{description}</p>
        <div className="mt-10 flex flex-col gap-3 sm:flex-row">
          <Link
            href={primaryAction.href}
            className="inline-flex h-12 items-center justify-center bg-brand-primary px-6 text-sm font-semibold text-white transition hover:bg-brand-primary-hover"
          >
            {primaryAction.label}
          </Link>
          {secondaryAction ? (
            <Link
              href={secondaryAction.href}
              className="inline-flex h-12 items-center justify-center border border-brand-border px-6 text-sm font-semibold text-brand-text transition hover:border-brand-text-muted hover:bg-white/5"
            >
              {secondaryAction.label}
            </Link>
          ) : null}
        </div>
      </div>

      <div className="relative min-h-72 border border-brand-border bg-brand-surface/80 p-6 shadow-2xl shadow-black/30">
        <div className="flex items-center justify-between border-b border-brand-border pb-5">
          <span className="font-mono text-sm text-brand-text-muted">status</span>
          <span className="font-mono text-sm text-brand-text-subtle">
            memberships.ubcesports.ca
          </span>
        </div>
        <div className="flex min-h-52 items-center justify-center">
          <p className="font-mono text-[7rem] font-semibold leading-none text-brand-text sm:text-[9rem]">
            {code}
          </p>
        </div>
        <div className="grid grid-cols-3 gap-3">
          <div className="h-2 bg-brand-primary" />
          <div className="h-2 bg-brand-primary-hover" />
          <div className="h-2 bg-brand-text-muted" />
        </div>
      </div>
    </section>
  );
}
