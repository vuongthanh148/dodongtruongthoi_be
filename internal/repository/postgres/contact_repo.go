package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ContactLinkRepository struct {
	pool *pgxpool.Pool
}

func NewContactLinkRepository(pool *pgxpool.Pool) *ContactLinkRepository {
	return &ContactLinkRepository{pool: pool}
}

func (r *ContactLinkRepository) List(ctx context.Context, includeInactive bool) ([]domain.ContactLink, error) {
	query := "SELECT id, platform, label, url, sort_order, is_active, created_at FROM contact_links"

	if !includeInactive {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY sort_order, created_at ASC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []domain.ContactLink
	for rows.Next() {
		var c domain.ContactLink
		err := rows.Scan(
			&c.ID, &c.Platform, &c.Label, &c.URL, &c.SortOrder, &c.IsActive, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}

func (r *ContactLinkRepository) Get(ctx context.Context, id string) (domain.ContactLink, bool, error) {
	query := "SELECT id, platform, label, url, sort_order, is_active, created_at FROM contact_links WHERE id = $1"

	var c domain.ContactLink
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Platform, &c.Label, &c.URL, &c.SortOrder, &c.IsActive, &c.CreatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.ContactLink{}, false, nil
		}
		return domain.ContactLink{}, false, err
	}
	return c, true, nil
}

func (r *ContactLinkRepository) Create(ctx context.Context, c domain.ContactLink) (domain.ContactLink, error) {
	query := `INSERT INTO contact_links (id, platform, label, url, sort_order, is_active, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, platform, label, url, sort_order, is_active, created_at`

	err := r.pool.QueryRow(ctx, query,
		c.ID, c.Platform, c.Label, c.URL, c.SortOrder, c.IsActive, c.CreatedAt,
	).Scan(&c.ID, &c.Platform, &c.Label, &c.URL, &c.SortOrder, &c.IsActive, &c.CreatedAt)

	return c, err
}

func (r *ContactLinkRepository) Update(ctx context.Context, id string, c domain.ContactLink) (domain.ContactLink, error) {
	if id == "" {
		return domain.ContactLink{}, errors.New("contact link id is required")
	}

	query := `UPDATE contact_links SET
		platform = COALESCE(NULLIF($2, ''), platform),
		label = COALESCE(NULLIF($3, ''), label),
		url = COALESCE(NULLIF($4, ''), url),
		sort_order = $5,
		is_active = $6
	WHERE id = $7
	RETURNING id, platform, label, url, sort_order, is_active, created_at`

	var result domain.ContactLink
	err := r.pool.QueryRow(ctx, query,
		id, c.Platform, c.Label, c.URL, c.SortOrder, c.IsActive, id,
	).Scan(
		&result.ID, &result.Platform, &result.Label, &result.URL, &result.SortOrder, &result.IsActive, &result.CreatedAt,
	)

	return result, err
}

func (r *ContactLinkRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "UPDATE contact_links SET is_active = false WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("contact link not found")
	}
	return nil
}
