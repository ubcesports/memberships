const GROUP_LABELS: Record<string, string> = {
    member: "Non-student",
    student: "Student",
    competitive_team: "Competitive team",
    executive: "Executive",
    director: "Director",
    board: "Board",
};

export function formatMembershipPrice(amountMinor: number, currency: string) {
    return new Intl.NumberFormat("en-CA", {
        style: "currency",
        currency: currency.toUpperCase(),
        minimumFractionDigits: 0,
        maximumFractionDigits: amountMinor % 100 === 0 ? 0 : 2,
    }).format(amountMinor / 100);
}

export function membershipGroupLabel(group: string) {
    return GROUP_LABELS[group] ?? group.replaceAll("_", " ");
}
