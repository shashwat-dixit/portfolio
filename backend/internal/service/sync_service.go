package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/cache"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/config"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"
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

func (s *SyncService) Sync(ctx context.Context) (*model.SyncResult, error) {
	repoDir, err := s.cloneOrPull(ctx)
	if err != nil {
		return nil, fmt.Errorf("git sync: %w", err)
	}

	mdFiles, err := s.walkMarkdown(repoDir)
	if err != nil {
		return nil, fmt.Errorf("walk markdown: %w", err)
	}

	shaMap, err := s.postRepo.GetSHAMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("get sha map: %w", err)
	}

	result := &model.SyncResult{}
	seenSlugs := make(map[string]bool)

	for _, path := range mdFiles {
		raw, err := os.ReadFile(path)
		if err != nil {
			slog.Warn("read file failed", "path", path, "error", err)
			continue
		}

		fm, htmlContent, wordCount, err := s.md.Parse(raw)
		if err != nil {
			slog.Warn("parse markdown failed", "path", path, "error", err)
			continue
		}

		if fm.Status == "writing" {
			continue
		}

		if fm.Slug == "" {
			base := filepath.Base(path)
			fm.Slug = strings.TrimSuffix(base, filepath.Ext(base))
		}

		seenSlugs[fm.Slug] = true
		contentHash := fmt.Sprintf("%x", sha256.Sum256(raw))

		if existingSHA, exists := shaMap[fm.Slug]; exists && existingSHA == contentHash {
			result.Synced++
			continue
		}

		post := &model.Post{
			Slug:        fm.Slug,
			Title:       fm.Title,
			Description: fm.Description,
			ContentMD:   string(raw),
			ContentHTML:  htmlContent,
			CoverImage:  fm.Cover,
			Status:      fm.Status,
			ReadingTime: (wordCount + 199) / 200,
			Author:      "Shashwat Dixit",
			GitLabSHA:   contentHash,
		}

		if fm.Date != "" {
			if t, err := time.Parse("2006-01-02", fm.Date); err == nil {
				post.PublishedAt = t
			}
		}
		if fm.Updated != "" {
			if t, err := time.Parse("2006-01-02", fm.Updated); err == nil {
				post.UpdatedAt = t
			}
		}
		if post.Status == "" {
			post.Status = "draft"
		}

		postID, err := s.postRepo.Upsert(ctx, post)
		if err != nil {
			slog.Warn("upsert post failed", "slug", fm.Slug, "error", err)
			continue
		}

		if len(fm.Tags) > 0 {
			tagIDs, err := s.tagRepo.UpsertMany(ctx, fm.Tags)
			if err != nil {
				slog.Warn("upsert tags failed", "slug", fm.Slug, "error", err)
			} else {
				if err := s.postRepo.SetTags(ctx, postID, tagIDs); err != nil {
					slog.Warn("set tags failed", "slug", fm.Slug, "error", err)
				}
			}
		}

		if _, exists := shaMap[fm.Slug]; exists {
			result.Updated++
		} else {
			result.Created++
		}
		result.Synced++
	}

	allSlugs, err := s.postRepo.AllSlugs(ctx)
	if err != nil {
		slog.Warn("get all slugs failed", "error", err)
	} else {
		for _, slug := range allSlugs {
			if !seenSlugs[slug] {
				if err := s.postRepo.Delete(ctx, slug); err != nil {
					slog.Warn("delete post failed", "slug", slug, "error", err)
				} else {
					result.Deleted++
				}
			}
		}
	}

	if err := s.cache.FlushBlog(ctx); err != nil {
		slog.Warn("flush cache failed", "error", err)
	}

	return result, nil
}

func (s *SyncService) cloneOrPull(ctx context.Context) (string, error) {
	dir := filepath.Join(os.TempDir(), "portfolio-blog-sync")

	repoURL := s.cfg.GitLabRepo
	if s.cfg.GitLabToken != "" && strings.HasPrefix(repoURL, "https://") {
		repoURL = strings.Replace(repoURL, "https://", fmt.Sprintf("https://oauth2:%s@", s.cfg.GitLabToken), 1)
	}

	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		cmd := exec.CommandContext(ctx, "git", "-C", dir, "pull", "--ff-only")
		cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		if output, err := cmd.CombinedOutput(); err != nil {
			slog.Warn("git pull failed, recloning", "error", err, "output", string(output))
			os.RemoveAll(dir)
		} else {
			return dir, nil
		}
	}

	os.RemoveAll(dir)
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth=1", repoURL, dir)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("git clone: %s: %w", string(output), err)
	}

	return dir, nil
}

func (s *SyncService) walkMarkdown(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(info.Name(), ".md") || strings.HasSuffix(info.Name(), ".markdown") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
