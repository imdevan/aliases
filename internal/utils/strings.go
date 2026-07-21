package utils

import (
	"path/filepath"
	"strings"
)

// GenerateAlias generates an alias name from a directory path.
func GenerateAlias(path string, separator string, lowercase bool, partLength int) string {
	cleaned := filepath.Clean(path)
	base := filepath.Base(cleaned)
	if base == "/" || base == "." {
		return "h"
	}

	var parts []string
	if separator != "" {
		parts = strings.Split(base, separator)
	} else {
		if strings.Contains(base, "-") {
			parts = strings.Split(base, "-")
		} else if strings.Contains(base, "_") {
			parts = strings.Split(base, "_")
		} else {
			parts = []string{base}
		}
	}

	var aliasParts []string
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		length := partLength
		if length > len(part) {
			length = len(part)
		}
		aliasParts = append(aliasParts, part[:length])
	}

	alias := strings.Join(aliasParts, separator)
	if lowercase {
		alias = strings.ToLower(alias)
	}
	return alias
}

// EscapeAliasValue escapes double quotes inside the alias value by replacing them with single quotes.
func EscapeAliasValue(val string) string {
	return strings.ReplaceAll(val, "\"", "'")
}
