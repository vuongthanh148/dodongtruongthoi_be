package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

const (
	lookupCodeLen      = 6
	lookupCodeAlphabet = "ABCDEFGHJKMNPQRSTUVWXYZ23456789" // no 0, O, 1, I, L
	lookupMaxAttempts  = 5
	lookupLockDuration = 15 * time.Minute
	lookupSessionTTL   = 15 * time.Minute
	lookupTokenBytes   = 32
)

// newLookupCode returns a random per-order code the buyer keeps to verify the order.
func newLookupCode() (string, error) {
	code := make([]byte, lookupCodeLen)
	max := big.NewInt(int64(len(lookupCodeAlphabet)))
	for i := range code {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("generate lookup code: %w", err)
		}
		code[i] = lookupCodeAlphabet[n.Int64()]
	}
	return string(code), nil
}

// normalizeLookupPhone removes spaces so "0988 123 456" and "0988123456" are the same phone.
func normalizeLookupPhone(phone string) string {
	return strings.ReplaceAll(phone, " ", "")
}

// OrderLookupLockedError reports that a phone is locked after too many failed codes.
// Handlers map it to 429.
type OrderLookupLockedError struct {
	LockedUntil time.Time
}

func (e OrderLookupLockedError) Error() string { return "locked" }

// OrderLookupInvalidCodeError reports a wrong code. Handlers map it to 401.
type OrderLookupInvalidCodeError struct {
	RemainingAttempts int
}

func (e OrderLookupInvalidCodeError) Error() string { return "invalid code" }

// OrderLookupSession is a verified lookup session, bound to the order it was verified for.
type OrderLookupSession struct {
	Token     string
	ExpiresAt time.Time
	OrderID   string
}

type OrderLookupUsecase struct {
	orderRepo   domain.OrderRepository
	attemptRepo domain.LookupAttemptRepository
	sessionRepo domain.LookupSessionRepository
	now         func() time.Time
}

func NewOrderLookupUsecase(
	orderRepo domain.OrderRepository,
	attemptRepo domain.LookupAttemptRepository,
	sessionRepo domain.LookupSessionRepository,
) *OrderLookupUsecase {
	return &OrderLookupUsecase{
		orderRepo:   orderRepo,
		attemptRepo: attemptRepo,
		sessionRepo: sessionRepo,
		now:         time.Now,
	}
}

// VerifyOrderLookup checks a phone and lookup code. A locked phone is rejected
// before the code is evaluated. A match against the most recent order of the phone
// resets the failure count and issues a session bound to that order. Any miss counts
// as a failure, whether or not the phone has orders, so the response never reveals
// which phones have orders.
func (u *OrderLookupUsecase) VerifyOrderLookup(ctx context.Context, phone, code string) (OrderLookupSession, error) {
	phone = normalizeLookupPhone(phone)
	code = strings.ToUpper(strings.TrimSpace(code))
	now := u.now()

	attempt, _, err := u.attemptRepo.Get(ctx, phone)
	if err != nil {
		return OrderLookupSession{}, err
	}
	if attempt.LockedUntil != nil && attempt.LockedUntil.After(now) {
		return OrderLookupSession{}, OrderLookupLockedError{LockedUntil: *attempt.LockedUntil}
	}

	orders, err := u.orderRepo.List(ctx, &phone, nil, 0, 0)
	if err != nil {
		return OrderLookupSession{}, err
	}
	// List returns newest first, so the most recent matching order wins.
	orderID := ""
	if code != "" {
		for _, ord := range orders {
			if ord.LookupCode == code {
				orderID = ord.ID
				break
			}
		}
	}

	if orderID == "" {
		failed, err := u.attemptRepo.IncrementFailed(ctx, phone)
		if err != nil {
			return OrderLookupSession{}, err
		}
		if failed >= lookupMaxAttempts {
			until := now.Add(lookupLockDuration)
			if err := u.attemptRepo.Lock(ctx, phone, until); err != nil {
				return OrderLookupSession{}, err
			}
			return OrderLookupSession{}, OrderLookupLockedError{LockedUntil: until}
		}
		return OrderLookupSession{}, OrderLookupInvalidCodeError{RemainingAttempts: lookupMaxAttempts - failed}
	}

	if err := u.attemptRepo.Reset(ctx, phone); err != nil {
		return OrderLookupSession{}, err
	}
	token, err := newLookupToken()
	if err != nil {
		return OrderLookupSession{}, err
	}
	expiresAt := now.Add(lookupSessionTTL)
	if err := u.sessionRepo.Create(ctx, domain.LookupSession{
		TokenHash: hashLookupToken(token),
		Phone:     phone,
		OrderID:   orderID,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}); err != nil {
		return OrderLookupSession{}, err
	}
	return OrderLookupSession{Token: token, ExpiresAt: expiresAt, OrderID: orderID}, nil
}

// AuthorizeOrder reports whether token is a live session bound to orderID. A token
// verified for another order, an expired token, and an empty or unknown token all
// return false.
func (u *OrderLookupUsecase) AuthorizeOrder(ctx context.Context, token, orderID string) (bool, error) {
	if token == "" || orderID == "" {
		return false, nil
	}
	s, ok, err := u.sessionRepo.Get(ctx, hashLookupToken(token))
	if err != nil || !ok {
		return false, err
	}
	return strings.EqualFold(s.OrderID, orderID) && s.ExpiresAt.After(u.now()), nil
}

func newLookupToken() (string, error) {
	b := make([]byte, lookupTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate lookup token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func hashLookupToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
