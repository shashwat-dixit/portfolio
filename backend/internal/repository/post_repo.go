package repository

import (
	"context"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepo struct {
	pool *pgxpool.Pool
}

func NewPostRepo(pool *pgxpool.Pool) *PostRepo {
	return &PostRepo{pool: pool}
}

// ListPublished returns paginated published posts, optionally filtered by tag.
func (r *PostRepo) ListPublished(ctx context.Context, tag string, page, limit int) ([]model.PostSummary, int, error) {
	// TODO: query posts JOIN post_tags/tags, filter by status='published',
	// optional tag filter, ORDER BY published_at DESC, with pagination.
	// Return posts slice, total count, error.
	return nil, 0, nil
}

// GetBySlug returns a single published post with full content.
func (r *PostRepo) GetBySlug(ctx context.Context, slug string) (*model.Post, []string, error) {
	// TODO: query post by slug, join tags, return post + tag names.
	return nil, nil, nil
}

// Upsert inserts or updates a post by slug.
func (r *PostRepo) Upsert(ctx context.Context, post *model.Post) error {
	// TODO: INSERT ... ON CONFLICT (slug) DO UPDATE
	return nil
}

// Delete removes a post by slug.
func (r *PostRepo) Delete(ctx context.Context, slug string) error {
	// TODO: DELETE FROM posts WHERE slug = $1
	return nil
}

// AllSlugs returns all slugs currently in the database.
func (r *PostRepo) AllSlugs(ctx context.Context) ([]string, error) {
	// TODO: SELECT slug FROM posts
	return nil, nil
}

// SetTags replaces the tags for a given post.
func (r *PostRepo) SetTags(ctx context.Context, postID string, tagIDs []int) error {
	// TODO: DELETE existing post_tags, INSERT new ones
	return nil
}
