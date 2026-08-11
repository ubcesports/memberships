type ToolbarRowProps = {
    children: React.ReactNode;
    justify?: boolean;
}

export function ToolbarRow({ children, justify = false }: ToolbarRowProps) {
    return (
        <div
            className={`flex flex-col gap-3 ${
                justify ? "xl:flex-row xl:items-end xl:justify-between" : "lg:flex-row lg:items-end"
            }`}
        >
            {children}
        </div>
    )
}