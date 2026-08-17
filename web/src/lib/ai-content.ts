export const MARKDOWN_CONTENT_TYPE = "text/markdown; charset=utf-8";
export const PLAIN_TEXT_CONTENT_TYPE = "text/plain; charset=utf-8";

const AI_USER_AGENTS = [
  "gptbot",
  "chatgpt-user",
  "oai-searchbot",
  "claudebot",
  "claude-user",
  "claude-searchbot",
  "anthropic-ai",
  "perplexitybot",
  "perplexity-user",
  "google-extended",
  "applebot-extended",
  "bytespider",
  "ccbot",
  "meta-externalagent",
  "cohere-ai",
];

export type LlmsPost = {
  slug: string;
  title: string;
  description: string;
  contentMd?: string;
};

export function wantsMarkdown(request: Request): boolean {
  const url = new URL(request.url);
  if (url.searchParams.get("format") === "md") {
    return true;
  }
  if (acceptsMarkdownType(request)) {
    return true;
  }

  const ua = (request.headers.get("User-Agent") ?? "").toLowerCase();
  return AI_USER_AGENTS.some((bot) => ua.includes(bot));
}

export function acceptsMarkdownType(request: Request): boolean {
  const accept = request.headers.get("Accept") ?? "";
  for (const part of accept.split(",")) {
    const mediaType = part.trim().split(";")[0]?.trim().toLowerCase();
    if (mediaType === "text/markdown" || mediaType === "text/x-markdown") {
      return true;
    }
  }
  return false;
}

export function markdownRewritePath(pathname: string, request: Request): string | null {
  if (!wantsMarkdown(request)) {
    return null;
  }

  if (
    pathname === "/llms.txt" ||
    pathname === "/llms-full.txt" ||
    pathname === "/robots.txt" ||
    pathname === "/sitemap.xml"
  ) {
    return null;
  }

  if (pathname.endsWith(".md") || pathname.endsWith(".txt") || pathname.endsWith(".xml")) {
    return null;
  }

  if (pathname.startsWith("/og/") || pathname.startsWith("/_")) {
    return null;
  }

  if (pathname.includes(".")) {
    return null;
  }

  if (pathname === "/") {
    return "/index.md";
  }

  if (pathname === "/blog") {
    return "/blog.md";
  }

  if (pathname.startsWith("/blog/") && !pathname.startsWith("/blog/tag/")) {
    return `${pathname.replace(/\/$/, "")}.md`;
  }

  return null;
}

export function markdownHeaders(request?: Request): HeadersInit {
  // ChatGPT's web reader rejects text/markdown with HTTP 400.
  // Serve plain text unless the client explicitly asked for markdown.
  const type =
    request && acceptsMarkdownType(request) ? MARKDOWN_CONTENT_TYPE : PLAIN_TEXT_CONTENT_TYPE;
  return {
    "Content-Type": type,
    "Cache-Control": "public, max-age=3600, stale-while-revalidate=86400",
    Vary: "Accept, User-Agent",
  };
}

export function plainTextHeaders(): HeadersInit {
  return {
    "Content-Type": PLAIN_TEXT_CONTENT_TYPE,
    "Cache-Control": "public, max-age=3600, stale-while-revalidate=86400",
  };
}

export function looksLikeMarkdown(body: string): boolean {
  const trimmed = body.trim();
  if (!trimmed) {
    return false;
  }
  return trimmed.startsWith("#") || trimmed.startsWith("---");
}

export function looksLikePostMarkdown(body: string): boolean {
  const trimmed = body.trim();
  if (!trimmed.startsWith("---")) {
    return false;
  }
  const closing = trimmed.indexOf("\n---", 3);
  if (closing === -1) {
    return false;
  }
  const frontmatter = trimmed.slice(3, closing);
  return /\btitle\s*:/.test(frontmatter) && /\bslug\s*:/.test(frontmatter);
}

