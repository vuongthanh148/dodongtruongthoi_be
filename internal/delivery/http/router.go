package http

import (
"net/http"
"time"

"github.com/go-chi/chi/v5"
"github.com/go-chi/chi/v5/middleware"
"github.com/vuongthanh148/dodongtruongthoi_be/internal/delivery/http/handler"
)

func NewRouter(healthHandler *handler.HealthHandler) http.Handler {
r := chi.NewRouter()

r.Use(middleware.RequestID)
r.Use(middleware.Recoverer)
r.Use(middleware.Logger)
r.Use(middleware.Timeout(30 * time.Second))

r.Get("/health", healthHandler.GetHealth)
r.Route("/api/v1", func(r chi.Router) {
r.Get("/health", healthHandler.GetHealth)
})

return r
}
