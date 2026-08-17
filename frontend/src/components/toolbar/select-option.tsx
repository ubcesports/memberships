type SelectFieldProps<T extends string> = {
  label: string;
  value: T | "";
  onChange: (value: string) => void;
  options: { value: T; label: string }[];
  allLabel?: string; // optional label for the "all" option
  allValue?: string;
  ariaLabel?: string;
};

export function SelectField<T extends string>({
  label,
  value,
  onChange,
  options,
  allLabel,
  allValue = "",
  ariaLabel,
}: SelectFieldProps<T>) {
  return (
    <label className="flex flex-col gap-1.5 text-sm text-brand-text-subtle">
      <span>{label}</span>
      <select
        value={value}
        onChange={(event) => onChange(event.target.value)}
        className="h-10 border border-brand-border bg-brand-surface px-3 text-sm text-brand-text"
        aria-label={ariaLabel}
      >
        {allLabel !== undefined && <option value={allValue}>{allLabel}</option>}
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </label>
  );
}
