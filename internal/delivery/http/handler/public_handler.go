package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/usecase"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

type PublicHandler struct {
	platform *usecase.PlatformUsecase
}

func NewPublicHandler(platform *usecase.PlatformUsecase) *PublicHandler {
	return &PublicHandler{platform: platform}
}

func (h *PublicHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListCategories(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *PublicHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	row, ok, err := h.platform.GetCategory(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		response.Error(w, http.StatusNotFound, "category not found")
		return
	}
	response.Success(w, http.StatusOK, row)
}

func (h *PublicHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	rows, err := h.platform.ListProducts(r.Context(), domain.ProductQuery{
		Category: r.URL.Query().Get("category"),
		Sort:     r.URL.Query().Get("sort"),
		Limit:    limit,
		Offset:   offset,
	}, false)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, rows)
}

func (h *PublicHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	row, ok, err := h.platform.GetProduct(r.Context(), id, false)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		response.Error(w, http.StatusNotFound, "product not found")
		return
	}
	response.Success(w, http.StatusOK, row)
}

func (h *PublicHandler) ListProductReviews(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	result, err := h.platform.ListReviews(r.Context(), id, false)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *PublicHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		response.Error(w, http.StatusBadRequest, "product id is required")
		return
	}

	var body struct {
		ReviewerName string `json:"reviewer_name"`
		Rating       int    `json:"rating"`
		Body         string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	if body.ReviewerName == "" {
		response.Error(w, http.StatusBadRequest, "reviewer_name is required")
		return
	}
	if body.Rating < 1 || body.Rating > 5 {
		response.Error(w, http.StatusBadRequest, "rating must be between 1 and 5")
		return
	}

	result, err := h.platform.CreateReview(r.Context(), productID, body.ReviewerName, body.Rating, ptrIfNotEmpty(body.Body))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, result)
}

func (h *PublicHandler) ListBanners(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListBanners(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *PublicHandler) ListCampaigns(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListAllCampaigns(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *PublicHandler) ListCustomerPhotos(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListCustomerPhotos(r.Context(), false)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *PublicHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.ListContacts(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *PublicHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	result, err := h.platform.GetPublicSettings(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *PublicHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Phone        string `json:"phone"`
		CustomerName string `json:"customerName"`
		Address      string `json:"address"`
		Note         string `json:"note"`
		Items        []struct {
			ProductID       string `json:"productId"`
			SizeCode        string `json:"sizeCode"`
			SizeLabel       string `json:"sizeLabel"`
			BGTone          string `json:"bgTone"`
			BGToneLabel     string `json:"bgToneLabel"`
			Frame           string `json:"frame"`
			FrameLabel      string `json:"frameLabel"`
			Quantity        int    `json:"quantity"`
			UnitPrice       int64  `json:"unitPrice"`
			VariantImageURL string `json:"variantImageUrl"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	items := make([]usecase.CreateOrderItem, 0, len(body.Items))
	for _, it := range body.Items {
		items = append(items, usecase.CreateOrderItem{
			ProductID:       it.ProductID,
			SizeCode:        ptrIfNotEmpty(it.SizeCode),
			SizeLabel:       ptrIfNotEmpty(it.SizeLabel),
			BGTone:          ptrIfNotEmpty(it.BGTone),
			BGToneLabel:     ptrIfNotEmpty(it.BGToneLabel),
			Frame:           ptrIfNotEmpty(it.Frame),
			FrameLabel:      ptrIfNotEmpty(it.FrameLabel),
			Quantity:        it.Quantity,
			UnitPrice:       it.UnitPrice,
			VariantImageURL: ptrIfNotEmpty(it.VariantImageURL),
		})
	}

	order, err := h.platform.CreateOrder(r.Context(), usecase.CreateOrderRequest{
		Phone:        body.Phone,
		CustomerName: ptrIfNotEmpty(body.CustomerName),
		Address:      ptrIfNotEmpty(body.Address),
		Note:         ptrIfNotEmpty(body.Note),
		Items:        items,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, order)
}

func (h *PublicHandler) ListOrdersByPhone(w http.ResponseWriter, r *http.Request) {
	phone := strings.TrimSpace(r.URL.Query().Get("phone"))
	if phone == "" {
		response.Error(w, http.StatusBadRequest, "phone is required")
		return
	}
	result, err := h.platform.ListOrdersByPhone(r.Context(), phone)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *PublicHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
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

func (h *PublicHandler) GetWishlistByPhone(w http.ResponseWriter, r *http.Request) {
	phone := strings.TrimSpace(r.URL.Query().Get("phone"))
	if phone == "" {
		response.Error(w, http.StatusBadRequest, "phone is required")
		return
	}
	result, err := h.platform.GetWishlist(r.Context(), phone)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *PublicHandler) SyncWishlist(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Phone      string   `json:"phone"`
		ProductIDs []string `json:"productIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if strings.TrimSpace(body.Phone) == "" {
		response.Error(w, http.StatusBadRequest, "phone is required")
		return
	}
	result, err := h.platform.SyncWishlist(r.Context(), body.Phone, body.ProductIDs)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, result)
}

func (h *PublicHandler) DeleteWishlistItem(w http.ResponseWriter, r *http.Request) {
	phone := chi.URLParam(r, "phone")
	productID := chi.URLParam(r, "productId")
	if phone == "" || productID == "" {
		response.Error(w, http.StatusBadRequest, "phone and productId are required")
		return
	}
	err := h.platform.DeleteWishlistItem(r.Context(), phone, productID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, map[string]string{"status": "deleted"})
}
