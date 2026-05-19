package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ImageRepository struct {
	pool *pgxpool.Pool
}

func NewImageRepository(pool *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{pool: pool}
}

func (r *ImageRepository) List(ctx context.Context) ([]domain.Image, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, url, cloudinary_public_id, created_at
		FROM images ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Image, 0)
	for rows.Next() {
		var img domain.Image
		if err := rows.Scan(&img.ID, &img.Name, &img.URL, &img.CloudinaryPublicID, &img.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, img)
	}
	return out, rows.Err()
}

func (r *ImageRepository) Get(ctx context.Context, id string) (domain.Image, bool, error) {
	var img domain.Image
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, url, cloudinary_public_id, created_at
		FROM images WHERE id = $1`, id).
		Scan(&img.ID, &img.Name, &img.URL, &img.CloudinaryPublicID, &img.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Image{}, false, nil
	}
	if err != nil {
		return domain.Image{}, false, err
	}
	return img, true, nil
}

func (r *ImageRepository) Create(ctx context.Context, img domain.Image) (domain.Image, error) {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO images (id, name, url, cloudinary_public_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, url, cloudinary_public_id, created_at`,
		img.ID, img.Name, img.URL, img.CloudinaryPublicID).
		Scan(&img.ID, &img.Name, &img.URL, &img.CloudinaryPublicID, &img.CreatedAt)
	return img, err
}

func (r *ImageRepository) Update(ctx context.Context, id string, name string) (domain.Image, error) {
	var img domain.Image
	err := r.pool.QueryRow(ctx, `
		UPDATE images SET name = $2 WHERE id = $1
		RETURNING id, name, url, cloudinary_public_id, created_at`,
		id, name).
		Scan(&img.ID, &img.Name, &img.URL, &img.CloudinaryPublicID, &img.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Image{}, errors.New("image not found")
	}
	return img, err
}

func (r *ImageRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM images WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("image not found")
	}
	return nil
}
