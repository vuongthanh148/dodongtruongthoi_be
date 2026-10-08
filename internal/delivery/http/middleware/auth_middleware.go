package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/usecase"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/response"
)

type contextKey string

const AdminUsernameKey contextKey = "admin_username"

func RequireAdminAuth(platform *usecase.PlatformUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if authorization == "" || !strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
				response.Error(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			token := strings.TrimSpace(authorization[len("Bearer "):])
			username, err := platform.VerifyTokenAndGetUsername(r.Context(), token)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), AdminUsernameKey, username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAdminUsername(r *http.Request) string {
	username, ok := r.Context().Value(AdminUsernameKey).(string)
	if !ok {
		return "admin"
	}
	return username
}
