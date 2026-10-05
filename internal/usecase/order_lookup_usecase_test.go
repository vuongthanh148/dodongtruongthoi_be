package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

// lookupOrderRepo serves orders by phone for lookup tests.
type lookupOrderRepo struct {
	domain.OrderRepository
	orders    []domain.Order // newest first, as the repo returns them
	listCalls int
}

func (f *lookupOrderRepo) List(ctx context.Context, phone *string, status *string, limit int, offset int) ([]domain.Order, error) {
	f.listCalls++
	var out []domain.Order
	for _, o := range f.orders {
		if phone == nil || o.Phone == *phone {
			out = append(out, o)
		}
	}
	return out, nil
}

type lookupAttemptRepo struct {
	domain.LookupAttemptRepository
	rows map[string]domain.LookupAttempt
}

func (f *lookupAttemptRepo) Get(ctx context.Context, phone string) (domain.LookupAttempt, bool, error) {
	a, ok := f.rows[phone]
	return a, ok, nil
}

func (f *lookupAttemptRepo) IncrementFailed(ctx context.Context, phone string) (int, error) {
	a := f.rows[phone]
	a.Phone = phone
	a.FailedCount++
	f.rows[phone] = a
	return a.FailedCount, nil
}

func (f *lookupAttemptRepo) Lock(ctx context.Context, phone string, until time.Time) error {
	a := f.rows[phone]
	a.Phone = phone
	a.FailedCount = 0
	a.LockedUntil = &until
	f.rows[phone] = a
	return nil
}

func (f *lookupAttemptRepo) Reset(ctx context.Context, phone string) error {
	a := f.rows[phone]
	a.Phone = phone
	a.FailedCount = 0
	a.LockedUntil = nil
	f.rows[phone] = a
	return nil
}

type lookupSessionRepo struct {
	domain.LookupSessionRepository
	rows map[string]domain.LookupSession
}

func (f *lookupSessionRepo) Create(ctx context.Context, s domain.LookupSession) error {
	f.rows[s.TokenHash] = s
	return nil
}

func (f *lookupSessionRepo) Get(ctx context.Context, tokenHash string) (domain.LookupSession, bool, error) {
	s, ok := f.rows[tokenHash]
	return s, ok, nil
}

const (
	lookupOrderID      = "d2715867-8857-4d2a-95b1-936a6414f040"
	lookupOlderOrderID = "e3826978-9968-4e3b-96c2-047b7525a151"
	lookupPhone        = "0988123456"
	otherLookupPhone   = "0977654321"
	lookupCode         = "K7QF3M"
	olderLookupCode    = "2XRN9H"
)

type lookupFixture struct {
	u        *OrderLookupUsecase
	orders   *lookupOrderRepo
	attempts *lookupAttemptRepo
	sessions *lookupSessionRepo
	now      time.Time
}

func newLookupFixture(orders ...domain.Order) *lookupFixture {
	f := &lookupFixture{
		orders:   &lookupOrderRepo{orders: orders},
		attempts: &lookupAttemptRepo{rows: map[string]domain.LookupAttempt{}},
		sessions: &lookupSessionRepo{rows: map[string]domain.LookupSession{}},
		now:      time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC),
	}
	f.u = NewOrderLookupUsecase(f.orders, f.attempts, f.sessions)
	f.u.now = func() time.Time { return f.now }
	return f
}

func lookupTestOrder() domain.Order {
	return domain.Order{ID: lookupOrderID, Phone: lookupPhone, Status: "pending_confirm", LookupCode: lookupCode}
}

func lookupOlderTestOrder() domain.Order {
	return domain.Order{ID: lookupOlderOrderID, Phone: lookupPhone, Status: "pending_confirm", LookupCode: olderLookupCode}
}

func TestNewLookupCode_LengthAndAlphabet(t *testing.T) {
	allowed := regexp.MustCompile(`^[ABCDEFGHJKMNPQRSTUVWXYZ23456789]{6}$`)
	for i := 0; i < 2000; i++ {
		code, err := newLookupCode()
		if err != nil {
			t.Fatalf("newLookupCode: %v", err)
		}
		if !allowed.MatchString(code) {
			t.Fatalf("newLookupCode() = %q, want 6 chars from the allowed alphabet", code)
		}
		if strings.ContainsAny(code, "01IOL") {
			t.Fatalf("newLookupCode() = %q contains an ambiguous character", code)
		}
	}
}

