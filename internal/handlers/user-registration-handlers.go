package handlers

import (
	"net/http"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	frameworkutils "github.com/geekible-ltd/serviceframework/framework-utils"
	"github.com/geekible-ltd/serviceframework/internal/middleware"
	"github.com/geekible-ltd/serviceframework/internal/services"
	"github.com/gin-gonic/gin"
)

type RegistrationHandlers struct {
	jwtSecret           string
	registrationService *services.UserRegistrationService
}

func NewRegistrationHandlers(jwtSecret string, registrationService *services.UserRegistrationService) *RegistrationHandlers {
	return &RegistrationHandlers{
		jwtSecret:           jwtSecret,
		registrationService: registrationService,
	}
}

func (h *RegistrationHandlers) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/registration")
	api.POST("/tenant", h.RegisterTenant)

	api.Use(middleware.BearerAuthMiddleware(h.jwtSecret))
	{
		api.POST("/user", h.AddUser)
	}
}

func (h *RegistrationHandlers) RegisterTenant(c *gin.Context) {
	tenantDTO := frameworkdto.TenantRegistrationDTO{}
	if err := c.ShouldBindJSON(&tenantDTO); err != nil {
		frameworkutils.ErrorResponse(c, err)
		return
	}

	err := h.registrationService.RegisterTenant(tenantDTO)
	if err != nil {
		if err == frameworkconstants.ErrTenantAlreadyExists {
			frameworkutils.ErrorResponse(c, frameworkutils.Conflict(err.Error()))
			return
		}

		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, nil, "Tenant registered successfully")
}

func (h *RegistrationHandlers) AddUser(c *gin.Context) {
	userDTO := frameworkdto.UserRegistrationDTO{}
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		frameworkutils.ErrorResponse(c, err)
		return
	}

	tokenDTO, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, err)
		return
	}

	err = h.registrationService.RegisterUser(tokenDTO.TenantID, userDTO)
	if err != nil {
		frameworkutils.ErrorResponse(c, err)
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, nil, "User added successfully")
}
