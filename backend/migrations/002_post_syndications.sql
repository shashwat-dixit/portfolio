CREATE TABLE IF NOT EXISTS post_syndications (
    post_id      UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    platform     TEXT NOT NULL CHECK (platform IN ('medium', 'substack')),
    remote_id    TEXT NOT NULL,
    remote_url   TEXT,
    content_hash TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    updated_at   TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (post_id, platform)
);

CREATE INDEX IF NOT EXISTS idx_post_syndications_platform ON post_syndications(platform);
