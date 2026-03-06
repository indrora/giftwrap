package internal

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func Slugify(s string) string {
	// NFD-normalize so accented characters decompose into base + combining marks,
	// then drop everything that isn't printable ASCII (removes combining marks,
	// emojis, and other non-ASCII runes).
	nfd := norm.NFD.String(s)
	var ascii strings.Builder
	for _, r := range nfd {
		if r <= unicode.MaxASCII && unicode.IsPrint(r) {
			ascii.WriteRune(r)
		}
	}

	slugRE := regexp.MustCompile(`[^a-z0-9]+`)
	lower := strings.ToLower(strings.TrimSpace(ascii.String()))
	return strings.Trim(slugRE.ReplaceAllLiteralString(lower, "-"), "-")
}
