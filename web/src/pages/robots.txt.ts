import type { APIRoute } from "astro";
import { CONFIG } from "@/data/config";
import { buildRobotsTxt } from "@/lib/ai-content";

export const GET: APIRoute = () => {
  return new Response(buildRobotsTxt(CONFIG.site.url), {
    headers: {
      "Content-Type": "text/plain; charset=utf-8",
      "Cache-Control": "public, max-age=86400",
    },
  });
};
