package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type AdminUserRepository struct {
	pool *pgxpool.Pool
}

func NewAdminUserRepository(pool *pgxpool.Pool) *AdminUserRepository {
	return &AdminUserRepository{pool: pool}
}

func (r *AdminUserRepository) GetByUsername(ctx context.Context, username string) (domain.AdminUser, bool, error) {
	query := "SELECT id, username, password_hash, display_name, is_active, created_at, last_login_at FROM admin_users WHERE username = $1"

	var user domain.AdminUser
	err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.IsActive, &user.CreatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.AdminUser{}, false, nil
		}
		return domain.AdminUser{}, false, err
	}
	return user, true, nil
}

func (r *AdminUserRepository) Get(ctx context.Context, id string) (domain.AdminUser, bool, error) {
	query := "SELECT id, username, password_hash, display_name, is_active, created_at, last_login_at FROM admin_users WHERE id = $1"

	var user domain.AdminUser
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.IsActive, &user.CreatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.AdminUser{}, false, nil
		}
		return domain.AdminUser{}, false, err
	}
	return user, true, nil
}

func (r *AdminUserRepository) Create(ctx context.Context, user domain.AdminUser) (domain.AdminUser, error) {
	query := `INSERT INTO admin_users (id, username, password_hash, display_name, is_active, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, username, password_hash, display_name, is_active, created_at, last_login_at`

	err := r.pool.QueryRow(ctx, query,
		user.ID, user.Username, user.PasswordHash, user.DisplayName, user.IsActive, user.CreatedAt,
	).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName, &user.IsActive, &user.CreatedAt, &user.LastLoginAt,
	)

	return user, err
}

func (r *AdminUserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "UPDATE admin_users SET last_login_at = now() WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("admin user not found")
	}
	return nil
}

func (r *AdminUserRepository) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	result, err := r.pool.Exec(ctx, "UPDATE admin_users SET password_hash = $1 WHERE id = $2", passwordHash, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("admin user not found")
	}
	return nil
}
