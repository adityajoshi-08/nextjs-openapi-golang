package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"nextjs-to-openapi/internal/models"
	"nextjs-to-openapi/internal/ollama"
	"nextjs-to-openapi/internal/scanner"

	"github.com/spf13/cobra"
)

// Simple OpenAPI structure
type OpenAPISpec struct {
	OpenAPI string                 `json:"openapi"`
	Info    map[string]interface{} `json:"info"`
	Paths   map[string]interface{} `json:"paths"`
}

func buildOpenAPISpec(client *ollama.Client, routes []models.APIRoute) OpenAPISpec {
	spec := OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: map[string]interface{}{
			"title":   "Next.js API Documentation",
			"version": "1.0.0",
		},
		Paths: make(map[string]interface{}),
	}

	fmt.Printf("\n🔄 Processing all %d routes...\n", len(routes))

	for i, route := range routes {
		fmt.Printf("Processing route %d/%d: %s\n", i+1, len(routes), route.FilePath)

		doc, err := client.DocumentRoute(route)
		if err != nil {
			fmt.Printf("⚠️ Error documenting %s: %v\n", route.FilePath, err)
			continue
		}

		// Convert to proper OpenAPI structure
		pathItem := make(map[string]interface{})
		for method, details := range doc.Methods {
			// Convert method to lowercase (OpenAPI requirement)
			methodLower := strings.ToLower(method)

			// Fix parameter structure
			var fixedParams []map[string]interface{}
			for _, param := range details.Parameters {
				fixedParam := map[string]interface{}{
					"name":     param.Name,
					"in":       param.In,
					"required": param.Required,
					"schema": map[string]interface{}{
						"type": param.Type,
					},
				}
				fixedParams = append(fixedParams, fixedParam)
			}

			// Add required responses section
			responses := map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Successful response",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type":        "object",
								"description": "Response data",
							},
						},
					},
				},
				"400": map[string]interface{}{
					"description": "Bad request",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"error": map[string]interface{}{
										"type": "string",
									},
								},
							},
						},
					},
				},
				"500": map[string]interface{}{
					"description": "Internal server error",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"error": map[string]interface{}{
										"type": "string",
									},
								},
							},
						},
					},
				},
			}

			pathItem[methodLower] = map[string]interface{}{
				"summary":     details.Summary,
				"description": details.Description,
				"parameters":  fixedParams,
				"responses":   responses, // ✅ Required responses section
			}
		}

		spec.Paths[doc.Path] = pathItem
	}

	return spec
}

func writeOpenAPIFile(filename string, spec OpenAPISpec) error {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

var (
	apiDir      string
	outputFile  string
	ollamaModel string
	workers     int
	ollamaURL   string
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var rootCmd = &cobra.Command{
	Use:   "nextjs-to-openapi",
	Short: "Convert Next.js API routes to OpenAPI specification",
	Long: `A CLI tool that scans your Next.js API routes and generates 
OpenAPI specification using Ollama for intelligent documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("🚀 Starting Next.js to OpenAPI conversion...\n")
		fmt.Printf("API Directory: %s\n", apiDir)
		fmt.Printf("Output File: %s\n", outputFile)
		fmt.Printf("Ollama Model: %s\n", ollamaModel)
		fmt.Printf("Workers: %d\n", workers)

		// Create scanner and scan for routes
		s := scanner.NewScanner(apiDir)
		routes, err := s.ScanRoutes()
		if err != nil {
			fmt.Printf("❌ Error scanning routes: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Found %d routes\n", len(routes))

		// Optional: Show route details (you can remove this debug section)
		if len(routes) > 0 {
			fmt.Printf("\n📋 Route Details:\n")
			for i, route := range routes {
				fmt.Printf("%d. File: %s\n", i+1, route.FilePath)
				fmt.Printf("   Type: %s\n", route.FileType)
				fmt.Printf("   Content preview (first 50 chars): %s...\n",
					route.Content[:min(50, len(route.Content))])
			}
		}

		if len(routes) == 0 {
			fmt.Printf("No routes found. Exiting.\n")
			return
		}

		// Create Ollama client
		client := ollama.NewClient(ollamaURL, ollamaModel)

		// Process all routes and build OpenAPI spec
		fmt.Printf("\n🤖 Generating documentation for all routes...\n")
		openAPISpec := buildOpenAPISpec(client, routes)

		// Write to file
		err = writeOpenAPIFile(outputFile, openAPISpec)
		if err != nil {
			fmt.Printf("❌ Error writing OpenAPI file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ OpenAPI specification written to: %s\n", outputFile)
		fmt.Printf("📁 File contains %d documented endpoints\n", len(openAPISpec.Paths))
	},
}

func init() {
	rootCmd.Flags().StringVarP(&apiDir, "api-dir", "d", "./api", "Directory containing Next.js API routes")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "openapi.json", "Output file for OpenAPI specification")
	rootCmd.Flags().StringVarP(&ollamaModel, "model", "m", "llama3.1", "Ollama model to use for documentation generation")
	rootCmd.Flags().IntVarP(&workers, "workers", "w", 3, "Number of worker goroutines")
	rootCmd.Flags().StringVar(&ollamaURL, "ollama-url", "http://localhost:11434", "Ollama server URL")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
