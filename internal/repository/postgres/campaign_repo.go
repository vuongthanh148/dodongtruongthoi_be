package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type CampaignRepository struct {
	pool *pgxpool.Pool
}

func NewCampaignRepository(pool *pgxpool.Pool) *CampaignRepository {
	return &CampaignRepository{pool: pool}
}

func (r *CampaignRepository) List(ctx context.Context, includeInactive bool) ([]domain.Campaign, error) {
	query := "SELECT id, name, description, discount_type, discount_value, starts_at, ends_at, is_active, created_at, updated_at FROM campaigns"

	if !includeInactive {
		query += " WHERE is_active = true"
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []domain.Campaign
	for rows.Next() {
		var c domain.Campaign
		err := rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.DiscountType, &c.DiscountValue,
			&c.StartsAt, &c.EndsAt, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, rows.Err()
}

func (r *CampaignRepository) Get(ctx context.Context, id string) (domain.Campaign, bool, error) {
	query := "SELECT id, name, description, discount_type, discount_value, starts_at, ends_at, is_active, created_at, updated_at FROM campaigns WHERE id = $1"

	var c domain.Campaign
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.DiscountType, &c.DiscountValue,
		&c.StartsAt, &c.EndsAt, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.Campaign{}, false, nil
		}
		return domain.Campaign{}, false, err
	}
	return c, true, nil
}

func (r *CampaignRepository) Create(ctx context.Context, c domain.Campaign) (domain.Campaign, error) {
	query := `INSERT INTO campaigns (id, name, description, discount_type, discount_value, starts_at, ends_at, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		c.ID, c.Name, c.Description, c.DiscountType, c.DiscountValue,
		c.StartsAt, c.EndsAt, c.IsActive, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.CreatedAt, &c.UpdatedAt)

	return c, err
}

func (r *CampaignRepository) Update(ctx context.Context, id string, c domain.Campaign) (domain.Campaign, error) {
	if id == "" {
		return domain.Campaign{}, errors.New("campaign id is required")
	}

	query := `UPDATE campaigns SET
		name = COALESCE(NULLIF($2, ''), name),
		description = $3,
		discount_type = COALESCE(NULLIF($4, ''), discount_type),
		discount_value = CASE WHEN $5 > 0 THEN $5 ELSE discount_value END,
		starts_at = CASE WHEN $6 IS NOT NULL THEN $6 ELSE starts_at END,
		ends_at = CASE WHEN $7 IS NOT NULL THEN $7 ELSE ends_at END,
		is_active = $8,
		updated_at = now()
	WHERE id = $9
	RETURNING id, name, description, discount_type, discount_value, starts_at, ends_at, is_active, created_at, updated_at`

	var result domain.Campaign
	err := r.pool.QueryRow(ctx, query,
		id, c.Name, c.Description, c.DiscountType, c.DiscountValue,
		c.StartsAt, c.EndsAt, c.IsActive, id,
	).Scan(
		&result.ID, &result.Name, &result.Description, &result.DiscountType, &result.DiscountValue,
		&result.StartsAt, &result.EndsAt, &result.IsActive, &result.CreatedAt, &result.UpdatedAt,
	)

	return result, err
}

func (r *CampaignRepository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, "UPDATE campaigns SET is_active = false, updated_at = now() WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("campaign not found")
	}
	return nil
}

func (r *CampaignRepository) SetProducts(ctx context.Context, campaignID string, productIDs []string) error {
	if campaignID == "" {
		return errors.New("campaign id is required")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM campaign_products WHERE campaign_id = $1", campaignID)
	if err != nil {
		return err
	}

	if len(productIDs) > 0 {
		stmt := "INSERT INTO campaign_products (campaign_id, product_id) VALUES "
		args := make([]interface{}, 0)
		for i, pid := range productIDs {
			if i > 0 {
				stmt += ", "
			}
			stmt += fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
			args = append(args, campaignID, pid)
		}
		_, err = tx.Exec(ctx, stmt, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *CampaignRepository) GetProductIDs(ctx context.Context, campaignID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, "SELECT product_id FROM campaign_products WHERE campaign_id = $1", campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productIDs []string
	for rows.Next() {
		var pid string
		if err := rows.Scan(&pid); err != nil {
			return nil, err
		}
		productIDs = append(productIDs, pid)
	}
	return productIDs, rows.Err()
}
