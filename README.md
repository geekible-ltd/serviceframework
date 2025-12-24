# Service Framework

A comprehensive Go service framework with multi-tenant support, user management, authentication, and authorization. Built on top of Gin and GORM, this framework provides a solid foundation for building secure, scalable SaaS applications.

[![Go Version](https://img.shields.io/badge/Go-1.24.5+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 📋 Table of Contents

- [Features](#-features)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Configuration](#-configuration)
- [API Documentation](#-api-documentation)
- [Authentication](#-authentication)
- [Available Endpoints](#-available-endpoints)
- [Database Support](#-database-support)
- [Development](#-development)
- [Testing](#-testing)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

- **Multi-Tenant Architecture** - Complete tenant isolation with tenant-specific data
- **User Management** - User registration, authentication, and role-based access control
- **JWT Authentication** - Secure token-based authentication
- **Role-Based Authorization** - Support for multiple user roles (Super Admin, Tenant Admin, Tenant User)
- **Email Verification** - Email verification workflow for new users
- **Password Reset** - Secure password reset functionality
- **Licence Management** - Tenant licence types and management
- **CORS Support** - Configurable CORS policies
- **Rate Limiting** - Built-in rate limiting middleware
- **API Documentation** - Full OpenAPI/Swagger documentation with Redoc support
- **Health Checks** - Built-in health check endpoints
- **Multiple Database Support** - PostgreSQL, MySQL, and SQLite

## 🔧 Prerequisites

Before using this framework, ensure you have:

- **Go 1.24.5 or higher** installed
- A supported database (PostgreSQL, MySQL, or SQLite)
- Basic understanding of Go and REST APIs

## 📦 Installation

### 1. Install the package

```bash
go get github.com/geekible-ltd/serviceframework
```

### 2. Install required dependencies

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get github.com/golang-jwt/jwt/v5
```

### 3. (Optional) Install Swagger tools for API documentation

```bash
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

## 🚀 Quick Start

### Basic Setup

Create a `main.go` file in your project:

```go
package main

import (
    "log"
    "net/http"
    "time"

    "github.com/geekible-ltd/serviceframework"
    frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
)

func main() {
    // Configure the framework
    cfg := frameworkdto.FrameworkConfig{
        Environment: frameworkdto.EnvDev,
        JWTSecret:   "your-super-secret-jwt-key-change-this-in-production", // 256 bit
        DBType:      frameworkdto.DatabaseTypePostgreSQL,
        DbCfg: frameworkdto.DatabaseConfig{
            Host:     "localhost",
            Port:     5432,
            Username: "postgres",
            Password: "your-database-password",
            Database: "your-database-name",
            SSLMode:  "disable",
        },
        CORSCfg: frameworkdto.CORSCfg{
            AllowedOrigins: []string{"*"},
            AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
            AllowedHeaders: []string{"*"},
        },
    }

    // Initialize the framework
    sf := serviceframework.NewServiceFramework(&cfg)
    
    // Get the configured router with rate limiting (5 requests per second, burst of 10)
    router := sf.GetRouter(5, 10)

    // Configure HTTP server
    server := &http.Server{
        Addr:           ":8080",
        Handler:        router,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
        IdleTimeout:    120 * time.Second,
    }

    log.Printf("Server is running on port 8080")
    log.Printf("API Documentation: http://localhost:8080/swagger/index.html")
    log.Printf("Redoc: http://localhost:8080/redoc")

    if err := server.ListenAndServe(); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### Run Your Application

```bash
go run main.go
```

Your API will be available at `http://localhost:8080`

## ⚙️ Configuration

### Framework Configuration

The `FrameworkConfig` struct provides comprehensive configuration options:

```go
type FrameworkConfig struct {
    Environment Environment          // EnvDev, EnvStaging, or EnvProduction
    JWTSecret   string               // Secret key for JWT token signing
    DBType      DatabaseType         // DatabaseTypePostgreSQL, DatabaseTypeMySQL, or DatabaseTypeSQLite
    DbCfg       DatabaseConfig       // Database connection details
    CORSCfg     CORSCfg             // CORS configuration
}
```

### Environment Types

```go
const (
    EnvDev        Environment = "development"
    EnvStaging    Environment = "staging"
    EnvProduction Environment = "production"
)
```

### Database Configuration

#### PostgreSQL Example

```go
DbCfg: frameworkdto.DatabaseConfig{
    Host:     "localhost",
    Port:     5432,
    Username: "postgres",
    Password: "your-password",
    Database: "your-database",
    SSLMode:  "disable",
}
```

#### MySQL Example

```go
DBType: frameworkdto.DatabaseTypeMySQL,
DbCfg: frameworkdto.DatabaseConfig{
    Host:     "localhost",
    Port:     3306,
    Username: "root",
    Password: "your-password",
    Database: "your-database",
}
```

#### SQLite Example

```go
DBType: frameworkdto.DatabaseTypeSQLite,
DbCfg: frameworkdto.DatabaseConfig{
    Database: "./data/app.db",
}
```

### CORS Configuration

```go
CORSCfg: frameworkdto.CORSCfg{
    AllowedOrigins: []string{"https://yourdomain.com", "https://app.yourdomain.com"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders: []string{"Origin", "Content-Type", "Authorization"},
}
```

### Rate Limiting

Configure rate limiting when getting the router:

```go
// 10 requests per second with a burst of 20
router := sf.GetRouter(10, 20)
```

## 📚 API Documentation

The framework automatically generates comprehensive API documentation using OpenAPI/Swagger.

### Accessing Documentation

Once your application is running:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **Redoc**: `http://localhost:8080/redoc`
- **OpenAPI JSON**: `http://localhost:8080/swagger/doc.json`

### Generating/Updating Documentation

If you extend the framework, regenerate the documentation:

```bash
# Using Make (recommended)
make swagger

# Or manually
swag init -g service-framework.go -o ./docs --parseDependency --parseInternal
```

## 🔐 Authentication

The framework uses JWT (JSON Web Tokens) for authentication.

### Login Flow

1. **Register a Tenant** (POST `/registration/tenant`)
2. **Login** (POST `/authentication/login`)
3. **Receive JWT Token**
4. **Use Token in Subsequent Requests**

### Using the Token

Include the JWT token in the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

### Example: Login Request

```bash
curl -X POST http://localhost:8080/authentication/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "your-password"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "message": "Login successful"
}
```

### Example: Authenticated Request

```bash
curl -X GET http://localhost:8080/user-maintenance/users/get-all \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## 🛣️ Available Endpoints

### Health Check

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/health` | Health check | No |

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/authentication/login` | User login | No |

### Registration

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/registration/tenant` | Register new tenant | No |
| POST | `/registration/user` | Add user to tenant | Yes |

### User Maintenance

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| DELETE | `/user-maintenance/user?userId={id}` | Delete user | Yes (Admin) |
| PUT | `/user-maintenance/user` | Update user | Yes (Admin/Self) |
| GET | `/user-maintenance/users/get-all` | Get all tenant users | Yes (Admin) |
| GET | `/user-maintenance/users/get-roles` | Get available roles | Yes (Admin) |
| POST | `/user-maintenance/reset-password-request` | Request password reset | No |
| POST | `/user-maintenance/reset-password` | Reset password | No |
| POST | `/user-maintenance/verify-email` | Verify email | No |

### Tenant Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/tenant/get-by-id` | Get tenant by ID | Yes |
| GET | `/tenant/get-all` | Get all tenants | Yes (Super Admin) |
| PUT | `/tenant/update` | Update tenant | Yes (Admin) |
| DELETE | `/tenant/delete` | Delete tenant | Yes (Admin) |

### Licence Type Management

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/licence-type/get-all` | Get all licence types | Yes (Super Admin) |
| GET | `/licence-type/get-by-id?id={id}` | Get licence type by ID | Yes (Super Admin) |
| POST | `/licence-type/create` | Create licence type | Yes (Super Admin) |
| PUT | `/licence-type/update` | Update licence type | Yes (Super Admin) |
| DELETE | `/licence-type/delete?id={id}` | Delete licence type | Yes (Super Admin) |

## 🗄️ Database Support

The framework automatically handles database migrations and supports:

### PostgreSQL (Recommended)

```go
DBType: frameworkdto.DatabaseTypePostgreSQL,
DbCfg: frameworkdto.DatabaseConfig{
    Host:     "localhost",
    Port:     5432,
    Username: "postgres",
    Password: "password",
    Database: "myapp",
    SSLMode:  "require", // Use "disable" for local development
}
```

### MySQL

```go
DBType: frameworkdto.DatabaseTypeMySQL,
DbCfg: frameworkdto.DatabaseConfig{
    Host:     "localhost",
    Port:     3306,
    Username: "root",
    Password: "password",
    Database: "myapp",
}
```

### SQLite (Development Only)

```go
DBType: frameworkdto.DatabaseTypeSQLite,
DbCfg: frameworkdto.DatabaseConfig{
    Database: "./data/app.db",
}
```

### Database Entities

The framework automatically creates and manages these tables:

- `tenants` - Tenant organizations
- `users` - User accounts
- `tenant_licences` - Tenant licence assignments
- `licence_types` - Available licence types

## 🔨 Development

### Project Structure

```
service-framework/
├── docs/                          # Generated API documentation
├── example/                       # Example implementation
│   └── main.go
├── framework-constants/           # Constants and error messages
├── framework-dto/                 # Data Transfer Objects
├── framework-utils/               # Utility functions
├── internal/
│   ├── config/                   # Internal configuration
│   ├── entities/                 # Database entities
│   ├── handlers/                 # HTTP handlers
│   ├── middleware/               # Middleware components
│   ├── repositories/             # Data access layer
│   └── services/                 # Business logic
├── Makefile                      # Build and development tasks
├── README.md                     # This file
├── go.mod
├── go.sum
└── service-framework.go          # Main framework entry point
```

### Using the Makefile

The project includes a comprehensive Makefile for common development tasks:

```bash
# Show all available commands
make help

# Install Swagger CLI tool
make install-swag

# Install all dependencies
make install-deps

# Generate API documentation
make swagger

# Run the example application
make run

# Build the project
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make format

# Run linter (requires golangci-lint)
make lint

# Clean build artifacts
make clean

# Complete setup (install tools, deps, generate docs)
make setup

# Quick start (setup + generate docs)
make quickstart
```

### Debugging

A VS Code debug configuration is included. To debug the example application:

1. Open the project in VS Code
2. Press `F5` or go to Run → Start Debugging
3. Select "Launch Example" from the dropdown
4. Set breakpoints and debug

### Extending the Framework

#### Adding a New Handler

1. Create your handler in `internal/handlers/`:

```go
package handlers

import (
    "github.com/gin-gonic/gin"
    frameworkutils "github.com/geekible-ltd/serviceframework/framework-utils"
)

type MyHandler struct {
    myService *services.MyService
}

func NewMyHandler(myService *services.MyService) *MyHandler {
    return &MyHandler{myService: myService}
}

func (h *MyHandler) RegisterRoutes(router *gin.Engine) {
    api := router.Group("/my-endpoint")
    {
        api.GET("/action", h.MyAction)
    }
}

// MyAction godoc
// @Summary My custom action
// @Description Detailed description of what this does
// @Tags MyTag
// @Accept json
// @Produce json
// @Success 200 {object} frameworkdto.SuccessResponseDTO
// @Router /my-endpoint/action [get]
func (h *MyHandler) MyAction(c *gin.Context) {
    result := h.myService.DoSomething()
    frameworkutils.SuccessResponse(c, 200, result, "Action completed")
}
```

2. Register your handler in `service-framework.go`:

```go
myService := services.NewMyService()
handlers.NewMyHandler(myService).RegisterRoutes(s.router)
```

3. Regenerate documentation:

```bash
make swagger
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/handlers/...
```

### Example Test

```go
package handlers_test

import (
    "testing"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
)

func TestHealthHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    
    handler := handlers.NewHealthHandlers()
    handler.RegisterRoutes(router)
    
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/api/health", nil)
    router.ServeHTTP(w, req)
    
    if w.Code != 200 {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

## 🔒 Security Best Practices

1. **Never commit secrets**: Use environment variables for sensitive data
2. **Use strong JWT secrets**: Generate secure, random secrets (minimum 32 characters)
3. **Enable SSL in production**: Set `SSLMode: "require"` for PostgreSQL
4. **Implement rate limiting**: Configure appropriate limits for your use case
5. **Validate input**: Always validate and sanitize user input
6. **Use HTTPS**: Always use HTTPS in production
7. **Keep dependencies updated**: Regularly update Go modules

### Example: Environment Variables

```go
import "os"

cfg := frameworkdto.FrameworkConfig{
    Environment: frameworkdto.Environment(os.Getenv("APP_ENV")),
    JWTSecret:   os.Getenv("JWT_SECRET"),
    DBType:      frameworkdto.DatabaseTypePostgreSQL,
    DbCfg: frameworkdto.DatabaseConfig{
        Host:     os.Getenv("DB_HOST"),
        Port:     5432,
        Username: os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        Database: os.Getenv("DB_NAME"),
        SSLMode:  "require",
    },
}
```

## 📖 Response Format

All API responses follow a consistent format:

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data here
  },
  "message": "Operation completed successfully"
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERR_001",
    "message": "Error description",
    "details": {}
  }
}
```

### HTTP Status Codes

- `200 OK` - Successful GET/PUT request
- `201 Created` - Successful POST (creation)
- `202 Accepted` - Accepted for processing
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., duplicate)
- `500 Internal Server Error` - Server error

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow standard Go conventions
- Run `gofmt` before committing
- Add comments for exported functions
- Include Swagger annotations for API endpoints
- Write tests for new features

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For issues, questions, or contributions:

- 📧 Email: support@geekible.com
- 🐛 Issues: [GitHub Issues](https://github.com/geekible-ltd/serviceframework/issues)
- 📖 Documentation: [API Docs](http://localhost:8080/swagger/index.html)

## 🗺️ Roadmap

- [ ] WebSocket support
- [ ] GraphQL support
- [ ] OAuth2 integration
- [ ] Audit logging
- [ ] Advanced search and filtering
- [ ] Caching layer (Redis)
- [ ] Background job processing
- [ ] File upload/storage
- [ ] Email service integration
- [ ] SMS service integration

---

**Built with ❤️ by Geekible Ltd**

