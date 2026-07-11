package repository

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/shashwat-dixit/portfolio/backend/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepo struct {
	pool *pgxpool.Pool
}

func NewTagRepo(pool *pgxpool.Pool) *TagRepo {
	return &TagRepo{pool: pool}
}

func (r *TagRepo) ListWithCounts(ctx context.Context) ([]model.Tag, error) {
	query := `
		SELECT t.id, t.name, t.slug, COUNT(pt.post_id) AS count
		FROM tags t
		LEFT JOIN post_tags pt ON pt.tag_id = t.id
		LEFT JOIN posts p ON p.id = pt.post_id AND p.status = 'published'
		GROUP BY t.id, t.name, t.slug
		HAVING COUNT(pt.post_id) > 0
		ORDER BY count DESC, t.name ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var t model.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.Count); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, t)
	}
	if tags == nil {
		tags = []model.Tag{}
	}

	return tags, nil
}

func slugify(name string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))
}

func (r *TagRepo) UpsertMany(ctx context.Context, names []string) ([]int, error) {
	if len(names) == 0 {
		return []int{}, nil
	}

	seen := make(map[string]string)
	dedupedNames := make([]string, 0, len(names))
	for _, name := range names {
		slug := slugify(name)
		if _, exists := seen[slug]; !exists {
			seen[slug] = strings.TrimSpace(name)
			dedupedNames = append(dedupedNames, strings.TrimSpace(name))
		}
	}

	var sb strings.Builder
	sb.WriteString(`INSERT INTO tags (name, slug) VALUES `)
	args := make([]any, 0, len(dedupedNames)*2)
	for i, name := range dedupedNames {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		args = append(args, name, slugify(name))
	}
	sb.WriteString(` ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name`)

	_, err := r.pool.Exec(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("upsert tags: %w", err)
	}

	slugs := make([]string, len(dedupedNames))
	slugArgs := make([]any, len(dedupedNames))
	for i, name := range dedupedNames {
		slugs[i] = fmt.Sprintf("$%d", i+1)
		slugArgs[i] = slugify(name)
	}

	selectQuery := fmt.Sprintf(
		`SELECT id FROM tags WHERE slug IN (%s) ORDER BY slug`,
		strings.Join(slugs, ", "),
	)

	rows, err := r.pool.Query(ctx, selectQuery, slugArgs...)
	if err != nil {
		return nil, fmt.Errorf("select tag ids: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan tag id: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
