import SectionLabel from "@/components/section/section-label";
import { DATA } from "@/data/resume";
import Markdown from "react-markdown";

export default function AboutSection() {
  return (
    <div className="flex min-h-0 flex-col gap-y-8">
      <SectionLabel>{DATA.sections.about.label}</SectionLabel>
      <div className="prose max-w-full text-pretty font-sans leading-relaxed text-muted-foreground dark:prose-invert">
        <Markdown>{DATA.summary}</Markdown>
      </div>
    </div>
  );
}
