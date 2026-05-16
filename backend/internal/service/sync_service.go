package service

import (
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/cache"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/config"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/repository"
)

type SyncService struct {
	cfg      *config.Config
	postRepo *repository.PostRepo
	tagRepo  *repository.TagRepo
	md       *MarkdownService
	cache    *cache.RedisCache
}

func NewSyncService(
	cfg *config.Config,
	postRepo *repository.PostRepo,
	tagRepo *repository.TagRepo,
	md *MarkdownService,
	cache *cache.RedisCache,
) *SyncService {
	return &SyncService{cfg: cfg, postRepo: postRepo, tagRepo: tagRepo, md: md, cache: cache}
}

// Sync pulls the GitLab repo, parses all markdown files, and upserts into the DB.
//
// Flow:
//  1. git clone --depth=1 (or git pull if already cloned) into a temp dir
//  2. Walk all *.md files
//  3. For each file:
//     a. Parse YAML frontmatter via MarkdownService.Parse()
//     b. Skip if status == "writing"
//     c. Compute content hash, compare with gitlab_sha in DB
//     d. If new/changed: render HTML, calc reading time, upsert post + tags
//  4. Delete posts whose slugs no longer exist in repo
//  5. Flush all Redis blog:* keys
//  6. Return SyncResult
//
// TODO: implement
