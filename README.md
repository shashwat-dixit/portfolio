# Portfolio — Shashwat Dixit

Personal portfolio and blog. Go backend syncs markdown from GitLab, stores in PostgreSQL with Redis caching, serves via REST API. Astro frontend renders the portfolio and blog at `shashwatdixit.com`.

## System Architecture

```
┌──────────────┐    weekly sync    ┌──────────────────────────────────────┐
│  GitLab Repo │ ───────────────── │           Go Backend (:8080)         │
│  (markdown)  │   git pull        │                                      │
└──────────────┘                   │  ┌────────┐  ┌───────────────────┐   │
                                   │  │ Sync   │  │ REST API          │   │
                                   │  │Service │  │ /api/posts        │   │
                                   │  └───┬────┘  │ /api/posts/:slug  │   │
                                   │      │       │ /api/tags         │   │
                                   │      ▼       │ /api/sync         │   │
                                   │  ┌────────┐  │ /api/feed.xml     │   │
                                   │  │  PG    │  │ /api/health       │   │
                                   │  │  DB    │◄─┤                   │   │
                                   │  └────────┘  └────────┬──────────┘   │
                                   │                       │              │
                                   │  ┌────────┐           │              │
                                   │  │ Redis  │◄──────────┘              │
                                   │  │ Cache  │  cache-aside             │
                                   │  └────────┘                          │
                                   └──────────────┬───────────────────────┘
                                                  │ JSON
                                                  ▼
                                   ┌──────────────────────────────────────┐
                                   │      Astro Frontend (:4321)          │
                                   │                                      │
                                   │  shashwatdixit.com      → portfolio  │
                                   │  shashwatdixit.com/blog → blog list  │
                                   │  shashwatdixit.com/blog/:slug → post │
                                   │                                      │
                                   │  Fuse.js client-side search          │
                                   │  localStorage reading progress       │
                                   └──────────────────────────────────────┘
```

## Tech Stack

| Layer | Choice |
| --- | --- |
| **Backend** | Go 1.26, chi router, pgx/v5, go-redis/v9, goldmark |
| **Database** | PostgreSQL (posts, tags) |
| **Cache** | Redis (response cache, 7-day TTL, flush on sync) |
| **Frontend** | Astro v6 (hybrid SSR), React 19 islands |
| **Styling** | Tailwind CSS v4, shadcn/ui, CVA |
| **Search** | Fuse.js (client-side fuzzy search) |
| **Content** | Markdown in GitLab repo, parsed by backend |
| **Fonts** | Fontsource variable — Outfit (sans), Geist Mono (mono) |
| **Animation** | Motion (Framer Motion successor) |
| **Theming** | next-themes (light/dark, system detection) |

## Directory Structure

```
portfolio/
├── README.md
├── docker-compose.yml
│
├── backend/                          # Go backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go              # entrypoint
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go            # env-based configuration
│   │   ├── handler/
│   │   │   ├── posts.go             # GET /api/posts, GET /api/posts/:slug
│   │   │   ├── tags.go              # GET /api/tags
│   │   │   ├── sync.go              # POST /api/sync
│   │   │   └── feed.go              # GET /api/feed.xml
│   │   ├── middleware/
│   │   │   ├── cache.go             # Redis response caching
│   │   │   └── cors.go              # CORS for frontend
│   │   ├── model/
│   │   │   └── post.go              # Post, Tag structs
│   │   ├── repository/
│   │   │   ├── post_repo.go         # PostgreSQL post queries
│   │   │   └── tag_repo.go          # PostgreSQL tag queries
│   │   ├── service/
│   │   │   ├── sync_service.go      # GitLab clone/pull, parse, upsert
│   │   │   ├── post_service.go      # Post retrieval logic
│   │   │   └── markdown.go          # frontmatter parsing + md→HTML
│   │   └── cache/
│   │       └── redis.go             # Redis client, key patterns, invalidation
│   ├── migrations/
│   │   └── 001_initial.sql          # DDL for posts, tags, post_tags
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
└── web/                              # Astro frontend
    ├── astro.config.mjs
    ├── package.json
    ├── tsconfig.json
    ├── public/                       # Static assets
    └── src/
        ├── data/
        │   ├── resume.tsx            # Portfolio data (identity, work, projects, etc.)
        │   └── config.ts             # Site settings, SEO, theme, API base URL
        ├── components/
        │   ├── HomePage.tsx          # Portfolio homepage sections
        │   ├── BlogList.tsx          # Blog listing with search + tag filter
        │   ├── BlogSearch.tsx        # Fuse.js search input (TODO)
        │   ├── ReadingProgress.tsx   # Progress bar + scroll restore (TODO)
        │   ├── TagList.tsx           # Tag chips with filter links (TODO)
        │   ├── section/              # Homepage sections
        │   ├── magicui/              # Animated UI components
        │   └── ui/                   # shadcn/ui primitives
        ├── hooks/
        │   └── useReadingProgress.ts # localStorage reading position (TODO)
        ├── lib/
        │   ├── api.ts               # Typed fetch wrapper for Go backend (TODO)
        │   ├── utils.ts
        │   └── pagination.ts
        ├── layouts/
        │   └── Layout.astro
        ├── pages/
        │   ├── index.astro           # Portfolio homepage
        │   ├── 404.astro
        │   └── blog/
        │       ├── index.astro       # Blog listing (fetches from API)
        │       ├── [slug].astro      # Blog post (fetches from API)
        │       └── tag/
        │           └── [tag].astro   # Tag-filtered listing (TODO)
        └── styles/
            └── global.css
```

