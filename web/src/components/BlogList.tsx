import { useCallback, useEffect, useRef, useState } from "react";
import BlurFade from "@/components/magicui/blur-fade";
import { CONFIG } from "@/data/config";
import type { PostSummary } from "@/lib/api";
import { ChevronRight } from "lucide-react";

const BLUR_FADE_DELAY = 0.04;

interface Post {
  id: string;
  title: string;
  publishedAt: string;
  isDraft?: boolean;
}

interface Pagination {
  page: number;
  totalPages: number;
  hasNextPage: boolean;
  hasPreviousPage: boolean;
}

interface BlogListProps {
  posts: Post[];
  allPostsCount: number;
  pagination: Pagination;
  pageSize: number;
  tag?: string;
}

function formatPost(post: PostSummary): Post {
  return {
    id: post.slug,
    title: post.title,
    publishedAt:
      post.status === "draft"
        ? "Coming Soon"
        : new Date(post.date).toLocaleDateString("en-US", {
            year: "numeric",
            month: "long",
            day: "numeric",
          }),
    isDraft: post.status === "draft",
  };
}

export default function BlogList({
  posts: initialPosts,
  allPostsCount,
  pagination: initialPagination,
  pageSize,
  tag,
}: BlogListProps) {
  const [posts, setPosts] = useState<Post[]>(initialPosts);
  const [page, setPage] = useState(initialPagination.page);
  const [hasNextPage, setHasNextPage] = useState(initialPagination.hasNextPage);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const sentinelRef = useRef<HTMLDivElement | null>(null);
  const loadingRef = useRef(false);

  const listSignature = `${tag ?? ""}:${initialPagination.page}:${initialPosts.map((p) => p.id).join(",")}`;

  // Reset when the server-rendered list changes (e.g. tag navigation via ClientRouter).
  useEffect(() => {
    setPosts(initialPosts);
    setPage(initialPagination.page);
    setHasNextPage(initialPagination.hasNextPage);
    setError(null);
  }, [listSignature, initialPosts, initialPagination.page, initialPagination.hasNextPage]);

  const loadMore = useCallback(async () => {
    if (loadingRef.current || !hasNextPage) return;

    loadingRef.current = true;
    setIsLoading(true);
    setError(null);

    const nextPage = page + 1;

    try {
      const params = new URLSearchParams({
        page: String(nextPage),
        limit: String(pageSize),
        drafts: "true",
      });
      if (tag) params.set("tag", tag);

      const res = await fetch(`${CONFIG.api.baseUrl}/api/posts?${params}`);
      if (!res.ok) {
        throw new Error(`Failed to load posts (${res.status})`);
      }

      const data = await res.json();
      const nextPosts: Post[] = (data.posts || []).map(formatPost);
      const totalPages = data.pagination?.totalPages ?? nextPage;

      setPosts((prev) => {
        const seen = new Set(prev.map((p) => p.id));
        return [...prev, ...nextPosts.filter((p) => !seen.has(p.id))];
      });
      setPage(nextPage);
      setHasNextPage(nextPage < totalPages);
    } catch {
      setError("Couldn't load more posts. Try again.");
    } finally {
      setIsLoading(false);
      loadingRef.current = false;
    }
  }, [hasNextPage, page, pageSize, tag]);

  useEffect(() => {
    const sentinel = sentinelRef.current;
    if (!sentinel || !hasNextPage || error) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries.some((entry) => entry.isIntersecting)) {
          void loadMore();
        }
      },
      { rootMargin: "200px 0px" }
    );

    observer.observe(sentinel);
    return () => observer.disconnect();
  }, [hasNextPage, error, loadMore]);

  return (
    <section id="blog">
      <BlurFade delay={BLUR_FADE_DELAY}>
        <h1 className="text-2xl font-semibold tracking-tight mb-4">
          Blog{" "}
          <span className="ml-1 bg-card border border-border rounded-md px-2 py-1 text-muted-foreground text-sm">
            {allPostsCount} posts
          </span>
        </h1>
        <p className="text-sm text-muted-foreground mb-8">
          My personal reflections about web development, life, and more.
        </p>
      </BlurFade>

      {posts.length > 0 ? (
        <>
          <BlurFade delay={BLUR_FADE_DELAY * 2}>
            <div className="flex flex-col gap-5">
              {posts.map((post, id) => {
                const Wrapper = post.isDraft ? "div" : "a";
                const wrapperProps = post.isDraft
                  ? {}
                  : { href: `/blog/${post.id}` };
                const staggerDelay =
                  id < initialPosts.length
                    ? BLUR_FADE_DELAY * 3 + id * 0.05
                    : 0;

                return (
                  <BlurFade delay={staggerDelay} key={post.id}>
                    <Wrapper
                      {...wrapperProps}
                      className={`flex items-start group focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 ${
                        post.isDraft ? "opacity-60 cursor-default" : "cursor-pointer"
                      }`}
                    >
                      <div className="flex flex-col gap-y-2 flex-1">
                        <p className="tracking-tight text-lg font-medium">
                          <span className={post.isDraft ? "" : "group-hover:text-foreground transition-colors"}>
                            {post.title}
                            {post.isDraft ? (
                              <span className="ml-2 inline-flex items-center text-[11px] font-medium px-2 py-0.5 rounded-full bg-muted text-muted-foreground border border-border">
                                Coming Soon
                              </span>
                            ) : (
                              <ChevronRight
                                className="ml-1 inline-block size-4 stroke-3 text-muted-foreground opacity-0 -translate-x-2 transition-all duration-200 group-hover:opacity-100 group-hover:translate-x-0"
                                aria-hidden
                              />
                            )}
                          </span>
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {post.publishedAt}
                        </p>
                      </div>
                    </Wrapper>
                  </BlurFade>
                );
              })}
            </div>
          </BlurFade>

          <div className="mt-8 flex flex-col items-center gap-3 min-h-8">
            {hasNextPage && (
              <div
                ref={sentinelRef}
                className="h-1 w-full"
                aria-hidden
              />
            )}

            {isLoading && (
              <div
                className="flex items-center gap-2 text-sm text-muted-foreground"
                role="status"
                aria-live="polite"
              >
                <div className="size-4 border-2 border-muted-foreground border-t-transparent rounded-full animate-spin" />
                Loading more posts…
              </div>
            )}

            {error && (
              <div className="flex flex-col items-center gap-2">
                <p className="text-sm text-muted-foreground">{error}</p>
                <button
                  type="button"
                  onClick={() => void loadMore()}
                  className="h-8 w-fit px-3 flex items-center justify-center text-sm border border-border rounded-lg hover:bg-accent/50 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                >
                  Try again
                </button>
              </div>
            )}

            {!hasNextPage && !isLoading && !error && posts.length > initialPosts.length && (
              <p className="text-sm text-muted-foreground">You've reached the end.</p>
            )}
          </div>
        </>
      ) : (
        <BlurFade delay={BLUR_FADE_DELAY * 2}>
          <div className="flex flex-col items-center justify-center py-12 px-4 border border-border rounded-xl">
            <p className="text-muted-foreground text-center">
              No blog posts yet. Check back soon!
            </p>
          </div>
        </BlurFade>
      )}
    </section>
  );
}
