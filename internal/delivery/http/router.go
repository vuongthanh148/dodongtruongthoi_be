package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/delivery/http/handler"
	authmiddleware "github.com/vuongthanh148/dodongtruongthoi_be/internal/delivery/http/middleware"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/usecase"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewRouter(
	healthHandler *handler.HealthHandler,
	publicHandler *handler.PublicHandler,
	adminHandler *handler.AdminHandler,
	platformUsecase *usecase.PlatformUsecase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(corsMiddleware)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", healthHandler.GetHealth)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler.GetHealth)

		r.Get("/categories", publicHandler.ListCategories)
		r.Get("/categories/{id}", publicHandler.GetCategory)

		r.Get("/products", publicHandler.ListProducts)
		r.Get("/products/{id}", publicHandler.GetProduct)
		r.Get("/products/{id}/reviews", publicHandler.ListProductReviews)
		r.Post("/products/{id}/reviews", publicHandler.SubmitReview)

		r.Get("/banners", publicHandler.ListBanners)
		r.Get("/campaigns", publicHandler.ListCampaigns)
		r.Get("/contacts", publicHandler.ListContacts)
		r.Get("/settings", publicHandler.GetSettings)

		r.Post("/orders", publicHandler.CreateOrder)
		r.Get("/orders", publicHandler.ListOrdersByPhone)
		r.Get("/orders/{id}", publicHandler.GetOrder)

		r.Get("/wishlists", publicHandler.GetWishlistByPhone)
		r.Post("/wishlists", publicHandler.SyncWishlist)
		r.Delete("/wishlists/{phone}/{productId}", publicHandler.DeleteWishlistItem)

		r.Route("/admin", func(r chi.Router) {
			r.Post("/login", adminHandler.Login)
			r.Post("/refresh", adminHandler.Refresh)

			r.Group(func(r chi.Router) {
				r.Use(authmiddleware.RequireAdminAuth(platformUsecase))

				r.Get("/products", adminHandler.ListProducts)
				r.Get("/products/{id}", adminHandler.GetProduct)
				r.Post("/products", adminHandler.CreateProduct)
				r.Put("/products/{id}", adminHandler.UpdateProduct)
				r.Delete("/products/{id}", adminHandler.DeleteProduct)
				r.Post("/products/{id}/images", adminHandler.UploadProductImages)
				r.Delete("/products/{id}/images/{imgId}", adminHandler.DeleteProductImage)
				r.Put("/products/{id}/sizes", adminHandler.SetProductSizes)

				r.Get("/categories", adminHandler.ListCategories)
				r.Post("/categories", adminHandler.CreateCategory)
				r.Put("/categories/{id}", adminHandler.UpdateCategory)
				r.Delete("/categories/{id}", adminHandler.DeleteCategory)

				r.Get("/campaigns", adminHandler.ListCampaigns)
				r.Post("/campaigns", adminHandler.CreateCampaign)
				r.Put("/campaigns/{id}", adminHandler.UpdateCampaign)
				r.Delete("/campaigns/{id}", adminHandler.DeleteCampaign)
				r.Put("/campaigns/{id}/products", adminHandler.SetCampaignProducts)

				r.Get("/banners", adminHandler.ListBanners)
				r.Post("/banners", adminHandler.CreateBanner)
				r.Put("/banners/{id}", adminHandler.UpdateBanner)
				r.Delete("/banners/{id}", adminHandler.DeleteBanner)

				r.Get("/contacts", adminHandler.ListContacts)
				r.Post("/contacts", adminHandler.CreateContact)
				r.Put("/contacts/{id}", adminHandler.UpdateContact)
				r.Delete("/contacts/{id}", adminHandler.DeleteContact)

				r.Get("/orders", adminHandler.ListOrders)
				r.Get("/orders/{id}", adminHandler.GetOrder)
				r.Put("/orders/{id}/status", adminHandler.UpdateOrderStatus)

				r.Get("/reviews", adminHandler.ListReviews)
				r.Put("/reviews/{id}", adminHandler.UpdateReview)
				r.Delete("/reviews/{id}", adminHandler.DeleteReview)

				r.Get("/settings", adminHandler.GetSettings)
				r.Put("/settings", adminHandler.UpdateSettings)
			})
		})
	})

	return r
}
