package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	entityType := strings.TrimSpace(r.URL.Query().Get("entity_type"))
	var entityTypePtr *string
	if entityType != "" {
		entityTypePtr = &entityType
	}

	limitParam := strings.TrimSpace(r.URL.Query().Get("limit"))
	limit := 50
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	offsetParam := strings.TrimSpace(r.URL.Query().Get("offset"))
	offset := 0
	if offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil && o >= 0 {
			offset = o
		}
	}

	result, err := h.platform.ListAuditLogs(r.Context(), limit, offset, entityTypePtr)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}
