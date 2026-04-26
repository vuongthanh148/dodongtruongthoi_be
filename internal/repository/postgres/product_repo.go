package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) List(ctx context.Context, category string, includeInactive bool) ([]domain.Product, error) {
	query := "SELECT id, title, subtitle, category_id, badge, base_price, description, meaning, default_bg, default_frame, bg_tones, frames, zodiac_ids, purpose_place, purpose_use, purpose_avoid, specs, requires_bg_tone, requires_frame, requires_size, is_active, sort_order, created_at, updated_at FROM products WHERE 1=1"

	args := []interface{}{}
	if !includeInactive {
		query += " AND is_active = true"
	}
	if category != "" {
		query += " AND category_id = $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, category)
	}
	query += " ORDER BY sort_order, created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		var specs map[string]string
		err := rows.Scan(
			&p.ID, &p.Title, &p.Subtitle, &p.CategoryID, &p.Badge, &p.BasePrice,
			&p.Description, &p.Meaning, &p.DefaultBG, &p.DefaultFrame,
			&p.BGTones, &p.Frames, &p.ZodiacIDs,
			&p.PurposePlace, &p.PurposeUse, &p.PurposeAvoid,
			&specs, &p.RequiresBGTone, &p.RequiresFrame, &p.RequiresSize,
			&p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		p.Specs = specs
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *ProductRepository) Get(ctx context.Context, id string, includeInactive bool) (domain.Product, bool, error) {
	query := "SELECT id, title, subtitle, category_id, badge, base_price, description, meaning, default_bg, default_frame, bg_tones, frames, zodiac_ids, purpose_place, purpose_use, purpose_avoid, specs, requires_bg_tone, requires_frame, requires_size, is_active, sort_order, created_at, updated_at FROM products WHERE id = $1"

	if !includeInactive {
		query += " AND is_active = true"
	}

	var p domain.Product
	var specs map[string]string
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Subtitle, &p.CategoryID, &p.Badge, &p.BasePrice,
		&p.Description, &p.Meaning, &p.DefaultBG, &p.DefaultFrame,
		&p.BGTones, &p.Frames, &p.ZodiacIDs,
		&p.PurposePlace, &p.PurposeUse, &p.PurposeAvoid,
		&specs, &p.RequiresBGTone, &p.RequiresFrame, &p.RequiresSize,
		&p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == nil {
		p.Specs = specs
	}

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.Product{}, false, nil
		}
		return domain.Product{}, false, err
	}
	return p, true, nil
}

func (r *ProductRepository) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	query := `INSERT INTO products (id, title, subtitle, category_id, badge, base_price, description, meaning, default_bg, default_frame, bg_tones, frames, zodiac_ids, purpose_place, purpose_use, purpose_avoid, specs, requires_bg_tone, requires_frame, requires_size, is_active, sort_order, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	RETURNING created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		p.ID, p.Title, p.Subtitle, p.CategoryID, p.Badge, p.BasePrice,
		p.Description, p.Meaning, p.DefaultBG, p.DefaultFrame,
		p.BGTones, p.Frames, p.ZodiacIDs,
		p.PurposePlace, p.PurposeUse, p.PurposeAvoid,
		p.Specs, p.RequiresBGTone, p.RequiresFrame, p.RequiresSize,
		p.IsActive, p.SortOrder, p.CreatedAt, p.UpdatedAt,
	).Scan(&p.CreatedAt, &p.UpdatedAt)

	return p, err
}

func (r *ProductRepository) Update(ctx context.Context, id string, p domain.Product) (domain.Product, error) {
	if id == "" {
		return domain.Product{}, errors.New("product id is required")
	}

	query := `UPDATE products SET
		title = COALESCE(NULLIF($2, ''), title),
		subtitle = COALESCE(NULLIF($3, ''), subtitle),
		category_id = COALESCE(NULLIF($4, ''), category_id),
		badge = COALESCE(NULLIF($5, ''), badge),
		base_price = CASE WHEN $6 > 0 THEN $6 ELSE base_price END,
		description = COALESCE(NULLIF($7, ''), description),
		meaning = COALESCE(NULLIF($8, ''), meaning),
		default_bg = COALESCE(NULLIF($9, ''), default_bg),
		default_frame = COALESCE(NULLIF($10, ''), default_frame),
		bg_tones = CASE WHEN array_length($11::text[], 1) > 0 THEN $11 ELSE bg_tones END,
		frames = CASE WHEN array_length($12::text[], 1) > 0 THEN $12 ELSE frames END,
		zodiac_ids = CASE WHEN array_length($13::text[], 1) > 0 THEN $13 ELSE zodiac_ids END,
		purpose_place = CASE WHEN array_length($14::text[], 1) > 0 THEN $14 ELSE purpose_place END,
		purpose_use = CASE WHEN array_length($15::text[], 1) > 0 THEN $15 ELSE purpose_use END,
		purpose_avoid = CASE WHEN array_length($16::text[], 1) > 0 THEN $16 ELSE purpose_avoid END,
		specs = COALESCE(NULLIF($17::jsonb, '{}'), specs),
		requires_bg_tone = $18,
		requires_frame = $19,
		requires_size = $20,
		is_active = $21,
		sort_order = $22,
		updated_at = now()
	WHERE id = $23
	RETURNING id, title, subtitle, category_id, badge, base_price, description, meaning, default_bg, default_frame, bg_tones, frames, zodiac_ids, purpose_place, purpose_use, purpose_avoid, specs, requires_bg_tone, requires_frame, requires_size, is_active, sort_order, created_at, updated_at`

	var result domain.Product
	err := r.pool.QueryRow(ctx, query,
		id, p.Title, p.Subtitle, p.CategoryID, p.Badge, p.BasePrice,
		p.Description, p.Meaning, p.DefaultBG, p.DefaultFrame,
		p.BGTones, p.Frames, p.ZodiacIDs,
		p.PurposePlace, p.PurposeUse, p.PurposeAvoid,
		p.Specs, p.RequiresBGTone, p.RequiresFrame, p.RequiresSize,
		p.IsActive, p.SortOrder, id,
	).Scan(
		&result.ID, &result.Title, &result.Subtitle, &result.CategoryID, &result.Badge, &result.BasePrice,
		&result.Description, &result.Meaning, &result.DefaultBG, &result.DefaultFrame,
		&result.BGTones, &result.Frames, &result.ZodiacIDs,
		&result.PurposePlace, &result.PurposeUse, &result.PurposeAvoid,
		&result.Specs, &result.RequiresBGTone, &result.RequiresFrame, &result.RequiresSize,
		&result.IsActive, &result.SortOrder, &result.CreatedAt, &result.UpdatedAt,
	)

	return result, err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "UPDATE products SET is_active = false, updated_at = now() WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}
