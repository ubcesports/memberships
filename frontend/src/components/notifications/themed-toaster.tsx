"use client";

import { Toaster } from "sonner";

export function ThemedToaster() {
    return (
        <Toaster
            position="top-right"
            closeButton
            toastOptions={{
                style: {
                    background: "var(--color-surface)",
                    border: "1px solid var(--color-border)",
                    color: "var(--color-text)",
                },
                classNames: {
                    description: "text-brand-text-muted",
                    actionButton:
                        "bg-brand-primary text-white hover:bg-brand-primary-hover",
                    cancelButton:
                        "border border-brand-border bg-transparent text-brand-text",
                },
            }}
        />
    );
}
