import { ChevronLeft, ChevronRight } from "lucide-react";
import { ActionButton } from "@/components/action-button";
import { PAGE_SIZE_OPTIONS } from "@/lib/admin-users.types";

const FIELD_CLASS_NAME =
  "h-10 border border-brand-border bg-brand-surface px-3 text-sm text-brand-text";

type UsersPaginationProps = {
  offset: number;
  limit: number;
  total: number;
  usersCount: number;
  onOffsetChange: (offset: number) => void;
  onLimitChange: (limit: number) => void;
};

function getRangeLabel(
  offset: number,
  usersCount: number,
  total: number,
): string {
  if (total === 0) {
    return "0 users";
  }

  if (usersCount === 0) {
    return `0 of ${total} users`;
  }

  const start = offset + 1;
  const end = offset + usersCount;
  return `Users ${start}–${end} of ${total}`;
}

export function UsersPagination({
  offset,
  limit,
  total,
  usersCount,
  onOffsetChange,
  onLimitChange,
}: UsersPaginationProps) {
  const canGoPrev = offset > 0;
  const canGoNext = offset + usersCount < total;

  return (
    <div className="flex flex-col gap-4 border-t border-brand-border px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
      <p className="text-sm text-brand-text-muted">
        {getRangeLabel(offset, usersCount, total)}
      </p>

      <div className="flex flex-wrap items-center gap-3">
        <label className="flex items-center gap-2 text-sm text-brand-text-subtle">
          <span>Page size</span>
          <select
            value={limit}
            onChange={(event) => onLimitChange(Number(event.target.value))}
            className={FIELD_CLASS_NAME}
            aria-label="Page size"
          >
            {PAGE_SIZE_OPTIONS.map((size) => (
              <option key={size} value={size}>
                {size}
              </option>
            ))}
          </select>
        </label>

        <div className="flex items-center gap-2">
          <ActionButton
            onClick={() => onOffsetChange(Math.max(0, offset - limit))}
            disabled={!canGoPrev}
            icon={<ChevronLeft aria-hidden="true" className="size-4" />}
            aria-label="Previous page"
          >
            Previous
          </ActionButton>
          <ActionButton
            onClick={() => onOffsetChange(offset + limit)}
            disabled={!canGoNext}
            icon={<ChevronRight aria-hidden="true" className="size-4" />}
            aria-label="Next page"
          >
            Next
          </ActionButton>
        </div>
      </div>
    </div>
  );
}
