package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/config"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/repository"
)

type crossPoster interface {
	Enabled() bool
	SupportsUpdate() bool
	Publish(ctx context.Context, post model.SyndicatablePost) (remoteID, remoteURL string, err error)
	Update(ctx context.Context, remoteID string, post model.SyndicatablePost) (remoteIDOut, remoteURL string, err error)
}

type postLister interface {
	ListPublishedFull(ctx context.Context, limit int) ([]model.SyndicatablePost, error)
}

type syndicationStore interface {
	List(ctx context.Context) (map[string]map[string]model.Syndication, error)
	Upsert(ctx context.Context, s model.Syndication) error
}

type SyndicationService struct {
	cfg      *config.Config
	postRepo postLister
	synRepo  syndicationStore
	medium   *MediumClient
	substack *SubstackClient
	pace     time.Duration
}

func NewSyndicationService(
	cfg *config.Config,
	postRepo *repository.PostRepo,
	synRepo *repository.SyndicationRepo,
) *SyndicationService {
	s := &SyndicationService{cfg: cfg, postRepo: postRepo, synRepo: synRepo, pace: 250 * time.Millisecond}
	if cfg.MediumToken != "" {
		s.medium = NewMediumClient(cfg)
	}
	if cfg.SubstackPublicationURL != "" && (cfg.SubstackSID != "" || cfg.SubstackCookies != "" || cfg.SubstackConnectSID != "") {
		s.substack = NewSubstackClient(cfg)
	}
	return s
}

func (s *SyndicationService) Enabled() bool {
	return s.medium.Enabled() || s.substack.Enabled()
}

func (s *SyndicationService) Syndicate(ctx context.Context) (medium, substack *model.CrossPostResult) {
	if !s.Enabled() {
		return nil, nil
	}

	posts, err := s.postRepo.ListPublishedFull(ctx, 0)
	if err != nil {
		slog.Error("list posts for syndication failed", "error", err)
		failed := &model.CrossPostResult{Failed: 1, Errors: []string{err.Error()}}
		if s.medium.Enabled() {
			medium = failed
		}
		if s.substack.Enabled() {
			substack = failed
		}
		return medium, substack
	}

	existing, err := s.synRepo.List(ctx)
	if err != nil {
		slog.Error("list syndications failed", "error", err)
		failed := &model.CrossPostResult{Failed: 1, Errors: []string{err.Error()}}
		if s.medium.Enabled() {
			medium = failed
		}
		if s.substack.Enabled() {
			substack = failed
		}
		return medium, substack
	}

	if s.medium.Enabled() {
		medium = s.syndicatePlatform(ctx, model.PlatformMedium, s.medium, posts, existing)
	}
	if s.substack.Enabled() {
		substack = s.syndicatePlatform(ctx, model.PlatformSubstack, s.substack, posts, existing)
	}
	return medium, substack
}

func (s *SyndicationService) syndicatePlatform(
	ctx context.Context,
	platform string,
	poster crossPoster,
	posts []model.SyndicatablePost,
	existing map[string]map[string]model.Syndication,
) *model.CrossPostResult {
	result := &model.CrossPostResult{}
	for i, post := range posts {
		if err := ctx.Err(); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
			return result
		}
		if i > 0 && s.pace > 0 {
			select {
			case <-ctx.Done():
				result.Failed++
				result.Errors = append(result.Errors, ctx.Err().Error())
				return result
			case <-time.After(s.pace):
			}
		}

		prev, ok := existing[post.ID][platform]
		if ok && prev.ContentHash == post.ContentHash && prev.RemoteID != "" {
			result.Skipped++
			continue
		}

		var remoteID, remoteURL string
		var err error
		created := !ok || prev.RemoteID == ""
		if created {
			remoteID, remoteURL, err = poster.Publish(ctx, post)
		} else if poster.SupportsUpdate() {
			remoteID, remoteURL, err = poster.Update(ctx, prev.RemoteID, post)
			if remoteID == "" {
				remoteID = prev.RemoteID
			}
			if remoteURL == "" {
				remoteURL = prev.RemoteURL
			}
		} else {
			result.Skipped++
			slog.Info("cross-post update skipped", "platform", platform, "slug", post.Slug, "reason", "platform cannot update")
			continue
		}
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", post.Slug, err))
			slog.Warn("cross-post failed", "platform", platform, "slug", post.Slug, "error", err)
			continue
		}

		if err := s.synRepo.Upsert(ctx, model.Syndication{
			PostID:      post.ID,
			Platform:    platform,
			RemoteID:    remoteID,
			RemoteURL:   remoteURL,
			ContentHash: post.ContentHash,
		}); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: save %v", post.Slug, err))
			continue
		}
		if created {
			result.Created++
		} else {
			result.Updated++
		}
		slog.Info("cross-posted", "platform", platform, "slug", post.Slug, "url", remoteURL)
	}
	return result
}
