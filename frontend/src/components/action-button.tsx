import type { ButtonHTMLAttributes, ReactNode } from "react";

export type ActionButtonProps = {
  children: ReactNode;
  className?: string;
  icon?: ReactNode;
  loading?: boolean;
  loadingIcon?: ReactNode;
} & Omit<ButtonHTMLAttributes<HTMLButtonElement>, "children" | "className">;

const BASE_CLASS_NAME =
  "inline-flex h-10 cursor-pointer items-center justify-center gap-2 border border-brand-border px-4 text-sm font-semibold text-brand-text transition hover:border-brand-text-muted hover:bg-white/5 disabled:cursor-not-allowed disabled:opacity-60";

export function ActionButton({
  children,
  className,
  icon,
  loading,
  loadingIcon,
  disabled,
  type,
  ...props
}: ActionButtonProps) {
  const content = (
    <>
      {loading ? loadingIcon : icon}
      <span>{children}</span>
    </>
  );

  return (
    <button
      {...props}
      type={type ?? "button"}
      disabled={disabled || loading}
      className={
        className ? `${BASE_CLASS_NAME} ${className}` : BASE_CLASS_NAME
      }
    >
      {content}
    </button>
  );
}
