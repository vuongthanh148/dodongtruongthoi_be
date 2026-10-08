package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type AuditLogRepository struct {
	pool *pgxpool.Pool
}

func NewAuditLogRepository(pool *pgxpool.Pool) *AuditLogRepository {
	return &AuditLogRepository{pool: pool}
}

func (r *AuditLogRepository) Create(ctx context.Context, entry domain.AuditLog) error {
	var beforeJSON, afterJSON *string
	if entry.Before != nil {
		b, err := json.Marshal(entry.Before)
		if err != nil {
			return err
		}
		s := string(b)
		beforeJSON = &s
	}
	if entry.After != nil {
		b, err := json.Marshal(entry.After)
		if err != nil {
			return err
		}
		s := string(b)
		afterJSON = &s
	}

	query := `INSERT INTO audit_logs (id, entity_type, entity_id, action, actor, before, after)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		entry.ID, entry.EntityType, entry.EntityID, entry.Action, entry.Actor, beforeJSON, afterJSON,
	)
	return err
}

func (r *AuditLogRepository) List(ctx context.Context, limit int, offset int, entityType *string) ([]domain.AuditLog, error) {
	query := "SELECT id, entity_type, entity_id, action, actor, before, after, created_at FROM audit_logs"
	var args []any
	argCount := 1

	if entityType != nil {
		query += fmt.Sprintf(" WHERE entity_type = $%d", argCount)
		args = append(args, *entityType)
		argCount++
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
		argCount++
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.AuditLog
	for rows.Next() {
		var entry domain.AuditLog
		var beforeJSON, afterJSON *string

		err := rows.Scan(&entry.ID, &entry.EntityType, &entry.EntityID, &entry.Action, &entry.Actor,
			&beforeJSON, &afterJSON, &entry.CreatedAt)
		if err != nil {
			return nil, err
		}

		if beforeJSON != nil {
			var before interface{}
			if err := json.Unmarshal([]byte(*beforeJSON), &before); err != nil {
				return nil, err
			}
			entry.Before = before
		}

		if afterJSON != nil {
			var after interface{}
			if err := json.Unmarshal([]byte(*afterJSON), &after); err != nil {
				return nil, err
			}
			entry.After = after
		}

		logs = append(logs, entry)
	}

	return logs, rows.Err()
}
