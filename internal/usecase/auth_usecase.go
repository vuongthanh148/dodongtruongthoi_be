package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type AuthUsecase struct {
	adminUserRepo domain.AdminUserRepository
	jwtSecret     string
}

func NewAuthUsecase(adminUserRepo domain.AdminUserRepository, jwtSecret string) *AuthUsecase {
	if jwtSecret == "" {
		jwtSecret = "dev-secret"
	}
	return &AuthUsecase{
		adminUserRepo: adminUserRepo,
		jwtSecret:     jwtSecret,
	}
}

func (u *AuthUsecase) Login(ctx context.Context, username, password string) (LoginResult, error) {
	user, ok, err := u.adminUserRepo.GetByUsername(ctx, username)
	if err != nil {
		return LoginResult{}, err
	}
	if !ok || !user.IsActive {
		return LoginResult{}, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return LoginResult{}, errors.New("invalid credentials")
	}

	// Update last login
	_ = u.adminUserRepo.UpdateLastLogin(ctx, user.ID)

	token, err := u.signAccessToken(username)
	if err != nil {
		return LoginResult{}, err
	}
	refreshToken, err := u.signRefreshToken(username)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: token, RefreshToken: refreshToken}, nil
}

func (u *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (LoginResult, error) {
	username, err := u.verifyTokenByType(ctx, refreshToken, "refresh")
	if err != nil {
		return LoginResult{}, err
	}

	token, err := u.signAccessToken(username)
	if err != nil {
		return LoginResult{}, err
	}
	rotatedRefreshToken, err := u.signRefreshToken(username)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: token, RefreshToken: rotatedRefreshToken}, nil
}

func (u *AuthUsecase) VerifyToken(ctx context.Context, token string) error {
	_, err := u.verifyTokenByType(ctx, token, "access")
	return err
}

func (u *AuthUsecase) verifyTokenByType(ctx context.Context, token, expectedType string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", errors.New("invalid token")
	}

	payloadRaw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", errors.New("invalid token")
	}
	sigRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.New("invalid token")
	}

	mac := hmac.New(sha256.New, []byte(u.jwtSecret))
	_, _ = mac.Write(payloadRaw)
	expected := mac.Sum(nil)
	if !hmac.Equal(sigRaw, expected) {
		return "", errors.New("invalid token")
	}

	var payload struct {
		Username string `json:"u"`
		Type     string `json:"t"`
		Exp      int64  `json:"exp"`
	}
	if err := json.Unmarshal(payloadRaw, &payload); err != nil {
		return "", errors.New("invalid token")
	}
	if payload.Type != expectedType {
		return "", errors.New("invalid token type")
	}
	if payload.Exp < time.Now().Unix() {
		return "", errors.New("token expired")
	}

	// Verify username exists
	_, ok, err := u.adminUserRepo.GetByUsername(ctx, payload.Username)
	if err != nil {
		return "", errors.New("invalid token user")
	}
	if !ok {
		return "", errors.New("invalid token user")
	}
	return payload.Username, nil
}

func (u *AuthUsecase) signAccessToken(username string) (string, error) {
	return u.signToken(username, "access", 24*time.Hour)
}

func (u *AuthUsecase) signRefreshToken(username string) (string, error) {
	return u.signToken(username, "refresh", 30*24*time.Hour)
}

func (u *AuthUsecase) signToken(username, tokenType string, ttl time.Duration) (string, error) {
	payload := struct {
		Username string `json:"u"`
		Type     string `json:"t"`
		Exp      int64  `json:"exp"`
	}{
		Username: username,
		Type:     tokenType,
		Exp:      time.Now().Add(ttl).Unix(),
	}
	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, []byte(u.jwtSecret))
	_, _ = mac.Write(payloadRaw)
	sig := mac.Sum(nil)

	return fmt.Sprintf("%s.%s",
		base64.RawURLEncoding.EncodeToString(payloadRaw),
		base64.RawURLEncoding.EncodeToString(sig),
	), nil
}
