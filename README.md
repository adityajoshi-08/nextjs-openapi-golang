# Next.js to OpenAPI Generator -- Under Construction and Not Complete

A powerful CLI tool written in Go that automatically converts Next.js App Router API routes to OpenAPI 3.0 specifications using AI-powered documentation generation.

## Features

🚀 **Automatic Discovery** - Recursively scans Next.js API routes (`route.js`, `route.ts`, `route.jsx`, `route.tsx`)  
🤖 **AI-Powered Documentation** - Uses Ollama to generate intelligent API documentation  
📝 **OpenAPI 3.0 Compliant** - Generates industry-standard OpenAPI specifications  
🔄 **Dynamic Route Support** - Converts `[id]` and `[...slug]` to OpenAPI path parameters  
⚡ **TypeScript & JavaScript** - Supports both TS and JS Next.js projects  
📊 **Swagger Compatible** - Generated specs work with Swagger UI, Postman, and other tools  

## Prerequisites

- **Go 1.21+** - [Install Go](https://golang.org/doc/install)
- **Ollama** - [Install Ollama](https://ollama.ai/download)
- **Ollama Model** - Download a model (e.g., `ollama pull gemma:2b`)

## Installation

### Option 1: Install from Source
```bash
git clone https://github.com/adityajoshi-08/nextjs-openapi-golang.git
cd nextjs-openapi-golang
go build -o nextjs-to-openapi cmd/main.go
```

### Option 2: Using Make
```bash
git clone https://github.com/adityajoshi-08/nextjs-openapi-golang.git
cd nextjs-openapi-golang
make build
```

## Quick Start

1. **Start Ollama** (if not already running):
   ```bash
   ollama serve
   ```

2. **Run the tool** on your Next.js project:
   ```bash
   ./nextjs-to-openapi --api-dir ./your-nextjs-app/app/api --model gemma:2b --output api-docs.json
   ```

3. **View your documentation** in Swagger UI:
   - Visit [editor.swagger.io](https://editor.swagger.io)
   - Copy-paste the contents of `api-docs.json`
   - Enjoy interactive API documentation!

## Usage

```bash
./nextjs-to-openapi [flags]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--api-dir` | `-d` | `./api` | Directory containing Next.js API routes |
| `--output` | `-o` | `openapi.json` | Output file for OpenAPI specification |
| `--model` | `-m` | `llama3.1` | Ollama model to use for documentation |
| `--workers` | `-w` | `3` | Number of worker goroutines (future feature) |
| `--ollama-url` | | `http://localhost:11434` | Ollama server URL |

### Examples

```bash
# Basic usage
./nextjs-to-openapi --api-dir ./app/api

# Custom model and output file
./nextjs-to-openapi -d ./src/app/api -m gemma:2b -o docs/openapi.json

# Remote Ollama instance
./nextjs-to-openapi --api-dir ./api --ollama-url http://192.168.1.100:11434
```

## Supported Next.js Patterns

### File Structure
```
app/api/
├── users/
│   ├── route.ts          ✅ /api/users
│   └── [id]/
│       └── route.js      ✅ /api/users/{id}
├── posts/
│   └── [...slug]/
│       └── route.tsx     ✅ /api/posts/{slug}
└── auth/
    └── login/
        └── route.jsx     ✅ /api/auth/login
```

### HTTP Methods
```typescript
// route.ts
export async function GET(request: Request) { /* ... */ }
export async function POST(request: Request) { /* ... */ }
export async function PUT(request: Request) { /* ... */ }
export async function DELETE(request: Request) { /* ... */ }
export async function PATCH(request: Request) { /* ... */ }
```

## Output Example

The tool generates OpenAPI 3.0 specifications like this:

```json
{
  "openapi": "3.0.0",
  "info": {
    "title": "Next.js API Documentation",
    "version": "1.0.0"
  },
  "paths": {
    "/api/users/{id}": {
      "GET": {
        "summary": "Get user by ID",
        "description": "Retrieves a specific user's information using their unique identifier",
        "parameters": [
          {
            "name": "id",
            "type": "string",
            "in": "path",
            "required": true
          }
        ]
      }
    }
  }
}
```

## Development

### Using Make
```bash
# Build the project
make build

# Build and run with default settings
make run

# Clean build artifacts
make clean
```

### Manual Build
```bash
# Build
go build -o nextjs-to-openapi cmd/main.go

# Run
./nextjs-to-openapi --api-dir ./app/api --model gemma:2b
```

## Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Scanner   │───▶│   Ollama    │───▶│   OpenAPI   │───▶│    File     │
│             │    │   Client    │    │   Builder   │    │   Output    │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
      │                    │                   │                 │
   Discovers           AI Analysis         Structures        Generates
  API routes         & Documentation      OpenAPI spec       JSON file
```

## Troubleshooting

### Common Issues

**"Ollama returned status 404"**
- Ensure Ollama is running: `ollama serve`
- Verify model exists: `ollama list`
- Check model name spelling

**"No routes found"**
- Verify API directory path
- Ensure route files are named `route.{js,ts,jsx,tsx}`
- Check file permissions

**"failed to parse JSON response"**
- Try a different Ollama model
- Ensure sufficient system resources for AI processing

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Roadmap

- [ ] **Parallel Processing** - Goroutines for concurrent route processing
- [ ] **Request/Response Schemas** - Generate complete data models
- [ ] **Authentication Documentation** - Support for auth schemes
- [ ] **Error Response Documentation** - Document error cases
- [ ] **Configuration File Support** - YAML/JSON config files
- [ ] **Multiple Output Formats** - YAML, Swagger UI HTML

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- **Next.js** - The React framework that makes this tool necessary
- **Ollama** - Local AI that powers intelligent documentation
- **OpenAPI Initiative** - Standard specification format
- **Cobra** - Excellent Go CLI framework

---

⭐ If this tool helped you, please consider starring the repository!

For questions, issues, or feature requests, please [open an issue](https://github.com/adityajoshi-08/nextjs-openapi-golang/issues).