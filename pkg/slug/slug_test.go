// pkg/slug/slug_test.go
package slug

import "testing"

func TestGenerate(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"vietnamese diacritics", "Tranh Đồng", "tranh-dong"},
		{"multi-word vietnamese", "Tượng Đồng Trang Trí", "tuong-dong-trang-tri"},
		{"punctuation and extra spaces", "  Hello,   World!!  ", "hello-world"},
		{"ampersand", "Café & Bar", "cafe-bar"},
		{"already a slug", "tranh-dong", "tranh-dong"},
		{"empty", "", ""},
		{"mixed case english", "Best Seller 2026", "best-seller-2026"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Generate(tc.in)
			if got != tc.want {
				t.Errorf("Generate(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
