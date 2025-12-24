# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Full OpenAPI/Swagger documentation with Redoc support
- Comprehensive README with detailed usage instructions
- Contributing guidelines
- MIT License
- Example environment configuration file
- Makefile with common development tasks
- VS Code debug configuration

### Changed
- Improved error handling across all handlers
- Enhanced response format consistency

### Fixed
- Query parameter handling in user deletion endpoint
- Swagger documentation generation compatibility

## [1.0.0] - 2025-01-XX

### Added
- Multi-tenant architecture with complete tenant isolation
- JWT-based authentication system
- User management with role-based access control (Super Admin, Tenant Admin, Tenant User)
- Tenant registration and management
- User registration and management
- Email verification workflow
- Password reset functionality
- Licence type management
- Health check endpoints
- CORS middleware with configurable policies
- Rate limiting middleware
- PostgreSQL, MySQL, and SQLite database support
- Automatic database migrations
- Bearer token authentication middleware

### Features

#### Authentication & Authorization
- JWT token generation and validation
- Role-based access control
- Bearer token middleware
- Secure password hashing

#### User Management
- User CRUD operations
- Email verification
- Password reset flow
- User role management
- Multi-tenant user isolation

#### Tenant Management
- Tenant registration
- Tenant CRUD operations
- Licence type assignment
- Tenant-specific data isolation

#### API Features
- RESTful API design
- Consistent response format
- Comprehensive error handling
- Input validation
- Request rate limiting

#### Developer Experience
- Swagger/OpenAPI documentation
- Redoc alternative UI
- Health check endpoints
- Example implementation
- Development Makefile
- Debug configuration

### Security
- JWT-based authentication
- Password hashing with bcrypt
- Role-based authorization
- CORS protection
- Rate limiting
- SQL injection protection via GORM
- Input validation

## Version History

### Version Numbering

This project uses semantic versioning:
- **MAJOR** version for incompatible API changes
- **MINOR** version for added functionality in a backward compatible manner
- **PATCH** version for backward compatible bug fixes

### Upgrade Guidelines

When upgrading between versions, please refer to the specific version notes above for:
- Breaking changes
- New features
- Deprecated features
- Migration steps (if any)

---

For more details on any release, please see the [GitHub Releases](https://github.com/geekible-ltd/serviceframework/releases) page.

