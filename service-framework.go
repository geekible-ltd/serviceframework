package serviceframework

import (
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/geekible-ltd/serviceframework/internal/config"
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
		cfg: cfg,
		fc:  fc,
		db:  gormDb,
	}
}

func (s *ServiceFramework) GetDatabase() *gorm.DB {
	return s.db
}

func (s *ServiceFramework) GetRouter() *gin.Engine {
	return s.router
}
