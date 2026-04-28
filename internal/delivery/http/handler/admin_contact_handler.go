package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.platform.ListContacts(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, contacts)
}

func (h *AdminHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var req domain.ContactLink
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	contact, err := h.platform.CreateContact(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, contact)
}

func (h *AdminHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "contactID")

	var req domain.ContactLink
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	contact, err := h.platform.UpdateContact(r.Context(), id, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, contact)
}

func (h *AdminHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "contactID")

	err := h.platform.DeleteContact(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusNoContent, nil)
}
