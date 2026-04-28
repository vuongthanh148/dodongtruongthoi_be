package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type BannerUsecase struct {
	bannerRepo domain.BannerRepository
}

func NewBannerUsecase(bannerRepo domain.BannerRepository) *BannerUsecase {
	return &BannerUsecase{
		bannerRepo: bannerRepo,
	}
}

func (u *BannerUsecase) ListBanners(ctx context.Context) ([]domain.Banner, error) {
	return u.bannerRepo.List(ctx, false)
}

func (u *BannerUsecase) ListAllBanners(ctx context.Context) ([]domain.Banner, error) {
	return u.bannerRepo.List(ctx, true)
}

func (u *BannerUsecase) CreateBanner(ctx context.Context, b domain.Banner) (domain.Banner, error) {
	if b.ID == "" || b.Title == nil || *b.Title == "" {
		return domain.Banner{}, errors.New("id and title are required")
	}

	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	b.IsActive = true

	return u.bannerRepo.Create(ctx, b)
}

func (u *BannerUsecase) UpdateBanner(ctx context.Context, id string, updates domain.Banner) (domain.Banner, error) {
	if id == "" {
		return domain.Banner{}, errors.New("banner id is required")
	}

	// Get existing banner
	existing, ok, err := u.bannerRepo.Get(ctx, id)
	if err != nil {
		return domain.Banner{}, err
	}
	if !ok {
		return domain.Banner{}, errors.New("banner not found")
	}

	if updates.Title != nil && *updates.Title != "" {
		existing.Title = updates.Title
	}
	if updates.Subtitle != nil && *updates.Subtitle != "" {
		existing.Subtitle = updates.Subtitle
	}
	if updates.ImageURL != nil && *updates.ImageURL != "" {
		existing.ImageURL = updates.ImageURL
	}
	if updates.LinkURL != nil && *updates.LinkURL != "" {
		existing.LinkURL = updates.LinkURL
	}
	existing.SortOrder = updates.SortOrder
	existing.IsActive = updates.IsActive
	existing.UpdatedAt = time.Now()

	return u.bannerRepo.Update(ctx, id, existing)
}

func (u *BannerUsecase) DeleteBanner(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("banner id is required")
	}
	return u.bannerRepo.Delete(ctx, id)
}
