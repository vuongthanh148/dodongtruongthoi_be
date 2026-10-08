package usecase

import (
	"context"
	"testing"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

// fakeAuditLogRepo is an in-memory domain.AuditLogRepository.
type fakeAuditLogRepo struct {
	entries []domain.AuditLog
}

func (f *fakeAuditLogRepo) Create(ctx context.Context, entry domain.AuditLog) error {
	f.entries = append(f.entries, entry)
	return nil
}

func (f *fakeAuditLogRepo) List(ctx context.Context, limit int, offset int, entityType *string) ([]domain.AuditLog, error) {
	return f.entries, nil
}

func TestRecordEntry_Success(t *testing.T) {
	repo := &fakeAuditLogRepo{}
	u := NewAuditLogUsecase(repo)

	before := map[string]string{"status": "pending_confirm"}
	after := map[string]string{"status": "confirmed"}

	u.RecordEntry(context.Background(), "order", "test-order-id", "status_change", "admin", before, after)

	if len(repo.entries) != 1 {
		t.Fatalf("Create called %d times, want 1", len(repo.entries))
	}

	entry := repo.entries[0]
	if entry.EntityType != "order" {
		t.Errorf("EntityType = %q, want %q", entry.EntityType, "order")
	}
	if entry.EntityID != "test-order-id" {
		t.Errorf("EntityID = %q, want %q", entry.EntityID, "test-order-id")
	}
	if entry.Action != "status_change" {
		t.Errorf("Action = %q, want %q", entry.Action, "status_change")
	}
	if entry.Actor != "admin" {
		t.Errorf("Actor = %q, want %q", entry.Actor, "admin")
	}
}

func TestListAuditLogs_Empty(t *testing.T) {
	repo := &fakeAuditLogRepo{}
	u := NewAuditLogUsecase(repo)

	logs, err := u.ListAuditLogs(context.Background(), 10, 0, nil)
	if err != nil {
		t.Fatalf("ListAuditLogs: unexpected error: %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("got %d logs, want 0", len(logs))
	}
}
