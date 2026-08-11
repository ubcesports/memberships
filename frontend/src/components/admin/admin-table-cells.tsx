import { formatTime } from "@/lib/utils/formatting";
import Link from "next/dist/client/link";
import Image from "next/image";

export function EmptyValue() {
  return <span className="text-brand-text-muted">—</span>;
}

export function formatOptionalTime(value: string | null) {
  if (!value) {
    return <EmptyValue />;
  }

  return formatTime(value);
}

type AvatarCellProps = {
    src: string | null; 
    alt?: string;
};

export function AvatarCell({ src, alt="" }: AvatarCellProps) {
    if (!src) return <EmptyValue />;

    return (
        <Image 
            src={src}
            alt={alt}
            width={32}
            height={32}
            className="size-8 border border-brand-border object-cover"
            unoptimized
        />
    );
}

type LinkCellProps = {
    href: string;
    children: React.ReactNode;
};

export function LinkCell({ href, children }: LinkCellProps) {
    return <Link href={href}>{children}</Link>
}