package usecase

import (
	"context"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

// validThemeIDs must match the theme ids the storefront knows (frontend src/lib/themes.ts).
var validThemeIDs = map[string]bool{
	"default":      true,
	"tet":          true,
	"independence": true,
	"labor-day":    true,
}

type unknownThemeError struct {
	theme string
}

func (e unknownThemeError) Error() string {
	return "unknown theme: " + e.theme
}

func (e unknownThemeError) Unwrap() error {
	return domain.ErrInvalidInput
}

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
	if theme, ok := in["active_theme"]; ok && !validThemeIDs[theme] {
		return nil, unknownThemeError{theme: theme}
	}
	if err := u.settingsRepo.SetBulk(ctx, in); err != nil {
		return nil, err
	}
	return u.settingsRepo.GetAll(ctx)
}
