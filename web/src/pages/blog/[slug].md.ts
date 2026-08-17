import type { APIRoute } from "astro";
import { getPostBySlug } from "@/lib/api";
import { looksLikePostMarkdown, markdownHeaders } from "@/lib/ai-content";

export const GET: APIRoute = async ({ params }) => {
  const slug = params.slug;
  if (!slug) {
    return new Response("Not found\n", { status: 404 });
  }

  try {
    const post = await getPostBySlug(slug);
    const body = post.contentMd?.trim()
      ? post.contentMd
      : `---\ntitle: "${post.title}"\nslug: ${post.slug}\n---\n\n${post.description}\n`;

    if (!looksLikePostMarkdown(body) && !body.trim().startsWith("#")) {
      console.warn("Post markdown did not look like markdown:", slug);
    }

    return new Response(body.endsWith("\n") ? body : `${body}\n`, {
      headers: markdownHeaders(),
    });
  } catch {
    return new Response("Not found\n", { status: 404 });
  }
};
