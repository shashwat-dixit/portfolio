export function blogPostOgImagePath(slug: string): string {
  return `/og/blog/${encodeURIComponent(slug)}.png`;
}

export const BLOG_INDEX_OG_PATH = "/og/blog.png";

export function blogOgImagePath(
  title?: string,
  description?: string
): string {
  if (!title && !description) return BLOG_INDEX_OG_PATH;

  const params = new URLSearchParams();
  if (title) params.set("title", title);
  if (description) params.set("description", description);
  return `${BLOG_INDEX_OG_PATH}?${params.toString()}`;
}
