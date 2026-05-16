package service

import (
	"gitlab.com/shashwat-dixit/portfolio/backend/internal/cache"
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

// TODO: implement List — check Redis cache, fall back to repo, cache result
// TODO: implement GetBySlug — check Redis cache, fall back to repo, cache result
