import type { APIRoute } from "astro";
import { CONFIG } from "@/data/config";
import { getAllPosts } from "@/lib/api";
import { buildBlogIndexMarkdown, markdownHeaders } from "@/lib/ai-content";

export const GET: APIRoute = async ({ request }) => {
  let posts: Awaited<ReturnType<typeof getAllPosts>> = [];
  try {
    posts = await getAllPosts();
  } catch (error) {
    console.error("Failed to load posts for blog.md:", error);
  }

  return new Response(buildBlogIndexMarkdown(CONFIG.site.url, posts), {
    headers: markdownHeaders(request),
  });
};
