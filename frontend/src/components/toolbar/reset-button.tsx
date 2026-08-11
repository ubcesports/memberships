import { RotateCcw } from "lucide-react";
import { ActionButton } from "../action-button";

type ResetButtonProps = {
  label?: string;
  onClick: () => void;
  disabled: boolean;
};

export function ResetButton({ label = "Reset", onClick, disabled }: ResetButtonProps) {
  return (
    <ActionButton
      onClick={onClick}
      disabled={disabled}
      icon={<RotateCcw aria-hidden="true" className="size-4" />}
    >
      {label}
    </ActionButton>
  );
}
