import assert from "node:assert/strict";
import { describe, it } from "node:test";
import {
  buildBlogIndexMarkdown,
  buildLlmsFullTxt,
  buildLlmsTxt,
  buildRobotsTxt,
  buildSitemapXml,
  looksLikeMarkdown,
  looksLikePostMarkdown,
  markdownHeaders,
  markdownRewritePath,
  plainTextHeaders,
  wantsMarkdown,
} from "./ai-content.ts";

function request(path: string, headers: Record<string, string> = {}): Request {
  return new Request(`https://shashwatdixit.com${path}`, { headers });
}

describe("wantsMarkdown", () => {
  it("is false for a normal browser request", () => {
    assert.equal(
      wantsMarkdown(
        request("/", {
          Accept: "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
          "User-Agent": "Mozilla/5.0",
        })
      ),
      false
    );
  });

  it("is true for ?format=md", () => {
    assert.equal(wantsMarkdown(request("/blog/some-post?format=md")), true);
  });

  it("is true for Accept: text/markdown", () => {
    assert.equal(
      wantsMarkdown(request("/", { Accept: "text/markdown; charset=utf-8" })),
      true
    );
  });

  it("does not treat */* as markdown", () => {
    assert.equal(wantsMarkdown(request("/", { Accept: "*/*" })), false);
  });

  it("is true for ChatGPT and GPTBot user agents", () => {
    assert.equal(
      wantsMarkdown(request("/", { "User-Agent": "ChatGPT-User/1.0" })),
      true
    );
    assert.equal(
      wantsMarkdown(request("/", { "User-Agent": "Mozilla/5.0 AppleWebKit/537.36 (compatible; GPTBot/1.2)" })),
      true
    );
  });

  it("does not rewrite Googlebot to markdown", () => {
    assert.equal(
      wantsMarkdown(request("/", { "User-Agent": "Mozilla/5.0 (compatible; Googlebot/2.1)" })),
      false
    );
  });
});

describe("markdownRewritePath", () => {
  const mdReq = request("/", { Accept: "text/markdown" });

  it("rewrites home, blog list, and post URLs", () => {
    assert.equal(markdownRewritePath("/", mdReq), "/index.md");
    assert.equal(markdownRewritePath("/blog", mdReq), "/blog.md");
    assert.equal(
      markdownRewritePath("/blog/django-orm-query-optimization", mdReq),
      "/blog/django-orm-query-optimization.md"
    );
  });

  it("leaves markdown, robots, and asset paths alone", () => {
    assert.equal(markdownRewritePath("/index.md", mdReq), null);
    assert.equal(markdownRewritePath("/llms.txt", mdReq), null);
    assert.equal(markdownRewritePath("/robots.txt", mdReq), null);
    assert.equal(markdownRewritePath("/favicon.ico", mdReq), null);
    assert.equal(markdownRewritePath("/og/blog.png", mdReq), null);
    assert.equal(markdownRewritePath("/blog/tag/code", mdReq), null);
  });

  it("does not rewrite HTML-only browser traffic", () => {
    assert.equal(
      markdownRewritePath("/", request("/", { Accept: "text/html" })),
      null
    );
  });
});

describe("markdown documents", () => {
  const posts = [
    {
      slug: "django-orm-query-optimization",
      title: "Optimizing Django ORM Queries",
      description: "Practical patterns for eliminating N+1 queries.",
      contentMd: `---
title: "Optimizing Django ORM Queries"
slug: django-orm-query-optimization
---

# Body
`,
    },
  ];

  it("builds llms.txt with .md links", () => {
    const body = buildLlmsTxt(
      "https://shashwatdixit.com",
      "Shashwat Dixit",
      "Software engineer.",
      posts
    );
    assert.equal(looksLikeMarkdown(body), true);
    assert.match(body, /^# Shashwat Dixit/m);
    assert.match(body, /https:\/\/shashwatdixit\.com\/index\.md/);
    assert.match(
      body,
      /https:\/\/shashwatdixit\.com\/blog\/django-orm-query-optimization\.md/
    );
  });

  it("keeps post frontmatter in the full dump", () => {
    const homepage = "# Shashwat Dixit\n\nHello.";
    const body = buildLlmsFullTxt(homepage, posts);
    assert.equal(looksLikeMarkdown(body), true);
    assert.equal(looksLikePostMarkdown(posts[0].contentMd ?? ""), true);
    assert.match(body, /^---$/m);
    assert.match(body, /slug: django-orm-query-optimization/);
  });

  it("builds a markdown blog index", () => {
    const body = buildBlogIndexMarkdown("https://shashwatdixit.com", posts);
    assert.equal(looksLikeMarkdown(body), true);
    assert.match(body, /^# Blog/m);
    assert.match(body, /django-orm-query-optimization\.md/);
  });

  it("allows GPTBot in robots.txt and points at llms.txt", () => {
    const body = buildRobotsTxt("https://shashwatdixit.com");
    assert.match(body, /User-agent: GPTBot/);
    assert.match(body, /Allow: \//);
    assert.match(body, /LLMs-Txt: https:\/\/shashwatdixit\.com\/llms\.txt/);
    assert.match(body, /Sitemap: https:\/\/shashwatdixit\.com\/sitemap\.xml/);
  });

  it("includes markdown alternates in the sitemap", () => {
    const xml = buildSitemapXml("https://shashwatdixit.com", [
      "django-orm-query-optimization",
    ]);
    assert.match(xml, /<loc>https:\/\/shashwatdixit\.com\/blog\/django-orm-query-optimization<\/loc>/);
    assert.match(
      xml,
      /type="text\/plain" href="https:\/\/shashwatdixit\.com\/blog\/django-orm-query-optimization\.md"/
    );
  });

  it("serves agent markdown as text/plain unless markdown is explicitly accepted", () => {
    const headers = new Headers(plainTextHeaders());
    assert.equal(headers.get("Content-Type"), "text/plain; charset=utf-8");
    assert.equal(
      new Headers(markdownHeaders()).get("Content-Type"),
      "text/plain; charset=utf-8"
    );
    assert.equal(
      new Headers(
        markdownHeaders(request("/blog/slug.md", { Accept: "text/markdown" }))
      ).get("Content-Type"),
      "text/markdown; charset=utf-8"
    );
  });
});
