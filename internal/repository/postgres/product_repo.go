package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

const productColumns = `
	p.id, p.title, p.subtitle, c.slug AS category_id, p.badge, p.base_price, p.description, p.meaning,
	p.variant_options, p.default_variant, p.zodiac_ids, p.purpose_place, p.purpose_use,
	p.purpose_avoid, p.specs, p.requires_size, p.is_active, p.sort_order, p.created_at, p.updated_at`

const productFrom = `FROM products p JOIN categories c ON c.id = p.category_id`

func scanProduct(row pgx.Row) (domain.Product, error) {
	var p domain.Product
	var variantOptionsRaw []byte
	var defaultVariantRaw []byte
	err := row.Scan(
		&p.ID, &p.Title, &p.Subtitle, &p.CategoryID, &p.Badge, &p.BasePrice,
		&p.Description, &p.Meaning,
		&variantOptionsRaw, &defaultVariantRaw,
		&p.ZodiacIDs, &p.PurposePlace, &p.PurposeUse, &p.PurposeAvoid,
		&p.Specs, &p.RequiresSize, &p.IsActive, &p.SortOrder, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return domain.Product{}, err
	}
	if len(variantOptionsRaw) > 0 {
		_ = json.Unmarshal(variantOptionsRaw, &p.VariantOptions)
	}
	if len(defaultVariantRaw) > 0 {
		_ = json.Unmarshal(defaultVariantRaw, &p.DefaultVariant)
	}
	if p.VariantOptions == nil {
		p.VariantOptions = []domain.VariantOption{}
	}
	if p.DefaultVariant == nil {
		p.DefaultVariant = map[string]string{}
	}
	return p, nil
}

func (r *ProductRepository) List(ctx context.Context, q domain.ProductQuery, includeInactive bool) ([]domain.Product, error) {
	query := "SELECT " + productColumns + " " + productFrom + " WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if !includeInactive {
		query += " AND p.is_active = true"
	}
	if q.Category != "" {
		query += fmt.Sprintf(" AND c.slug = $%d", argCount)
		args = append(args, q.Category)
		argCount++
	}

	orderBy := "p.sort_order, p.created_at DESC"
	switch q.Sort {
	case "price_asc":
		orderBy = "p.base_price ASC, p.sort_order"
	case "price_desc":
		orderBy = "p.base_price DESC, p.sort_order"
	case "newest":
		orderBy = "p.created_at DESC, p.sort_order"
	}
	query += " ORDER BY " + orderBy

	if q.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, q.Limit)
		argCount++
	}
	if q.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, q.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *ProductRepository) Get(ctx context.Context, id string, includeInactive bool) (domain.Product, bool, error) {
	query := "SELECT " + productColumns + " " + productFrom + " WHERE p.id = $1"
	if !includeInactive {
		query += " AND p.is_active = true"
	}

	p, err := scanProduct(r.pool.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, false, nil
	}
	if err != nil {
		return domain.Product{}, false, err
	}
	return p, true, nil
}

func (r *ProductRepository) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	if p.ZodiacIDs == nil {
		p.ZodiacIDs = []string{}
	}
	if p.PurposePlace == nil {
		p.PurposePlace = []string{}
	}
	if p.PurposeUse == nil {
		p.PurposeUse = []string{}
	}
	if p.PurposeAvoid == nil {
		p.PurposeAvoid = []string{}
	}
	variantOptionsJSON, _ := json.Marshal(p.VariantOptions)
	defaultVariantJSON, _ := json.Marshal(p.DefaultVariant)

	query := `INSERT INTO products (
		id, title, subtitle, category_id, badge, base_price, description, meaning,
		variant_options, default_variant, zodiac_ids, purpose_place, purpose_use,
		purpose_avoid, specs, requires_size, is_active, sort_order, created_at, updated_at
	) VALUES (
		$1,$2,$3,(SELECT id FROM categories WHERE slug = $4),$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20
	) RETURNING created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		p.ID, p.Title, p.Subtitle, p.CategoryID, p.Badge, p.BasePrice,
		p.Description, p.Meaning,
		variantOptionsJSON, defaultVariantJSON,
		p.ZodiacIDs, p.PurposePlace, p.PurposeUse, p.PurposeAvoid,
		p.Specs, p.RequiresSize, p.IsActive, p.SortOrder, p.CreatedAt, p.UpdatedAt,
	).Scan(&p.CreatedAt, &p.UpdatedAt)

	return p, err
}

func (r *ProductRepository) Update(ctx context.Context, id string, p domain.Product) (domain.Product, error) {
	if id == "" {
		return domain.Product{}, errors.New("product id is required")
	}

	variantOptionsJSON, _ := json.Marshal(p.VariantOptions)
	defaultVariantJSON, _ := json.Marshal(p.DefaultVariant)

	query := `UPDATE products SET
		title          = COALESCE(NULLIF($2, ''), title),
		subtitle       = $3,
		category_id    = COALESCE((SELECT id FROM categories WHERE slug = NULLIF($4, '')), category_id),
		badge          = $5,
		base_price     = CASE WHEN $6 > 0 THEN $6 ELSE base_price END,
		description    = $7,
		meaning        = $8,
		variant_options = $9,
		default_variant = $10,
		zodiac_ids     = CASE WHEN array_length($11::text[], 1) > 0 THEN $11 ELSE zodiac_ids END,
		purpose_place  = $12,
		purpose_use    = $13,
		purpose_avoid  = $14,
		specs          = COALESCE(NULLIF($15::jsonb, '{}'), specs),
		requires_size  = $16,
		is_active      = $17,
		sort_order     = $18,
		updated_at     = now()
	WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		id, p.Title, p.Subtitle, p.CategoryID, p.Badge, p.BasePrice,
		p.Description, p.Meaning,
		variantOptionsJSON, defaultVariantJSON,
		p.ZodiacIDs, p.PurposePlace, p.PurposeUse, p.PurposeAvoid,
		p.Specs, p.RequiresSize, p.IsActive, p.SortOrder,
	)
	if err != nil {
		return domain.Product{}, err
	}
	if tag.RowsAffected() == 0 {
		return domain.Product{}, errors.New("product not found")
	}

	result, _, err := r.Get(ctx, id, true)
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
