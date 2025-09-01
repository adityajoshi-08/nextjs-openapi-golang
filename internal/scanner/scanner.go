package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"nextjs-to-openapi/internal/models"
)

type Scanner struct {
	rootDir string
}

func NewScanner(rootDir string) *Scanner {
	return &Scanner{rootDir: rootDir}
}

// Simplified scanner - just find files and read content
func (s *Scanner) ScanRoutes() ([]models.APIRoute, error) {
	var routes []models.APIRoute

	err := filepath.WalkDir(s.rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		// Just check if it's a route file
		if isRouteFile(d.Name()) {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil // Skip problematic files
			}

			// Minimal processing - let Ollama figure out the rest
			route := models.APIRoute{
				FilePath: path,
				FileType: strings.TrimPrefix(filepath.Ext(path), "."),
				Content:  string(content),
				// Let Ollama determine: Path, Method, Parameters, etc.
			}

			routes = append(routes, route)
		}
		return nil
	})

	return routes, err
}

func isRouteFile(filename string) bool {
	// Match: route.js, route.ts, route.jsx, route.tsx
	matched, _ := regexp.MatchString(`^route\.(js|ts|jsx|tsx)$`, filename)
	return matched
}
