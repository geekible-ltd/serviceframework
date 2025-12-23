package serviceframework

import (
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/geekible-ltd/serviceframework/internal/config"
	"github.com/geekible-ltd/serviceframework/internal/handlers"
	"github.com/geekible-ltd/serviceframework/internal/middleware"
	"github.com/geekible-ltd/serviceframework/internal/repositories"
	"github.com/geekible-ltd/serviceframework/internal/services"
	"github.com/gin-gonic/gin"
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
