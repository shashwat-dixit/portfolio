import type { APIRoute } from "astro";
import { getPostBySlug } from "@/lib/api";
import {
  fallbackOgResponse,
  pngResponse,
  renderOgImage,
} from "@/lib/og-image";
import { DATA } from "@/data/resume";

export const GET: APIRoute = async ({ params }) => {
  const slug = params.slug;
  if (!slug) {
    return new Response("Not found", { status: 404 });
  }

  try {
    const post = await getPostBySlug(slug);
    const png = await renderOgImage({
      title: post.title,
      description: post.description || DATA.description,
      eyebrow: "Blog",
      author: post.author || DATA.name,
    });
    return pngResponse(png);
  } catch (error) {
    if (error instanceof Error && /404|not found/i.test(error.message)) {
      return new Response("Not found", { status: 404 });
    }
    console.error("Failed to render blog OG image:", error);
    return fallbackOgResponse();
  }
};
