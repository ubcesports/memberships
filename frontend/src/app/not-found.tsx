import { BasePage } from "@/components/layout/base-page";
import { StatusPage } from "../components/status-page";

export default function NotFound() {
    return (
        <BasePage>
            <div className="flex flex-1 items-center py-16">
                <StatusPage
                    code="404"
                    eyebrow="Page not found"
                    title="Page not found"
                    description="The page you are looking for does not exist."
                    primaryAction={{
                        href: "/",
                        label: "Return home",
                    }}
                />
            </div>
        </BasePage>
    );
}
