package rootpath

import (
	"path/filepath"
	"strings"
	"unicode"
)

// Windows-specific constants
const (
	WindowsMaxPath   = 260 // Standard MAX_PATH
	WindowsMaxName   = 255 // Maximum filename length
	WindowsExtMaxLen = 4   // Including the dot
)

func isValidWindowsPath(path string) bool {
	// Check for invalid Windows characters
	if strings.ContainsAny(path, "<>:\"|?*") {
		return false
	}

	// Try to get absolute path
	_, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Check Windows path length limit
	if len(path) > 259 {
		return false
	}

	// Check drive letter format if present
	if len(path) >= 2 && path[1] == ':' {
		driveLetter := path[0]
		if !((driveLetter >= 'A' && driveLetter <= 'Z') || (driveLetter >= 'a' && driveLetter <= 'z')) {
			return false
		}
	}

	return true
}

func sanitizeWindowsFile(filename string) string {
	// Step 1: Convert to clean filename without path
	filename = filepath.Base(filename)

	return sanitizeWindowsFilename(filename)
}

func sanitizeWindowsFilename(filename string) string {
	// Step 2: Handle empty filename
	if filename == "" {
		return "_"
	}

	// Step 3: Replace invalid characters
	// Windows invalid: < > : " / \ | ? *
	invalidChars := []string{
		"<", ">", ":", "\"", "/", "\\", "|", "?", "*",
		"\x00", "\x01", "\x02", "\x03", "\x04", "\x05", "\x06", "\x07",
		"\x08", "\x09", "\x0a", "\x0b", "\x0c", "\x0d", "\x0e", "\x0f",
		"\x10", "\x11", "\x12", "\x13", "\x14", "\x15", "\x16", "\x17",
		"\x18", "\x19", "\x1a", "\x1b", "\x1c", "\x1d", "\x1e", "\x1f",
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

	// Step 6: Handle Windows reserved names
	result = handleWindowsReservedNames(result)

	// Step 7: Ensure filename isn't empty after all replacements
	if result == "" {
		return "_"
	}

	// Step 8: Handle maximum length
	return truncateWindowsFilename(result, WindowsMaxName)
}

func handleWindowsReservedNames(filename string) string {
	// Windows reserved names (case-insensitive)
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	// Get the name without extension
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Check if the name matches any reserved name (case-insensitive)
	upperName := strings.ToUpper(nameWithoutExt)
	for _, reserved := range reservedNames {
		if upperName == reserved {
			return "_" + filename
		}
	}

	return filename
}

func truncateWindowsFilename(filename string, maxLength int) string {
	if len(filename) <= maxLength {
		return filename
	}

	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// If extension is longer than 4 chars (including dot), truncate it
	if len(ext) > WindowsExtMaxLen {
		ext = ext[:WindowsExtMaxLen]
	}

	// Calculate available space for the name
	maxNameLength := maxLength - len(ext)
	if maxNameLength < 1 {
		maxNameLength = 1
	}

	return nameWithoutExt[:maxNameLength] + ext
}

func sanitizeWindowsPath(path string) string {
	// Handle full paths
	if len(path) > WindowsMaxPath {
		// If path is too long, truncate it while preserving the drive letter and separator
		if len(path) >= 2 && path[1] == ':' {
			return path[:2] + sanitizeWindowsFile(path[2:])
		}
		return sanitizeWindowsFile(path)
	}

	// Handle drive letter
	if len(path) >= 2 && path[1] == ':' {
		driveLetter := strings.ToUpper(path[:2])
		remainingPath := path[2:]
		parts := strings.Split(remainingPath, "\\")

		// Sanitize each part
		for i, part := range parts {
			if part != "" {
				parts[i] = sanitizeWindowsFile(part)
			}
		}

		return driveLetter + strings.Join(parts, "\\")
	}

	// Handle UNC paths
	if strings.HasPrefix(path, "\\\\") {
		parts := strings.Split(path[2:], "\\")
		for i, part := range parts {
			if i < 2 { // Server and share name
				continue
			}
			if part != "" {
				parts[i] = sanitizeWindowsFile(part)
			}
		}
		return "\\\\" + strings.Join(parts, "\\")
	}

	// Regular path
	parts := strings.Split(path, "\\")
	for i, part := range parts {
		if part != "" {
			parts[i] = sanitizeWindowsFile(part)
		}
	}
	return strings.Join(parts, "\\")
}
