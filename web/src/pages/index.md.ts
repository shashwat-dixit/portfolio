import type { APIRoute } from "astro";
import { markdownHeaders } from "@/lib/ai-content";
import { homepageMarkdown } from "@/lib/homepage-md";

export const GET: APIRoute = ({ request }) => {
  return new Response(homepageMarkdown(), { headers: markdownHeaders(request) });
};
