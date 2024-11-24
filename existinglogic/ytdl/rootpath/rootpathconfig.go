package rootpath

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// This functionality will abstract the full path out,
// so that the save will go to an expected location

// Save Root Path configuration
func SaveRootPath(rootPath string) {

}

func GetRootPath() string {
	return ""
}

// Creates directory with full path if it does not exist (including hierarchy of nested dirs)
func CreateDirectoryIfNotExists(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// 1. Removes all OS-bound invalid file name chars
// 2. Removes all file name chars that may be valid, but annoying to deal with
// 3. ONLY APPLIES THIS TO DIR/FILE AT THE END AND DOES NOT TOUCH THE EXTENSION IF FILE
// 4. NON-PATHS CAN BE PLACED AS WELL
func RemoveInvalidFileNameChars(path string) string {

	if !isValidPath(path) {
		return sanitizeFilename(path)
	}

	dir := filepath.Dir(path)
	file := filepath.Base(path)
	ext := filepath.Ext(path)
	fileWithoutExt := strings.TrimSuffix(file, ext)

	sanitizedFile := sanitizeFilename(fileWithoutExt)
	sanitizedFileWithExt := sanitizedFile + ext

	pathWithSanitizedFileName := filepath.Join(dir, sanitizedFileWithExt)

	return pathWithSanitizedFileName
}

// / UTILITIES
func isValidPath(path string) bool {
	// Check if path is empty
	if path == "" {
		return false
	}

	// Different validation rules based on OS
	switch runtime.GOOS {
	case "windows":
		return isValidWindowsPath(path)
	case "linux", "darwin":
		return isValidUnixPath(path)
	default:
		return false
	}
}

// already assumes file name split from path
func sanitizeFilename(filename string) string {
	if filename == "" {
		return ""
	}

	// Different validation rules based on OS
	switch runtime.GOOS {
	case "windows":
		return sanitizeWindowsFilename(filename)
	case "linux", "darwin":
		return sanitizeUnixFilename(filename)
	default:
		return ""
	}
}

// sanitizes file name with path
func sanitizeFile(path string) string {
	if path == "" {
		return ""
	}

	// Different validation rules based on OS
	switch runtime.GOOS {
	case "windows":
		return sanitizeWindowsFile(path)
	case "linux", "darwin":
		return sanitizeUnixFile(path)
	default:
		return ""
	}
}

// sanitizes full path
func sanitizePath(path string) string {
	if path == "" {
		return ""
	}

	// Different validation rules based on OS
	switch runtime.GOOS {
	case "windows":
		return sanitizeWindowsPath(path)
	case "linux", "darwin":
		return sanitizeUnixPath(path)
	default:
		return ""
	}
}

// Check if path exists on filesystem
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Combined check for validity and existence
func isValidAndExists(path string) bool {
	return isValidPath(path) && pathExists(path)
}
