package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/usecase"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

type AdminHandler struct {
	platform *usecase.PlatformUsecase
}

func NewAdminHandler(platform *usecase.PlatformUsecase) *AdminHandler {
	return &AdminHandler{platform: platform}
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.platform.Login(r.Context(), body.Username, body.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if body.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	result, err := h.platform.Refresh(r.Context(), body.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	result, err := h.platform.ListProducts(r.Context(), usecase.ProductQuery{
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

// Helper function to convert empty string to nil pointer
func ptrIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
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

func (h *AdminHandler) ListCampaigns(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListAllCampaigns(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Description   string `json:"description"`
		DiscountType  string `json:"discount_type"`
		DiscountValue int64  `json:"discount_value"`
		StartsAt      string `json:"starts_at"`
		EndsAt        string `json:"ends_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	startsAt, err := time.Parse(time.RFC3339, body.StartsAt)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid starts_at format")
		return
	}
	endsAt, err := time.Parse(time.RFC3339, body.EndsAt)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid ends_at format")
		return
	}

	campaign := domain.Campaign{
		ID:            body.ID,
		Name:          body.Name,
		Description:   ptrIfNotEmpty(body.Description),
		DiscountType:  body.DiscountType,
		DiscountValue: body.DiscountValue,
		StartsAt:      startsAt,
		EndsAt:        endsAt,
	}

	result, err := h.platform.CreateCampaign(r.Context(), campaign)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, result)
}

func (h *AdminHandler) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "campaign id is required")
		return
	}

	var body struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		DiscountType  string `json:"discount_type"`
		DiscountValue int64  `json:"discount_value"`
		StartsAt      string `json:"starts_at"`
		EndsAt        string `json:"ends_at"`
		IsActive      bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	startsAt, err := time.Parse(time.RFC3339, body.StartsAt)
	if err != nil && body.StartsAt != "" {
		response.Error(w, http.StatusBadRequest, "invalid starts_at format")
		return
	}
	endsAt, err := time.Parse(time.RFC3339, body.EndsAt)
	if err != nil && body.EndsAt != "" {
		response.Error(w, http.StatusBadRequest, "invalid ends_at format")
		return
	}

	updates := domain.Campaign{
		Name:          body.Name,
		Description:   ptrIfNotEmpty(body.Description),
		DiscountType:  body.DiscountType,
		DiscountValue: body.DiscountValue,
		StartsAt:      startsAt,
		EndsAt:        endsAt,
		IsActive:      body.IsActive,
	}

	result, err := h.platform.UpdateCampaign(r.Context(), id, updates)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "campaign id is required")
		return
	}

	err := h.platform.DeleteCampaign(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "campaign deleted"})
}

func (h *AdminHandler) SetCampaignProducts(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "campaign id is required")
		return
	}

	var body struct {
		ProductIDs []string `json:"product_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	err := h.platform.SetCampaignProducts(r.Context(), id, body.ProductIDs)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "campaign products updated"})
}

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

func (h *AdminHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListContacts(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID        string `json:"id"`
		Platform  string `json:"platform"`
		Label     string `json:"label"`
		URL       string `json:"url"`
		SortOrder int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	contact := domain.ContactLink{
		ID:        body.ID,
		Platform:  body.Platform,
		Label:     body.Label,
		URL:       body.URL,
		SortOrder: body.SortOrder,
	}

	result, err := h.platform.CreateContact(r.Context(), contact)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, result)
}

func (h *AdminHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "contact id is required")
		return
	}

	var body struct {
		Platform  string `json:"platform"`
		Label     string `json:"label"`
		URL       string `json:"url"`
		SortOrder int    `json:"sort_order"`
		IsActive  bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	updates := domain.ContactLink{
		Platform:  body.Platform,
		Label:     body.Label,
		URL:       body.URL,
		SortOrder: body.SortOrder,
		IsActive:  body.IsActive,
	}

	result, err := h.platform.UpdateContact(r.Context(), id, updates)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "contact id is required")
		return
	}

	err := h.platform.DeleteContact(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"message": "contact deleted"})
}

func (h *AdminHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
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

func (h *AdminHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("productId")
	approved := r.URL.Query().Get("approved")

	if productID != "" {
		result, err := h.platform.ListReviews(r.Context(), productID, true)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.Success(w, http.StatusOK, result)
		return
	}

	includeUnapproved := approved == "" || approved == "false"
	result, err := h.platform.ListAllReviews(r.Context(), includeUnapproved)
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

func (h *AdminHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.GetAdminSettings(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *AdminHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.platform.UpdateSettings(r.Context(), body)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}
