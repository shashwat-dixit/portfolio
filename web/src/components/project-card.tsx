/* eslint-disable @next/next/no-img-element */

import { Badge } from "@/components/ui/badge";
import { ProjectShareButton } from "@/components/project-share-button";
import { cn } from "@/lib/utils";
import { getGithubUrlFromLinks } from "@/lib/share";
import { ArrowUpRight } from "lucide-react";
import { useState } from "react";
import Markdown from "react-markdown";

function ProjectMedia({
  src,
  alt,
  video,
}: {
  src?: string;
  alt: string;
  video?: string;
}) {
  const [imageError, setImageError] = useState(false);

  if (video) {
    return (
      <video
        src={video}
        autoPlay
        loop
        muted
        playsInline
        className="h-full w-full object-cover"
      />
    );
  }

  if (!src || imageError) {
    return null;
  }

  return (
    <img
      src={src}
      alt={alt}
      className="h-full w-full object-cover"
      onError={() => setImageError(true)}
    />
  );
}

interface Props {
  title: string;
  href?: string;
  description: string;
  dates: string;
  tags: readonly string[];
  image?: string;
  video?: string;
  index?: number;
  links?: readonly {
    icon: React.ReactNode;
    type: string;
    href: string;
  }[];
  className?: string;
}

export function ProjectCard({
  title,
  href,
  description,
  dates,
  tags,
  image,
  video,
  index,
  links,
  className,
}: Props) {
  const githubUrl = getGithubUrlFromLinks(links);
  const hasMedia = Boolean(video || image);
  const number =
    typeof index === "number" ? String(index + 1).padStart(2, "0") : undefined;

  return (
    <article
      className={cn(
        "group flex flex-col gap-4 rounded-xl border border-border p-5 transition-all duration-200 hover:ring-2 hover:ring-muted sm:flex-row sm:items-start",
        className
      )}
    >
      <div className="flex min-w-0 flex-1 flex-col gap-3">
        <div className="flex items-start justify-between gap-3">
          <div className="flex min-w-0 items-start gap-3">
            {number ? (
              <span className="mt-0.5 font-mono text-xs tabular-nums text-muted-foreground">
                {number}
              </span>
            ) : null}
            <div className="flex min-w-0 flex-col gap-1">
              <h3 className="font-semibold leading-none">{title}</h3>
              <time className="text-xs text-muted-foreground">{dates}</time>
            </div>
          </div>
          <a
            href={href || "#"}
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-sm text-muted-foreground transition-colors hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            aria-label={`Open ${title}`}
          >
            <ArrowUpRight className="h-4 w-4" aria-hidden />
          </a>
        </div>

        <div className="prose max-w-full text-pretty font-sans text-xs leading-relaxed text-muted-foreground dark:prose-invert">
          <Markdown>{description}</Markdown>
        </div>

        {tags && tags.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {tags.map((tag) => (
              <Badge
                key={tag}
                className="h-6 w-fit px-2 text-[11px] font-medium"
                variant="outline"
              >
                {tag}
              </Badge>
            ))}
          </div>
        )}

        <div className="relative z-10 mt-auto flex flex-wrap items-center gap-2 overflow-visible">
          {links?.map((link, idx) => (
            <a
              href={link.href}
              key={idx}
              target="_blank"
              rel="noopener noreferrer"
            >
              <Badge
                className="flex items-center gap-1.5 bg-black text-xs text-white hover:bg-black/90"
                variant="default"
              >
                {link.icon}
                {link.type}
              </Badge>
            </a>
          ))}
          <ProjectShareButton
            title={title}
            description={description}
            websiteUrl={href}
            githubUrl={githubUrl}
          />
        </div>
      </div>

      {hasMedia ? (
        <a
          href={href || "#"}
          target="_blank"
          rel="noopener noreferrer"
          className="block h-36 w-full shrink-0 overflow-hidden rounded-lg bg-muted sm:h-28 sm:w-40"
          aria-label={`${title} preview`}
        >
          <ProjectMedia src={image} alt={title} video={video} />
        </a>
      ) : null}
    </article>
  );
}
