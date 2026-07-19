import { BasePage } from "@/components/layout/base-page";
import { CheckoutResult } from "@/components/membership/checkout-result";

type CheckoutPageProps = {
  searchParams: Promise<{
    success?: string | string[];
  }>;
};

export default async function CheckoutPage({ searchParams }: CheckoutPageProps) {
  const { success } = await searchParams;
  const successful = Array.isArray(success) ? success.includes("true") : success === "true";

  return (
    <BasePage>
      <div className="flex flex-1 items-center py-12 sm:py-16">
        <CheckoutResult successful={successful} />
      </div>
    </BasePage>
  );
}
