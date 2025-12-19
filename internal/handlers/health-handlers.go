package handlers

import (
	"net/http"

	frameworkutils "github.com/geekible-ltd/serviceframework/framework-utils"
	"github.com/gin-gonic/gin"
)

type HealthHandlers struct {
}

func NewHealthHandlers() *HealthHandlers {
	return &HealthHandlers{}
}

func (h *HealthHandlers) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/health", h.GetHealth)
}

func (h *HealthHandlers) GetHealth(c *gin.Context) {
	frameworkutils.SuccessResponse(c, http.StatusOK, "OK", "Health check successful")
}
