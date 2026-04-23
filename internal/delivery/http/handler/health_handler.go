package handler

import (
"net/http"

"github.com/vuongthanh148/dodongtruongthoi_be/internal/usecase"
"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

type HealthHandler struct {
healthUsecase *usecase.HealthUsecase
}

func NewHealthHandler(healthUsecase *usecase.HealthUsecase) *HealthHandler {
return &HealthHandler{healthUsecase: healthUsecase}
}

func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
status := h.healthUsecase.GetStatus(r.Context())
response.Success(w, http.StatusOK, status)
}
