package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	// Accept a short alias for pending status so CMS clients can use either
	// "pending" or full database value "pending_confirm".
	if status == "pending" {
		status = "pending_confirm"
	}
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	result, err := h.platform.ListOrders(r.Context(), statusPtr)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	row, ok, err := h.platform.GetOrder(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		response.Error(w, http.StatusNotFound, "order not found")
		return
	}
	response.Success(w, http.StatusOK, row)
}

func (h *AdminHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body struct {
		Status    string `json:"status"`
		AdminNote string `json:"adminNote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if body.Status == "" {
		response.Error(w, http.StatusBadRequest, "status is required")
		return
	}
	err := h.platform.UpdateOrderStatus(r.Context(), id, body.Status, body.AdminNote)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	// Return updated order
	order, _, err := h.platform.GetOrder(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, order)
}