## Blog Frontmatter Schema

Markdown files in the GitLab repo use this YAML frontmatter:

```yaml
---
title: "What I am learning in 2026"
slug: what-i-am-learning-2026
date: 2025-01-25
updated: 2025-02-01            # optional
tags: [code, engineering, personal]
description: "Short summary for listings and SEO"
cover: /images/learning-2026.jpg  # optional
status: published              # published | draft | writing
---
```

The backend auto-generates OG/Twitter/SEO metadata from these fields. `status: writing` posts are ignored during sync. `draft` posts are stored but not served publicly.

## Database Schema

```sql
CREATE TABLE posts (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug         TEXT UNIQUE NOT NULL,
    title        TEXT NOT NULL,
    description  TEXT,
    content_md   TEXT NOT NULL,
    content_html TEXT NOT NULL,
    cover_image  TEXT,
    status       TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('published', 'draft')),
    reading_time INT  NOT NULL DEFAULT 0,
    author       TEXT NOT NULL DEFAULT 'Shashwat Dixit',
    published_at TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ,
    gitlab_sha   TEXT,
    created_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tags (
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    slug TEXT UNIQUE NOT NULL
);

CREATE TABLE post_tags (
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  INT  REFERENCES tags(id)  ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX idx_posts_status       ON posts(status);
CREATE INDEX idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX idx_posts_slug         ON posts(slug);
CREATE INDEX idx_tags_slug          ON tags(slug);
```

## REST API

```
GET  /api/posts                   List published posts
     ?tag=engineering             Filter by tag
     ?page=1&limit=10            Pagination
     ?q=search+term              Search (title + description)
     Response: { posts: [...], pagination: { page, limit, total, totalPages } }

GET  /api/posts/:slug             Single post with full HTML
     Response: { slug, title, description, contentHtml, tags, date, updated, cover, readingTime }

GET  /api/tags                    All tags with post counts
     Response: { tags: [{ name, slug, count }] }

POST /api/sync                    Trigger GitLab sync (protected)
     Header: X-API-Key: <secret>
     Response: { synced, created, updated, deleted }

GET  /api/feed.xml                RSS 2.0 feed

GET  /api/health                  Health check
```

## Redis Caching

| Key Pattern | Content | TTL |
| --- | --- | --- |
| `blog:post:{slug}` | Single post JSON | 7 days |
| `blog:posts:list:{tag}:{page}:{limit}` | Post list JSON | 7 days |
| `blog:tags:all` | Tags list JSON | 7 days |

On `POST /api/sync`, all `blog:*` keys are flushed. Cache-aside pattern: handler checks Redis first, falls back to PG, writes to Redis on miss.

HTTP response headers for CDN caching:

- Single post: `Cache-Control: public, max-age=86400, stale-while-revalidate=604800`
- List/tags: `Cache-Control: public, max-age=3600, stale-while-revalidate=86400`

## GitLab Sync Flow

