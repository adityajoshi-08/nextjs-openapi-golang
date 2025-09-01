package models

// APIRoute represents a discovered API route in Next.js
type APIRoute struct {
	Path       string   `json:"path"`
	Method     string   `json:"method"`
	FilePath   string   `json:"file_path"`
	FileType   string   `json:"file_type"` // "ts", "js", "tsx", "jsx"
	Parameters []string `json:"parameters,omitempty"`
	Content    string   `json:"content"`
}

// DocumentedRoute represents an API route with generated documentation
type DocumentedRoute struct {
	Route       APIRoute `json:"route"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Error       error    `json:"-"`
}

// Config holds CLI configuration
type Config struct {
	APIDir      string `json:"api_dir"`
	OutputFile  string `json:"output_file"`
	OllamaModel string `json:"ollama_model"`
	Workers     int    `json:"workers"`
	OllamaURL   string `json:"ollama_url"`
}
