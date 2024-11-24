package rootpath

import (
	"path/filepath"
	"strings"
	"unicode"
)

func isValidUnixPath(path string) bool {
	// Check for null characters (invalid in Unix paths)
	if strings.Contains(path, "\x00") {
		return false
	}

	// Try to get absolute path
	_, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Unix systems typically have a path length limit of 4096 characters
	if len(path) > 4096 {
		return false
	}

	return true
}

func sanitizeUnixFile(filename string) string {
	// Step 1: Convert to clean filename without path
	filename = filepath.Base(filename)

	return sanitizeUnixFilename(filename)
}

func sanitizeUnixFilename(filename string) string {
	// Step 2: Handle empty or dot-only filenames
	if filename == "" || strings.TrimLeft(filename, ".") == "" {
		return "_"
	}

	// Step 3: Replace invalid characters
	invalidChars := []string{
		"/",    // Path separator
		"\x00", // Null byte
		"\\",   // Backslash (for compatibility)
	}

	result := filename
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}

	// Step 4: Handle control characters and non-printing characters
	result = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || !unicode.IsPrint(r) {
			return '_'
		}
		return r
	}, result)

	// Step 5: Trim spaces and dots from ends
	result = strings.Trim(result, " .")

	// Step 6: Handle special cases
	if result == "." || result == ".." {
		return "_" + result
	}

	// Step 7: Ensure filename isn't empty after all replacements
	if result == "" {
		return "_"
	}

	// Step 8: Handle maximum length (255 bytes is common max for many Unix filesystems)
	return truncateUnixFilename(result, 255)
}

func truncateUnixFilename(filename string, maxBytes int) string {
	if len(filename) <= maxBytes {
		return filename
	}

	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]

	// Calculate available space for the name
	maxNameBytes := maxBytes - len(ext)
	if maxNameBytes < 1 {
		// If extension is too long, truncate it
		return filename[:maxBytes]
	}

	// Truncate the name part while preserving the extension
	return nameWithoutExt[:maxNameBytes] + ext
}

// Optional: Additional utility functions that might be useful

func isValidUnixFilename(filename string) bool {
	// Check if filename is valid for Unix systems
	if filename == "" || len(filename) > 255 {
		return false
	}

	// Check for invalid characters
	if strings.ContainsAny(filename, "/\x00") {
		return false
	}

	// Check if filename is "." or ".."
	if filename == "." || filename == ".." {
		return false
	}

	// Check for control characters
	for _, r := range filename {
		if unicode.IsControl(r) {
			return false
		}
	}

	return true
}

func sanitizeUnixPath(path string) string {
	// Handle full paths by sanitizing each component
	parts := strings.Split(filepath.Clean(path), string(filepath.Separator))
	for i, part := range parts {
		if part != "" {
			parts[i] = sanitizeUnixFilename(part)
		}
	}

	// Reconstruct the path
	if strings.HasPrefix(path, "/") {
		return "/" + strings.Join(parts, "/")
	}
	return strings.Join(parts, "/")
}
