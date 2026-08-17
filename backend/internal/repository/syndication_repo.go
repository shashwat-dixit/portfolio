package repository

import (
	"context"
	"fmt"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SyndicationRepo struct {
	pool *pgxpool.Pool
}

func NewSyndicationRepo(pool *pgxpool.Pool) *SyndicationRepo {
	return &SyndicationRepo{pool: pool}
}

func (r *SyndicationRepo) List(ctx context.Context) (map[string]map[string]model.Syndication, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT post_id, platform, remote_id, COALESCE(remote_url, ''), content_hash
		FROM post_syndications`)
	if err != nil {
		return nil, fmt.Errorf("list syndications: %w", err)
	}
	defer rows.Close()

	out := make(map[string]map[string]model.Syndication)
	for rows.Next() {
		var s model.Syndication
		if err := rows.Scan(&s.PostID, &s.Platform, &s.RemoteID, &s.RemoteURL, &s.ContentHash); err != nil {
			return nil, fmt.Errorf("scan syndication: %w", err)
		}
		if out[s.PostID] == nil {
			out[s.PostID] = make(map[string]model.Syndication)
		}
		out[s.PostID][s.Platform] = s
	}
	return out, nil
}

func (r *SyndicationRepo) Upsert(ctx context.Context, s model.Syndication) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO post_syndications (post_id, platform, remote_id, remote_url, content_hash, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (post_id, platform) DO UPDATE SET
			remote_id = EXCLUDED.remote_id,
			remote_url = EXCLUDED.remote_url,
			content_hash = EXCLUDED.content_hash,
			updated_at = NOW()`,
		s.PostID, s.Platform, s.RemoteID, s.RemoteURL, s.ContentHash,
	)
	if err != nil {
		return fmt.Errorf("upsert syndication: %w", err)
	}
	return nil
}
