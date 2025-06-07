# Portfolio + Personal Blog System

A Markdown + Go + Next.js-based blog platform with likes, comments, and strong developer ergonomics. Built for flexibility, performance, and long-term ownership.

---

## 📁 Project Structure

```
portfolio/
├── backend/
├── next-app/            # Markdown blog posts (linked to Obsidian vault)
├── scripts/             # Deploy/backup utilities
└── README.md
```

---

## ✅ Tech Stack Overview

| Layer             | Technology                       | Status |
|------------------|----------------------------------|--------|
| Frontend          | Next.js + Tailwind              | [ ] |
| Data Fetching     | TanStack Query                  | [ ] |
| Backend           | Go (REST API)                   | [ ] |
| DB                | SQLite                          | [ ] |
| Image Hosting     | AWS S3                          | [ ] |
| Blog Content      | GitHub Obsidian Vault           | [ ] |
| RSS Feed          | Markdown → XML generation       | [ ] |
| CDN               | Not used                        | [ ] |
| Caching           | Enabled (local or nginx)        | [ ] |
| TUI (optional)    | BubbleTea (Go)                  | [ ] |
| Rate Limiting     | In Go backend                   | [ ] |
| Monitoring        | Basic logging                   | [ ] |
| Backups           | Weekly SQLite → GitHub          | [ ] |

---

## 🛠️ Backend (Go) Features

| Feature                              | Status |
|--------------------------------------|--------|
| REST API for likes/comments          | [ ] |
| SQLite schema for posts + feedback   | [ ] |
| Rate limiting per IP                 | [ ] |
| Markdown metadata parser             | [ ] |
| RSS Feed generator                   | [ ] |
| GitHub webhook endpoint (`/sync`)    | [ ] |
| Markdown sync logic (git pull/parse) | [ ] |
| Dockerfile                           | [ ] |
| Systemd/runner script                | [ ] |

---

## 🎯 Frontend (Next.js)

| Feature                              | Status |
|--------------------------------------|--------|
| Static blog post rendering           | [ ] |
| MDX support                          | [ ] |
| Likes/comments (API connected)       | [ ] |
| TanStack Query integration           | [ ] |
| RSS link                             | [ ] |
| S3 image rendering                   | [ ] |
| `.env.local` config for backend URL  | [ ] |
| Vercel deployment (initial)          | [ ] |
| `next export` (optional portability) | [ ] |

---

## 🔁 GitHub Webhook Sync

| Task                                        | Status |
|--------------------------------------------|--------|
| Webhook trigger on `push`                  | [ ] |
| Go backend endpoint to receive webhook     | [ ] |
| Git pull & post parsing                    | [ ] |
| Slug/title extraction from Markdown files  | [ ] |
| Update SQLite with new post metadata       | [ ] |

---

## 📦 Deployment

| Component       | Method              | Status |
|------------------|---------------------|--------|
| Backend (Go)     | Docker + EC2        | [ ] |
| SQLite backup    | Cron to GitHub repo | [ ] |
| Frontend (Next)  | Vercel              | [ ] |
| Images           | AWS S3              | [ ] |
| Pre-process imgs | Local compression   | [ ] |

---

## 🧪 Local Dev Setup

| Task                                 | Status |
|--------------------------------------|--------|
| Run backend (`:8080`)                | [ ] |
| Run frontend (`:3000`)               | [ ] |
| Connect frontend → backend via `.env`| [ ] |
| Simulate webhook with curl           | [ ] |

---

## ✍️ Blog Authoring Flow

| Step                                  | Status |
|---------------------------------------|--------|
| Write in Obsidian, push to GitHub     | [ ] |
| Webhook triggers post sync in backend | [ ] |
| Frontend uses updated post slugs      | [ ] |
| Comments & likes linked to post slugs | [ ] |

---

## 🔧 Optional Enhancements

| Feature                               | Status |
|----------------------------------------|--------|
| Admin dashboard (views/moderation)     | [ ] |
| Full-text blog search                  | [ ] |
| BubbleTea-based CLI/TUI                | [ ] |
| Auto-optimize images before S3 upload  | [ ] |
| Obsidian preview support via S3 URLs   | [ ] |

---

## 📌 To-Do Highlights

- [ ] Implement GitHub webhook for vault sync
- [ ] Markdown parsing script in Go
- [ ] SQLite weekly backup automation
- [ ] Build script for `next export`
- [ ] Optional: admin tools or CLI

---

## 🧠 Inspiration & Philosophy

Built for:
- Long-term content ownership
- No vendor lock-in (can leave Vercel/SaaS)
- Markdown-first, git-based writing
- Powerful Go backend + static-friendly frontend
---