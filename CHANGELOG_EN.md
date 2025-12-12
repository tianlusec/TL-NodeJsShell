# Changelog

English | [简体中文](CHANGELOG.md)

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- WebSocket support for real-time communication
- Plugin system for extensibility
- Multi-user support with authentication
- Enhanced logging and audit trails
- Docker deployment support

## [1.0.0] - 2024-12-12

### Added
- Initial release of TL-NodeJsShell
- Memory shell injection capabilities
  - Express middleware injection
  - Koa middleware injection
  - Prototype pollution techniques
- Multiple encoding support (Base64, XOR, AES)
- Interactive virtual terminal with xterm.js
- Comprehensive file management system
  - File browser with directory navigation
  - Upload/download with chunked transfer
  - File preview and editing with Monaco editor
- Proxy support (HTTP/HTTPS/SOCKS5)
- Custom HTTP headers configuration
- Real-time shell status monitoring
- System information collection
- Command history tracking
- Modern Vue 3 + TypeScript frontend
- Go backend with Gin framework
- SQLite database for data persistence
- RESTful API design
- Responsive UI with Element Plus

### Security
- Password protection for shells
- Multiple encoding methods for payload obfuscation
- Proxy support for anonymity
- Input validation and sanitization

### Documentation
- Comprehensive README with English and Chinese versions
- API documentation
- Contributing guidelines
- Code of conduct
- MIT License

---

[Unreleased]: https://github.com/tianlusec/TL-NodeJsShell/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/tianlusec/TL-NodeJsShell/releases/tag/v1.0.0