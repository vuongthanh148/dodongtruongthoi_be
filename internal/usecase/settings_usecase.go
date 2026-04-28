package usecase

import (
	"context"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type SettingsUsecase struct {
	settingsRepo domain.SiteSettingsRepository
}

func NewSettingsUsecase(settingsRepo domain.SiteSettingsRepository) *SettingsUsecase {
	return &SettingsUsecase{
		settingsRepo: settingsRepo,
	}
}

func (u *SettingsUsecase) GetPublicSettings(ctx context.Context) (map[string]string, error) {
	return u.settingsRepo.GetAll(ctx)
}

func (u *SettingsUsecase) GetAdminSettings(ctx context.Context) (map[string]string, error) {
	return u.settingsRepo.GetAll(ctx)
}

func (u *SettingsUsecase) UpdateSettings(ctx context.Context, in map[string]string) (map[string]string, error) {
	if err := u.settingsRepo.SetBulk(ctx, in); err != nil {
		return nil, err
	}
	return u.settingsRepo.GetAll(ctx)
}
