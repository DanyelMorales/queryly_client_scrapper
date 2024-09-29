package queryl

import (
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

func removeDiacritics(s string) string {
	// Normalize the string to decompose combined characters
	t := norm.NFD.String(s)

	// Use a builder to create the final string without diacritics
	var b strings.Builder
	for _, r := range t {
		// Filter out all combining characters (diacritics)
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func ProcessString(input string) string {
	// Step 1: Convert to lowercase
	lower := strings.ToLower(input)

	// Step 2: Remove diacritical marks (accents)
	normalized := removeDiacritics(lower)

	// Step 3: Replace single quotes with escaped version
	result := strings.ReplaceAll(normalized, "'", "\\'")

	return result
}

func BuildQueryParams(queryKey, endIndex, batchSize string) *map[string]string {
	return &map[string]string{
		"queryly_key":        queryKey,
		"endindex":           endIndex,
		"batchsize":          batchSize,
		"query":              "Noticia",
		"showfaceted":        "true",
		"extendeddatafields": "creator,subheadline,pubdateunix,link,creator",
		"timezoneoffset":     "480",
	}
}
