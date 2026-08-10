import { useEffect, useId, useRef, useState, type MouseEvent as ReactMouseEvent } from "react";
import { Check, Copy, Share2 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import {
  buildShareText,
  copyShareText,
  getShareTargets,
  shareProject,
  type ProjectSharePayload,
} from "@/lib/share";

type Props = {
  title: string;
  description: string;
  websiteUrl?: string;
  githubUrl?: string;
  className?: string;
};

export function ProjectShareButton({
  title,
  description,
  websiteUrl,
  githubUrl,
  className,
}: Props) {
  const [open, setOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const menuId = useId();

  const payload: ProjectSharePayload = {
    title,
    description,
    websiteUrl,
    githubUrl,
  };
  const shareText = buildShareText(payload);
  const targets = getShareTargets(shareText, websiteUrl);

  useEffect(() => {
    if (!open) return;

    const onPointerDown = (event: globalThis.MouseEvent) => {
      if (!menuRef.current?.contains(event.target as Node)) {
        setOpen(false);
      }
    };

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setOpen(false);
    };

    document.addEventListener("mousedown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("mousedown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [open]);

  useEffect(() => {
    if (!copied) return;
    const timer = window.setTimeout(() => setCopied(false), 2000);
    return () => window.clearTimeout(timer);
  }, [copied]);

  const handleShareClick = async (event: ReactMouseEvent) => {
    event.preventDefault();
    event.stopPropagation();

    const result = await shareProject(payload);
    if (result === "fallback") {
      setOpen((prev) => !prev);
    } else {
      setOpen(false);
    }
  };

  const handleCopy = async (event: ReactMouseEvent) => {
    event.preventDefault();
    event.stopPropagation();
    const ok = await copyShareText(shareText);
    if (ok) setCopied(true);
  };

  return (
    <div ref={menuRef} className={cn("relative", className)}>
      <button
        type="button"
        onClick={handleShareClick}
        aria-haspopup="menu"
        aria-expanded={open}
        aria-controls={menuId}
        aria-label={`Share ${title}`}
        className="rounded-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
      >
        <Badge
          className="flex items-center gap-1.5 text-xs bg-black text-white hover:bg-black/90"
          variant="default"
        >
          <Share2 className="size-3" aria-hidden />
          Share
        </Badge>
      </button>

      {open && (
        <div
          id={menuId}
          role="menu"
          aria-label={`Share ${title} options`}
          className="absolute right-0 top-full z-20 mt-2 w-48 rounded-lg border border-border bg-background p-1.5 shadow-lg"
          onClick={(event) => event.stopPropagation()}
        >
          <p className="px-2 py-1.5 text-[11px] font-medium text-muted-foreground">
            Share via
          </p>
          {targets.map((target) => (
            <a
              key={target.id}
              href={target.href}
              target="_blank"
              rel="noopener noreferrer"
              role="menuitem"
              className="flex w-full items-center rounded-md px-2 py-1.5 text-sm text-foreground hover:bg-accent"
              onClick={() => setOpen(false)}
            >
              {target.label}
            </a>
          ))}
          <button
            type="button"
            role="menuitem"
            onClick={handleCopy}
            className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm text-foreground hover:bg-accent"
          >
            {copied ? (
              <Check className="size-3.5 text-green-600" aria-hidden />
            ) : (
              <Copy className="size-3.5" aria-hidden />
            )}
            {copied ? "Copied" : "Copy text"}
          </button>
        </div>
      )}
    </div>
  );
}
