package usecase

import (
	"context"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

const (
	contactNameMaxLen    = 100
	contactMessageMinLen = 10
	contactMessageMaxLen = 2000
	contactPhoneMinLen   = 9
	contactPhoneMaxLen   = 15
)

// contactPhonePattern allows digits, spaces and '+' only.
var contactPhonePattern = regexp.MustCompile(`^[0-9+ ]+$`)

// contactValidationError is a user-facing input error. It unwraps to
// domain.ErrInvalidInput so handlers can map it to 400.
type contactValidationError struct {
	msg string
}

func (e contactValidationError) Error() string { return e.msg }

func (e contactValidationError) Unwrap() error { return domain.ErrInvalidInput }

type ContactMessageUsecase struct {
	contactMessageRepo domain.ContactMessageRepository
}

func NewContactMessageUsecase(contactMessageRepo domain.ContactMessageRepository) *ContactMessageUsecase {
	return &ContactMessageUsecase{contactMessageRepo: contactMessageRepo}
}

func (u *ContactMessageUsecase) SubmitContactMessage(ctx context.Context, name, phone, message string) (domain.ContactMessage, error) {
	name = strings.TrimSpace(name)
	phone = strings.TrimSpace(phone)
	message = strings.TrimSpace(message)

	if name == "" {
		return domain.ContactMessage{}, contactValidationError{"name is required"}
	}
	if utf8.RuneCountInString(name) > contactNameMaxLen {
		return domain.ContactMessage{}, contactValidationError{"name must be at most 100 characters"}
	}

	if phone == "" {
		return domain.ContactMessage{}, contactValidationError{"phone is required"}
	}
	if !contactPhonePattern.MatchString(phone) {
		return domain.ContactMessage{}, contactValidationError{"phone may only contain digits, spaces and +"}
	}
	digits := strings.ReplaceAll(phone, " ", "")
	if len(digits) < contactPhoneMinLen || len(digits) > contactPhoneMaxLen {
		return domain.ContactMessage{}, contactValidationError{"phone must be 9 to 15 characters after removing spaces"}
	}

	msgLen := utf8.RuneCountInString(message)
	if message == "" {
		return domain.ContactMessage{}, contactValidationError{"message is required"}
	}
	if msgLen < contactMessageMinLen {
		return domain.ContactMessage{}, contactValidationError{"message must be at least 10 characters"}
	}
	if msgLen > contactMessageMaxLen {
		return domain.ContactMessage{}, contactValidationError{"message must be at most 2000 characters"}
	}

	return u.contactMessageRepo.Create(ctx, domain.ContactMessage{
		ID:      uuid.NewString(),
		Name:    name,
		Phone:   phone,
		Message: message,
		Handled: false,
	})
}

func (u *ContactMessageUsecase) ListContactMessages(ctx context.Context, handled *bool) ([]domain.ContactMessage, error) {
	return u.contactMessageRepo.List(ctx, handled)
}

// SetContactMessageHandled returns domain.ErrNotFound for unknown or malformed ids.
func (u *ContactMessageUsecase) SetContactMessageHandled(ctx context.Context, id string, handled bool) (domain.ContactMessage, error) {
	if _, err := uuid.Parse(id); err != nil {
		return domain.ContactMessage{}, domain.ErrNotFound
	}
	return u.contactMessageRepo.SetHandled(ctx, id, handled)
}
