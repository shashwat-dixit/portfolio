import BlurFade from "@/components/magicui/blur-fade";

interface Tag {
  id: number;
  name: string;
  slug: string;
  count: number;
}

interface TagListProps {
  tags: Tag[];
  activeTag?: string;
}

export default function TagList({ tags, activeTag }: TagListProps) {
  return (
    <BlurFade delay={0.04}>
      <div className="flex flex-wrap gap-2">
        <a
          href="/blog"
          className={`text-xs px-3 py-1.5 rounded-full border transition-colors ${
            !activeTag
              ? "bg-primary text-primary-foreground border-primary"
              : "border-border text-muted-foreground hover:bg-accent hover:text-accent-foreground"
          }`}
        >
          All
        </a>
        {tags.map((tag) => (
          <a
            key={tag.id}
            href={`/blog/tag/${tag.slug}`}
            className={`text-xs px-3 py-1.5 rounded-full border transition-colors inline-flex items-center gap-1.5 ${
              activeTag === tag.slug
                ? "bg-primary text-primary-foreground border-primary"
                : "border-border text-muted-foreground hover:bg-accent hover:text-accent-foreground"
            }`}
          >
            {tag.name}
            <span className="opacity-60">({tag.count})</span>
          </a>
        ))}
      </div>
    </BlurFade>
  );
}
