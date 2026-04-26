package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SiteSettingsRepository struct {
	pool *pgxpool.Pool
}

func NewSiteSettingsRepository(pool *pgxpool.Pool) *SiteSettingsRepository {
	return &SiteSettingsRepository{pool: pool}
}

func (r *SiteSettingsRepository) Get(ctx context.Context, key string) (string, error) {
	var value string
	err := r.pool.QueryRow(ctx, "SELECT value FROM site_settings WHERE key = $1", key).Scan(&value)
	return value, err
}

func (r *SiteSettingsRepository) GetAll(ctx context.Context) (map[string]string, error) {
	query := "SELECT key, value FROM site_settings"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		settings[key] = value
	}
	return settings, rows.Err()
}

func (r *SiteSettingsRepository) Set(ctx context.Context, key string, value string) error {
	query := `INSERT INTO site_settings (key, value, updated_at) VALUES ($1, $2, now())
	ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = now()`

	_, err := r.pool.Exec(ctx, query, key, value)
	return err
}

func (r *SiteSettingsRepository) SetBulk(ctx context.Context, settings map[string]string) error {
	for key, value := range settings {
		if err := r.Set(ctx, key, value); err != nil {
			return err
		}
	}
	return nil
}
