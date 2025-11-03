package slug

import (
	"regexp"
	"strings"
)

// GenerateSlug creates a clean URL-safe slug from a string
func GenerateSlug(name string) (string, error) {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace common symbols
	slug = strings.ReplaceAll(slug, "&", "and")

	// Replace spaces with dashes
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove all non-alphanumeric and non-dash characters
	re := regexp.MustCompile(`[^a-z0-9\-]+`)
	slug = re.ReplaceAllString(slug, "")

	// Remove multiple dashes
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	// Trim leading/trailing dashes
	slug = strings.Trim(slug, "-")

	return slug, nil
}
