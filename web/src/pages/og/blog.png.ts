import type { APIRoute } from "astro";
import {
  fallbackOgResponse,
  pngResponse,
  renderOgImage,
} from "@/lib/og-image";
import { DATA } from "@/data/resume";

const DEFAULT_TITLE = "Blog";
const DEFAULT_DESCRIPTION =
  "Thoughts on software development, life, and more.";

export const GET: APIRoute = async ({ url }) => {
  const title = url.searchParams.get("title")?.trim() || DEFAULT_TITLE;
  const description =
    url.searchParams.get("description")?.trim() || DEFAULT_DESCRIPTION;

  try {
    const png = await renderOgImage({
      title,
      description,
      eyebrow: "Blog",
      author: DATA.name,
    });
    return pngResponse(png);
  } catch (error) {
    console.error("Failed to render blog index OG image:", error);
    return fallbackOgResponse();
  }
};