func TestVerifyOrderLookup_CorrectCodeIssuesHashedTokenBoundToOrder(t *testing.T) {
	f := newLookupFixture(lookupTestOrder())

	got, err := f.u.VerifyOrderLookup(context.Background(), "0988 123 456", " k7qf3m ")
	if err != nil {
		t.Fatalf("VerifyOrderLookup: unexpected error: %v", err)
	}
	if !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(got.Token) {
		t.Errorf("token = %q, want 32 random bytes hex", got.Token)
	}
	if got.OrderID != lookupOrderID {
		t.Errorf("OrderID = %q, want %q", got.OrderID, lookupOrderID)
	}
	if want := f.now.Add(15 * time.Minute); !got.ExpiresAt.Equal(want) {
		t.Errorf("ExpiresAt = %v, want %v", got.ExpiresAt, want)
	}
	if _, raw := f.sessions.rows[got.Token]; raw {
		t.Error("raw token stored as session key, want only its hash")
	}
	sess, ok := f.sessions.rows[hashLookupToken(got.Token)]
	if !ok {
		t.Fatal("no session stored under the SHA-256 of the token")
	}
	if sess.OrderID != lookupOrderID {
		t.Errorf("stored session OrderID = %q, want %q", sess.OrderID, lookupOrderID)
	}
}

func TestVerifyOrderLookup_CodeResolvesToItsOwnOrder(t *testing.T) {
	// Two orders on one phone: each code opens only the order it was issued for.
	f := newLookupFixture(lookupTestOrder(), lookupOlderTestOrder())

	got, err := f.u.VerifyOrderLookup(context.Background(), lookupPhone, olderLookupCode)
	if err != nil {
		t.Fatalf("VerifyOrderLookup with older order's code: %v", err)
	}
	if got.OrderID != lookupOlderOrderID {
		t.Errorf("OrderID = %q, want the older order %q", got.OrderID, lookupOlderOrderID)
	}
}

func TestVerifyOrderLookup_CodeFromAnotherPhoneRejected(t *testing.T) {
	f := newLookupFixture(lookupTestOrder())

	_, err := f.u.VerifyOrderLookup(context.Background(), otherLookupPhone, lookupCode)
	var invalid OrderLookupInvalidCodeError
	if !errors.As(err, &invalid) || invalid.RemainingAttempts != 4 {
		t.Fatalf("error = %v, want invalid code with 4 remaining", err)
	}
}

func TestVerifyOrderLookup_WrongCodeDecrementsRemainingAttempts(t *testing.T) {
	f := newLookupFixture(lookupTestOrder())

	for i, want := range []int{4, 3, 2, 1} {
		_, err := f.u.VerifyOrderLookup(context.Background(), lookupPhone, "ZZZZZZ")
		var invalid OrderLookupInvalidCodeError
		if !errors.As(err, &invalid) {
			t.Fatalf("attempt %d: error = %v, want OrderLookupInvalidCodeError", i+1, err)
		}
		if invalid.RemainingAttempts != want {
			t.Errorf("attempt %d: RemainingAttempts = %d, want %d", i+1, invalid.RemainingAttempts, want)
		}
	}
}

func TestVerifyOrderLookup_UnknownPhoneIsIndistinguishableFromWrongCode(t *testing.T) {
	f := newLookupFixture(lookupTestOrder())

	_, err := f.u.VerifyOrderLookup(context.Background(), otherLookupPhone, lookupCode)
	var invalid OrderLookupInvalidCodeError
	if !errors.As(err, &invalid) || invalid.RemainingAttempts != 4 {
		t.Fatalf("error = %v, want invalid code with 4 remaining", err)
	}
}

func TestVerifyOrderLookup_FifthWrongCodeLocksPhone(t *testing.T) {
	f := newLookupFixture(lookupTestOrder())

	for i := 0; i < lookupMaxAttempts-1; i++ {
		if _, err := f.u.VerifyOrderLookup(context.Background(), lookupPhone, "ZZZZZZ"); err == nil {
			t.Fatalf("attempt %d: want error", i+1)
		}
	}
	_, err := f.u.VerifyOrderLookup(context.Background(), lookupPhone, "ZZZZZZ")
	var locked OrderLookupLockedError
	if !errors.As(err, &locked) {
		t.Fatalf("5th attempt: error = %v, want OrderLookupLockedError", err)
	}
	if want := f.now.Add(15 * time.Minute); !locked.LockedUntil.Equal(want) {
		t.Errorf("LockedUntil = %v, want %v", locked.LockedUntil, want)
	}
	row := f.attempts.rows[lookupPhone]
	if row.FailedCount != 0 || row.LockedUntil == nil {
		t.Errorf("attempt row = %+v, want reset count and lock set", row)
	}
}

