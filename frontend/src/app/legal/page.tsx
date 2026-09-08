import type { Metadata } from "next";
import { BasePage } from "@/components/layout/base-page";

const CONTACT_EMAIL = "communications@ubcesports.ca";

export const metadata: Metadata = {
  title: "Privacy Policy and Terms of Use",
  description: "Privacy Policy and Terms of Use for the UBC Esports membership portal.",
  alternates: {
    canonical: "/legal",
  },
};

const headingClassName = "text-lg font-semibold text-brand-text";
const paragraphClassName = "mt-2 text-sm leading-6 text-brand-text-muted";
const listClassName = "mt-3 list-disc space-y-2 pl-5 text-sm leading-6 text-brand-text-muted";
const contactClassName =
  "font-medium text-brand-primary underline underline-offset-4 hover:text-brand-primary-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand-primary";

export default function LegalPage() {
  return (
    <BasePage>
      <div className="mx-auto grid w-full max-w-4xl gap-6 py-10 sm:py-14">
        <article
          id="privacy"
          aria-labelledby="privacy-title"
          className="scroll-mt-28 border border-brand-border bg-brand-surface/80"
        >
          <header className="border-b border-brand-border px-6 py-6 sm:px-8">
            <h1 id="privacy-title" className="mt-2 text-2xl font-semibold text-brand-text">
              Privacy Policy
            </h1>
          </header>

          <div className="grid gap-8 px-6 py-8 sm:px-8">
            <section>
              <h3 className={headingClassName}>Who we are</h3>
              <p className={paragraphClassName}>
                This portal is operated by the UBC Esports Association (&ldquo;UBCEA&rdquo;), a
                student club operating under the Alma Mater Society of UBC Vancouver. UBCEA is
                responsible for the personal information it handles through this portal.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Information we collect</h3>
              <ul className={listClassName}>
                <li>
                  Account information provided through Zetrova, such as your name, email address,
                  account identifier, and profile image.
                </li>
                <li>
                  Information you provide to UBCEA, including student status, student ID, and
                  membership eligibility information.
                </li>
                <li>
                  Membership and payment records, including your tier, groups, dates, amount paid,
                  payment status, purchase type, and Stripe transaction identifiers.
                </li>
                <li>
                  Essential session, request, security, error, and administrative audit information
                  generated while operating the portal.
                </li>
              </ul>
            </section>

            <section>
              <h3 className={headingClassName}>How we use information</h3>
              <p className={paragraphClassName}>
                We use this information to authenticate accounts, determine eligibility and pricing,
                process and administer memberships, communicate with members, prevent misuse,
                support users, investigate errors, maintain records, and comply with legal
                obligations. We do not sell personal information or use it for behavioural
                advertising.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Service providers</h3>
              <p className={paragraphClassName}>
                We use Zetrova for authentication, Stripe for payment processing, Resend for email
                delivery, and OVHcloud in Canada for application and database hosting. These
                providers process information as needed to provide their services and may process
                some information outside Canada. Payment-card details are submitted to Stripe and
                are not stored by this portal.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Cookies and sessions</h3>
              <p className={paragraphClassName}>
                The portal uses cookies or similar storage that are necessary to keep you signed in
                and protect your account. We do not currently use advertising or behavioural
                analytics cookies. Third-party services may use their own necessary cookies when you
                visit their pages.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Retention and security</h3>
              <p className={paragraphClassName}>
                We retain personal information for as long as reasonably necessary to administer
                memberships, maintain financial, security, and audit records, resolve disputes, and
                meet legal obligations. Some records may therefore be kept for an extended period.
                When identifiable information is no longer reasonably required, we will delete or
                anonymize it. We use reasonable safeguards, but no online service can guarantee
                absolute security.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Your choices and contact</h3>
              <p className={paragraphClassName}>
                You may ask to access or correct your personal information, withdraw consent where
                applicable, request deletion subject to necessary retention, or raise a privacy
                concern. Contact the UBCEA Communications Team at{" "}
                <a className={contactClassName} href={`mailto:${CONTACT_EMAIL}`}>
                  {CONTACT_EMAIL}
                </a>
                . We may need to verify your identity before completing a request.
              </p>
            </section>
          </div>
        </article>

        <article
          id="terms"
          aria-labelledby="terms-title"
          className="scroll-mt-28 border border-brand-border bg-brand-surface/80"
        >
          <header className="border-b border-brand-border px-6 py-6 sm:px-8">
            <h2 id="terms-title" className="mt-2 text-2xl font-semibold text-brand-text">
              Terms of Use
            </h2>
          </header>

          <div className="grid gap-8 px-6 py-8 sm:px-8">
            <section>
              <h3 className={headingClassName}>Agreement and eligibility</h3>
              <p className={paragraphClassName}>
                By using the portal or purchasing a membership, you agree to these terms. You must
                be at least 13 years old. If you are under 19, you must have permission from a
                parent or legal guardian before purchasing a membership.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Accounts</h3>
              <p className={paragraphClassName}>
                You must provide accurate information, keep your account secure, and use only your
                own account. Memberships and accounts are personal and may not be sold, transferred,
                or shared. UBCEA may correct eligibility or membership errors and may suspend access
                for fraud, abuse, chargebacks, or violations of UBCEA or AMS rules.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Memberships and payment</h3>
              <p className={paragraphClassName}>
                Membership benefits, eligibility, start dates, and expiry dates are described on the
                pricing and checkout pages. Unless clearly stated otherwise, prices are in Canadian
                dollars and payments are one-time charges processed by Stripe. Upgrade prices may
                account for the amount paid for an eligible current membership.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Refunds</h3>
              <p className={paragraphClassName}>
                Membership purchases are final and non-refundable except where a duplicate charge,
                technical malfunction, or other system error affected the purchase, or where a
                refund is required by law. Send requests to{" "}
                <a className={contactClassName} href={`mailto:${CONTACT_EMAIL}`}>
                  {CONTACT_EMAIL}
                </a>{" "}
                with reasonable supporting information, such as your account email, transaction
                details, screenshots, and any available request ID. UBCEA will review each request
                based on the available records and circumstances.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Availability and changes</h3>
              <p className={paragraphClassName}>
                We may maintain, change, suspend, or discontinue parts of the portal or membership
                offerings. We will not knowingly remove a paid benefit without a reasonable
                operational reason. The portal is provided on an &ldquo;as available&rdquo; basis,
                and interruptions or errors may occur.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Responsibility</h3>
              <p className={paragraphClassName}>
                To the extent permitted by law, UBCEA is not responsible for indirect or incidental
                losses caused by service interruptions, third-party services, or unauthorized
                account use. Nothing in these terms excludes rights or liabilities that cannot
                legally be excluded.
              </p>
            </section>

            <section>
              <h3 className={headingClassName}>Changes and governing law</h3>
              <p className={paragraphClassName}>
                We may update these terms and will post the revised date on this page. Material
                changes will apply prospectively where required. These terms are governed by the
                laws of British Columbia and Canada. Questions can be sent to{" "}
                <a className={contactClassName} href={`mailto:${CONTACT_EMAIL}`}>
                  {CONTACT_EMAIL}
                </a>
                .
              </p>
            </section>
          </div>
        </article>
      </div>
    </BasePage>
  );
}
