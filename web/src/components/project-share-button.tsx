import {
  useCallback,
  useEffect,
  useId,
  useLayoutEffect,
  useRef,
  useState,
  type ComponentType,
  type CSSProperties,
  type MouseEvent as ReactMouseEvent,
} from "react";
import { createPortal } from "react-dom";
import { Check, Copy, Share2 } from "lucide-react";
import {
  FaEnvelope,
  FaFacebook,
  FaLinkedinIn,
  FaRedditAlien,
  FaTelegram,
  FaWhatsapp,
  FaXTwitter,
} from "react-icons/fa6";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
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

const SHARE_ICONS: Record<
  string,
  ComponentType<{ className?: string }>
> = {
  x: FaXTwitter,
  linkedin: FaLinkedinIn,
  facebook: FaFacebook,
  reddit: FaRedditAlien,
  whatsapp: FaWhatsapp,
  telegram: FaTelegram,
  email: FaEnvelope,
};

const MENU_GAP = 8;
const VIEWPORT_PADDING = 8;

function positionMenu(
  anchor: DOMRect,
  menu: HTMLElement
): CSSProperties {
  const menuWidth = menu.offsetWidth || 208;
  const menuHeight = menu.offsetHeight || 0;
  const maxLeft = window.innerWidth - menuWidth - VIEWPORT_PADDING;
  const left = Math.max(
    VIEWPORT_PADDING,
    Math.min(anchor.right - menuWidth, maxLeft)
  );

  const spaceBelow = window.innerHeight - anchor.bottom - MENU_GAP;
  const spaceAbove = anchor.top - MENU_GAP;
  const openUpward =
    menuHeight > 0 && spaceBelow < menuHeight && spaceAbove > spaceBelow;

  const top = openUpward
    ? Math.max(VIEWPORT_PADDING, anchor.top - menuHeight - MENU_GAP)
    : Math.min(
        anchor.bottom + MENU_GAP,
        window.innerHeight - menuHeight - VIEWPORT_PADDING
      );

  return {
    position: "fixed",
    top,
    left,
    zIndex: 50,
  };
}

export function ProjectShareButton({
  title,
  description,
  websiteUrl,
  githubUrl,
  className,
}: Props) {
  const [open, setOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  const [menuStyle, setMenuStyle] = useState<CSSProperties>({
    position: "fixed",
    top: 0,
    left: 0,
    zIndex: 50,
  });
  const rootRef = useRef<HTMLDivElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);
  const menuId = useId();

  const setMenuRef = useCallback((node: HTMLDivElement | null) => {
    menuRef.current = node;
    const anchor = rootRef.current?.getBoundingClientRect();
    if (!node || !anchor) return;
    setMenuStyle(positionMenu(anchor, node));
  }, []);

  const payload: ProjectSharePayload = {
    title,
    description,
    websiteUrl,
    githubUrl,
  };
  const shareText = buildShareText(payload);
  const targets = getShareTargets(shareText, websiteUrl);

  const updateMenuPosition = useCallback(() => {
    const anchor = rootRef.current?.getBoundingClientRect();
    const menu = menuRef.current;
    if (!anchor || !menu) return;
    setMenuStyle(positionMenu(anchor, menu));
  }, []);

  useLayoutEffect(() => {
    if (!open) return;
    updateMenuPosition();
  }, [open, updateMenuPosition, targets.length, copied]);

  useEffect(() => {
    if (!open) return;

    const onPointerDown = (event: globalThis.MouseEvent) => {
      const target = event.target as Node;
      if (
        rootRef.current?.contains(target) ||
        menuRef.current?.contains(target)
      ) {
        return;
      }
      setOpen(false);
    };

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setOpen(false);
    };

    document.addEventListener("mousedown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    window.addEventListener("resize", updateMenuPosition);
    window.addEventListener("scroll", updateMenuPosition, true);
    return () => {
      document.removeEventListener("mousedown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
      window.removeEventListener("resize", updateMenuPosition);
      window.removeEventListener("scroll", updateMenuPosition, true);
    };
  }, [open, updateMenuPosition]);

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

  const itemClassName =
    "flex h-9 w-full shrink-0 items-center gap-2.5 rounded-md px-2.5 text-sm text-foreground outline-none hover:bg-accent focus-visible:bg-accent";

  const menu =
    open && typeof document !== "undefined"
      ? createPortal(
          <div
            ref={setMenuRef}
            id={menuId}
            role="menu"
            aria-label={`Share ${title} options`}
            style={menuStyle}
            className="flex max-h-[min(24rem,calc(100vh-1rem))] w-52 flex-col overflow-y-auto overscroll-contain rounded-xl border border-border bg-popover p-1 text-popover-foreground shadow-lg"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="px-2.5 py-1.5 text-[11px] font-medium text-muted-foreground">
              Share via
            </div>
            {targets.map((target) => {
              const Icon = SHARE_ICONS[target.id];
              return (
                <a
                  key={target.id}
                  href={target.href}
                  target="_blank"
                  rel="noopener noreferrer"
                  role="menuitem"
                  className={itemClassName}
                  onClick={() => setOpen(false)}
                >
                  {Icon ? (
                    <span
                      className="inline-flex size-3.5 shrink-0 items-center justify-center text-muted-foreground"
                      aria-hidden
                    >
                      <Icon className="size-3.5" />
                    </span>
                  ) : null}
                  {target.label}
                </a>
              );
            })}
            <Separator className="my-1" />
            <button
              type="button"
              role="menuitem"
              onClick={handleCopy}
              className={itemClassName}
            >
              {copied ? (
                <Check className="size-3.5 shrink-0 text-green-600" aria-hidden />
              ) : (
                <Copy className="size-3.5 shrink-0 text-muted-foreground" aria-hidden />
              )}
              {copied ? "Copied" : "Copy text"}
            </button>
          </div>,
          document.body
        )
      : null;

  return (
    <div ref={rootRef} className={cn("relative inline-flex shrink-0", className)}>
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
      {menu}
    </div>
  );
}
