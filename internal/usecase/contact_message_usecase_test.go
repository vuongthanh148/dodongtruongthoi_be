package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

// fakeContactMessageRepo is an in-memory domain.ContactMessageRepository.
type fakeContactMessageRepo struct {
	created []domain.ContactMessage
}

func (f *fakeContactMessageRepo) Create(ctx context.Context, m domain.ContactMessage) (domain.ContactMessage, error) {
	f.created = append(f.created, m)
	return m, nil
}

func (f *fakeContactMessageRepo) List(ctx context.Context, handled *bool) ([]domain.ContactMessage, error) {
	return f.created, nil
}

func (f *fakeContactMessageRepo) SetHandled(ctx context.Context, id string, handled bool) (domain.ContactMessage, error) {
	return domain.ContactMessage{}, domain.ErrNotFound
}

const validContactMessage = "Tôi muốn hỏi về trống đồng cỡ lớn."

func TestSubmitContactMessage_Validation(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		phone   string
		message string
		wantMsg string
	}{
		{name: "empty name", in: "   ", phone: "0912345678", message: validContactMessage, wantMsg: "name is required"},
		{name: "name too long", in: strings.Repeat("a", 101), phone: "0912345678", message: validContactMessage, wantMsg: "name must be at most 100 characters"},
		{name: "empty phone", in: "An", phone: "", message: validContactMessage, wantMsg: "phone is required"},
		{name: "phone with letters", in: "An", phone: "09123abc78", message: validContactMessage, wantMsg: "phone may only contain digits, spaces and +"},
		{name: "phone too short", in: "An", phone: "09 1234", message: validContactMessage, wantMsg: "phone must be 9 to 15 characters after removing spaces"},
		{name: "phone too long", in: "An", phone: "0912345678901234", message: validContactMessage, wantMsg: "phone must be 9 to 15 characters after removing spaces"},
		{name: "empty message", in: "An", phone: "0912345678", message: "  ", wantMsg: "message is required"},
		{name: "short message", in: "An", phone: "0912345678", message: "too short", wantMsg: "message must be at least 10 characters"},
		{name: "message too long", in: "An", phone: "0912345678", message: strings.Repeat("x", 2001), wantMsg: "message must be at most 2000 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeContactMessageRepo{}
			u := NewContactMessageUsecase(repo)

			_, err := u.SubmitContactMessage(context.Background(), tt.in, tt.phone, tt.message)
			if err == nil {
				t.Fatalf("SubmitContactMessage returned nil error, want %q", tt.wantMsg)
			}
			if err.Error() != tt.wantMsg {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantMsg)
			}
			if !errors.Is(err, domain.ErrInvalidInput) {
				t.Errorf("error does not wrap domain.ErrInvalidInput: %v", err)
			}
			if len(repo.created) != 0 {
				t.Errorf("Create called %d times for rejected input, want 0", len(repo.created))
			}
		})
	}
}

func TestSubmitContactMessage_Success(t *testing.T) {
	repo := &fakeContactMessageRepo{}
	u := NewContactMessageUsecase(repo)

	got, err := u.SubmitContactMessage(context.Background(), "  Nguyễn An  ", " +84 912 345 678 ", "  "+validContactMessage+"  ")
	if err != nil {
		t.Fatalf("SubmitContactMessage: unexpected error: %v", err)
	}
	if got.ID == "" {
		t.Error("ID is empty, want generated uuid")
	}
	if got.Name != "Nguyễn An" {
		t.Errorf("Name = %q, want trimmed %q", got.Name, "Nguyễn An")
	}
	if got.Phone != "+84 912 345 678" {
		t.Errorf("Phone = %q, want trimmed value", got.Phone)
	}
	if got.Message != validContactMessage {
		t.Errorf("Message = %q, want trimmed value", got.Message)
	}
	if got.Handled {
		t.Error("Handled = true, want false for new message")
	}
	if len(repo.created) != 1 {
		t.Fatalf("Create called %d times, want 1", len(repo.created))
	}
}

func TestSetContactMessageHandled_MalformedIDIsNotFound(t *testing.T) {
	u := NewContactMessageUsecase(&fakeContactMessageRepo{})

	_, err := u.SetContactMessageHandled(context.Background(), "not-a-uuid", true)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("error = %v, want domain.ErrNotFound", err)
	}
}
