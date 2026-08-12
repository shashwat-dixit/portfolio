import { Badge } from "@/components/ui/badge";
import SectionLabel from "@/components/section/section-label";
import { DATA } from "@/data/resume";
import { Timeline, TimelineItem, TimelineConnectItem } from "@/components/timeline";

export default function WorkTimelineSection() {
  return (
    <div className="flex min-h-0 flex-col gap-y-8 w-full">
      <SectionLabel>{DATA.sections.work.label}</SectionLabel>
      <Timeline>
        {DATA.work.map((job) => (
          <TimelineItem
            key={job.company + job.start}
            className="grid w-full grid-cols-[auto_minmax(0,1fr)] items-start gap-x-3 gap-y-4 sm:flex sm:items-start sm:justify-between sm:gap-10"
          >
            <TimelineConnectItem className="flex items-start justify-center max-sm:self-start max-sm:[&_[data-timeline-line]]:hidden">
              {job.logoUrl ? (
                <img
                  src={job.logoUrl}
                  alt={job.company}
                  className="size-8 sm:size-10 bg-card z-10 shrink-0 overflow-hidden p-1 border rounded-full shadow ring-2 ring-border object-contain flex-none"
                />
              ) : (
                <div className="size-8 sm:size-10 bg-card z-10 shrink-0 overflow-hidden p-1 border rounded-full shadow ring-2 ring-border flex-none" />
              )}
            </TimelineConnectItem>
            <div className="flex min-w-0 flex-1 flex-col justify-start gap-2">
              <time className="text-xs text-muted-foreground">
                {job.start} — {job.end || "Present"}
              </time>
              <h3 className="font-semibold leading-none">{job.title}</h3>
              <div className="flex items-center gap-2">
                <a
                  href={job.href}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                >
                  {job.company}
                </a>
                {job.badges && job.badges.length > 0 && job.badges.map((badge) => (
                  <Badge key={badge} variant="secondary" className="text-[10px]">
                    {badge}
                  </Badge>
                ))}
              </div>
              <ul className="mt-1 hidden space-y-1.5 text-sm leading-relaxed text-muted-foreground sm:block">
                {job.bullets.map((bullet, idx) => (
                  <li key={idx} className="flex gap-2">
                    <span className="mt-1.5 shrink-0 text-muted-foreground/60">•</span>
                    <span>{bullet}</span>
                  </li>
                ))}
              </ul>
            </div>
            <ul className="col-span-2 list-disc space-y-3 pl-5 text-sm leading-relaxed text-muted-foreground marker:text-muted-foreground/70 sm:hidden">
              {job.bullets.map((bullet, idx) => (
                <li key={idx} className="pl-1 text-pretty">
                  {bullet}
                </li>
              ))}
            </ul>
          </TimelineItem>
        ))}
      </Timeline>
    </div>
  );
}
