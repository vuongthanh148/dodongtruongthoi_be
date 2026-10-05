package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type LookupAttemptRepository struct {
	pool *pgxpool.Pool
}

func NewLookupAttemptRepository(pool *pgxpool.Pool) *LookupAttemptRepository {
	return &LookupAttemptRepository{pool: pool}
}

func (r *LookupAttemptRepository) Get(ctx context.Context, phone string) (domain.LookupAttempt, bool, error) {
	query := `SELECT phone, failed_count, locked_until, updated_at FROM lookup_attempts WHERE phone = $1`

	var a domain.LookupAttempt
	err := r.pool.QueryRow(ctx, query, phone).Scan(&a.Phone, &a.FailedCount, &a.LockedUntil, &a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.LookupAttempt{}, false, nil
	}
	if err != nil {
		return domain.LookupAttempt{}, false, err
	}
	return a, true, nil
}

func (r *LookupAttemptRepository) IncrementFailed(ctx context.Context, phone string) (int, error) {
	query := `INSERT INTO lookup_attempts (phone, failed_count, updated_at)
	VALUES ($1, 1, now())
	ON CONFLICT (phone) DO UPDATE
	SET failed_count = lookup_attempts.failed_count + 1, updated_at = now()
	RETURNING failed_count`

	var count int
	err := r.pool.QueryRow(ctx, query, phone).Scan(&count)
	return count, err
}

func (r *LookupAttemptRepository) Lock(ctx context.Context, phone string, until time.Time) error {
	query := `INSERT INTO lookup_attempts (phone, failed_count, locked_until, updated_at)
	VALUES ($1, 0, $2, now())
	ON CONFLICT (phone) DO UPDATE
	SET failed_count = 0, locked_until = $2, updated_at = now()`

	_, err := r.pool.Exec(ctx, query, phone, until)
	return err
}

func (r *LookupAttemptRepository) Reset(ctx context.Context, phone string) error {
	query := `INSERT INTO lookup_attempts (phone, failed_count, locked_until, updated_at)
	VALUES ($1, 0, NULL, now())
	ON CONFLICT (phone) DO UPDATE
	SET failed_count = 0, locked_until = NULL, updated_at = now()`

	_, err := r.pool.Exec(ctx, query, phone)
	return err
}

type LookupSessionRepository struct {
	pool *pgxpool.Pool
}

func NewLookupSessionRepository(pool *pgxpool.Pool) *LookupSessionRepository {
	return &LookupSessionRepository{pool: pool}
}

func (r *LookupSessionRepository) Create(ctx context.Context, s domain.LookupSession) error {
	query := `INSERT INTO lookup_sessions (token_hash, phone, order_id, expires_at) VALUES ($1, $2, $3, $4)`

	_, err := r.pool.Exec(ctx, query, s.TokenHash, s.Phone, s.OrderID, s.ExpiresAt)
	return err
}

func (r *LookupSessionRepository) Get(ctx context.Context, tokenHash string) (domain.LookupSession, bool, error) {
	query := `SELECT token_hash, phone, order_id, expires_at, created_at FROM lookup_sessions WHERE token_hash = $1`

	var s domain.LookupSession
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(&s.TokenHash, &s.Phone, &s.OrderID, &s.ExpiresAt, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.LookupSession{}, false, nil
	}
	if err != nil {
		return domain.LookupSession{}, false, err
	}
	return s, true, nil
}
