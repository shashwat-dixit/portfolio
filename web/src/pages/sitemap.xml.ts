import type { APIRoute } from "astro";
import { CONFIG } from "@/data/config";
import { getAllPosts } from "@/lib/api";
import { buildSitemapXml } from "@/lib/ai-content";

export const GET: APIRoute = async () => {
  let slugs: string[] = [];
  try {
    slugs = (await getAllPosts()).map((post) => post.slug);
  } catch (error) {
    console.error("Failed to load posts for sitemap:", error);
  }

  const body = buildSitemapXml(CONFIG.site.url, slugs, new Date().toISOString().slice(0, 10));
  return new Response(body, {
    headers: {
      "Content-Type": "application/xml; charset=utf-8",
      "Cache-Control": "public, max-age=3600, stale-while-revalidate=86400",
    },
  });
};
