package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ProductSizeRepository struct {
	pool *pgxpool.Pool
}

func NewProductSizeRepository(pool *pgxpool.Pool) *ProductSizeRepository {
	return &ProductSizeRepository{pool: pool}
}

func (r *ProductSizeRepository) ListByProduct(ctx context.Context, productID string) ([]domain.ProductSize, error) {
	query := "SELECT id, product_id, size_label, size_code, price, sort_order FROM product_sizes WHERE product_id = $1 ORDER BY sort_order, size_code ASC"

	rows, err := r.pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sizes []domain.ProductSize
	for rows.Next() {
		var s domain.ProductSize
		err := rows.Scan(&s.ID, &s.ProductID, &s.SizeLabel, &s.SizeCode, &s.Price, &s.SortOrder)
		if err != nil {
			return nil, err
		}
		sizes = append(sizes, s)
	}
	return sizes, rows.Err()
}

func (r *ProductSizeRepository) Get(ctx context.Context, id string) (domain.ProductSize, bool, error) {
	query := "SELECT id, product_id, size_label, size_code, price, sort_order FROM product_sizes WHERE id = $1"

	var s domain.ProductSize
	err := r.pool.QueryRow(ctx, query, id).Scan(&s.ID, &s.ProductID, &s.SizeLabel, &s.SizeCode, &s.Price, &s.SortOrder)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.ProductSize{}, false, nil
		}
		return domain.ProductSize{}, false, err
	}
	return s, true, nil
}

func (r *ProductSizeRepository) Create(ctx context.Context, s domain.ProductSize) (domain.ProductSize, error) {
	query := `INSERT INTO product_sizes (id, product_id, size_label, size_code, price, sort_order)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, product_id, size_label, size_code, price, sort_order`

	err := r.pool.QueryRow(ctx, query,
		s.ID, s.ProductID, s.SizeLabel, s.SizeCode, s.Price, s.SortOrder,
	).Scan(&s.ID, &s.ProductID, &s.SizeLabel, &s.SizeCode, &s.Price, &s.SortOrder)

	return s, err
}

func (r *ProductSizeRepository) Update(ctx context.Context, id string, s domain.ProductSize) (domain.ProductSize, error) {
	query := `UPDATE product_sizes SET
		size_label = COALESCE(NULLIF($2, ''), size_label),
		price = CASE WHEN $3 > 0 THEN $3 ELSE price END,
		sort_order = $4
	WHERE id = $5
	RETURNING id, product_id, size_label, size_code, price, sort_order`

	var result domain.ProductSize
	err := r.pool.QueryRow(ctx, query,
		id, s.SizeLabel, s.Price, s.SortOrder, id,
	).Scan(&result.ID, &result.ProductID, &result.SizeLabel, &result.SizeCode, &result.Price, &result.SortOrder)

	return result, err
}

func (r *ProductSizeRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM product_sizes WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("product size not found")
	}
	return nil
}

func (r *ProductSizeRepository) DeleteByProduct(ctx context.Context, productID string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM product_sizes WHERE product_id = $1", productID)
	return err
}
