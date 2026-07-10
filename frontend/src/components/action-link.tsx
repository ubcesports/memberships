import type { AnchorHTMLAttributes, ReactNode } from "react";

export type ActionLinkProps = {
  children: ReactNode;
  className?: string;
  icon?: ReactNode;
} & Omit<AnchorHTMLAttributes<HTMLAnchorElement>, "children" | "className">;

const BASE_CLASS_NAME =
  "inline-flex h-10 items-center justify-center gap-2 border border-brand-border px-4 text-sm font-semibold text-brand-text transition hover:border-brand-text-muted hover:bg-white/5";

export function ActionLink({ children, className, icon, ...props }: ActionLinkProps) {
  return (
    <a
      {...props}
      className={
        className ? `${BASE_CLASS_NAME} ${className}` : BASE_CLASS_NAME
      }
    >
      {icon}
      <span>{children}</span>
    </a>
  );
}
