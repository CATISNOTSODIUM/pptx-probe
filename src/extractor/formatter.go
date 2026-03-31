package extractor

import (
	"regexp"
	"strings"
)

func FixPowerPointQuotes(input string) string {
	// Define the pairs of (old, new)
	replacer := strings.NewReplacer(
		"“", "\"", // Opening curly quote
		"”", "\"", // Closing curly quote
	)

	return replacer.Replace(input)
}

func StripHeaderMarker(input string) string {
	parts := strings.SplitN(input, "\n", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func ExtractFileName(input string) string {
	// Character class breakdown:
	// \w   : Alphanumeric + underscore
	// \.   : Literal dot
	// \/   : Forward slash
	// \\   : Literal backslash (escaped)
	// -    : Literal hyphen
	re := regexp.MustCompile(`#!([\w\.\/\\-]+)`)

	match := re.FindStringSubmatch(input)

	if len(match) > 1 {
		return match[1]
	}
	return ""
}
