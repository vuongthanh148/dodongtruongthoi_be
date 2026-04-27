package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type CustomerPhotoRepository struct {
	pool *pgxpool.Pool
}

func NewCustomerPhotoRepository(pool *pgxpool.Pool) *CustomerPhotoRepository {
	return &CustomerPhotoRepository{pool: pool}
}

func (r *CustomerPhotoRepository) List(ctx context.Context, includeInactive bool) ([]domain.CustomerPhoto, error) {
	query := "SELECT id, image_url, caption, sort_order, is_active, created_at FROM customer_photos"
	if !includeInactive {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY sort_order, created_at DESC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []domain.CustomerPhoto
	for rows.Next() {
		var p domain.CustomerPhoto
		if err := rows.Scan(&p.ID, &p.ImageURL, &p.Caption, &p.SortOrder, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

func (r *CustomerPhotoRepository) Get(ctx context.Context, id string) (domain.CustomerPhoto, bool, error) {
	query := "SELECT id, image_url, caption, sort_order, is_active, created_at FROM customer_photos WHERE id = $1"

	var p domain.CustomerPhoto
	err := r.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.ImageURL, &p.Caption, &p.SortOrder, &p.IsActive, &p.CreatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.CustomerPhoto{}, false, nil
		}
		return domain.CustomerPhoto{}, false, err
	}
	return p, true, nil
}

func (r *CustomerPhotoRepository) Create(ctx context.Context, p domain.CustomerPhoto) (domain.CustomerPhoto, error) {
	query := `INSERT INTO customer_photos (id, image_url, caption, sort_order, is_active, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, image_url, caption, sort_order, is_active, created_at`

	err := r.pool.QueryRow(ctx, query,
		p.ID, p.ImageURL, p.Caption, p.SortOrder, p.IsActive, p.CreatedAt,
	).Scan(&p.ID, &p.ImageURL, &p.Caption, &p.SortOrder, &p.IsActive, &p.CreatedAt)

	return p, err
}

func (r *CustomerPhotoRepository) Update(ctx context.Context, id string, p domain.CustomerPhoto) (domain.CustomerPhoto, error) {
	if id == "" {
		return domain.CustomerPhoto{}, errors.New("customer photo id is required")
	}

	query := `UPDATE customer_photos
	SET caption = $2, sort_order = $3, is_active = $4
	WHERE id = $1
	RETURNING id, image_url, caption, sort_order, is_active, created_at`

	err := r.pool.QueryRow(ctx, query, id, p.Caption, p.SortOrder, p.IsActive).
		Scan(&p.ID, &p.ImageURL, &p.Caption, &p.SortOrder, &p.IsActive, &p.CreatedAt)

	return p, err
}

func (r *CustomerPhotoRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("customer photo id is required")
	}
	_, err := r.pool.Exec(ctx, "DELETE FROM customer_photos WHERE id = $1", id)
	return err
}
