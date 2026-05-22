import { useState, useEffect, useRef, useCallback } from "react";

interface TOCItem {
  id: string;
  text: string;
  level: number;
}

interface BlogTOCProps {
  contentHtml: string;
}

function parseHeadings(html: string): TOCItem[] {
  const parser = new DOMParser();
  const doc = parser.parseFromString(html, "text/html");
  const headings = doc.querySelectorAll("h2, h3");
  const items: TOCItem[] = [];

  headings.forEach((el) => {
    const id = el.getAttribute("id");
    const text = el.textContent?.trim();
    if (id && text) {
      items.push({
        id,
        text,
        level: parseInt(el.tagName[1], 10),
      });
    }
  });

  return items;
}

export default function BlogTOC({ contentHtml }: BlogTOCProps) {
  const [headings, setHeadings] = useState<TOCItem[]>([]);
  const [activeId, setActiveId] = useState<string>("");
  const observerRef = useRef<IntersectionObserver | null>(null);

  useEffect(() => {
    setHeadings(parseHeadings(contentHtml));
  }, [contentHtml]);

  const handleObserver = useCallback((entries: IntersectionObserverEntry[]) => {
    const visibleEntries = entries.filter((e) => e.isIntersecting);
    if (visibleEntries.length > 0) {
      const topMost = visibleEntries.reduce((a, b) =>
        a.boundingClientRect.top < b.boundingClientRect.top ? a : b
      );
      setActiveId(topMost.target.id);
    }
  }, []);

  useEffect(() => {
    if (headings.length === 0) return;

    observerRef.current = new IntersectionObserver(handleObserver, {
      rootMargin: "-64px 0px -75% 0px",
      threshold: 0,
    });

    headings.forEach(({ id }) => {
      const el = document.getElementById(id);
      if (el) observerRef.current?.observe(el);
    });

    return () => observerRef.current?.disconnect();
  }, [headings, handleObserver]);

  const handleClick = (e: React.MouseEvent<HTMLAnchorElement>, id: string) => {
    e.preventDefault();
    const el = document.getElementById(id);
    if (el) {
      el.scrollIntoView({ behavior: "smooth", block: "start" });
      setActiveId(id);
    }
  };

  if (headings.length === 0) return null;

  return (
    <nav className="hidden xl:block fixed top-1/2 -translate-y-1/2 left-[max(1rem,calc(50%-400px-14rem))] w-52">
      <ul className="flex flex-col gap-1 text-[13px] leading-snug">
        {headings.map((h) => (
          <li key={h.id}>
            <a
              href={`#${h.id}`}
              onClick={(e) => handleClick(e, h.id)}
              className={`
                block py-1 transition-colors duration-200 border-l-2
                ${h.level === 3 ? "pl-5" : "pl-3"}
                ${
                  activeId === h.id
                    ? "border-primary text-foreground font-medium"
                    : "border-transparent text-muted-foreground hover:text-foreground hover:border-muted-foreground/40"
                }
              `}
            >
              {h.text}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
}