func TestVerifyOrderLookup_LockedPhoneRejectsWithoutCheckingCode(t *testing.T) {
	f := newLookupFixture(lookupTestOrder())
	until := f.now.Add(10 * time.Minute)
	f.attempts.rows[lookupPhone] = domain.LookupAttempt{Phone: lookupPhone, LockedUntil: &until}

	_, err := f.u.VerifyOrderLookup(context.Background(), lookupPhone, lookupCode)
	var locked OrderLookupLockedError
	if !errors.As(err, &locked) {
		t.Fatalf("error = %v, want OrderLookupLockedError even with the right code", err)
	}
	if !locked.LockedUntil.Equal(until) {
		t.Errorf("LockedUntil = %v, want %v", locked.LockedUntil, until)
	}
	if f.orders.listCalls != 0 {
		t.Errorf("order list called %d times while locked, want 0", f.orders.listCalls)
	}
}

func TestVerifyOrderLookup_LockExpiresAfterFifteenMinutes(t *testing.T) {
	f := newLookupFixture(lookupTestOrder())
	until := f.now.Add(-time.Second)
	f.attempts.rows[lookupPhone] = domain.LookupAttempt{Phone: lookupPhone, LockedUntil: &until}

	if _, err := f.u.VerifyOrderLookup(context.Background(), lookupPhone, lookupCode); err != nil {
		t.Fatalf("expired lock: unexpected error: %v", err)
	}
}

func TestAuthorizeOrder_TokenOpensOnlyItsOwnOrder(t *testing.T) {
	f := newLookupFixture(lookupTestOrder(), lookupOlderTestOrder())
	sess, err := f.u.VerifyOrderLookup(context.Background(), lookupPhone, lookupCode)
	if err != nil {
		t.Fatalf("VerifyOrderLookup: %v", err)
	}

	ok, err := f.u.AuthorizeOrder(context.Background(), sess.Token, lookupOrderID)
	if err != nil || !ok {
		t.Fatalf("AuthorizeOrder for the verified order = (%v, %v), want allowed", ok, err)
	}
	ok, err = f.u.AuthorizeOrder(context.Background(), sess.Token, lookupOlderOrderID)
	if err != nil || ok {
		t.Errorf("AuthorizeOrder for a different order of the same phone = (%v, %v), want rejected", ok, err)
	}
}

func TestAuthorizeOrder_ExpiredTokenRejected(t *testing.T) {
	f := newLookupFixture(lookupTestOrder())
	sess, err := f.u.VerifyOrderLookup(context.Background(), lookupPhone, lookupCode)
	if err != nil {
		t.Fatalf("VerifyOrderLookup: %v", err)
	}

	f.now = sess.ExpiresAt
	if ok, _ := f.u.AuthorizeOrder(context.Background(), sess.Token, lookupOrderID); ok {
		t.Error("token authorized at its expiry instant, want rejected")
	}
}

func TestAuthorizeOrder_UnknownAndEmptyTokensRejected(t *testing.T) {
	f := newLookupFixture(lookupTestOrder())

	for _, token := range []string{"", "not-a-real-token"} {
		if ok, err := f.u.AuthorizeOrder(context.Background(), token, lookupOrderID); ok || err != nil {
			t.Errorf("AuthorizeOrder(%q) = ok %v, err %v; want rejected", token, ok, err)
		}
	}
}

// Cancel uses the same gate as detail: a cancel without the token for that exact
// order must not reach the order, so the status change never runs.
func TestAuthorizeOrder_CancelRequiresTokenForThatOrder(t *testing.T) {
	f := newLookupFixture(lookupTestOrder(), lookupOlderTestOrder())
	sess, err := f.u.VerifyOrderLookup(context.Background(), lookupPhone, lookupCode)
	if err != nil {
		t.Fatalf("VerifyOrderLookup: %v", err)
	}

	if ok, _ := f.u.AuthorizeOrder(context.Background(), "", lookupOrderID); ok {
		t.Error("cancel without a token authorized, want rejected")
	}
	if ok, _ := f.u.AuthorizeOrder(context.Background(), sess.Token, lookupOlderOrderID); ok {
		t.Error("cancel of another order with this token authorized, want rejected")
	}
	if ok, _ := f.u.AuthorizeOrder(context.Background(), sess.Token, lookupOrderID); !ok {
		t.Error("cancel with the matching token rejected, want allowed")
	}
}
