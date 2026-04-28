package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListCategories(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Tone        string `json:"tone"`
		ImageURL    string `json:"image_url"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	category := domain.Category{
		ID:          body.ID,
		Name:        body.Name,
		Slug:        body.Slug,
		Description: ptrIfNotEmpty(body.Description),
		Tone:        body.Tone,
		ImageURL:    ptrIfNotEmpty(body.ImageURL),
		SortOrder:   body.SortOrder,
	}

	result, err := h.platform.CreateCategory(r.Context(), category)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, result)
}

func (h *AdminHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "category id is required")
		return
	}

	var body struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Tone        string `json:"tone"`
		ImageURL    string `json:"image_url"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	updates := domain.Category{
		Name:        body.Name,
		Slug:        body.Slug,
		Description: ptrIfNotEmpty(body.Description),
		Tone:        body.Tone,
		ImageURL:    ptrIfNotEmpty(body.ImageURL),
		SortOrder:   body.SortOrder,
		IsActive:    body.IsActive,
	}

	result, err := h.platform.UpdateCategory(r.Context(), id, updates)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "category id is required")
		return
	}

	err := h.platform.DeleteCategory(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "category deleted"})
}
