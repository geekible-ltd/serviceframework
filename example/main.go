package main

import (
	"log"
	"net/http"
	"time"

	"github.com/geekible-ltd/serviceframework"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
)

func main() {
	cfg := frameworkdto.FrameworkConfig{
		Environment: frameworkdto.EnvDev,
		JWTSecret:   "sjkskjd8fkjnd89w3rmnasdiofsdfsadf",
		DBType:      frameworkdto.DatabaseTypePostgreSQL,
		DbCfg: frameworkdto.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "SuperSecretPassword",
			Database: "mydatabase",
			SSLMode:  "disable",
		},
		CORSCfg: frameworkdto.CORSCfg{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowedHeaders: []string{"*"},
		},
	}

	sf := serviceframework.NewServiceFramework(&cfg)
	router := sf.GetRouter(5, 10)

	server := &http.Server{
		Addr:           ":8000",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    120 * time.Second,
	}

	log.Printf("Server is running on port %s", "8000")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
