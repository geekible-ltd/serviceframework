package handlers

import (
	"net/http"

	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	frameworkutils "github.com/geekible-ltd/serviceframework/framework-utils"
	"github.com/geekible-ltd/serviceframework/internal/services"
	"github.com/gin-gonic/gin"
)

type LoginHandlers struct {
	loginService *services.LoginService
}

func NewLoginHandlers(loginService *services.LoginService) *LoginHandlers {
	return &LoginHandlers{loginService: loginService}
}

func (h *LoginHandlers) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/authentication")
	api.POST("/login", h.Login)
}

func (h *LoginHandlers) Login(c *gin.Context) {
	var loginRequest frameworkdto.LoginDTO
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loginResponse, err := h.loginService.Login(loginRequest, c.ClientIP())
	if err != nil {
		frameworkutils.ErrorResponse(c, err)
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, loginResponse, "Login successful")
}
