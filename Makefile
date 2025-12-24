.PHONY: help install-swag swagger-generate swagger-clean run build test clean tidy lint format

# Default target - show help
help:
	@echo "Available targets:"
	@echo "  make install-swag      - Install swag CLI tool for generating Swagger docs"
	@echo "  make swagger-generate  - Generate Swagger documentation"
	@echo "  make swagger-clean     - Clean generated Swagger files"
	@echo "  make swagger           - Clean and regenerate Swagger docs"
	@echo "  make run              - Run the example application"
	@echo "  make build            - Build the project"
	@echo "  make test             - Run tests"
	@echo "  make tidy             - Run go mod tidy"
	@echo "  make format           - Format code with gofmt"
	@echo "  make lint             - Run golangci-lint (if installed)"
	@echo "  make clean            - Clean build artifacts and docs"

# Install swag CLI tool
install-swag:
	@echo "Installing swag CLI tool..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Swag installed successfully!"
	@echo "Make sure $(go env GOPATH)/bin is in your PATH"

# Install swagger dependencies
install-deps:
	@echo "Installing Swagger dependencies..."
	go get -u github.com/swaggo/swag/cmd/swag
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files
	@echo "Dependencies installed successfully!"

# Generate Swagger documentation
swagger-generate:
	@echo "Generating Swagger documentation..."
	swag init -g service-framework.go -o ./docs --parseDependency --parseInternal
	@echo "Fixing compatibility issues..."
	@sed -i.bak '/LeftDelim:/d' ./docs/docs.go && rm -f ./docs/docs.go.bak || true
	@sed -i.bak '/RightDelim:/d' ./docs/docs.go && rm -f ./docs/docs.go.bak || true
	@echo "Swagger documentation generated successfully!"
	@echo "Available at: http://localhost:8080/swagger/index.html"
	@echo "Redoc available at: http://localhost:8080/redoc"

# Clean Swagger generated files
swagger-clean:
	@echo "Cleaning Swagger documentation..."
	rm -rf ./docs/docs.go ./docs/swagger.json ./docs/swagger.yaml
	@echo "Swagger documentation cleaned!"

# Clean and regenerate Swagger docs
swagger: swagger-clean swagger-generate

# Run the example application
run:
	@echo "Running example application..."
	cd example && go run main.go

# Build the project
build:
	@echo "Building project..."
	go build -v ./...
	@echo "Build complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "Dependencies tidied!"

# Format code
format:
	@echo "Formatting code..."
	gofmt -w -s .
	@echo "Code formatted!"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf ./docs/docs.go ./docs/swagger.json ./docs/swagger.yaml
	rm -f coverage.out coverage.html
	go clean
	@echo "Clean complete!"

# Setup project (install tools and dependencies)
setup: install-swag install-deps tidy
	@echo "Project setup complete!"
	@echo "Run 'make swagger' to generate API documentation"

# Quick start - setup and generate docs
quickstart: setup swagger
	@echo "Quick start complete!"
	@echo "Run 'make run' to start the application"

