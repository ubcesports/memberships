type SearchFieldProps = {
    label: string;
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    ariaLabel?: string;
    type?: "text" | "search";
    className?: string;
};

export function SearchField({
    label,
    value,
    onChange,
    placeholder,
    ariaLabel,
    type = "text",
    className,
}: SearchFieldProps) {
    return (
        <label className={`flex min-w-40 flex-col gap-1.5 text-sm text-brand-text-subtle ${className}`}>
          <span>{label}</span>
            <input
                type={type}
                placeholder={placeholder}
                className="h-10 border border-brand-border bg-brand-surface px-3 text-sm text-brand-text"
                value={value}
                onChange={(e) => onChange(e.target.value)}
                aria-label={ariaLabel}
            />
        </label>
    )
}