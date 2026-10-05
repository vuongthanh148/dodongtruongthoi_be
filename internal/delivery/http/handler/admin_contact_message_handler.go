package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListContactMessages(w http.ResponseWriter, r *http.Request) {
	handledParam := strings.TrimSpace(r.URL.Query().Get("handled"))

	var handledFilter *bool
	if handledParam != "" {
		v := strings.EqualFold(handledParam, "true")
		if !v && !strings.EqualFold(handledParam, "false") {
			response.Error(w, http.StatusBadRequest, "handled must be true or false")
			return
		}
		handledFilter = &v
	}

	result, err := h.platform.ListContactMessages(r.Context(), handledFilter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) SetContactMessageHandled(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var body struct {
		Handled *bool `json:"handled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Handled == nil {
		response.Error(w, http.StatusBadRequest, "body must be {\"handled\": true|false}")
		return
	}

	result, err := h.platform.SetContactMessageHandled(r.Context(), id, *body.Handled)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "contact message not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}
