CREATE TABLE IF NOT EXISTS posts (
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

CREATE TABLE IF NOT EXISTS tags (
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    slug TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS post_tags (
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  INT  REFERENCES tags(id)  ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_posts_status       ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_slug         ON posts(slug);
CREATE INDEX IF NOT EXISTS idx_tags_slug          ON tags(slug);
