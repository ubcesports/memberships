import { AlertTriangle, Loader2, RefreshCw } from "lucide-react";
import { ActionButton } from "@/components/action-button";
import { SurfacePanel } from "@/components/surface-panel";

type MembershipLoadErrorProps = {
  isRetrying: boolean;
  onRetry: () => void;
};

export function MembershipLoadError({ isRetrying, onRetry }: MembershipLoadErrorProps) {
  return (
    <SurfacePanel
      className="flex flex-col gap-4 border-red-400/35 bg-red-400/10 px-5 py-5 sm:flex-row sm:items-center sm:justify-between"
      role="alert"
    >
      <div className="flex items-start gap-3">
        <AlertTriangle aria-hidden="true" className="mt-0.5 size-5 shrink-0 text-red-200" />
        <div>
          <h3 className="text-sm font-semibold text-red-100">Unable to load memberships</h3>
          <p className="mt-1 text-sm leading-6 text-brand-text-muted">
            We could not confirm the membership records. Try loading them again.
          </p>
        </div>
      </div>
      <ActionButton
        className="shrink-0 border-red-300/40 text-red-100 hover:border-red-200/60 hover:bg-red-300/10"
        onClick={onRetry}
        loading={isRetrying}
        icon={<RefreshCw aria-hidden="true" className="size-4" />}
        loadingIcon={<Loader2 aria-hidden="true" className="size-4 animate-spin" />}
      >
        {isRetrying ? "Retrying" : "Try again"}
      </ActionButton>
    </SurfacePanel>
  );
}
