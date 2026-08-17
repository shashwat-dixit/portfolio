import type { APIRoute } from "astro";
import { getAllPosts, getPostBySlug } from "@/lib/api";
import { buildLlmsFullTxt, markdownHeaders } from "@/lib/ai-content";
import { homepageMarkdown } from "@/lib/homepage-md";

export const GET: APIRoute = async () => {
  let posts: Awaited<ReturnType<typeof getAllPosts>> = [];
  try {
    posts = await getAllPosts();
  } catch (error) {
    console.error("Failed to load posts for llms-full.txt:", error);
  }

  const withBodies = await Promise.all(
    posts.map(async (post) => {
      try {
        const full = await getPostBySlug(post.slug);
        return {
          slug: post.slug,
          title: post.title,
          description: post.description,
          contentMd: full.contentMd,
        };
      } catch (error) {
        console.error("Failed to load post markdown:", post.slug, error);
        return post;
      }
    })
  );

  const body = buildLlmsFullTxt(homepageMarkdown(), withBodies);
  return new Response(body, { headers: markdownHeaders() });
};
