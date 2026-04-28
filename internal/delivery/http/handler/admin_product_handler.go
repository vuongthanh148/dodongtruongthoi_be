package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func (h *AdminHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	result, err := h.platform.ListProducts(r.Context(), domain.ProductQuery{
		Category: r.URL.Query().Get("category"),
		Sort:     r.URL.Query().Get("sort"),
		Limit:    limit,
		Offset:   offset,
	}, true)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	product, ok, err := h.platform.GetProduct(r.Context(), id, true)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		response.Error(w, http.StatusNotFound, "product not found")
		return
	}
	response.Success(w, http.StatusOK, product)
}

func (h *AdminHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID             string            `json:"id"`
		Title          string            `json:"title"`
		Subtitle       string            `json:"subtitle"`
		CategoryID     string            `json:"category_id"`
		Badge          string            `json:"badge"`
		BasePrice      int64             `json:"base_price"`
		Description    string            `json:"description"`
		Meaning        string            `json:"meaning"`
		DefaultBG      string            `json:"default_bg"`
		DefaultFrame   string            `json:"default_frame"`
		BGTones        []string          `json:"bg_tones"`
		Frames         []string          `json:"frames"`
		ZodiacIDs      []string          `json:"zodiac_ids"`
		PurposePlace   []string          `json:"purpose_place"`
		PurposeUse     []string          `json:"purpose_use"`
		PurposeAvoid   []string          `json:"purpose_avoid"`
		Specs          map[string]string `json:"specs"`
		RequiresBGTone bool              `json:"requires_bg_tone"`
		RequiresFrame  bool              `json:"requires_frame"`
		RequiresSize   bool              `json:"requires_size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if body.ID == "" || body.Title == "" || body.CategoryID == "" {
		response.Error(w, http.StatusBadRequest, "id, title, category_id are required")
		return
	}

	product := domain.Product{
		ID:             body.ID,
		Title:          body.Title,
		Subtitle:       ptrIfNotEmpty(body.Subtitle),
		CategoryID:     body.CategoryID,
		Badge:          ptrIfNotEmpty(body.Badge),
		BasePrice:      body.BasePrice,
		Description:    ptrIfNotEmpty(body.Description),
		Meaning:        ptrIfNotEmpty(body.Meaning),
		DefaultBG:      body.DefaultBG,
		DefaultFrame:   body.DefaultFrame,
		BGTones:        body.BGTones,
		Frames:         body.Frames,
		ZodiacIDs:      body.ZodiacIDs,
		PurposePlace:   body.PurposePlace,
		PurposeUse:     body.PurposeUse,
		PurposeAvoid:   body.PurposeAvoid,
		Specs:          body.Specs,
		RequiresBGTone: body.RequiresBGTone,
		RequiresFrame:  body.RequiresFrame,
		RequiresSize:   body.RequiresSize,
	}

	result, err := h.platform.CreateProduct(r.Context(), product)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, result)
}

func (h *AdminHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	var body struct {
		Title          string            `json:"title"`
		Subtitle       string            `json:"subtitle"`
		CategoryID     string            `json:"category_id"`
		Badge          string            `json:"badge"`
		BasePrice      int64             `json:"base_price"`
		Description    string            `json:"description"`
		Meaning        string            `json:"meaning"`
		DefaultBG      string            `json:"default_bg"`
		DefaultFrame   string            `json:"default_frame"`
		BGTones        []string          `json:"bg_tones"`
		Frames         []string          `json:"frames"`
		ZodiacIDs      []string          `json:"zodiac_ids"`
		PurposePlace   []string          `json:"purpose_place"`
		PurposeUse     []string          `json:"purpose_use"`
		PurposeAvoid   []string          `json:"purpose_avoid"`
		Specs          map[string]string `json:"specs"`
		RequiresBGTone bool              `json:"requires_bg_tone"`
		RequiresFrame  bool              `json:"requires_frame"`
		RequiresSize   bool              `json:"requires_size"`
		IsActive       bool              `json:"is_active"`
		SortOrder      int               `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	updates := domain.Product{
		Title:          body.Title,
		Subtitle:       ptrIfNotEmpty(body.Subtitle),
		CategoryID:     body.CategoryID,
		Badge:          ptrIfNotEmpty(body.Badge),
		BasePrice:      body.BasePrice,
		Description:    ptrIfNotEmpty(body.Description),
		Meaning:        ptrIfNotEmpty(body.Meaning),
		DefaultBG:      body.DefaultBG,
		DefaultFrame:   body.DefaultFrame,
		BGTones:        body.BGTones,
		Frames:         body.Frames,
		ZodiacIDs:      body.ZodiacIDs,
		PurposePlace:   body.PurposePlace,
		PurposeUse:     body.PurposeUse,
		PurposeAvoid:   body.PurposeAvoid,
		Specs:          body.Specs,
		RequiresBGTone: body.RequiresBGTone,
		RequiresFrame:  body.RequiresFrame,
		RequiresSize:   body.RequiresSize,
		IsActive:       body.IsActive,
		SortOrder:      body.SortOrder,
	}

	result, err := h.platform.UpdateProduct(r.Context(), id, updates)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	err := h.platform.DeleteProduct(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "product deleted"})
}

func (h *AdminHandler) UploadProductImages(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	// Parse multipart form data (32MB max)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "file parameter required")
		return
	}
	defer file.Close()

	// Get optional parameters
	bgTone := r.FormValue("bgTone")
	frame := r.FormValue("frame")

	// Upload to Cloudinary and save to database
	img, err := h.platform.UploadProductImage(r.Context(), productID, file, header.Filename,
		ptrIfNotEmpty(bgTone), ptrIfNotEmpty(frame))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "upload failed: "+err.Error())
		return
	}

	response.Success(w, http.StatusCreated, img)
}

func (h *AdminHandler) DeleteProductImage(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	imageID := chi.URLParam(r, "imgId")
	if productID == "" || imageID == "" {
		response.Error(w, http.StatusBadRequest, "product id and image id are required")
		return
	}

	err := h.platform.DeleteProductImage(r.Context(), productID, imageID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "image deleted"})
}

func (h *AdminHandler) SetProductSizes(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	var body struct {
		Sizes []domain.ProductSize `json:"sizes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.platform.SetProductSizes(r.Context(), id, body.Sizes)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}
