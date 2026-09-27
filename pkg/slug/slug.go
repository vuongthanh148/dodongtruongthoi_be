// pkg/slug/slug.go
package slug

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	nonAlnum  = regexp.MustCompile(`[^a-z0-9]+`)
	trimDash  = regexp.MustCompile(`^-+|-+$`)
	stripMarks = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
)

// Generate converts a display name (Vietnamese or English) into a
// lowercase, URL-safe, latin slug: diacritics stripped, everything
// that isn't a letter or digit collapsed to a single hyphen, leading
// and trailing hyphens trimmed. "Tranh Đồng" -> "tranh-dong".
func Generate(name string) string {
	// đ/Đ don't decompose under NFD (they're distinct letters, not a
	// base letter + combining mark), so fold them explicitly first.
	folded := strings.NewReplacer("đ", "d", "Đ", "D").Replace(name)

	ascii, _, err := transform.String(stripMarks, folded)
	if err != nil {
		ascii = folded
	}

	lower := strings.ToLower(ascii)
	dashed := nonAlnum.ReplaceAllString(lower, "-")
	return trimDash.ReplaceAllString(dashed, "")
}
