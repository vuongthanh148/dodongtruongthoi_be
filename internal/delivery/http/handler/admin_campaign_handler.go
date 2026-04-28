package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListCampaigns(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListAllCampaigns(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Description   string `json:"description"`
		DiscountType  string `json:"discount_type"`
		DiscountValue int64  `json:"discount_value"`
		StartsAt      string `json:"starts_at"`
		EndsAt        string `json:"ends_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	startsAt, err := time.Parse(time.RFC3339, body.StartsAt)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid starts_at format")
		return
	}
	endsAt, err := time.Parse(time.RFC3339, body.EndsAt)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid ends_at format")
		return
	}

	campaign := domain.Campaign{
		ID:            body.ID,
		Name:          body.Name,
		Description:   ptrIfNotEmpty(body.Description),
		DiscountType:  body.DiscountType,
		DiscountValue: body.DiscountValue,
		StartsAt:      startsAt,
		EndsAt:        endsAt,
	}

	result, err := h.platform.CreateCampaign(r.Context(), campaign)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, result)
}

func (h *AdminHandler) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "campaign id is required")
		return
	}

	var body struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		DiscountType  string `json:"discount_type"`
		DiscountValue int64  `json:"discount_value"`
		StartsAt      string `json:"starts_at"`
		EndsAt        string `json:"ends_at"`
		IsActive      bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	startsAt, err := time.Parse(time.RFC3339, body.StartsAt)
	if err != nil && body.StartsAt != "" {
		response.Error(w, http.StatusBadRequest, "invalid starts_at format")
		return
	}
	endsAt, err := time.Parse(time.RFC3339, body.EndsAt)
	if err != nil && body.EndsAt != "" {
		response.Error(w, http.StatusBadRequest, "invalid ends_at format")
		return
	}

	updates := domain.Campaign{
		Name:          body.Name,
		Description:   ptrIfNotEmpty(body.Description),
		DiscountType:  body.DiscountType,
		DiscountValue: body.DiscountValue,
		StartsAt:      startsAt,
		EndsAt:        endsAt,
		IsActive:      body.IsActive,
	}

	result, err := h.platform.UpdateCampaign(r.Context(), id, updates)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "campaign id is required")
		return
	}

	err := h.platform.DeleteCampaign(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "campaign deleted"})
}

func (h *AdminHandler) SetCampaignProducts(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "campaign id is required")
		return
	}

	var body struct {
		ProductIDs []string `json:"product_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	err := h.platform.SetCampaignProducts(r.Context(), id, body.ProductIDs)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "campaign products updated"})
}
