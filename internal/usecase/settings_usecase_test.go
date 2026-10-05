package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

// fakeSettingsRepo is an in-memory domain.SiteSettingsRepository for
// usecase-level unit tests — no database needed.
type fakeSettingsRepo struct {
	values   map[string]string
	setBulks int
}

func newFakeSettingsRepo() *fakeSettingsRepo {
	return &fakeSettingsRepo{values: map[string]string{"active_theme": "default"}}
}

func (f *fakeSettingsRepo) Get(ctx context.Context, key string) (string, error) {
	return f.values[key], nil
}

func (f *fakeSettingsRepo) GetAll(ctx context.Context) (map[string]string, error) {
	out := make(map[string]string, len(f.values))
	for k, v := range f.values {
		out[k] = v
	}
	return out, nil
}

func (f *fakeSettingsRepo) Set(ctx context.Context, key string, value string) error {
	f.values[key] = value
	return nil
}

func (f *fakeSettingsRepo) SetBulk(ctx context.Context, settings map[string]string) error {
	f.setBulks++
	for k, v := range settings {
		f.values[k] = v
	}
	return nil
}

func TestUpdateSettings_ActiveTheme(t *testing.T) {
	tests := []struct {
		name    string
		in      map[string]string
		wantErr bool
		wantMsg string
	}{
		{name: "known theme default", in: map[string]string{"active_theme": "default"}},
		{name: "known theme tet", in: map[string]string{"active_theme": "tet"}},
		{name: "known theme independence", in: map[string]string{"active_theme": "independence"}},
		{name: "known theme labor-day", in: map[string]string{"active_theme": "labor-day"}},
		{name: "unknown theme dark", in: map[string]string{"active_theme": "dark"}, wantErr: true, wantMsg: "unknown theme: dark"},
		{name: "empty theme rejected", in: map[string]string{"active_theme": ""}, wantErr: true, wantMsg: "unknown theme: "},
		{name: "no active_theme key is allowed", in: map[string]string{"hotline": "0899012288"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeSettingsRepo()
			u := NewSettingsUsecase(repo)

			_, err := u.UpdateSettings(context.Background(), tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("UpdateSettings(%v) returned nil error, want error", tt.in)
				}
				if err.Error() != tt.wantMsg {
					t.Errorf("error = %q, want %q", err.Error(), tt.wantMsg)
				}
				if !errors.Is(err, domain.ErrInvalidInput) {
					t.Errorf("error does not wrap domain.ErrInvalidInput: %v", err)
				}
				if repo.setBulks != 0 {
					t.Errorf("SetBulk called %d times for rejected input, want 0", repo.setBulks)
				}
				if repo.values["active_theme"] != "default" {
					t.Errorf("active_theme changed to %q after rejected input", repo.values["active_theme"])
				}
				return
			}
			if err != nil {
				t.Fatalf("UpdateSettings(%v) returned error: %v", tt.in, err)
			}
		})
	}
}
