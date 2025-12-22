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
	//tenantLicenceRepo := repositories.NewTenantLicenceRepository(s.db)

	// Register Services
	loginService := services.NewLoginService(s.cfg, userRepo, tenantRepo)

	// Register login handlers
	handlers.NewLoginHandlers(loginService).RegisterRoutes(s.router)

	return s.router
}
