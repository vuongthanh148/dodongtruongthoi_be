package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ReviewRepository struct {
	pool *pgxpool.Pool
}

func NewReviewRepository(pool *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{pool: pool}
}

func (r *ReviewRepository) List(ctx context.Context, productID string) ([]domain.Review, error) {
	query := "SELECT id, product_id, reviewer_name, rating, body, is_approved, created_at FROM reviews WHERE product_id = $1 AND is_approved = true ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var rev domain.Review
		err := rows.Scan(
			&rev.ID, &rev.ProductID, &rev.ReviewerName, &rev.Rating, &rev.Body, &rev.IsApproved, &rev.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}

func (r *ReviewRepository) ListAll(ctx context.Context, approved *bool) ([]domain.Review, error) {
	query := "SELECT id, product_id, reviewer_name, rating, body, is_approved, created_at FROM reviews"

	if approved != nil {
		if *approved {
			query += " WHERE is_approved = true"
		} else {
			query += " WHERE is_approved = false"
		}
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var rev domain.Review
		err := rows.Scan(
			&rev.ID, &rev.ProductID, &rev.ReviewerName, &rev.Rating, &rev.Body, &rev.IsApproved, &rev.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}

func (r *ReviewRepository) Create(ctx context.Context, rev domain.Review) (domain.Review, error) {
	query := `INSERT INTO reviews (id, product_id, reviewer_name, rating, body, is_approved, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, product_id, reviewer_name, rating, body, is_approved, created_at`

	err := r.pool.QueryRow(ctx, query,
		rev.ID, rev.ProductID, rev.ReviewerName, rev.Rating, rev.Body, rev.IsApproved, rev.CreatedAt,
	).Scan(&rev.ID, &rev.ProductID, &rev.ReviewerName, &rev.Rating, &rev.Body, &rev.IsApproved, &rev.CreatedAt)

	return rev, err
}

func (r *ReviewRepository) UpdateApproval(ctx context.Context, id string, isApproved bool) error {
	result, err := r.pool.Exec(ctx, "UPDATE reviews SET is_approved = $1 WHERE id = $2", isApproved, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("review not found")
	}
	return nil
}

func (r *ReviewRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "DELETE FROM reviews WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("review not found")
	}
	return nil
}
