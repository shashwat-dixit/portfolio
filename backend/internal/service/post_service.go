package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/cache"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/repository"
)

type PostService struct {
	postRepo *repository.PostRepo
	tagRepo  *repository.TagRepo
	cache    *cache.RedisCache
}

func NewPostService(postRepo *repository.PostRepo, tagRepo *repository.TagRepo, cache *cache.RedisCache) *PostService {
	return &PostService{postRepo: postRepo, tagRepo: tagRepo, cache: cache}
}

func (s *PostService) List(ctx context.Context, tag string, page, limit int) (*model.PostListResponse, error) {
	cacheKey := cache.PostListKey(tag, page, limit)

	cached, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		slog.Warn("cache get error", "key", cacheKey, "error", err)
	}
	if cached != "" {
		var resp model.PostListResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	}

	posts, total, err := s.postRepo.ListPublished(ctx, tag, page, limit)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}

	totalPages := (total + limit - 1) / limit
	resp := &model.PostListResponse{
		Posts: posts,
		Pagination: model.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	data, err := json.Marshal(resp)
	if err == nil {
		if err := s.cache.Set(ctx, cacheKey, string(data)); err != nil {
			slog.Warn("cache set error", "key", cacheKey, "error", err)
		}
	}

	return resp, nil
}

func (s *PostService) GetBySlug(ctx context.Context, slug string) (*model.Post, []string, error) {
	cacheKey := cache.PostKey(slug)

	cached, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		slog.Warn("cache get error", "key", cacheKey, "error", err)
	}
	if cached != "" {
		var entry cachedPost
		if err := json.Unmarshal([]byte(cached), &entry); err == nil {
			return &entry.Post, entry.Tags, nil
		}
	}

	post, tags, err := s.postRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, nil, fmt.Errorf("get post: %w", err)
	}
	if post == nil {
		return nil, nil, nil
	}

	entry := cachedPost{Post: *post, Tags: tags}
	data, err := json.Marshal(entry)
	if err == nil {
		if err := s.cache.Set(ctx, cacheKey, string(data)); err != nil {
			slog.Warn("cache set error", "key", cacheKey, "error", err)
		}
	}

	return post, tags, nil
}

type cachedPost struct {
	Post model.Post `json:"post"`
	Tags []string   `json:"tags"`
}
