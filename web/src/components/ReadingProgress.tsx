import { useReadingProgress } from "@/hooks/useReadingProgress";

interface ReadingProgressProps {
  slug: string;
}

export default function ReadingProgress({ slug }: ReadingProgressProps) {
  const { progress } = useReadingProgress(slug);

  return (
    <div className="fixed top-0 left-0 right-0 z-50 h-1 bg-muted">
      <div
        className="h-full bg-primary transition-[width] duration-150 ease-out"
        style={{ width: `${progress}%` }}
      />
    </div>
  );
}
