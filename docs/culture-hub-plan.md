# Culture Hub Plan

Saved so this can be designed and built later without losing the shape of the idea.

## Why this exists

The site currently has two tents:

1. **Career home** (`/`) — about, work, education, skills, projects, contact
2. **Blog** (`/blog`) — long-form writing synced from markdown

Personal culture does not fit either cleanly:

- watching lots of movies → Letterboxd
- eating out and writing about it → Google Maps + blog essays
- art interest → not yet a first-class surface
- books → coming slowly

Forcing Letterboxd ratings, Maps pins, and book shelves into blog tags will feel like clutter. Those are different content shapes from essays.

## The big tent

Create a **third surface**: a Culture hub.

- One nav item, not four dock icons
- Suggested routes: `/life` or `/taste` (name TBD — something personal like `Life`, `Notes`, or `Atelier`)
- Homepage stays career-forward; culture gets a small “Lately” teaser at most

### Shelves inside the hub

| Shelf | What lives there | Source of truth |
| --- | --- | --- |
| **Writing** | essays, restaurant writeups, film notes | existing blog |
| **Films** | ratings, short takes, watchlist | Letterboxd (embed / RSS / curated picks) |
| **Places** | restaurants, trips, maps | Google Maps + blog posts |
| **Books** | shelf + notes (when ready) | static list first → later Goodreads/etc. |

Suggested routes later:

- `/life` — hub overview
- `/life/films`
- `/life/places`
- `/life/books`
- Writing can deep-link into `/blog` filtered by tags like `film`, `food`, `books`

## The tent rule

- **Short logs** live in each shelf (rating, pin, shelf entry)
- **Blog** is only for when there is something longer to say
- Cross-link, do not duplicate:
  - film page → “wrote about this” blog post
  - Maps pin → restaurant essay
  - book note → longer review when it exists

## How this fits the current stack

Relevant existing pieces:

- Nav / IA: `web/src/data/resume.tsx` (`navbar`, `sections`)
- Homepage composition: `web/src/components/HomePage.tsx`, `web/src/components/section/*`
- Blog API + SSR pages: `web/src/lib/api.ts`, `web/src/pages/blog/*`
- Backend post sync is markdown-post-specific; do not stretch it for ratings/pins yet
- Disabled travel photos section already exists as a soft precedent (`sections.photos`, `enabled: false`)

### Recommended build approach

1. **IA first**
   - Add one Culture nav entry
   - Add `/life` (or chosen name) with four sub-hub stubs
2. **Homepage restraint**
   - Optional “Lately” strip: latest film, place, post
   - Do not dump culture into the career narrative
3. **Start thin**
   - Curated JSON / manual picks in `web/src/data/`
   - Outbound links to Letterboxd + Google Maps profiles
   - Reuse blog posts via tags for writing
4. **Automate later only if volume hurts**
   - Letterboxd RSS/API sync
   - Maps export / manual review log
   - Books import

## Content model sketch (later)

Keep hubs typed and separate from `Post` until there is a real need for a shared backend.

Possible lightweight local data shapes:

```ts
// films
{ title, year?, rating?, letterboxdUrl, note?, watchedAt?, relatedPostSlug? }

// places
{ name, city?, mapsUrl, cuisine?, note?, visitedAt?, relatedPostSlug? }

// books
{ title, author, status: "reading" | "read" | "want", note?, relatedPostSlug? }
```

Blog remains the essay engine. Hubs remain the catalogs.

## Open decisions

- Final hub name and URL (`/life` vs `/taste` vs something else)
- Launch with Films + Places first, or all four stubs immediately
- How prominently culture appears on the homepage
- Whether art becomes its own shelf or an ambient theme across the hub
- Whether photos/travel (`sections.photos`) folds into Places

## Suggested first implementation slice

When ready to build:

1. Pick the name
2. Add nav item + `/life` overview page
3. Stub Films / Places / Books pages with profile links + empty states
4. Tag existing blog posts (`film`, `food`, etc.) and surface recent writing on the hub
5. Add a tiny curated “currently into” list by hand
6. Only then consider syncing Letterboxd/Maps

## Non-goals for v1

- Replacing Letterboxd or Google Maps
- Building a full review CMS
- Merging culture content into the career homepage sections
- Overloading the existing markdown sync pipeline with non-essay content
