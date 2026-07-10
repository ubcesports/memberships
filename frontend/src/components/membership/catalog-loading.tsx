export function CatalogLoading() {
  return (
    <div
      className="grid gap-5 lg:grid-cols-2"
      aria-label="Loading membership passes"
    >
      {[0, 1].map((item) => (
        <div
          key={item}
          className="h-124 animate-pulse border border-brand-border bg-brand-surface/60"
        />
      ))}
    </div>
  );
}
