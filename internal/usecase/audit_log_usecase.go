package usecase

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type AuditLogUsecase struct {
	auditLogRepo domain.AuditLogRepository
}

func NewAuditLogUsecase(auditLogRepo domain.AuditLogRepository) *AuditLogUsecase {
	return &AuditLogUsecase{auditLogRepo: auditLogRepo}
}

// RecordEntry creates an audit log entry. It does not fail the operation if logging fails;
// errors are logged server-side only.
func (u *AuditLogUsecase) RecordEntry(ctx context.Context, entityType, entityID, action, actor string, before, after interface{}) {
	entry := domain.AuditLog{
		ID:         uuid.NewString(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Actor:      actor,
		Before:     before,
		After:      after,
	}

	if err := u.auditLogRepo.Create(ctx, entry); err != nil {
		log.Printf("failed to record audit log: %v", err)
	}
}

func (u *AuditLogUsecase) ListAuditLogs(ctx context.Context, limit int, offset int, entityType *string) ([]domain.AuditLog, error) {
	return u.auditLogRepo.List(ctx, limit, offset, entityType)
}
