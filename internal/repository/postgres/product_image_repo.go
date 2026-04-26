package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ProductImageRepository struct {
	pool *pgxpool.Pool
}

func NewProductImageRepository(pool *pgxpool.Pool) *ProductImageRepository {
	return &ProductImageRepository{pool: pool}
}

func (r *ProductImageRepository) ListByProduct(ctx context.Context, productID string) ([]domain.ProductImage, error) {
	query := `SELECT id, product_id, bg_tone, frame, url, alt_text, sort_order, created_at
	FROM product_images WHERE product_id = $1 ORDER BY sort_order, created_at ASC`

	rows, err := r.pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []domain.ProductImage
	for rows.Next() {
		var img domain.ProductImage
		err := rows.Scan(
			&img.ID, &img.ProductID, &img.BGTone, &img.Frame, &img.URL, &img.AltText, &img.SortOrder, &img.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

func (r *ProductImageRepository) Get(ctx context.Context, id string) (domain.ProductImage, bool, error) {
	query := `SELECT id, product_id, bg_tone, frame, url, alt_text, sort_order, created_at
	FROM product_images WHERE id = $1`

	var img domain.ProductImage
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&img.ID, &img.ProductID, &img.BGTone, &img.Frame, &img.URL, &img.AltText, &img.SortOrder, &img.CreatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.ProductImage{}, false, nil
		}
		return domain.ProductImage{}, false, err
	}
	return img, true, nil
}

func (r *ProductImageRepository) Create(ctx context.Context, img domain.ProductImage) (domain.ProductImage, error) {
	query := `INSERT INTO product_images (product_id, bg_tone, frame, url, alt_text, sort_order)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, product_id, bg_tone, frame, url, alt_text, sort_order, created_at`

	err := r.pool.QueryRow(ctx, query,
		img.ProductID, img.BGTone, img.Frame, img.URL, img.AltText, img.SortOrder,
	).Scan(
		&img.ID, &img.ProductID, &img.BGTone, &img.Frame, &img.URL, &img.AltText, &img.SortOrder, &img.CreatedAt,
	)

	return img, err
}

func (r *ProductImageRepository) Update(ctx context.Context, id string, img domain.ProductImage) (domain.ProductImage, error) {
	query := `UPDATE product_images SET
		alt_text = COALESCE(NULLIF($2, ''), alt_text),
		sort_order = $3
	WHERE id = $4
	RETURNING id, product_id, bg_tone, frame, url, alt_text, sort_order, created_at`

	var result domain.ProductImage
	err := r.pool.QueryRow(ctx, query,
		id, img.AltText, img.SortOrder, id,
	).Scan(
		&result.ID, &result.ProductID, &result.BGTone, &result.Frame, &result.URL, &result.AltText, &result.SortOrder, &result.CreatedAt,
	)

	return result, err
}

func (r *ProductImageRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM product_images WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("product image not found")
	}
	return nil
}

func (r *ProductImageRepository) DeleteByProduct(ctx context.Context, productID string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM product_images WHERE product_id = $1", productID)
	return err
}
