export default function SectionLabel({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex items-center w-full">
      <div className="flex-1 h-px bg-linear-to-r from-transparent from-5% via-border via-95% to-transparent" />
      <div className="border bg-primary z-10 rounded-xl px-4 py-1">
        <h2 className="text-background text-sm font-medium">{children}</h2>
      </div>
      <div className="flex-1 h-px bg-linear-to-l from-transparent from-5% via-border via-95% to-transparent" />
    </div>
  );
}
