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
	query := "SELECT id, name, slug, description, tone, image_url, sort_order, is_active, created_at, updated_at FROM categories"

	if !includeInactive {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY sort_order, created_at DESC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		err := rows.Scan(
			&c.ID, &c.Name, &c.Slug, &c.Description, &c.Tone, &c.ImageURL,
			&c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) Get(ctx context.Context, id string, includeInactive bool) (domain.Category, bool, error) {
	query := "SELECT id, name, slug, description, tone, image_url, sort_order, is_active, created_at, updated_at FROM categories WHERE id = $1"

	if !includeInactive {
		query += " AND is_active = true"
	}

	var c domain.Category
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Slug, &c.Description, &c.Tone, &c.ImageURL,
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
	query := `INSERT INTO categories (id, name, slug, description, tone, image_url, sort_order, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		c.ID, c.Name, c.Slug, c.Description, c.Tone, c.ImageURL,
		c.SortOrder, c.IsActive, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.CreatedAt, &c.UpdatedAt)

	return c, err
}

func (r *CategoryRepository) Update(ctx context.Context, id string, c domain.Category) (domain.Category, error) {
	if id == "" {
		return domain.Category{}, errors.New("category id is required")
	}

	query := `UPDATE categories SET
		name = COALESCE(NULLIF($2, ''), name),
		slug = COALESCE(NULLIF($3, ''), slug),
		description = $4,
		tone = COALESCE(NULLIF($5, ''), tone),
		image_url = COALESCE(NULLIF($6, ''), image_url),
		sort_order = $7,
		is_active = $8,
		updated_at = now()
	WHERE id = $9
	RETURNING id, name, slug, description, tone, image_url, sort_order, is_active, created_at, updated_at`

	var result domain.Category
	err := r.pool.QueryRow(ctx, query,
		id, c.Name, c.Slug, c.Description, c.Tone, c.ImageURL,
		c.SortOrder, c.IsActive, id,
	).Scan(
		&result.ID, &result.Name, &result.Slug, &result.Description, &result.Tone, &result.ImageURL,
		&result.SortOrder, &result.IsActive, &result.CreatedAt, &result.UpdatedAt,
	)

	return result, err
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "UPDATE categories SET is_active = false, updated_at = now() WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("category not found")
	}
	return nil
}
