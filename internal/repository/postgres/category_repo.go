package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) List(ctx context.Context, includeInactive bool) ([]domain.Category, error) {
	query := `SELECT c.slug, c.name, c.description, c.tone, c.image_url, c.sort_order, c.is_active,
		COUNT(p.id) FILTER (WHERE p.is_active = true) AS product_count,
		c.created_at, c.updated_at
		FROM categories c
		LEFT JOIN products p ON p.category_id = c.id`

	if !includeInactive {
		query += " WHERE c.is_active = true"
	}
	query += " GROUP BY c.id ORDER BY c.sort_order, c.created_at DESC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		err := rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.Tone, &c.ImageURL,
			&c.SortOrder, &c.IsActive, &c.ProductCount, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) Get(ctx context.Context, slug string, includeInactive bool) (domain.Category, bool, error) {
	query := "SELECT slug, name, description, tone, image_url, sort_order, is_active, created_at, updated_at FROM categories WHERE slug = $1"

	if !includeInactive {
		query += " AND is_active = true"
	}

	var c domain.Category
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&c.ID, &c.Name, &c.Description, &c.Tone, &c.ImageURL,
		&c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.Category{}, false, nil
		}
		return domain.Category{}, false, err
	}
	return c, true, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c domain.Category) (domain.Category, error) {
	query := `INSERT INTO categories (slug, name, description, tone, image_url, sort_order, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		c.ID, c.Name, c.Description, c.Tone, c.ImageURL,
		c.SortOrder, c.IsActive, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.CreatedAt, &c.UpdatedAt)

	return c, err
}

func (r *CategoryRepository) Update(ctx context.Context, slug string, c domain.Category) (domain.Category, error) {
	if slug == "" {
		return domain.Category{}, errors.New("category slug is required")
	}

	query := `UPDATE categories SET
		name = COALESCE(NULLIF($2, ''), name),
		description = $3,
		tone = COALESCE(NULLIF($4, ''), tone),
		image_url = COALESCE(NULLIF($5, ''), image_url),
		sort_order = $6,
		is_active = $7,
		updated_at = now()
	WHERE slug = $1
	RETURNING slug, name, description, tone, image_url, sort_order, is_active, created_at, updated_at`

	var result domain.Category
	err := r.pool.QueryRow(ctx, query,
		slug, c.Name, c.Description, c.Tone, c.ImageURL,
		c.SortOrder, c.IsActive,
	).Scan(
		&result.ID, &result.Name, &result.Description, &result.Tone, &result.ImageURL,
		&result.SortOrder, &result.IsActive, &result.CreatedAt, &result.UpdatedAt,
	)

	return result, err
}

func (r *CategoryRepository) Delete(ctx context.Context, slug string) error {
	result, err := r.pool.Exec(ctx, "UPDATE categories SET is_active = false, updated_at = now() WHERE slug = $1", slug)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("category not found")
	}
	return nil
}
