// Package serviceframework provides a comprehensive service framework with multi-tenant support
package serviceframework

// @title Service Framework API
// @version 1.0
// @description A comprehensive service framework API with multi-tenant support, user management, authentication, and authorization.
// @description
// @description ## Authentication
// @description Most endpoints require Bearer token authentication. Include the token in the Authorization header:
// @description `Authorization: Bearer <your_jwt_token>`
// @description
// @description ## Error Responses
// @description All error responses follow a standard format with success, error code, and message fields.
//
// @contact.name API Support
// @contact.email support@geekible.com
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8000
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
//
// @tag.name Health
// @tag.description Health check endpoints
//
// @tag.name Authentication
// @tag.description User authentication and login
//
// @tag.name Registration
// @tag.description Tenant and user registration
//
// @tag.name User Maintenance
// @tag.description User management, password reset, and email verification
//
// @tag.name Tenant
// @tag.description Tenant management operations
//
// @tag.name Licence Type
// @tag.description Licence type management (Super Admin only)

import (
	_ "github.com/geekible-ltd/serviceframework/docs"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/geekible-ltd/serviceframework/internal/config"
	"github.com/geekible-ltd/serviceframework/internal/handlers"
	"github.com/geekible-ltd/serviceframework/internal/middleware"
	"github.com/geekible-ltd/serviceframework/internal/repositories"
	"github.com/geekible-ltd/serviceframework/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type ServiceFramework struct {
	cfg    *frameworkdto.FrameworkConfig
	fc     *config.FrameworkConfiguration
	db     *gorm.DB
	router *gin.Engine
}

func NewServiceFramework(cfg *frameworkdto.FrameworkConfig) *ServiceFramework {
	fc := config.NewFrameworkConfig(cfg)
	gormDb := fc.GetDatabase()

	return &ServiceFramework{
		cfg:    cfg,
		fc:     fc,
		db:     gormDb,
		router: fc.GetRouter(),
	}
}

func (s *ServiceFramework) GetDatabase() *gorm.DB {
	return s.db
}

func (s *ServiceFramework) GetRouter(requestPerSecond, burst int) *gin.Engine {
	if s.router == nil {
		panic("router is not initialized")
	}

	s.router.Use(middleware.CORSMiddleware(s.cfg.CORSCfg))
	s.router.Use(middleware.RateLimitMiddleware(requestPerSecond, burst))

	if s.cfg.Environment == frameworkdto.EnvDev {
		s.router.Use(gin.Logger())
	}

	handlers.NewHealthHandlers().RegisterRoutes(s.router)

	// Register Swagger/Redoc documentation routes

	if s.cfg.Environment == frameworkdto.EnvDev {
		s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		s.router.GET("/docs", func(c *gin.Context) {
			c.Redirect(302, "/swagger/index.html")
		})
	}
	// Redoc endpoint - provides alternative API documentation UI
	s.router.GET("/redoc", func(c *gin.Context) {
		html := `<!DOCTYPE html>
<html>
  <head>
    <title>API Documentation - ReDoc</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='/swagger/doc.json'></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
  </body>
</html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, html)
	})

	// Register Repos
	userRepo := repositories.NewUserRepository(s.db)
	tenantRepo := repositories.NewTenantRepository(s.db)
	tenantLicenceRepo := repositories.NewTenantLicenceRepository(s.db)
	licenceTypeRepo := repositories.NewLicenceTypeRepository(s.db)

	// Register Services
	loginService := services.NewLoginService(s.cfg, userRepo, tenantRepo)
	licenceTypeService := services.NewLicenceTypeService(licenceTypeRepo)
	registrationService := services.NewUserRegistrationService(userRepo, tenantRepo, tenantLicenceRepo, licenceTypeRepo)
	tenantService := services.NewTenantService(tenantRepo)
	userMaintenanceService := services.NewUserMaintenanceService(userRepo)

	// Register login handlers
	handlers.NewLoginHandlers(loginService).RegisterRoutes(s.router)
	handlers.NewLicenceTypeHandler(s.cfg.JWTSecret, licenceTypeService).RegisterRoutes(s.router)
	handlers.NewRegistrationHandlers(s.cfg.JWTSecret, registrationService).RegisterRoutes(s.router)
	handlers.NewUserMaintenanceHandler(s.cfg.JWTSecret, userMaintenanceService).RegisterRoutes(s.router)
	handlers.NewTenantHandler(s.cfg.JWTSecret, tenantService).RegisterRoutes(s.router)

	return s.router
}