```
POST /api/sync (protected by X-API-Key)
  │
  ├─ 1. git clone --depth=1 (or git pull) the GitLab blog repo
  ├─ 2. Walk all *.md files
  ├─ 3. For each file:
  │     ├─ Parse YAML frontmatter
  │     ├─ Skip if status == "writing"
  │     ├─ Compute content hash, compare with gitlab_sha in DB
  │     ├─ If new/changed:
  │     │     ├─ Convert markdown → HTML (goldmark + syntax highlighting)
  │     │     ├─ Calculate reading time (words / 200 wpm)
  │     │     └─ Upsert post + tags in PostgreSQL
  │     └─ If unchanged: skip
  ├─ 4. Delete DB posts whose slugs no longer exist in repo
  ├─ 5. Flush all Redis blog:* keys
  └─ 6. Return sync summary
```

Automated via cron (Tuesday 9:00 AM IST / 3:30 AM UTC):

```
30 3 * * 2 curl -X POST -H "X-API-Key: $SYNC_KEY" https://api.shashwatdixit.com/api/sync
```

## Request Flow

```
Browser → GET shashwatdixit.com/blog/some-post
  │
  ├─ Caddy/Nginx reverse proxy → Astro (:4321)
  │
  ├─ Astro SSR page → fetch GET /api/posts/some-post → Go backend (:8080)
  │     │
  │     ├─ Redis cache hit? → return cached JSON
  │     └─ Cache miss → query PG → cache in Redis → return JSON
  │
  ├─ Astro renders HTML with post content
  └─ Browser: localStorage saves reading progress on scroll
```

## Frontend: Client-Side Search (Fuse.js)

The blog list page loads a lightweight index of all posts (title, description, tags, slug) from `GET /api/posts?limit=all`. Fuse.js performs fuzzy search client-side — no backend search endpoint needed for the expected post volume. The search index is fetched once and cached in the React component state.

## Frontend: Reading Progress (Device-Level)

No backend involvement. Pure localStorage:

```
localStorage key: "blog-read-progress"
Value: {
  [slug]: {
    scrollPercent: 0-100,
    lastPosition: scrollY pixels,
    contentHeight: document height at save time,
    lastRead: unix timestamp
  }
}
```

- `useReadingProgress(slug)` hook: restores scroll on mount, debounced save on scroll (300ms)
- `ReadingProgress` component: thin fixed bar at top of post page
- Entries older than 30 days auto-cleaned

## Local Development

### Prerequisites

- Docker & Docker Compose
- Node.js ≥ 22
- pnpm

### 1. Configure environment

```bash
cp .env.example .env
```

Fill in the GitLab credentials in `.env`:

```env
GITLAB_REPO=https://gitlab.com/shashwat-dixit/blog.git
GITLAB_TOKEN=glpat-xxxxxxxxxxxxxxxxxxxx
SYNC_API_KEY=any-secret-string-for-local
```

The remaining values (PostgreSQL, Redis, CORS) are overridden by `docker-compose.local.yml` for local development — no need to change them.

### 2. Start the backend

```bash
docker compose -f docker-compose.local.yml up -d
```

This starts three containers:

| Service    | Port   | Details                          |
| ---------- | ------ | -------------------------------- |
| PostgreSQL | `5432` | User `postgres`, password `postgres`, DB `portfolio` |
| Redis      | `6379` | No password                      |
| Backend    | `8080` | Go API server                    |

The backend auto-runs database migrations on startup.

### 3. Sync blog posts

Pull markdown posts from the GitLab repo into the local database:

```bash
curl -X POST http://localhost:8080/api/sync \
  -H "X-API-Key: <your SYNC_API_KEY from .env>"
```

You should see a response like:

```json
{"synced":97,"created":97,"updated":0,"deleted":0}
```

Re-run this command anytime you want to pull the latest posts.

### 4. Start the frontend

```bash
cd web
pnpm install
pnpm dev
```

Open <http://localhost:4321>. The frontend fetches data from the backend at `http://localhost:8080`.

### Verify everything works

- **Blog listing**: <http://localhost:4321/blog>
- **Health check**: <http://localhost:8080/api/health>
- **API posts**: <http://localhost:8080/api/posts>

### Stopping

```bash
docker compose -f docker-compose.local.yml down     # stop and remove containers
docker compose -f docker-compose.local.yml down -v   # also delete database volume
```

### Production (Docker)

The production `docker-compose.yml` runs only the backend (PostgreSQL and Redis are external services). Caddy handles reverse proxy and TLS:

```bash
docker compose up -d
```

## Publishing Workflow

1. Write markdown in GitLab repo (`status: writing`)
2. When ready: set `status: published`, set `date`, commit and push
3. Tuesday 9 AM IST: cron triggers `POST /api/sync`
4. Backend syncs: pulls repo, parses, upserts DB, flushes Redis cache
5. Blog is live immediately (Astro SSR fetches fresh data)

