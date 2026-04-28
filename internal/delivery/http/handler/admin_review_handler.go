package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("productId")
	approved := strings.TrimSpace(r.URL.Query().Get("approved"))

	var approvedFilter *bool
	if approved != "" {
		v := strings.EqualFold(approved, "true")
		if !v && !strings.EqualFold(approved, "false") {
			response.Error(w, http.StatusBadRequest, "approved must be true or false")
			return
		}
		approvedFilter = &v
	}

	if productID != "" {
		includeUnapproved := approvedFilter == nil || !*approvedFilter
		result, err := h.platform.ListReviews(r.Context(), productID, includeUnapproved)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.Success(w, http.StatusOK, result)
		return
	}

	result, err := h.platform.ListAllReviews(r.Context(), approvedFilter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "review id is required")
		return
	}

	var body struct {
		IsApproved bool `json:"is_approved"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	err := h.platform.UpdateReviewApproval(r.Context(), id, body.IsApproved)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	// Return updated review if needed
	response.Success(w, http.StatusOK, map[string]string{"message": "review updated"})
}

func (h *AdminHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "review id is required")
		return
	}

	err := h.platform.DeleteReview(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "review deleted"})
}
