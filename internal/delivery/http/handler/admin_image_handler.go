package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	images, err := h.platform.ListImages(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, images)
}

func (h *AdminHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
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

	name := r.FormValue("name")
	if name == "" {
		name = header.Filename
	}

	img, err := h.platform.UploadImage(r.Context(), file, header.Filename, name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "upload failed: "+err.Error())
		return
	}
	response.Success(w, http.StatusCreated, img)
}

func (h *AdminHandler) RenameImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	img, err := h.platform.RenameImage(r.Context(), id, body.Name)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, img)
}

func (h *AdminHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}
	if err := h.platform.DeleteImage(r.Context(), id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "image deleted"})
}
