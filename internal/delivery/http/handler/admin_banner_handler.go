package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListBanners(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListBanners(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Subtitle  string `json:"subtitle"`
		ImageURL  string `json:"image_url"`
		LinkURL   string `json:"link_url"`
		SortOrder int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	banner := domain.Banner{
		ID:        body.ID,
		Title:     ptrIfNotEmpty(body.Title),
		Subtitle:  ptrIfNotEmpty(body.Subtitle),
		ImageURL:  ptrIfNotEmpty(body.ImageURL),
		LinkURL:   ptrIfNotEmpty(body.LinkURL),
		SortOrder: body.SortOrder,
	}

	result, err := h.platform.CreateBanner(r.Context(), banner)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, result)
}

func (h *AdminHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "banner id is required")
		return
	}

	var body struct {
		Title     string `json:"title"`
		Subtitle  string `json:"subtitle"`
		ImageURL  string `json:"image_url"`
		LinkURL   string `json:"link_url"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	updates := domain.Banner{
		Title:     ptrIfNotEmpty(body.Title),
		Subtitle:  ptrIfNotEmpty(body.Subtitle),
		ImageURL:  ptrIfNotEmpty(body.ImageURL),
		LinkURL:   ptrIfNotEmpty(body.LinkURL),
		SortOrder: body.SortOrder,
		IsActive:  body.IsActive,
	}

	result, err := h.platform.UpdateBanner(r.Context(), id, updates)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "banner id is required")
		return
	}

	err := h.platform.DeleteBanner(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "banner deleted"})
}
