package middleware

import (
	"net/http"
	"strings"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/usecase"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

func RequireAdminAuth(platform *usecase.PlatformUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" || !strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
				response.Error(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			token := strings.TrimSpace(authorization[len("Bearer "):])
			if err := platform.VerifyToken(r.Context(), token); err != nil {
				response.Error(w, http.StatusUnauthorized, err.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
