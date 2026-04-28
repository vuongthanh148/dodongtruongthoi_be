package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type CampaignUsecase struct {
	campaignRepo domain.CampaignRepository
}

func NewCampaignUsecase(campaignRepo domain.CampaignRepository) *CampaignUsecase {
	return &CampaignUsecase{
		campaignRepo: campaignRepo,
	}
}

func (u *CampaignUsecase) CreateCampaign(ctx context.Context, c domain.Campaign) (domain.Campaign, error) {
	if c.ID == "" || c.Name == "" {
		return domain.Campaign{}, errors.New("id and name are required")
	}

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.IsActive == false && c.ID != "" {
		// Keep explicit false
	} else {
		c.IsActive = true
	}

	return u.campaignRepo.Create(ctx, c)
}

func (u *CampaignUsecase) UpdateCampaign(ctx context.Context, id string, updates domain.Campaign) (domain.Campaign, error) {
	if id == "" {
		return domain.Campaign{}, errors.New("campaign id is required")
	}

	// Get existing campaign
	existing, ok, err := u.campaignRepo.Get(ctx, id)
	if err != nil {
		return domain.Campaign{}, err
	}
	if !ok {
		return domain.Campaign{}, errors.New("campaign not found")
	}

	if updates.Name != "" {
		existing.Name = updates.Name
	}
	if updates.Description != nil {
		existing.Description = updates.Description
	}
	if updates.DiscountType != "" {
		existing.DiscountType = updates.DiscountType
	}
	if updates.DiscountValue > 0 {
		existing.DiscountValue = updates.DiscountValue
	}
	if !updates.StartsAt.IsZero() {
		existing.StartsAt = updates.StartsAt
	}
	if !updates.EndsAt.IsZero() {
		existing.EndsAt = updates.EndsAt
	}
	existing.IsActive = updates.IsActive
	existing.UpdatedAt = time.Now()

	return u.campaignRepo.Update(ctx, id, existing)
}

func (u *CampaignUsecase) DeleteCampaign(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("campaign id is required")
	}
	return u.campaignRepo.Delete(ctx, id)
}

func (u *CampaignUsecase) SetCampaignProducts(ctx context.Context, campaignID string, productIDs []string) error {
	if campaignID == "" {
		return errors.New("campaign id is required")
	}
	return u.campaignRepo.SetProducts(ctx, campaignID, productIDs)
}

func (u *CampaignUsecase) ListAllCampaigns(ctx context.Context) ([]domain.Campaign, error) {
	return u.campaignRepo.List(ctx, true)
}
