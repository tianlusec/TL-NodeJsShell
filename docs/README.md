# Documentation

Welcome to the TL-NodeJsShell documentation!

## Available Documentation

### Getting Started
- **[Installation Guide](INSTALLATION.md)** - Complete installation instructions for all platforms
- **[Quick Start](../README.md#quick-start)** - Get up and running quickly

### Usage
- **[API Documentation](API.md)** - Complete API reference
- **[User Guide](../README.md#usage)** - How to use TL-NodeJsShell

### Development
- **[Contributing Guide](../CONTRIBUTING.md)** - How to contribute to the project
- **[Architecture Overview](#architecture)** - System architecture and design

### Security
- **[Security Policy](SECURITY.md)** - Security guidelines and vulnerability reporting
- **[Best Practices](SECURITY.md#security-best-practices)** - Security recommendations

### Project Information
- **[Changelog](../CHANGELOG.md)** - Version history and updates
- **[License](../LICENSE)** - MIT License

---

## Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────┐
│      Web Browser       │
│     (Vue 3 Frontend)      │
└────────────────────┬────────────────────────────────────┘
      │ HTTP/WebSocket
      │
┌────────────────────▼────────────────────────────────────┐
│     Backend Server       │
│     (Go + Gin)        │
│ ┌──────────────────────────────────────────────────┐ │
│ │ API Layer (Handlers)       │ │
│ └──────────────┬───────────────────────────────────┘ │
│     │          │
│ ┌──────────────▼───────────────────────────────────┐ │
│ │ Core Services         │ │
│ │ - Payload Generator        │ │
│ │ - Crypto (Base64/XOR/AES)      │ │
│ │ - Transport (HTTP/Multipart)     │ │
│ │ - Proxy Manager         │ │
│ └──────────────┬───────────────────────────────────┘ │
│     │          │
│ ┌──────────────▼───────────────────────────────────┐ │
│ │ Database Layer (SQLite + GORM)     │ │
│ └──────────────────────────────────────────────────┘ │
└────────────────────┬────────────────────────────────────┘
      │ HTTP/HTTPS
      │ (with optional proxy)
┌────────────────────▼────────────────────────────────────┐
│    Target Node.js Server      │
│    (WebShell Endpoint)       │
└─────────────────────────────────────────────────────────┘
```

### Component Details

#### Frontend (Vue 3 + TypeScript)
- **Framework**: Vue 3 with Composition API
- **UI Library**: Element Plus
- **State Management**: Pinia
- **Routing**: Vue Router
- **Terminal**: xterm.js
- **Editor**: Monaco Editor
- **Build Tool**: Vite

#### Backend (Go)
- **Framework**: Gin Web Framework
- **ORM**: GORM
- **Database**: SQLite
- **HTTP Client**: Custom with proxy support

#### Core Modules

**1. Payload Generator**
- Template-based payload generation
- Multiple encoding support
- Memory shell injection templates

**2. Crypto Module**
- Base64 encoding/decoding
- XOR encryption
- AES-256 encryption

**3. Transport Layer**
- HTTP/HTTPS client
- Multipart form data
- Custom headers support
- Proxy integration

**4. Proxy Manager**
- HTTP/HTTPS proxy
- SOCKS5 proxy
- Authentication support

---

## Development Workflow

### Backend Development

```bash
# Install dependencies
cd backend
go mod download

# Run tests
go test ./...

# Run with hot reload (using air)
go install github.com/cosmtrek/air@latest
air

# Build
go build -o NodeJsshell main.go
```

### Frontend Development

```bash
# Install dependencies
cd frontend
npm install

# Development server
npm run dev

# Type checking
npm run type-check

# Build for production
npm run build

# Preview production build
npm run preview
```

### Code Structure

```
backend/
├── app/
│ ├── app.go    # Application initialization
│ ├── middleware/   # HTTP middleware
│ └── routes/    # Route definitions
├── config/
│ └── config.go   # Configuration management
├── core/
│ ├── crypto/    # Encryption utilities
│ ├── exploit/   # Exploit modules
│ ├── payload/   # Payload generation
│ ├── proxy/    # Proxy management
│ └── transport/   # HTTP transport
├── database/
│ ├── db.go    # Database connection
│ └── shell.go   # Data models
├── handlers/
│ ├── shellHandler.go  # Shell management
│ ├── fileHandler.go  # File operations
│ ├── cmdHandler.go  # Command execution
│ └── payloadHandler.go # Payload generation
└── main.go     # Entry point

frontend/
├── src/
│ ├── api/    # API client functions
│ ├── components/   # Reusable components
│ ├── views/    # Page components
│ ├── stores/    # Pinia stores
│ ├── router/    # Route configuration
│ ├── types/    # TypeScript types
│ ├── App.vue    # Root component
│ └── main.ts    # Entry point
├── public/     # Static assets
└── index.html    # HTML template
```

---

## API Integration

### Example: Creating a Shell

```typescript
// TypeScript (Frontend)
import axios from 'axios'

interface ShellConfig {
 url: string
 password: string
 encode_type: stringmethod: string
}

async function createShell(config: ShellConfig) {
 const response = await axios.post('/api/shells', config)
 return response.data
}
```

```go
// Go (Backend)
func (h *ShellHandler) Create(c *gin.Context) {
 var req struct {
  URL  string `json:"url" binding:"required"`
  Password string `json:"password" binding:"required"`
  EncodeType string `json:"encode_type"`
  Method  string `json:"method"`
 }
 
 if err := c.ShouldBindJSON(&req); err != nil {
  c.JSON(400, gin.H{"error": err.Error()})
  return
 }
 
 shell := database.Shell{
  URL:  req.URL,
  Password: req.Password,
  EncodeType: req.EncodeType,
  Method:  req.Method,
  Status:  "offline",
 }
 
 h.db.Create(&shell)
 c.JSON(201, shell)
}
```

---

## Testing

### Backend Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./core/crypto/...

# Verbose output
go test -v ./...
```

### Frontend Tests

```bash
# Run unit tests
npm run test

# Run with coverage
npm run test:coverage

# Run in watch mode
npm run test:watch
```

---

## Deployment

### Production Build

```bash
# Build backend
cd backend
go build -ldflags="-s -w" -o NodeJsshell main.go

# Build frontend
cd frontend
npm run build

# The frontend build output will be in frontend/dist/
# Configure backend to serve these static files
```

### Environment Variables

```bash
# Backend
export PORT=8080
export HOST=0.0.0.0
export DB_PATH=./data/shells.db

# Frontend (build time)
export VITE_API_BASE_URL=http://localhost:8080
```

---

## Troubleshooting

### Common Development Issues

**1. CORS Errors**
- Ensure backend CORS middleware is properly configured
- Check frontend proxy settings in `vite.config.ts`

**2. Database Locked**
- Close other connections to the SQLite database
- Ensure proper transaction handling

**3. Module Not Found**
- Run `go mod tidy` in backend
- Run `npm install` in frontend

---

## Resources

### External Documentation
- [Go Documentation](https://go.dev/doc/)
- [Vue 3 Documentation](https://vuejs.org/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [Element Plus](https://element-plus.org/)
- [GORM](https://gorm.io/docs/)

### Community
- [GitHub Issues](https://github.com/tianlusec/TL-NodeJsShell/issues)
- [GitHub Discussions](https://github.com/tianlusec/TL-NodeJsShell/discussions)

---

## Contributing

We welcome contributions! Please see our [Contributing Guide](../CONTRIBUTING.md) for details.

---

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.