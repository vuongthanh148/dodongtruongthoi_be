package handler

import (
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/usecase"
)

type AdminHandler struct {
	platform *usecase.PlatformUsecase
}

func NewAdminHandler(platform *usecase.PlatformUsecase) *AdminHandler {
	return &AdminHandler{platform: platform}
}

// Helper function to convert empty string to nil pointer
func ptrIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
