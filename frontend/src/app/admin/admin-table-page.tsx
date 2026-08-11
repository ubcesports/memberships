import { BasePage } from "@/components/layout/base-page";
import { Loader2 } from "lucide-react";

type AdminTablePageProps = {
  title: string;
  description: string;
  isLoading: boolean;
  loadingLabel: string;
  toolbar: React.ReactNode;
  table: React.ReactNode;
  pagination: React.ReactNode;
};

export function AdminTablePage({
  title,
  description,
  isLoading,
  loadingLabel,
  toolbar,
  table,
  pagination,
}: AdminTablePageProps) {
  return (
    <BasePage>
      <div className="flex flex-1 items-center py-6">
        <section className="mx-auto flex min-h-[85vh] w-full max-h-[calc(100vh-3rem)] flex-col">
          {isLoading ? (
            <div className="flex items-center gap-3 text-brand-text-muted">
              <Loader2 aria-hidden="true" className="size-5 animate-spin" />
              <span>{loadingLabel}</span>
            </div>
          ) : (
            <div className="flex min-h-[85vh] flex-1 flex-col border border-brand-border bg-brand-surface/80 shadow-2xl shadow-black/25">
              <div className="shrink-0 border-b border-brand-border px-5 py-5 sm:px-6">
                <h1 className="text-lg font-semibold text-brand-text">{title}</h1>
                <p className="mt-1 text-sm text-brand-text-subtle">{description}</p>
              </div>

              <div className="shrink-0">{toolbar}</div>
              <div className="flex min-h-0 flex-1 flex-col p-5 sm:p-6">{table}</div>
              <div className="shrink-0">{pagination}</div>
            </div>
          )}
        </section>
      </div>
    </BasePage>
  );
}
