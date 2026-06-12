import { BasePage } from "@/components/layout/base-page";
import { StatusPage } from "../../components/status-page";

export default function ForbiddenPage() {
    return (
        <BasePage>
            <div className="flex flex-1 items-center py-16">
                <StatusPage
                    code="403"
                    eyebrow="Forbidden"
                    title="Access denied"
                    description="You do not have permission to view this page."
                    primaryAction={{
                        href: "/",
                        label: "Return home",
                    }}
                    secondaryAction={{
                        href: "/login",
                        label: "Switch account",
                    }}
                />
            </div>
        </BasePage>
    );
}
