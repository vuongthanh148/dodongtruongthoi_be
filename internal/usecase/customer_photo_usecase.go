package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type CustomerPhotoUsecase struct {
	customerPhotoRepo domain.CustomerPhotoRepository
	imageUploader     ImageUploader
}

func NewCustomerPhotoUsecase(
	customerPhotoRepo domain.CustomerPhotoRepository,
	imageUploader ImageUploader,
) *CustomerPhotoUsecase {
	return &CustomerPhotoUsecase{
		customerPhotoRepo: customerPhotoRepo,
		imageUploader:     imageUploader,
	}
}

func (u *CustomerPhotoUsecase) ListCustomerPhotos(ctx context.Context, includeInactive bool) ([]domain.CustomerPhoto, error) {
	if u.customerPhotoRepo == nil {
		return []domain.CustomerPhoto{}, nil
	}
	return u.customerPhotoRepo.List(ctx, includeInactive)
}

func (u *CustomerPhotoUsecase) CreateCustomerPhoto(ctx context.Context, file interface{}, filename string, caption *string, isActive bool, sortOrder int) (domain.CustomerPhoto, error) {
	if u.customerPhotoRepo == nil {
		return domain.CustomerPhoto{}, errors.New("customer photo feature not configured")
	}
	if u.imageUploader == nil {
		return domain.CustomerPhoto{}, errors.New("image uploads disabled; cloudinary not configured")
	}

	url, err := u.imageUploader.UploadImage(ctx, file, filename, "customer-photos")
	if err != nil {
		return domain.CustomerPhoto{}, fmt.Errorf("cloudinary upload failed: %w", err)
	}

	photo := domain.CustomerPhoto{
		ID:        uuid.NewString(),
		ImageURL:  url,
		Caption:   caption,
		SortOrder: sortOrder,
		IsActive:  isActive,
		CreatedAt: time.Now(),
	}

	result, err := u.customerPhotoRepo.Create(ctx, photo)
	if err != nil {
		return domain.CustomerPhoto{}, fmt.Errorf("failed to save customer photo: %w", err)
	}

	return result, nil
}

func (u *CustomerPhotoUsecase) UpdateCustomerPhoto(ctx context.Context, id string, caption *string, isActive bool, sortOrder int) (domain.CustomerPhoto, error) {
	if u.customerPhotoRepo == nil {
		return domain.CustomerPhoto{}, errors.New("customer photo feature not configured")
	}

	photo := domain.CustomerPhoto{
		Caption:   caption,
		SortOrder: sortOrder,
		IsActive:  isActive,
	}

	updated, err := u.customerPhotoRepo.Update(ctx, id, photo)
	if err != nil {
		return domain.CustomerPhoto{}, err
	}

	return updated, nil
}

func (u *CustomerPhotoUsecase) DeleteCustomerPhoto(ctx context.Context, id string) error {
	if u.customerPhotoRepo == nil {
		return errors.New("customer photo feature not configured")
	}

	return u.customerPhotoRepo.Delete(ctx, id)
}
