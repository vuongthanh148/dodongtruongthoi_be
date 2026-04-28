package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListCustomerPhotos(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListCustomerPhotos(r.Context(), true)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) UploadCustomerPhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "file parameter required")
		return
	}
	defer file.Close()

	caption := ptrIfNotEmpty(r.FormValue("caption"))
	isActive := true
	if raw := r.FormValue("isActive"); raw != "" {
		isActive = raw == "true" || raw == "1"
	}
	sortOrder := 0
	if raw := r.FormValue("sortOrder"); raw != "" {
		sortOrder, _ = strconv.Atoi(raw)
	}

	photo, err := h.platform.CreateCustomerPhoto(r.Context(), file, header.Filename, caption, isActive, sortOrder)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "upload failed: "+err.Error())
		return
	}

	response.Success(w, http.StatusCreated, photo)
}

func (h *AdminHandler) UpdateCustomerPhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "customer photo id is required")
		return
	}

	var body struct {
		Caption   string `json:"caption"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	updated, err := h.platform.UpdateCustomerPhoto(r.Context(), id, ptrIfNotEmpty(body.Caption), body.IsActive, body.SortOrder)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusOK, updated)
}

func (h *AdminHandler) DeleteCustomerPhoto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "customer photo id is required")
		return
	}

	err := h.platform.DeleteCustomerPhoto(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "customer photo deleted"})
}
