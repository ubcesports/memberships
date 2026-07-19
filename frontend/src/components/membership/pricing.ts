import type {
  EligibleMembershipTier,
  MembershipTier,
  MembershipTierPrice,
} from "@/lib/membership.hook";

const STUDENT_LABELS: Record<string, string> = {
  true: "Student",
  false: "Community",
  null: "Assigned",
};

export function formatMembershipPrice(price: number) {
  return new Intl.NumberFormat("en-CA", {
    style: "currency",
    currency: "CAD",
    minimumFractionDigits: 0,
    maximumFractionDigits: Math.round(price * 100) % 100 === 0 ? 0 : 2,
  }).format(price);
}

export function membershipPriceLabel(price: MembershipTierPrice) {
  return STUDENT_LABELS[String(price.is_student_required)] ?? "Member";
}

export function getPriceByStudentStatus(tier: MembershipTier, isStudent: boolean) {
  return tier.prices.find((price) => price.is_student_required === isStudent);
}

export function getFallbackPrice(tier: MembershipTier) {
  return getPriceByStudentStatus(tier, true) ?? tier.prices[0];
}

export function purchaseLabel(tier: EligibleMembershipTier) {
  if (tier.purchase_type === "upgrade") {
    return "Upgrade to Premium";
  }

  return "Choose this pass";
}

export function isMembershipTierPrice(
  price: MembershipTierPrice | undefined,
): price is MembershipTierPrice {
  return Boolean(price);
}
