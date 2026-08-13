export type ProjectSharePayload = {
  title: string;
  description: string;
  websiteUrl?: string;
  githubUrl?: string;
};

export function buildShareText({
  title,
  description,
  websiteUrl,
  githubUrl,
}: ProjectSharePayload): string {
  const lines = [title, "", description.trim()];

  if (websiteUrl) {
    lines.push("", `Website: ${websiteUrl}`);
  }

  if (githubUrl) {
    lines.push(`GitHub: ${githubUrl}`);
  }

  return lines.join("\n");
}

export function getGithubUrlFromLinks(
  links?: readonly { type: string; href: string }[]
): string | undefined {
  if (!links?.length) return undefined;

  const sourceLink = links.find(
    (link) =>
      link.type.toLowerCase() === "source" ||
      link.type.toLowerCase() === "github" ||
      link.href.includes("github.com")
  );

  return sourceLink?.href;
}

export type ShareTarget = {
  id: string;
  label: string;
  href: string;
};

export function getShareTargets(text: string, url?: string): ShareTarget[] {
  const encodedText = encodeURIComponent(text);
  const shareUrl = url || (typeof window !== "undefined" ? window.location.href : "");
  const encodedUrl = encodeURIComponent(shareUrl);
  const encodedTitle = encodeURIComponent(text.split("\n")[0] || "Project");

  const targets: ShareTarget[] = [
    {
      id: "x",
      label: "X",
      href: `https://twitter.com/intent/tweet?text=${encodedText}`,
    },
    {
      id: "whatsapp",
      label: "WhatsApp",
      href: `https://wa.me/?text=${encodedText}`,
    },
    {
      id: "telegram",
      label: "Telegram",
      href: `https://t.me/share/url?url=${encodedUrl}&text=${encodedText}`,
    },
    {
      id: "email",
      label: "Email",
      href: `mailto:?subject=${encodedTitle}&body=${encodedText}`,
    },
  ];

  if (shareUrl) {
    targets.splice(
      1,
      0,
      {
        id: "linkedin",
        label: "LinkedIn",
        href: `https://www.linkedin.com/sharing/share-offsite/?url=${encodedUrl}`,
      },
      {
        id: "facebook",
        label: "Facebook",
        href: `https://www.facebook.com/sharer/sharer.php?u=${encodedUrl}&quote=${encodedText}`,
      },
      {
        id: "reddit",
        label: "Reddit",
        href: `https://www.reddit.com/submit?url=${encodedUrl}&title=${encodedTitle}`,
      }
    );
  }

  return targets;
}

function isHandheldDevice(): boolean {
  if (typeof navigator === "undefined") return false;

  const ua = navigator.userAgent;
  if (/Android|iPhone|iPod|iPad/i.test(ua)) return true;

  // iPadOS 13+ Safari reports as Macintosh with touch points.
  return navigator.platform === "MacIntel" && navigator.maxTouchPoints > 1;
}

export function shouldUseNativeShare(): boolean {
  return (
    typeof navigator !== "undefined" &&
    typeof navigator.share === "function" &&
    isHandheldDevice()
  );
}

function buildNativeShareData(payload: ProjectSharePayload): ShareData {
  const url = payload.websiteUrl || payload.githubUrl;
  return {
    title: payload.title,
    text: payload.description.trim(),
    ...(url ? { url } : {}),
  };
}

export async function shareProject(
  payload: ProjectSharePayload
): Promise<"shared" | "fallback" | "cancelled"> {
  // Desktop share sheets (Safari/macOS especially) scrape URLs out of `text`
  // and present them as link-preview images — e.g. "2 Images" / "Phabric and
  // Phabric" — instead of the project write-up. The in-page menu is the
  // better desktop UX anyway.
  if (shouldUseNativeShare()) {
    const data = buildNativeShareData(payload);
    if (!navigator.canShare || navigator.canShare(data)) {
      try {
        await navigator.share(data);
        return "shared";
      } catch (error) {
        if (error instanceof DOMException && error.name === "AbortError") {
          return "cancelled";
        }
        // Fall through to custom menu when native share fails
      }
    }
  }

  return "fallback";
}

export async function copyShareText(text: string): Promise<boolean> {
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text);
      return true;
    }
  } catch {
    // Fall through to legacy copy
  }

  try {
    const textarea = document.createElement("textarea");
    textarea.value = text;
    textarea.setAttribute("readonly", "");
    textarea.style.position = "fixed";
    textarea.style.left = "-9999px";
    document.body.appendChild(textarea);
    textarea.select();
    const ok = document.execCommand("copy");
    document.body.removeChild(textarea);
    return ok;
  } catch {
    return false;
  }
}
