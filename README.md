# Portfolio — Shashwat Dixit

Personal portfolio and blog, built with Astro 6, React 19, Tailwind CSS 4, and shadcn/ui.

## Tech Stack

| Layer | Choice |
| --- | --- |
| Framework | [Astro v6](https://astro.build) — static site generator |
| UI / Interactivity | [React 19](https://react.dev) via Astro islands (`client:load`, `client:only`) |
| Content | [MDX](https://mdxjs.com/) with content collections, remark-gfm, rehype-pretty-code + Shiki |
| Styling | [Tailwind CSS v4](https://tailwindcss.com), shadcn/ui primitives (Radix), CVA, tailwind-merge |
| Fonts | Fontsource variable — Outfit (sans), Geist Mono (mono) |
| Animation | [Motion](https://motion.dev/) (Framer Motion successor) |
| Theming | next-themes (light/dark, system detection) |
| Icons | lucide-react, @radix-ui/react-icons, custom SVGs |
| Package Manager | pnpm |
| Node | >= 22.12.0 |
| Deployment | TBD — no adapter configured; currently builds to static |

## Architecture

### Data-Driven Design

The entire portfolio is driven by two data files. Components never contain hardcoded personal data — they read from the exported objects:

```
src/data/
├── resume.tsx   ← identity, work, education, projects, skills, achievements, contact
└── config.ts    ← site URL, SEO, theme colors, typography, blog settings
```

`resume.tsx` exports a single `DATA` object consumed by every page and section component.
`config.ts` exports a single `CONFIG` object consumed by the layout, Astro config, and SEO tags.

### Directory Structure

```
web/
├── astro.config.mjs          # Astro config (site URL, integrations, markdown pipeline)
├── components.json            # shadcn/ui component registry config
├── package.json
├── tsconfig.json
├── public/                    # Static assets (avatar, OG image, photos, project screenshots)
│   └── favicon.svg
└── src/
    ├── content.config.ts      # Content collection schema (blog posts)
    ├── middleware.ts           # Security headers (active only in server mode)
    ├── content/
    │   └── blog/              # MDX blog posts
    │       ├── building-restful-apis.mdx
    │       └── ...
    ├── data/
    │   ├── config.ts          # Site settings, SEO, theme
    │   └── resume.tsx         # All personal/portfolio data
    ├── components/
    │   ├── HomePage.tsx        # Main page — renders all sections from DATA
    │   ├── NavbarIsland.tsx    # Navigation island (React)
    │   ├── navbar.tsx          # Navbar layout
    │   ├── icons.tsx           # Icon registry (GitHub, LinkedIn, etc.)
    │   ├── project-card.tsx    # Project card component
    │   ├── timeline.tsx        # Work/education/hackathon timeline
    │   ├── mode-toggle.tsx     # Light/dark theme toggle
    │   ├── theme-provider.tsx  # next-themes provider wrapper
    │   ├── section/            # Page sections (contact, work, projects, hackathons, photos)
    │   ├── magicui/            # Animated UI primitives (blur-fade, dock, flickering-grid)
    │   ├── mdx/                # MDX rendering components (code blocks, media)
    │   └── ui/                 # shadcn/ui primitives + custom SVG icons
    ├── layouts/
    │   └── Layout.astro        # HTML shell — meta tags, OG, theme injection
    ├── lib/
    │   ├── utils.ts            # cn() helper (clsx + tailwind-merge)
    │   ├── pagination.ts       # Blog pagination logic
    │   └── remark-code-meta.ts # Custom remark plugin for code block metadata
    ├── pages/
    │   ├── index.astro         # Homepage — renders <HomePage /> island
    │   ├── 404.astro
    │   └── blog/
    │       ├── index.astro     # Blog listing page
    │       └── [slug].astro    # Individual blog post (content collection)
    └── styles/
        └── global.css          # Font imports, Tailwind base, theme CSS variables
```

### Page Rendering Flow

```
index.astro
  └─ Layout.astro              # Reads CONFIG + DATA for meta/OG tags
       └─ <HomePage />          # React island (client:only="react")
            ├─ About section     # DATA.summary rendered via react-markdown
            ├─ Work section      # DATA.work[] → Timeline component
            ├─ Education section # DATA.education[] → Timeline component
            ├─ Skills section    # DATA.skills[] → icon grid
            ├─ Projects section  # DATA.projects[] → ProjectCard grid
            ├─ Photos section    # DATA.photos[] → image gallery (disabled)
            ├─ Achievements      # DATA.hackathons[] → Timeline component
            └─ Contact section   # DATA.contact → social links + dock
```

### Blog Pipeline

Blog posts are `.mdx` files in `src/content/blog/`. The content collection schema (`content.config.ts`) validates frontmatter with Zod. Posts are rendered by `[slug].astro` which generates JSON-LD structured data using `DATA.name` and `CONFIG.site.url`.

Markdown is processed through: remark-gfm → remark-code-meta → rehype-pretty-code (Shiki, dual GitHub themes).

## Quick Start

```bash
cd web
pnpm install
pnpm dev
```

Open http://localhost:4321.

## Customizing Content

| File | Controls |
| --- | --- |
| `src/data/resume.tsx` | Name, bio, work, education, projects, skills, achievements, contact |
| `src/data/config.ts` | Site URL, SEO, theme colors, typography |
| `src/content/blog/*.mdx` | Blog posts |

## Commands

| Command | Action |
| --- | --- |
| `pnpm install` | Install dependencies |
| `pnpm dev` | Start dev server at localhost:4321 |
| `pnpm build` | Build for production (static output to `dist/`) |
| `pnpm preview` | Preview production build locally |

## Deployment

No deployment adapter is currently configured. The site builds to static HTML in `dist/`.

To deploy, add an Astro adapter for your platform of choice:
- **Vercel**: `@astrojs/vercel`
- **Netlify**: `@astrojs/netlify`
- **Cloudflare**: `@astrojs/cloudflare`
- **Node**: `@astrojs/node`
- **Static hosting** (GitHub Pages, S3, etc.): no adapter needed, just serve `dist/`

See [Astro deployment docs](https://docs.astro.build/en/guides/deploy/) for details.

## Manual Data TODO

The following items need to be provided manually (not extractable from the resume):

- [ ] **Avatar image** — add personal photo as `web/public/picofme.png`
- [ ] **OG image** — add or generate `web/public/og_image.png` for social sharing
- [ ] **GitHub profile URL** — update in `resume.tsx` → `contact.social.GitHub.url`
- [ ] **LinkedIn profile URL** — update in `resume.tsx` → `contact.social.LinkedIn.url`
- [ ] **Twitter/X handle** — if applicable, update in `resume.tsx` and `config.ts`
- [ ] **Personal domain** — update `url` in `resume.tsx` and `site.url` in `config.ts`
- [ ] **Project live URLs** — add URLs for Jamin and Zort in `resume.tsx` → `projects[].href`
- [ ] **Project source code URLs** — add GitHub repo links in `resume.tsx` → `projects[].links[]`
- [ ] **Project images/screenshots** — add to `web/public/` and reference in `resume.tsx`
- [ ] **Company logos** — optionally replace favicon-based logos for Instahyre and Pummyz Foods
- [ ] **Photo gallery** — add photos to `web/public/photos/` and enable photos section in `resume.tsx`
- [ ] **Blog posts** — current posts are template placeholders; write your own or remove them
