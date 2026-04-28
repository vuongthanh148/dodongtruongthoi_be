package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type BannerRepository struct {
	pool *pgxpool.Pool
}

func NewBannerRepository(pool *pgxpool.Pool) *BannerRepository {
	return &BannerRepository{pool: pool}
}

func (r *BannerRepository) List(ctx context.Context, includeInactive bool) ([]domain.Banner, error) {
	query := "SELECT id, title, subtitle, image_url, link_url, sort_order, is_active, created_at, updated_at FROM banners"

	if !includeInactive {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY sort_order, created_at DESC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []domain.Banner
	for rows.Next() {
		var b domain.Banner
		err := rows.Scan(
			&b.ID, &b.Title, &b.Subtitle, &b.ImageURL, &b.LinkURL, &b.SortOrder, &b.IsActive, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		banners = append(banners, b)
	}
	return banners, rows.Err()
}

func (r *BannerRepository) Get(ctx context.Context, id string) (domain.Banner, bool, error) {
	query := "SELECT id, title, subtitle, image_url, link_url, sort_order, is_active, created_at, updated_at FROM banners WHERE id = $1"

	var b domain.Banner
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.Title, &b.Subtitle, &b.ImageURL, &b.LinkURL, &b.SortOrder, &b.IsActive, &b.CreatedAt, &b.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.Banner{}, false, nil
		}
		return domain.Banner{}, false, err
	}
	return b, true, nil
}

func (r *BannerRepository) Create(ctx context.Context, b domain.Banner) (domain.Banner, error) {
	query := `INSERT INTO banners (id, title, subtitle, image_url, link_url, sort_order, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, title, subtitle, image_url, link_url, sort_order, is_active, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		b.ID, b.Title, b.Subtitle, b.ImageURL, b.LinkURL, b.SortOrder, b.IsActive, b.CreatedAt, b.UpdatedAt,
	).Scan(&b.ID, &b.Title, &b.Subtitle, &b.ImageURL, &b.LinkURL, &b.SortOrder, &b.IsActive, &b.CreatedAt, &b.UpdatedAt)

	return b, err
}

func (r *BannerRepository) Update(ctx context.Context, id string, b domain.Banner) (domain.Banner, error) {
	if id == "" {
		return domain.Banner{}, errors.New("banner id is required")
	}

	query := `UPDATE banners SET
		title = COALESCE(NULLIF($1, ''), title),
		subtitle = COALESCE(NULLIF($2, ''), subtitle),
		image_url = COALESCE(NULLIF($3, ''), image_url),
		link_url = COALESCE(NULLIF($4, ''), link_url),
		sort_order = $5,
		is_active = $6,
		updated_at = now()
	WHERE id = $7
	RETURNING id, title, subtitle, image_url, link_url, sort_order, is_active, created_at, updated_at`

	var result domain.Banner
	err := r.pool.QueryRow(ctx, query,
		b.Title, b.Subtitle, b.ImageURL, b.LinkURL, b.SortOrder, b.IsActive, id,
	).Scan(
		&result.ID, &result.Title, &result.Subtitle, &result.ImageURL, &result.LinkURL, &result.SortOrder, &result.IsActive, &result.CreatedAt, &result.UpdatedAt,
	)

	return result, err
}

func (r *BannerRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "UPDATE banners SET is_active = false, updated_at = now() WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("banner not found")
	}
	return nil
}
