# Contributing to Service Framework

Thank you for your interest in contributing to the Service Framework! This document provides guidelines and instructions for contributing.

## 🤝 How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs actual behavior
- **Environment details** (Go version, OS, database type)
- **Code samples** or test cases if applicable
- **Error messages** or logs

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Clear description** of the proposed feature
- **Use cases** and examples
- **Potential implementation approach** (optional)
- **Alternative solutions** you've considered

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Follow the coding standards** outlined below
3. **Add tests** for new functionality
4. **Update documentation** including README and API docs
5. **Ensure all tests pass** before submitting
6. **Update the CHANGELOG** (if applicable)

## 💻 Development Setup

### Prerequisites

- Go 1.24.5 or higher
- Git
- PostgreSQL, MySQL, or SQLite
- golangci-lint (for linting)

### Setup Steps

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/serviceframework.git
cd serviceframework

# Install dependencies
make install-deps

# Install development tools
make install-swag

# Run tests
make test

# Generate API documentation
make swagger
```

## 📋 Coding Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Use `golint` and `go vet` for code quality
- Keep functions small and focused
- Write clear, descriptive variable names

### Code Formatting

Before committing:

```bash
# Format code
make format

# Run linter
make lint
```

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `handlers`, `services`)
- **Files**: lowercase with hyphens (e.g., `user-handler.go`)
- **Types**: PascalCase (e.g., `UserHandler`, `LoginService`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Constants**: PascalCase or SCREAMING_SNAKE_CASE for enums

### Comments and Documentation

- Add comments for all exported functions, types, and constants
- Use godoc-style comments
- Include Swagger/OpenAPI annotations for all API endpoints

Example:

```go
// UserHandler handles user-related HTTP requests.
type UserHandler struct {
    userService *services.UserService
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieves user information by user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} frameworkdto.SuccessResponseDTO
// @Failure 404 {object} frameworkdto.ErrorResponseDTO
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    // Implementation
}
```

## 🧪 Testing

### Writing Tests

- Write tests for all new functionality
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for >80% code coverage

### Test Structure

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   UserDTO
        want    *User
        wantErr bool
    }{
        {
            name: "valid user",
            input: UserDTO{
                Email: "test@example.com",
                Name:  "Test User",
            },
            want: &User{
                Email: "test@example.com",
                Name:  "Test User",
            },
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test -v ./internal/services/...
```

## 📝 Documentation

### README Updates

When adding new features, update the README to include:

- Feature description
- Configuration options
- Usage examples
- API endpoints (if applicable)

### API Documentation

When adding or modifying API endpoints:

1. Add Swagger annotations to the handler function
2. Regenerate documentation: `make swagger`
3. Verify the documentation in Swagger UI
4. Update the endpoint table in README if necessary

### Code Comments

- Comment complex logic
- Explain "why" not just "what"
- Keep comments up-to-date with code changes

## 🔄 Git Workflow

### Branch Naming

- `feature/description` - New features
- `bugfix/description` - Bug fixes
- `hotfix/description` - Urgent fixes
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:

```bash
feat(auth): add OAuth2 support

Add OAuth2 authentication provider support with Google and GitHub

Closes #123

fix(handlers): resolve null pointer in user update

The user update handler was not checking for nil values before
accessing nested properties.

Fixes #456

docs(readme): update installation instructions

Add missing step for database migration
```

### Pull Request Process

1. **Update your fork** with the latest changes from main:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Create a pull request** with:
   - Clear title and description
   - Reference to related issues
   - List of changes made
   - Screenshots (if UI changes)
   - Test results

3. **Respond to feedback** from reviewers promptly

4. **Squash commits** if requested before merging

## 🏗️ Architecture Guidelines

### Layered Architecture

The framework follows a layered architecture:

```
Handlers (HTTP Layer)
    ↓
Services (Business Logic)
    ↓
Repositories (Data Access)
    ↓
Database
```

### Dependency Injection

- Use constructor injection
- Keep dependencies explicit
- Avoid global state

Example:

```go
type UserService struct {
    userRepo   repositories.UserRepository
    emailService services.EmailService
}

func NewUserService(
    userRepo repositories.UserRepository,
    emailService services.EmailService,
) *UserService {
    return &UserService{
        userRepo:   userRepo,
        emailService: emailService,
    }
}
```

### Error Handling

- Use custom error types for domain errors
- Wrap errors with context
- Return errors, don't panic
- Handle errors at appropriate levels

Example:

```go
if err := h.userService.CreateUser(dto); err != nil {
    if errors.Is(err, ErrUserAlreadyExists) {
        frameworkutils.ErrorResponse(c, frameworkutils.Conflict("User already exists"))
        return
    }
    frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
    return
}
```

## 📋 Checklist

Before submitting a pull request, ensure:

- [ ] Code follows project style guidelines
- [ ] Code is properly formatted (`make format`)
- [ ] All tests pass (`make test`)
- [ ] New tests added for new functionality
- [ ] Documentation updated (README, API docs, comments)
- [ ] No unnecessary dependencies added
- [ ] Commit messages follow conventional commits
- [ ] PR description is clear and complete
- [ ] Related issues are referenced
- [ ] Code has been self-reviewed
- [ ] No debug code or console.logs left

## 🎯 Areas for Contribution

We particularly welcome contributions in these areas:

- **Testing**: Improve test coverage
- **Documentation**: Improve examples and guides
- **Performance**: Optimization and benchmarking
- **Security**: Security audits and improvements
- **Features**: New features from the roadmap
- **Bug Fixes**: Fix reported issues
- **Refactoring**: Code quality improvements

## 💬 Community

- Be respectful and constructive
- Help others in issues and discussions
- Share your use cases and experiences
- Provide constructive feedback on PRs

## 📧 Contact

For questions or discussions:

- Email: support@geekible.com
- GitHub Issues: For bug reports and feature requests
- Pull Requests: For code contributions

Thank you for contributing to Service Framework! 🎉

