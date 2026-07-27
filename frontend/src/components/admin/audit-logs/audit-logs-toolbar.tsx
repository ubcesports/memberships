import { ActionButton } from "@/components/action-button";
import { RotateCcw } from "lucide-react";

const FIELD_CLASS_NAME =
  "h-10 border border-brand-border bg-brand-surface px-3 text-sm text-brand-text";

type AuditLogsToolbarProps = {
  searchInput: string;
  total: number;
  onSearchInputChange: (value: string) => void;
  onResetSearch: () => void;
};

export function AuditLogsToolbar({
  searchInput,
  onSearchInputChange,
  onResetSearch,
}: AuditLogsToolbarProps) {
  return (
    <div className="flex flex-col gap-5 border-b border-brand-border px-5 py-5 sm:px-6">
      <div className="flex flex-col gap-3 lg:flex-row lg:items-end">
        <label className="flex min-w-40 flex-col gap-1.5 text-sm text-brand-text-subtle">
          <span>Search by actor name</span>
          <input
            type="text"
            placeholder="Enter actor name"
            className={FIELD_CLASS_NAME}
            value={searchInput}
            onChange={(e) => onSearchInputChange(e.target.value)}
          />
        </label>
        <div className="flex flex-wrap gap-2">
          <ActionButton
            onClick={onResetSearch}
            disabled={searchInput === ""}
            icon={<RotateCcw aria-hidden="true" className="size-4" />}
          >
            Reset search
          </ActionButton>
        </div>
      </div>
    </div>
  );
}
