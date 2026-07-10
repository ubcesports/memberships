import type { HTMLAttributes, ReactNode } from "react";

export type SurfacePanelProps = {
  children: ReactNode;
  className?: string;
} & HTMLAttributes<HTMLDivElement>;

const BASE_CLASS_NAME = "border border-brand-border bg-white/[0.03]";

export function SurfacePanel({
  children,
  className,
  ...props
}: SurfacePanelProps) {
  return (
    <div
      {...props}
      className={
        className ? `${BASE_CLASS_NAME} ${className}` : BASE_CLASS_NAME
      }
    >
      {children}
    </div>
  );
}
