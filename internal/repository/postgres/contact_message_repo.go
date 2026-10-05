package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ContactMessageRepository struct {
	pool *pgxpool.Pool
}

func NewContactMessageRepository(pool *pgxpool.Pool) *ContactMessageRepository {
	return &ContactMessageRepository{pool: pool}
}

func (r *ContactMessageRepository) Create(ctx context.Context, m domain.ContactMessage) (domain.ContactMessage, error) {
	query := `INSERT INTO contact_messages (id, name, phone, message, handled)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, name, phone, message, handled, created_at`

	err := r.pool.QueryRow(ctx, query,
		m.ID, m.Name, m.Phone, m.Message, m.Handled,
	).Scan(&m.ID, &m.Name, &m.Phone, &m.Message, &m.Handled, &m.CreatedAt)

	return m, err
}

func (r *ContactMessageRepository) List(ctx context.Context, handled *bool) ([]domain.ContactMessage, error) {
	query := "SELECT id, name, phone, message, handled, created_at FROM contact_messages"
	var args []any
	if handled != nil {
		query += " WHERE handled = $1"
		args = append(args, *handled)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []domain.ContactMessage{}
	for rows.Next() {
		var m domain.ContactMessage
		if err := rows.Scan(&m.ID, &m.Name, &m.Phone, &m.Message, &m.Handled, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (r *ContactMessageRepository) SetHandled(ctx context.Context, id string, handled bool) (domain.ContactMessage, error) {
	query := `UPDATE contact_messages SET handled = $2 WHERE id = $1
	RETURNING id, name, phone, message, handled, created_at`

	var m domain.ContactMessage
	err := r.pool.QueryRow(ctx, query, id, handled).
		Scan(&m.ID, &m.Name, &m.Phone, &m.Message, &m.Handled, &m.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ContactMessage{}, domain.ErrNotFound
	}
	return m, err
}
