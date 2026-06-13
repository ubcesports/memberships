import { BasePage } from "@/components/layout/base-page";

export default function HomePage() {
  return (
    <BasePage>
      <div className="flex flex-1 flex-col items-center justify-center gap-6">
        <h1 className="text-2xl font-semibold text-brand-text">
          UBCEA Memberships
        </h1>
      </div>
    </BasePage>
  );
}
