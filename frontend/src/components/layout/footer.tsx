import Image from "next/image";
import Link from "next/link";

const primaryLinks = [
  { href: "/", label: "Home" },
  { href: "/pricing", label: "Pricing" },
];

const legalLinks = [
  { href: "/privacy", label: "Privacy Policy" },
  { href: "/terms", label: "Terms" },
  { href: "/contact", label: "Contact" },
];

const linkClassName =
  "text-sm text-brand-text-muted transition hover:text-brand-text focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-brand-primary";

export function Footer() {
  return (
    <footer className="border-t border-brand-border bg-brand-surface">
      <div className="mx-auto w-full max-w-7xl px-5 py-8 sm:px-8 lg:px-10">
        <div className="flex flex-col gap-6 pb-8 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-4">
            <Image
              src="/ubcea_logo.jpg"
              alt=""
              width={44}
              height={44}
              className="size-11 shrink-0 border border-brand-border object-cover"
            />
            <div>
              <p className="text-base font-bold tracking-[0.12em] text-brand-text">
                UBC Esports Portal
              </p>
              <p className="mt-1 text-sm text-brand-text-subtle">app.ubcesports.ca</p>
            </div>
          </div>

          <nav aria-label="Footer navigation">
            <ul className="flex items-center gap-6">
              {primaryLinks.map((link) => (
                <li key={link.href}>
                  <Link href={link.href} className={linkClassName}>
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </nav>
        </div>

        <div className="flex flex-col-reverse gap-4 border-t border-brand-border/70 pt-5 sm:flex-row sm:items-center sm:justify-between">
          <p className="text-xs text-brand-text-subtle">
            © {new Date().getFullYear()} UBC Esports Association. All rights reserved.
          </p>

          <nav aria-label="Legal links">
            <ul className="flex flex-wrap items-center gap-x-5 gap-y-2">
              {legalLinks.map((link) => (
                <li key={link.href}>
                  <Link href={link.href} className={linkClassName}>
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </nav>
        </div>
      </div>
    </footer>
  );
}
