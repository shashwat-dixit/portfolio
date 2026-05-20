import { useState, useEffect, useMemo } from "react";
import Fuse from "fuse.js";
import { CONFIG } from "@/data/config";
import type { PostSummary } from "@/lib/api";
import BlurFade from "@/components/magicui/blur-fade";

export default function BlogSearch() {
  const [query, setQuery] = useState("");
  const [posts, setPosts] = useState<PostSummary[]>([]);
  const [results, setResults] = useState<PostSummary[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const fetchPosts = async () => {
      setIsLoading(true);
      try {
        const res = await fetch(`${CONFIG.api.baseUrl}/api/posts?limit=100`);
        if (res.ok) {
          const data = await res.json();
          setPosts(data.posts || []);
        }
      } catch {
        // Silently fail - search just won't work
      } finally {
        setIsLoading(false);
      }
    };
    fetchPosts();
  }, []);

  const fuse = useMemo(
    () =>
      new Fuse(posts, {
        keys: [
          { name: "title", weight: 0.4 },
          { name: "description", weight: 0.3 },
          { name: "tags", weight: 0.3 },
        ],
        threshold: 0.4,
        includeScore: true,
      }),
    [posts]
  );

  useEffect(() => {
    if (!query.trim()) {
      setResults([]);
      return;
    }
    const searchResults = fuse.search(query).map((r) => r.item);
    setResults(searchResults);
  }, [query, fuse]);

  const showResults = query.trim().length > 0 && results.length > 0;

  return (
    <BlurFade delay={0.02}>
      <div className="relative">
        <div className="relative">
          <svg
            className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          >
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.3-4.3" />
          </svg>
          <input
            type="text"
            placeholder="Search posts..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 text-sm border border-border rounded-lg bg-background text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 transition-shadow"
          />
          {isLoading && (
            <div className="absolute right-3 top-1/2 -translate-y-1/2">
              <div className="size-4 border-2 border-muted-foreground border-t-transparent rounded-full animate-spin" />
            </div>
          )}
        </div>
        {showResults && (
          <div className="absolute top-full left-0 right-0 mt-2 bg-card border border-border rounded-lg shadow-lg z-50 max-h-80 overflow-y-auto">
            {results.map((post) => (
              <a
                key={post.slug}
                href={`/blog/${post.slug}`}
                className="flex flex-col gap-1 px-4 py-3 hover:bg-accent/50 transition-colors border-b border-border last:border-b-0"
              >
                <span className="text-sm font-medium">{post.title}</span>
                {post.description && (
                  <span className="text-xs text-muted-foreground line-clamp-1">
                    {post.description}
                  </span>
                )}
              </a>
            ))}
          </div>
        )}
        {query.trim().length > 0 && results.length === 0 && !isLoading && (
          <div className="absolute top-full left-0 right-0 mt-2 bg-card border border-border rounded-lg shadow-lg z-50 p-4">
            <p className="text-sm text-muted-foreground text-center">
              No posts found for "{query}"
            </p>
          </div>
        )}
      </div>
    </BlurFade>
  );
}
