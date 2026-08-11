import { ResetButton } from "@/components/toolbar/reset-button";
import { SearchField } from "@/components/toolbar/search-field";
import { ToolbarContainer } from "@/components/toolbar/toolbar-container";
import { ToolbarRow } from "@/components/toolbar/toolbar-row";

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
    <ToolbarContainer>
      <ToolbarRow>
        <SearchField 
          label="Search by actor name"
          placeholder="Enter actor name"
          value={searchInput}
          onChange={onSearchInputChange}
        />
        <div className="flex flex-wrap gap-2">
          <ResetButton 
            label="Reset search"
            onClick={onResetSearch}
            disabled={searchInput === ""}
          />
        </div>
      </ToolbarRow>
    </ToolbarContainer>
  )
}
