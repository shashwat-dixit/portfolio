import type { APIRoute } from "astro";
import { DATA } from "@/data/resume";
import { CONFIG } from "@/data/config";
import { getAllPosts } from "@/lib/api";
import { buildLlmsTxt, plainTextHeaders } from "@/lib/ai-content";

export const GET: APIRoute = async () => {
  let posts: Awaited<ReturnType<typeof getAllPosts>> = [];
  try {
    posts = await getAllPosts();
  } catch (error) {
    console.error("Failed to load posts for llms.txt:", error);
  }

  const body = buildLlmsTxt(CONFIG.site.url, DATA.name, DATA.description, posts);
  return new Response(body, { headers: plainTextHeaders() });
};
