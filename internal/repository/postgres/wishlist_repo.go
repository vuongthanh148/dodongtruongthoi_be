package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type WishlistRepository struct {
	pool *pgxpool.Pool
}

func NewWishlistRepository(pool *pgxpool.Pool) *WishlistRepository {
	return &WishlistRepository{pool: pool}
}

func (r *WishlistRepository) GetByPhone(ctx context.Context, phone string) ([]domain.WishlistItem, error) {
	query := "SELECT id, phone, product_id, created_at FROM wishlists WHERE phone = $1 ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, phone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []domain.WishlistItem
	for rows.Next() {
		var w domain.WishlistItem
		err := rows.Scan(&w.ID, &w.Phone, &w.ProductID, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		wishlists = append(wishlists, w)
	}
	return wishlists, rows.Err()
}

func (r *WishlistRepository) Add(ctx context.Context, w domain.WishlistItem) (domain.WishlistItem, error) {
	query := `INSERT INTO wishlists (id, phone, product_id, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id, phone, product_id, created_at`

	err := r.pool.QueryRow(ctx, query, w.ID, w.Phone, w.ProductID, w.CreatedAt).Scan(
		&w.ID, &w.Phone, &w.ProductID, &w.CreatedAt,
	)

	return w, err
}

func (r *WishlistRepository) Remove(ctx context.Context, phone string, productID string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM wishlists WHERE phone = $1 AND product_id = $2", phone, productID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("wishlist item not found")
	}
	return nil
}

func (r *WishlistRepository) RemoveByPhone(ctx context.Context, phone string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM wishlists WHERE phone = $1", phone)
	return err
}
