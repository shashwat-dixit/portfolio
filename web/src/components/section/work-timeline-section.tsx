import { Badge } from "@/components/ui/badge";
import { DATA } from "@/data/resume";
import { Timeline, TimelineItem, TimelineConnectItem } from "@/components/timeline";

export default function WorkTimelineSection() {
  return (
    <div className="flex min-h-0 flex-col gap-y-8 w-full">
      <div className="flex flex-col gap-y-4 items-center justify-center">
        <div className="flex items-center w-full">
          <div className="flex-1 h-px bg-linear-to-r from-transparent from-5% via-border via-95% to-transparent" />
          <div className="border bg-primary z-10 rounded-xl px-4 py-1">
            <span className="text-background text-sm font-medium">
              {DATA.sections.work.label}
            </span>
          </div>
          <div className="flex-1 h-px bg-linear-to-l from-transparent from-5% via-border via-95% to-transparent" />
        </div>
        <div className="flex flex-col gap-y-3 items-center justify-center">
          <h2 className="text-3xl font-bold tracking-tighter sm:text-4xl">
            {DATA.sections.work.heading}
          </h2>
          <p className="text-muted-foreground md:text-lg/relaxed lg:text-base/relaxed xl:text-lg/relaxed text-balance text-center">
            {DATA.sections.work.text}
          </p>
        </div>
      </div>
      <Timeline>
        {DATA.work.map((job) => (
          <TimelineItem key={job.company + job.start} className="w-full flex items-start justify-between gap-10">
            <TimelineConnectItem className="flex items-start justify-center">
              {job.logoUrl ? (
                <img
                  src={job.logoUrl}
                  alt={job.company}
                  className="size-10 bg-card z-10 shrink-0 overflow-hidden p-1 border rounded-full shadow ring-2 ring-border object-contain flex-none"
                />
              ) : (
                <div className="size-10 bg-card z-10 shrink-0 overflow-hidden p-1 border rounded-full shadow ring-2 ring-border flex-none" />
              )}
            </TimelineConnectItem>
            <div className="flex flex-1 flex-col justify-start gap-2 min-w-0">
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
              <ul className="mt-1 space-y-1.5 text-sm text-muted-foreground leading-relaxed">
                {job.bullets.map((bullet, idx) => (
                  <li key={idx} className="flex gap-2">
                    <span className="text-muted-foreground/60 mt-1.5 shrink-0">•</span>
                    <span>{bullet}</span>
                  </li>
                ))}
              </ul>
            </div>
          </TimelineItem>
        ))}
      </Timeline>
    </div>
  );
}
