package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ContactUsecase struct {
	contactRepo domain.ContactLinkRepository
}

func NewContactUsecase(contactRepo domain.ContactLinkRepository) *ContactUsecase {
	return &ContactUsecase{
		contactRepo: contactRepo,
	}
}

func (u *ContactUsecase) ListContacts(ctx context.Context) ([]domain.ContactLink, error) {
	return u.contactRepo.List(ctx, false)
}

func (u *ContactUsecase) ListAllContacts(ctx context.Context) ([]domain.ContactLink, error) {
	return u.contactRepo.List(ctx, true)
}

func (u *ContactUsecase) CreateContact(ctx context.Context, cl domain.ContactLink) (domain.ContactLink, error) {
	if cl.ID == "" || cl.Platform == "" || cl.URL == "" {
		return domain.ContactLink{}, errors.New("id, platform, and url are required")
	}

	now := time.Now()
	cl.CreatedAt = now
	cl.IsActive = true

	return u.contactRepo.Create(ctx, cl)
}

func (u *ContactUsecase) UpdateContact(ctx context.Context, id string, updates domain.ContactLink) (domain.ContactLink, error) {
	if id == "" {
		return domain.ContactLink{}, errors.New("contact id is required")
	}

	// Get existing contact
	existing, ok, err := u.contactRepo.Get(ctx, id)
	if err != nil {
		return domain.ContactLink{}, err
	}
	if !ok {
		return domain.ContactLink{}, errors.New("contact not found")
	}

	if updates.Platform != "" {
		existing.Platform = updates.Platform
	}
	if updates.Label != "" {
		existing.Label = updates.Label
	}
	if updates.URL != "" {
		existing.URL = updates.URL
	}
	existing.SortOrder = updates.SortOrder
	existing.IsActive = updates.IsActive

	return u.contactRepo.Update(ctx, id, existing)
}

func (u *ContactUsecase) DeleteContact(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("contact id is required")
	}
	return u.contactRepo.Delete(ctx, id)
}
