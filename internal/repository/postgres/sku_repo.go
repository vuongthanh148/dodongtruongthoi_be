package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ProductSKURepository struct {
	pool *pgxpool.Pool
}

func NewProductSKURepository(pool *pgxpool.Pool) *ProductSKURepository {
	return &ProductSKURepository{pool: pool}
}

func (r *ProductSKURepository) ListByProduct(ctx context.Context, productID string) ([]domain.ProductSKU, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, product_id, size_code, attrs, price, sort_order, created_at
		FROM product_skus
		WHERE product_id = $1
		ORDER BY sort_order, created_at`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.ProductSKU, 0)
	for rows.Next() {
		var s domain.ProductSKU
		var attrsJSON []byte
		if err := rows.Scan(&s.ID, &s.ProductID, &s.SizeCode, &attrsJSON, &s.Price, &s.SortOrder, &s.CreatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(attrsJSON, &s.Attrs); err != nil {
			s.Attrs = map[string]string{}
		}
		if s.Attrs == nil {
			s.Attrs = map[string]string{}
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *ProductSKURepository) SetByProduct(ctx context.Context, productID string, skus []domain.ProductSKU) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if _, err := tx.Exec(ctx, `DELETE FROM product_skus WHERE product_id = $1`, productID); err != nil {
		return err
	}

	for i, s := range skus {
		attrsJSON, err := json.Marshal(s.Attrs)
		if err != nil {
			return err
		}
		id := fmt.Sprintf("sku_%s_%d_%d", productID, time.Now().UnixNano(), i)
		if _, err := tx.Exec(ctx,
			`INSERT INTO product_skus (id, product_id, size_code, attrs, price, sort_order, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, now())`,
			id, productID, s.SizeCode, attrsJSON, s.Price, i,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
