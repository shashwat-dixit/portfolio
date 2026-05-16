package repository

import (
	"context"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepo struct {
	pool *pgxpool.Pool
}

func NewTagRepo(pool *pgxpool.Pool) *TagRepo {
	return &TagRepo{pool: pool}
}

// ListWithCounts returns all tags with the number of published posts using each.
func (r *TagRepo) ListWithCounts(ctx context.Context) ([]model.Tag, error) {
	// TODO: SELECT t.*, COUNT(pt.post_id) FROM tags t
	// JOIN post_tags pt ... JOIN posts p WHERE p.status='published'
	// GROUP BY t.id ORDER BY count DESC
	return nil, nil
}

// UpsertMany ensures all given tag names exist and returns their IDs.
func (r *TagRepo) UpsertMany(ctx context.Context, names []string) ([]int, error) {
	// TODO: INSERT INTO tags (name, slug) VALUES ... ON CONFLICT DO NOTHING
	// then SELECT id FROM tags WHERE name = ANY($1)
	return nil, nil
}
