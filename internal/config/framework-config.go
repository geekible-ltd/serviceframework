package config

import (
	"fmt"

	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/geekible-ltd/serviceframework/internal/entities"
	"github.com/geekible-ltd/serviceframework/internal/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FrameworkConfiguration struct {
	db     *gorm.DB
	router *gin.Engine
}

func NewFrameworkConfig(cfg *frameworkdto.FrameworkConfig) *FrameworkConfiguration {
	var fc FrameworkConfiguration

	switch cfg.DBType {
	case frameworkdto.DatabaseTypeMySQL:
		fc.db = connectToMySQL(cfg)
	case frameworkdto.DatabaseTypePostgreSQL:
		fc.db = connectToPostgreSQL(cfg)
	case frameworkdto.DatabaseTypeSQLite:
		fc.db = connectToSQLite(cfg)
	}

	err := fc.db.AutoMigrate(&entities.Tenant{}, &entities.User{}, &entities.TenantLicence{}, &entities.LicenceType{})
	if err != nil {
		panic(err)
	}

	licenceTypeRepo := repositories.NewLicenceTypeRepository(fc.db)
	licenceTypeRepo.Create(entities.LicenceType{
		Name:        "Free",
		Description: "Free licence type",
		MaxSeats:    1,
	}, true)

	fc.router = buildGinEngine()

	return &fc
}

func connectToMySQL(cfg *frameworkdto.FrameworkConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DbCfg.Username,
		cfg.DbCfg.Password,
		cfg.DbCfg.Host,
		cfg.DbCfg.Port)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sql := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		cfg.DbCfg.Database,
	)

	err = db.Exec(sql).Error
	if err != nil {
		panic(err)
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DbCfg.Username,
		cfg.DbCfg.Password,
		cfg.DbCfg.Host,
		cfg.DbCfg.Port,
		cfg.DbCfg.Database)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func connectToPostgreSQL(cfg *frameworkdto.FrameworkConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		cfg.DbCfg.Host,
		cfg.DbCfg.Port,
		cfg.DbCfg.Username,
		cfg.DbCfg.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)`
	if err := db.Raw(query, cfg.DbCfg.Database).Scan(&exists).Error; err != nil {
		panic(err)
	}

	if !exists {
		query = `CREATE DATABASE ` + cfg.DbCfg.Database
		if err := db.Exec(query).Error; err != nil {
			panic(err)
		}
	}

	dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DbCfg.Host,
		cfg.DbCfg.Port,
		cfg.DbCfg.Username,
		cfg.DbCfg.Password,
		cfg.DbCfg.Database,
		cfg.DbCfg.SSLMode)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func connectToSQLite(cfg *frameworkdto.FrameworkConfig) *gorm.DB {
	dbName := fmt.Sprintf("%s.db", cfg.DbCfg.Database)
	dbPath := fmt.Sprintf("%s/%s", cfg.DbCfg.Database, dbName)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func buildGinEngine() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	return router
}

func (c *FrameworkConfiguration) GetDatabase() *gorm.DB {
	if c.db == nil {
		panic("database connection is not established")
	}
	return c.db
}

func (c *FrameworkConfiguration) GetRouter() *gin.Engine {
	if c.router == nil {
		panic("router is not initialized")
	}
	return c.router
}
