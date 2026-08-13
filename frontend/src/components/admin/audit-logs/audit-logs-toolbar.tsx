import { ActionButton } from "@/components/action-button";
import { ResetButton } from "@/components/toolbar/reset-button";
import { SearchField } from "@/components/toolbar/search-field";
import { ToolbarContainer } from "@/components/toolbar/toolbar-container";
import { ToolbarRow } from "@/components/toolbar/toolbar-row";
import { Download, Loader2 } from "lucide-react";

type AuditLogsToolbarProps = {
  searchInput: string;
  total: number;
  isExporting: boolean;
  onSearchInputChange: (value: string) => void;
  onResetSearch: () => void;
  onExport: () => void;
};

export function AuditLogsToolbar({
  searchInput,
  total,
  isExporting,
  onSearchInputChange,
  onResetSearch,
  onExport,
}: AuditLogsToolbarProps) {
  return (
    <ToolbarContainer>
      <ToolbarRow>
        <SearchField
          label="Search by actor name"
          placeholder="Enter actor name"
          value={searchInput}
          onChange={onSearchInputChange}
        />
        <div className="flex flex-wrap gap-2">
          <ResetButton label="Reset search" onClick={onResetSearch} disabled={searchInput === ""} />
          <ActionButton
            onClick={onExport}
            disabled={total === 0 || isExporting}
            loading={isExporting}
            icon={<Download aria-hidden="true" className="size-4" />}
            loadingIcon={<Loader2 aria-hidden="true" className="size-4 animate-spin" />}
          >
            {isExporting ? "Exporting" : "Export CSV"}
          </ActionButton>
        </div>
      </ToolbarRow>
    </ToolbarContainer>
  );
}
