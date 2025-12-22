package handlers

import (
	"net/http"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	frameworkutils "github.com/geekible-ltd/serviceframework/framework-utils"
	"github.com/geekible-ltd/serviceframework/internal/services"
	"github.com/gin-gonic/gin"
)

type RegistrationHandlers struct {
	registrationService *services.UserRegistrationService
}

func NewRegistrationHandlers(registrationService *services.UserRegistrationService) *RegistrationHandlers {
	return &RegistrationHandlers{registrationService: registrationService}
}

func (h *RegistrationHandlers) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/registration")
	api.POST("/tenant", h.RegisterTenant)
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
