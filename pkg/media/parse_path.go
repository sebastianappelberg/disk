package media

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Regex pathPatterns for different season formats
var pathPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^(?P<base>.+)\\(?P<title>[^\\]+)\\Season[ _]?(?P<season>\d+)\\(?P<episode>.+)$`),
	regexp.MustCompile(`(?i)^(?P<base>.+)\\(?P<title>[^\\]+)\\(?P<title2>[^\\]+)[ _-]Season[ _]?(?P<season>\d+)\\(?P<episode>.+)$`),
	regexp.MustCompile(`(?i)^(?P<base>.+)\\(?P<title>[^\\]+)\\(?P<episode>.+)$`), // No season folder
}

type ParsedPath struct {
	Title   string
	Season  int
	Episode string
}

func ParsePath(path string) ParsedPath {
	path = strings.TrimSpace(path) // Remove leading/trailing spaces
	path = filepath.Clean(path)    // Normalize path separators

	for _, pattern := range pathPatterns {
		match := pattern.FindStringSubmatch(path)
		if match != nil {
			result := make(map[string]string)
			for i, name := range pattern.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}

			// Extracted values
			title := result["title"]
			// If no season and if no siblings, assume that it's a movie.
			// If no season and siblings, assume that it's season 1 and that the episodes are in lexicographical order.
			season := 0
			if result["season"] != "" {
				season, _ = strconv.Atoi(strings.TrimLeft(result["season"], "0"))
			}
			//season := strings.TrimLeft(result["season"], "0") // This might be empty if no season info
			episode := result["episode"]

			// If a second title exists, use it as the title instead
			if val, exists := result["title2"]; exists {
				title = val
			}

			return ParsedPath{Title: title, Season: season, Episode: episode}
		}
	}
	return ParsedPath{}
}
