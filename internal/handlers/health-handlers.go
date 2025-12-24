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

// GetHealth godoc
// @Summary Health check
// @Description Check if the API is running and healthy
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} frameworkdto.SuccessResponseDTO{data=string} "Service is healthy"
// @Router /api/health [get]
func (h *HealthHandlers) GetHealth(c *gin.Context) {
	frameworkutils.SuccessResponse(c, http.StatusOK, "OK", "Health check successful")
}
