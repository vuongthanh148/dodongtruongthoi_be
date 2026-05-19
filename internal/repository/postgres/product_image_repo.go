package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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
	rows, err := r.pool.Query(ctx, `
		SELECT pi.id, pi.product_id, pi.image_id, i.url, i.name, pi.sort_order, pi.created_at
		FROM product_images pi
		JOIN images i ON i.id = pi.image_id
		WHERE pi.product_id = $1
		ORDER BY pi.sort_order, pi.created_at`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.ProductImage, 0)
	for rows.Next() {
		var img domain.ProductImage
		if err := rows.Scan(&img.ID, &img.ProductID, &img.ImageID, &img.URL, &img.Name, &img.SortOrder, &img.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, img)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load attrs for each product_image
	for i := range out {
		attrs, err := r.loadAttrs(ctx, out[i].ID)
		if err != nil {
			return nil, err
		}
		out[i].Attrs = attrs
	}
	return out, nil
}

func (r *ProductImageRepository) Get(ctx context.Context, id string) (domain.ProductImage, bool, error) {
	var img domain.ProductImage
	err := r.pool.QueryRow(ctx, `
		SELECT pi.id, pi.product_id, pi.image_id, i.url, i.name, pi.sort_order, pi.created_at
		FROM product_images pi
		JOIN images i ON i.id = pi.image_id
		WHERE pi.id = $1`, id).
		Scan(&img.ID, &img.ProductID, &img.ImageID, &img.URL, &img.Name, &img.SortOrder, &img.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ProductImage{}, false, nil
	}
	if err != nil {
		return domain.ProductImage{}, false, err
	}
	attrs, err := r.loadAttrs(ctx, img.ID)
	if err != nil {
		return domain.ProductImage{}, false, err
	}
	img.Attrs = attrs
	return img, true, nil
}

func (r *ProductImageRepository) Attach(ctx context.Context, img domain.ProductImage) (domain.ProductImage, error) {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO product_images (id, product_id, image_id, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, product_id, image_id, sort_order, created_at`,
		img.ID, img.ProductID, img.ImageID, img.SortOrder).
		Scan(&img.ID, &img.ProductID, &img.ImageID, &img.SortOrder, &img.CreatedAt)
	if err != nil {
		return domain.ProductImage{}, err
	}

	if len(img.Attrs) > 0 {
		if err := r.SetAttrs(ctx, img.ID, img.Attrs); err != nil {
			return domain.ProductImage{}, err
		}
	}

	// Re-read to get url + name from images join
	result, _, err := r.Get(ctx, img.ID)
	return result, err
}

func (r *ProductImageRepository) SetAttrs(ctx context.Context, productImageID string, attrs []domain.VariantAttr) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM product_image_attrs WHERE product_image_id = $1", productImageID); err != nil {
		return err
	}

	for _, a := range attrs {
		if _, err := tx.Exec(ctx,
			"INSERT INTO product_image_attrs (product_image_id, attr_key, attr_value) VALUES ($1, $2, $3)",
			productImageID, a.Key, a.Value); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *ProductImageRepository) Detach(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM product_images WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("product image not found")
	}
	return nil
}

func (r *ProductImageRepository) DetachByProduct(ctx context.Context, productID string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM product_images WHERE product_id = $1", productID)
	return err
}

func (r *ProductImageRepository) loadAttrs(ctx context.Context, productImageID string) ([]domain.VariantAttr, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT attr_key, attr_value FROM product_image_attrs WHERE product_image_id = $1 ORDER BY attr_key",
		productImageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attrs := make([]domain.VariantAttr, 0)
	for rows.Next() {
		var a domain.VariantAttr
		if err := rows.Scan(&a.Key, &a.Value); err != nil {
			return nil, err
		}
		attrs = append(attrs, a)
	}
	return attrs, rows.Err()
}
