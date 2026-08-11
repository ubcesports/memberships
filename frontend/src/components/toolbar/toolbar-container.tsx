export function ToolbarContainer({ children }: { children: React.ReactNode }) {
    return (
        <div className="flex flex-col gap-5 border-b border-brand-border px-5 py-5 sm:px-6">
            {children}
        </div>
    )
}