## Deployment

```
┌─────────────────────────────────────────────────┐
│                   VPS / Cloud                    │
│                                                  │
│  Caddy (reverse proxy, auto-TLS)                 │
│    ├─ shashwatdixit.com     → Astro (:4321)      │
│    └─ api.shashwatdixit.com → Go backend (:8080) │
│                                                  │
│  PostgreSQL (:5432)                              │
│  Redis (:6379)                                   │
│                                                  │
│  All services via docker-compose                 │
└─────────────────────────────────────────────────┘
```

## TODO

### Portfolio (manual data)

- [ ] Avatar image — add photo as `web/public/picofme.png`
- [ ] OG image — add or generate `web/public/og_image.png`
- [ ] GitHub profile URL — update in `resume.tsx` → `contact.social.GitHub.url`
- [ ] LinkedIn profile URL — update in `resume.tsx` → `contact.social.LinkedIn.url`
- [ ] Twitter/X handle — if applicable, update in `resume.tsx` and `config.ts`
- [ ] Personal domain — update `url` in `resume.tsx` and `site.url` in `config.ts`
- [ ] Project live URLs — Jamin and Zort in `resume.tsx` → `projects[].href`
- [ ] Project source code URLs — GitHub repos in `resume.tsx` → `projects[].links[]`
- [ ] Project images — add to `web/public/`, reference in `resume.tsx`
- [ ] Company logos — replace favicon-based logos for Instahyre and Pummyz Foods
- [ ] Photo gallery — add photos to `web/public/photos/`, enable section in `resume.tsx`

### Backend

- [x] Config loader — env vars for PG, Redis, GitLab token, API key, port
- [x] PostgreSQL connection pool (pgx)
- [x] Redis client (go-redis)
- [x] Database migrations runner
- [x] Post repository — CRUD + list with tag filter + pagination
- [x] Tag repository — CRUD + counts
- [x] Markdown service — frontmatter parsing (go-yaml) + goldmark HTML rendering
- [x] Sync service — git clone/pull, walk .md files, diff + upsert
- [x] POST /api/sync handler (API key protected)
- [x] GET /api/posts handler (pagination, tag filter)
- [x] GET /api/posts/:slug handler
- [x] GET /api/tags handler
- [x] GET /api/feed.xml handler (RSS 2.0)
- [x] GET /api/health handler
- [x] Redis cache-aside middleware (blog:* keys, 7-day TTL)
- [x] CORS middleware
- [x] HTTP cache headers (Cache-Control, ETag)
- [x] Dockerfile
- [x] .env.example

### Frontend

- [x] Switch Astro to hybrid output + Node adapter
- [x] Add `api.baseUrl` to `web/src/data/config.ts`
- [x] `web/src/lib/api.ts` — typed fetch wrapper for Go backend
- [x] Rewrite `blog/index.astro` to fetch from API instead of content collections
- [x] Rewrite `blog/[slug].astro` to fetch from API
- [x] Add `blog/tag/[tag].astro` — tag-filtered listing
- [x] `BlogSearch.tsx` — Fuse.js client-side fuzzy search on blog list
- [x] `TagList.tsx` — tag chips with filter links
- [x] `ReadingProgress.tsx` — progress bar at top of post
- [x] `useReadingProgress.ts` — localStorage scroll position hook
- [x] Remove placeholder .mdx blog posts from `src/content/blog/`
- [x] Remove `content.config.ts` (content collections no longer needed)
- [x] Convert project grid into carousel
- [x] Remove the achievements section instead use that timeline to show the current work experience
- [x] Display Draft Blogs as Coming Soon
- [X] Render the headings and subheading as a table of content on the left side

### Misc

- [x] Integrate Posthog for analytics
- [ ] Add a newsletter functionality
- [ ] Cross-post to Medium (REST API integration in sync service)
- [ ] Cross-post to Substack (RSS feed import or API when available)
- [ ] Cache the landing page and utilize bf cache when navigating back from blog
- [x] Fix the click on dark mode button

### DevOps

- [x] `docker-compose.yml` — Go backend, Astro, PostgreSQL, Redis
- [x] Caddy / Nginx reverse proxy config
- [x] Option To Trigger Deploy and Trigger Gitlab Repo Pull (via POST /api/sync + GitLab CI)
- [x] Cron job for weekly sync (GitLab scheduled pipeline)
- [x] CI/CD pipeline (GitLab CI)
