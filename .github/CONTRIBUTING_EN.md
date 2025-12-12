# Contributing to TL-NodeJsShell

English | [简体中文](CONTRIBUTING.md)

---

Thank you for your interest in contributing to TL-NodeJsShell! This document provides guidelines for contributing to the project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates.

**When reporting a bug, include:**
- Clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Screenshots (if applicable)
- Environment details (OS, Go version, Node.js version)
- Error messages or logs

### Suggesting Enhancements

Enhancement suggestions are welcome! Please provide:
- Clear description of the enhancement
- Use cases and benefits
- Possible implementation approach
- Any relevant examples or mockups

### Pull Requests

1. **Fork the repository**
 ```bash
 git clone https://github.com/YOUR_USERNAME/TL-NodeJsShell.git
 cd TL-NodeJsShell
 ```

2. **Create a feature branch**
 ```bash
 git checkout -b feature/your-feature-name
 ```

3. **Make your changes**
 - Follow the coding standards
 - Write clear commit messages
 - Add tests if applicable
 - Update documentation

4. **Test your changes**
 ```bash
 # Backend tests
 cd backend
 go test ./...
 
 # Frontend tests
 cd frontend
 npm run test
 ```

5. **Commit your changes**
 ```bash
 git add .
 git commit -m "feat: add your feature description"
 ```

6. **Push to your fork**
 ```bash
 git push origin feature/your-feature-name
 ```

7. **Create a Pull Request**
 - Provide a clear description of changes
 - Reference any related issues
 - Ensure all tests pass

## Coding Standards

### Go (Backend)

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for code formatting
- Write meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and concise

**Example:**
```go
// GetShellByID retrieves a shell configuration by its ID
func (h *ShellHandler) GetShellByID(id uint) (*database.Shell, error) {
 var shell database.Shell
 if err := h.db.First(&shell, id).Error; err != nil {
  return nil, fmt.Errorf("shell not found: %w", err)
 }
 return &shell, nil
}
```

### TypeScript/Vue (Frontend)

- Follow [Vue.js Style Guide](https://vuejs.org/style-guide/)
- Use TypeScript for type safety
- Use composition API for Vue components
- Follow ESLint rules
- Write self-documenting code

**Example:**
```typescript
// Define clear interfaces
interface ShellConfig {
 id: number
 url: string
 password: string
 encodeType: string
}

// Use composition API
const { shells, loading, error } = useShells()
```

### Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, etc.)
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `chore:` Maintenance tasks

**Examples:**
```
feat: add proxy authentication support
fix: resolve file upload chunking issue
docs: update API documentation
refactor: improve error handling in shell handler
```

## Project Structure

```
TL-NodeJsShell/
├── backend/
│ ├── app/    # Application initialization
│ ├── config/   # Configuration management
│ ├── core/    # Core business logic
│ ├── database/   # Database models and operations
│ ├── handlers/   # HTTP request handlers
│ └── main.go   # Entry point
├── frontend/
│ ├── src/
│ │ ├── api/   # API client functions
│ │ ├── components/ # Reusable Vue components
│ │ ├── views/  # Page components
│ │ ├── stores/  # Pinia state management
│ │ └── router/  # Vue Router configuration
│ └── public/   # Static assets
└── docs/    # Documentation
```

## Development Setup

### Backend Development

```bash
cd backend
go mod download
go run main.go
```

### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

## Testing

### Backend Tests

```bash
cd backend
go test ./... -v
go test -cover ./...
```

### Frontend Tests

```bash
cd frontend
npm run test
npm run test:coverage
```

## Documentation

- Update README.md for user-facing changes
- Update API documentation in docs/API.md
- Add inline comments for complex code
- Update CHANGELOG.md for notable changes

## Security

- **Never commit sensitive data** (passwords, keys, tokens)
- Report security vulnerabilities privately
- Follow secure coding practices
- Validate all user inputs

## Questions?

Feel free to:
- Open an issue for discussion
- Join our community discussions
- Contact maintainers directly

---

<div align="center">

**Thank you for contributing!**

</div>