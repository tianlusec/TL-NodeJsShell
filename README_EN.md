# TL-NodeJsShell<div align="center">

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Vue Version](https://img.shields.io/badge/Vue-3.3+-4FC08D?logo=vue.js)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)

**A Modern Node.js WebShell Management Platform**

**Developed by Tianlu Laboratory**

English | [简体中文](README.md)

</div>

---

## Legal Disclaimer

**This tool is intended for authorized security testing and educational purposes only.**

- Users must obtain explicit permission before testing any systems
- Unauthorized access to computer systems is illegal
- Users are solely responsible for their actions
- The authors assume no liability for misuse or damage

## Overview

TL-NodeJsShell is a comprehensive WebShell management platform designed for security professionals and penetration testers. It provides a modern web interface for managing Node.js-based shells with advanced features including memory shell injection, command execution, file management, and proxy support.

## Key Features

- **Memory Shell Injection**
 - Express middleware injection
 - Koa middleware injection
 - Prototype pollution techniques
 - Multiple encoding methods (Base64, XOR, AES)

- **Interactive Terminal**
 - Real-time command execution
 - Virtual terminal with xterm.js
 - Command history tracking
 - Multi-shell management

- **File Management**
 - File browser with directory navigation
 - Upload/download files with chunked transfer
 - File preview and editing
 - Monaco editor integration

- **Security Features**
 - Multiple encoding types support
 - Custom HTTP headers
 - Proxy support (HTTP/HTTPS/SOCKS5)
 - Password protection

- **Modern UI**
 - Vue 3 + TypeScript frontend
 - Element Plus components
 - Responsive design
 - Real-time status monitoring

## Screenshots

<div align="center">
<img src="docs/images/image-20251212115456-c1eetoy.png" width="45%" />
<img src="docs/images/image-20251212115520-tdfnt8f.png" width="45%" />
<img src="docs/images/image-20251212115541-ee6pqdg.png" width="45%" />
<img src="docs/images/image-20251212115610-lv6r139.png" width="45%" />
<img src="docs/images/image-20251212115646-80qohh6.png" width="45%" />
<img src="docs/images/image-20251212115655-4286u25.png" width="45%" />
<img src="docs/images/image-20251212115705-fuy8qqu.png" width="45%" />
<img src="docs/images/image-20251212115734-uif94za.png" width="45%" />
<img src="docs/images/image-20251212115805-83cnawp.png" width="45%" />
<img src="docs/images/image-20251212115820-dlh0rqy.png" width="45%" />
</div>

## Architecture

```
TL-NodeJsShell/
├── server/    # Go backend server
│ ├── cmd/    # Entry points
│ ├── internal/  # Internal logic
│ │ ├── app/   # Application core
│ │ ├── config/  # Configuration
│ │ ├── core/  # Core functionality
│ │ ├── database/ # Database models
│ │ └── handlers/ # HTTP handlers
│ └── go.mod
├── web/     # Vue.js frontend
│ ├── src/
│ │ ├── api/  # API clients
│ │ ├── components/ # Vue components
│ │ ├── router/  # Vue router
│ │ ├── stores/  # Pinia stores
│ │ ├── types/  # TypeScript types
│ │ └── views/  # Page views
│ └── public/   # Static assets
└── docs/    # Documentation assets
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 16 or higher
- npm or yarn

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/tianlusec/TL-NodeJsShell.git
cd TL-NodeJsShell
```

2. **Build and run backend**
```bash
cd server
go mod download
go build -o NodeJsshell cmd/api/main.go
./NodeJsshell
```

3. **Build and run frontend**
```bash
cd web
npm install
npm run build
```

4. **Access the application**
```
Open your browser and navigate to: http://localhost:8080
```

### Development Mode

**Backend:**
```bash
cd server
go run cmd/api/main.go
```

**Frontend:**
```bash
cd web
npm run dev
```

## Usage

### 1. Add a Shell

- Navigate to Shell Manager
- Click "Add Shell"
- Configure:
 - Target URL
 - Password
 - Encoding type (Base64/XOR/AES)
 - HTTP method (GET/POST)
 - Optional: Proxy settings
 - Optional: Custom headers

### 2. Manage Shells

- View all connected shells
- Test connection status
- Check system information
- Monitor latency

### 3. Execute Commands

- Select a shell
- Use the virtual terminal
- Execute commands in real-time
- View command history

### 4. File Operations

- Browse remote file system
- Upload files (with chunked transfer for large files)
- Download files
- Preview and edit files

### 5. Payload Generation

- Select template type
- Configure encoding
- Generate payload
- Copy and deploy

## Configuration

Backend configuration is located in `server/internal/config/config.go`:

```go
type Config struct {
 Port string // Default: "8080"
 Host string // Default: "0.0.0.0"
}
```

## Technology Stack

**Backend:**
- Go 1.21+
- Gin Web Framework
- GORM (SQLite)
- Gorilla WebSocket

**Frontend:**
- Vue 3
- TypeScript
- Element Plus
- Vite
- Pinia (State Management)
- Vue Router
- Axios
- xterm.js
- Monaco Editor

## Documentation

- [Installation Guide](docs/INSTALLATION.md)
- [API Documentation](docs/API.md)
- [Security Policy](docs/SECURITY.md)
- [Contributing Guide](.github/CONTRIBUTING.md)
- [Project Structure](PROJECT_STRUCTURE.md)

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](.github/CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to all contributors
- Inspired by various WebShell management tools
- Built with modern web technologies

## Contact

- GitHub: [@tianlusec](https://github.com/tianlusec)
- Issues: [GitHub Issues](https://github.com/tianlusec/TL-NodeJsShell/issues)

---

<div align="center">

** If this project helps you, please give it a star!**

</div>
