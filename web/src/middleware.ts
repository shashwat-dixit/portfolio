import { defineMiddleware } from "astro:middleware";
import { markdownRewritePath } from "@/lib/ai-content";

export const onRequest = defineMiddleware(async (ctx, next) => {
  const rewriteTo = ctx.isPrerendered
    ? null
    : markdownRewritePath(ctx.url.pathname, ctx.request);
  const response = rewriteTo ? await ctx.rewrite(rewriteTo) : await next();

  response.headers.set("X-Content-Type-Options", "nosniff");
  response.headers.set("X-Frame-Options", "DENY");
  response.headers.set("Referrer-Policy", "strict-origin-when-cross-origin");
  response.headers.set("Permissions-Policy", "camera=(), microphone=(), geolocation=()");
  if (rewriteTo) {
    response.headers.set("Vary", "Accept, User-Agent");
  }
  return response;
});
