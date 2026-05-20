package repository

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepo struct {
	pool *pgxpool.Pool
}

func NewPostRepo(pool *pgxpool.Pool) *PostRepo {
	return &PostRepo{pool: pool}
}

func (r *PostRepo) ListPublished(ctx context.Context, tag string, page, limit int, includeDrafts bool) ([]model.PostSummary, int, error) {
	offset := (page - 1) * limit

	statusFilter := "p.status = 'published'"
	if includeDrafts {
		statusFilter = "p.status IN ('published', 'draft')"
	}

	var countQuery string
	var listQuery string
	var args []any

	if tag != "" {
		countQuery = fmt.Sprintf(`
			SELECT COUNT(DISTINCT p.id)
			FROM posts p
			JOIN post_tags pt ON pt.post_id = p.id
			JOIN tags t ON t.id = pt.tag_id
			WHERE %s AND t.slug = $1`, statusFilter)
		args = append(args, tag)

		listQuery = fmt.Sprintf(`
			SELECT p.slug, p.title, p.description, p.cover_image, p.status, p.reading_time, p.published_at,
				COALESCE(
					(SELECT array_agg(t2.name) FROM tags t2 JOIN post_tags pt2 ON pt2.tag_id = t2.id WHERE pt2.post_id = p.id),
					'{}'
				) AS tags
			FROM posts p
			JOIN post_tags pt ON pt.post_id = p.id
			JOIN tags t ON t.id = pt.tag_id
			WHERE %s AND t.slug = $1
			ORDER BY p.published_at DESC NULLS LAST
			LIMIT $2 OFFSET $3`, statusFilter)
		args = append(args, limit, offset)
	} else {
		countQuery = fmt.Sprintf(`SELECT COUNT(*) FROM posts p WHERE %s`, statusFilter)

		listQuery = fmt.Sprintf(`
			SELECT p.slug, p.title, p.description, p.cover_image, p.status, p.reading_time, p.published_at,
				COALESCE(
					(SELECT array_agg(t.name) FROM tags t JOIN post_tags pt ON pt.tag_id = t.id WHERE pt.post_id = p.id),
					'{}'
				) AS tags
			FROM posts p
			WHERE %s
			ORDER BY p.published_at DESC NULLS LAST
			LIMIT $1 OFFSET $2`, statusFilter)
		args = append(args, limit, offset)
	}

	var total int
	if tag != "" {
		err := r.pool.QueryRow(ctx, countQuery, tag).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("count posts: %w", err)
		}
	} else {
		err := r.pool.QueryRow(ctx, countQuery).Scan(&total)
		if err != nil {
			return nil, 0, fmt.Errorf("count posts: %w", err)
		}
	}

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	var posts []model.PostSummary
	for rows.Next() {
		var p model.PostSummary
		var tags []string
		if err := rows.Scan(&p.Slug, &p.Title, &p.Description, &p.CoverImage, &p.Status, &p.ReadingTime, &p.PublishedAt, &tags); err != nil {
			return nil, 0, fmt.Errorf("scan post: %w", err)
		}
		p.Tags = tags
		posts = append(posts, p)
	}
	if posts == nil {
		posts = []model.PostSummary{}
	}

	return posts, total, nil
}

func (r *PostRepo) GetBySlug(ctx context.Context, slug string) (*model.Post, []string, error) {
	query := `
		SELECT id, slug, title, description, content_html, cover_image, status,
			reading_time, author, published_at, updated_at
		FROM posts
		WHERE slug = $1 AND status = 'published'`

	var post model.Post
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&post.ID, &post.Slug, &post.Title, &post.Description,
		&post.ContentHTML, &post.CoverImage, &post.Status,
		&post.ReadingTime, &post.Author, &post.PublishedAt, &post.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("get post by slug: %w", err)
	}

	tagQuery := `
		SELECT t.name FROM tags t
		JOIN post_tags pt ON pt.tag_id = t.id
		WHERE pt.post_id = $1
		ORDER BY t.name`

	rows, err := r.pool.Query(ctx, tagQuery, post.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("get post tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, name)
	}
	if tags == nil {
		tags = []string{}
	}

	return &post, tags, nil
}

func (r *PostRepo) Upsert(ctx context.Context, post *model.Post) (string, error) {
	query := `
		INSERT INTO posts (slug, title, description, content_md, content_html, cover_image, status, reading_time, author, published_at, updated_at, gitlab_sha)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (slug) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			content_md = EXCLUDED.content_md,
			content_html = EXCLUDED.content_html,
			cover_image = EXCLUDED.cover_image,
			status = EXCLUDED.status,
			reading_time = EXCLUDED.reading_time,
			author = EXCLUDED.author,
			published_at = EXCLUDED.published_at,
			updated_at = EXCLUDED.updated_at,
			gitlab_sha = EXCLUDED.gitlab_sha
		RETURNING id`

	var id string
	err := r.pool.QueryRow(ctx, query,
		post.Slug, post.Title, post.Description, post.ContentMD, post.ContentHTML,
		post.CoverImage, post.Status, post.ReadingTime, post.Author,
		post.PublishedAt, post.UpdatedAt, post.GitLabSHA,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("upsert post: %w", err)
	}
	return id, nil
}

func (r *PostRepo) Delete(ctx context.Context, slug string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM posts WHERE slug = $1`, slug)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	return nil
}

func (r *PostRepo) AllSlugs(ctx context.Context) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT slug FROM posts`)
	if err != nil {
		return nil, fmt.Errorf("all slugs: %w", err)
	}
	defer rows.Close()

	var slugs []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, fmt.Errorf("scan slug: %w", err)
		}
		slugs = append(slugs, s)
	}
	return slugs, nil
}

func (r *PostRepo) GetSHAMap(ctx context.Context) (map[string]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT slug, COALESCE(gitlab_sha, '') FROM posts`)
	if err != nil {
		return nil, fmt.Errorf("get sha map: %w", err)
	}
	defer rows.Close()

	m := make(map[string]string)
	for rows.Next() {
		var slug, sha string
		if err := rows.Scan(&slug, &sha); err != nil {
			return nil, fmt.Errorf("scan sha: %w", err)
		}
		m[slug] = sha
	}
	return m, nil
}

func (r *PostRepo) SetTags(ctx context.Context, postID string, tagIDs []int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM post_tags WHERE post_id = $1`, postID)
	if err != nil {
		return fmt.Errorf("delete old tags: %w", err)
	}

	if len(tagIDs) > 0 {
		var sb strings.Builder
		sb.WriteString(`INSERT INTO post_tags (post_id, tag_id) VALUES `)
		args := []any{postID}
		for i, tagID := range tagIDs {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("($1, $%d)", i+2))
			args = append(args, tagID)
		}
		sb.WriteString(` ON CONFLICT DO NOTHING`)

		_, err = tx.Exec(ctx, sb.String(), args...)
		if err != nil {
			return fmt.Errorf("insert tags: %w", err)
		}
	}

	return tx.Commit(ctx)
}