export function buildLlmsTxt(siteUrl: string, name: string, description: string, posts: LlmsPost[]): string {
  const origin = siteUrl.replace(/\/$/, "");
  const lines = [
    `# ${name}`,
    "",
    `> ${description}`,
    "",
    "This file is an index of markdown pages for language models and coding agents.",
    "Prefer the `.md` URLs below over HTML when you need the source content.",
    "",
    "## Pages",
    "",
    `- [Home](${origin}/index.md): About, work, education, skills, projects, and contact.`,
    `- [Blog](${origin}/blog.md): Index of published writing.`,
    `- [Full content](${origin}/llms-full.txt): Homepage plus every published post in one markdown file.`,
    "",
    "## Blog",
    "",
  ];

  if (posts.length === 0) {
    lines.push("- No published posts yet.");
  } else {
    for (const post of posts) {
      const desc = post.description?.trim() ? `: ${post.description.trim()}` : "";
      lines.push(`- [${post.title}](${origin}/blog/${post.slug}.md)${desc}`);
    }
  }

  lines.push("");
  return lines.join("\n");
}

export function buildLlmsFullTxt(homepage: string, posts: LlmsPost[]): string {
  const parts = [homepage.trim()];
  for (const post of posts) {
    const body = (post.contentMd ?? "").trim() || `# ${post.title}\n\n${post.description}`;
    parts.push(body);
  }
  return `${parts.join("\n\n---\n\n")}\n`;
}

export function buildBlogIndexMarkdown(siteUrl: string, posts: LlmsPost[]): string {
  const origin = siteUrl.replace(/\/$/, "");
  const lines = [
    "# Blog",
    "",
    "Published posts in markdown. Each link is the raw source an agent should read.",
    "",
  ];

  if (posts.length === 0) {
    lines.push("No published posts yet.");
  } else {
    for (const post of posts) {
      lines.push(`## [${post.title}](${origin}/blog/${post.slug}.md)`);
      if (post.description?.trim()) {
        lines.push("");
        lines.push(post.description.trim());
      }
      lines.push("");
    }
  }

  return `${lines.join("\n").trim()}\n`;
}

export function buildRobotsTxt(siteUrl: string): string {
  const origin = siteUrl.replace(/\/$/, "");
  const aiAgents = [
    "GPTBot",
    "ChatGPT-User",
    "OAI-SearchBot",
    "ClaudeBot",
    "Claude-User",
    "Claude-SearchBot",
    "PerplexityBot",
    "Perplexity-User",
    "Google-Extended",
    "Applebot-Extended",
  ];

  const lines = [
    "# Allow search and AI crawlers. Markdown copies live at *.md, /llms.txt, and /llms-full.txt.",
    "User-agent: *",
    "Allow: /",
    "",
  ];

  for (const agent of aiAgents) {
    lines.push(`User-agent: ${agent}`);
    lines.push("Allow: /");
    lines.push("");
  }

  lines.push(`Sitemap: ${origin}/sitemap.xml`);
  lines.push(`LLMs-Txt: ${origin}/llms.txt`);
  lines.push("");
  return lines.join("\n");
}

export function buildSitemapXml(siteUrl: string, slugs: string[], lastmod?: string): string {
  const origin = siteUrl.replace(/\/$/, "");
  const urls = [
    { loc: `${origin}/`, extra: [`${origin}/index.md`] },
    { loc: `${origin}/blog`, extra: [`${origin}/blog.md`] },
    { loc: `${origin}/llms.txt`, extra: [] },
    { loc: `${origin}/llms-full.txt`, extra: [] },
    ...slugs.map((slug) => ({
      loc: `${origin}/blog/${slug}`,
      extra: [`${origin}/blog/${slug}.md`],
    })),
  ];

  const lastmodTag = lastmod ? `\n    <lastmod>${lastmod}</lastmod>` : "";
  const body = urls
    .map(({ loc, extra }) => {
      const alts = extra
        .map(
          (href) =>
            `    <xhtml:link rel="alternate" type="text/plain" href="${escapeXml(href)}" />`
        )
        .join("\n");
      return `  <url>\n    <loc>${escapeXml(loc)}</loc>${lastmodTag}${alts ? `\n${alts}` : ""}\n  </url>`;
    })
    .join("\n");

  return `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:xhtml="http://www.w3.org/1999/xhtml">
${body}
</urlset>
`;
}

function escapeXml(value: string): string {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&apos;");
}
