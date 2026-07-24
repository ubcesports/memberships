import {
  ArrowRight,
  ArrowUpRight,
  CalendarDays,
  Clock3,
  CreditCard,
  TicketPercent,
  UsersRound,
} from "lucide-react";
import Link from "next/link";
import { BasePage } from "@/components/layout/base-page";

const memberPerks = [
  {
    title: "Legion Lounge access",
    description: "Daily access to the Legion Gaming Lounge",
    icon: Clock3,
  },
  {
    title: "Event ticket savings",
    description: "Discounted UBCEA raffle & event tickets.",
    icon: TicketPercent,
  },
  {
    title: "Member events",
    description: "Regular access to events created for the UBCEA member community.",
    icon: UsersRound,
  },
  {
    title: "Flexible upgrades",
    description: "Move from Basic to Lounge tier anytime and pay only the difference.",
    icon: ArrowUpRight,
  },
];

const primaryLinkClass =
  "inline-flex h-11 items-center justify-center gap-2 border border-brand-primary bg-brand-primary px-5 text-sm font-semibold text-white transition hover:border-brand-primary-hover hover:bg-brand-primary-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand-text";

export default function HomePage() {
  return (
    <BasePage>
      <div className="flex flex-1 items-center py-10 md:py-16">
        <div className="mx-auto w-full max-w-6xl">
          <section className="relative overflow-hidden border border-brand-border bg-brand-surface/85 shadow-2xl shadow-black/25">
            <div
              aria-hidden="true"
              className="absolute top-0 right-0 h-full w-2/5 border-l border-brand-primary/20 bg-brand-primary/8 [clip-path:polygon(30%_0,100%_0,100%_100%,0_100%)]"
            />

            <div className="relative max-w-3xl px-6 py-12 sm:px-10 sm:py-16 lg:px-14 lg:py-20">
              <h1 className="mt-5 max-w-2xl text-4xl leading-tight font-semibold tracking-[-0.03em] text-brand-text sm:text-5xl">
                Join UBC’s largest gaming community.
              </h1>
              <p className="mt-5 max-w-xl text-base leading-7 text-brand-text-muted sm:text-lg">
                Become a member to join a community of 3,000+ gamers. Gain special perks such as
                Legion Lounge access, discounted event tickets, member events, and more - all
                managed from one account.
              </p>

              <div className="mt-8 flex flex-wrap gap-3">
                <Link href="/pricing" className={primaryLinkClass}>
                  Become a member
                  <ArrowRight aria-hidden="true" className="size-4" />
                </Link>
              </div>
            </div>

            <div className="relative border-t border-brand-border bg-black/15">
              <div className="border-b border-brand-border px-6 py-4 sm:px-8">
                <h2 className="text-xs font-bold tracking-[0.16em] text-brand-text-muted uppercase">
                  Why become a member?
                </h2>
              </div>

              <ul className="grid sm:grid-cols-2 lg:grid-cols-4">
                {memberPerks.map((perk) => {
                  const Icon = perk.icon;

                  return (
                    <li
                      key={perk.title}
                      className="border-b border-brand-border p-6 last:border-b-0 sm:odd:border-r sm:nth-last-[-n+2]:border-b-0 lg:border-r lg:border-b-0 lg:last:border-r-0"
                    >
                      <span className="flex size-9 items-center justify-center border border-brand-primary/40 bg-brand-primary/15 text-brand-text">
                        <Icon aria-hidden="true" className="size-4" />
                      </span>
                      <h3 className="mt-5 text-base font-semibold text-brand-text">{perk.title}</h3>
                      <p className="mt-2 text-sm leading-6 text-brand-text-muted">
                        {perk.description}
                      </p>
                    </li>
                  );
                })}
              </ul>
            </div>
          </section>

          <section aria-label="Membership services" className="mt-6 grid gap-6 lg:grid-cols-2">
            <div className="border border-brand-border bg-brand-surface/75 p-6 sm:p-8">
              <CreditCard aria-hidden="true" className="size-6 text-brand-primary" />
              <h2 className="mt-5 text-xl font-semibold text-brand-text">
                Find the membership that fits you.
              </h2>
              <p className="mt-2 max-w-lg text-sm leading-6 text-brand-text-muted">
                Compare available options, eligibility, and pricing before you join.
              </p>
              <Link
                href="/pricing"
                className="mt-6 inline-flex items-center gap-2 text-sm font-semibold text-brand-text hover:text-brand-text-subtle focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-brand-primary"
              >
                Explore membership pricing
                <ArrowRight aria-hidden="true" className="size-4" />
              </Link>
            </div>

            <div className="relative overflow-hidden border border-brand-border bg-brand-surface/75 p-6 sm:p-8">
              <CalendarDays aria-hidden="true" className="size-6 text-brand-primary" />
              <h2 className="mt-5 text-xl font-semibold text-brand-text">Online lounge booking</h2>
              <p className="mt-2 max-w-lg text-sm leading-6 text-brand-text-muted">
                Check availability and reserve the Legion Lounge online.
              </p>
              <span className="mt-6 inline-flex min-h-6 items-center border border-brand-primary/40 bg-brand-primary/15 px-2 text-xs font-semibold text-brand-text p-1">
                Coming soon!
              </span>
            </div>
          </section>
        </div>
      </div>
    </BasePage>
  );
}
