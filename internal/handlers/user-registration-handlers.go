package handlers

import (
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

	protected := api.Use(middleware.BearerAuthMiddleware(h.jwtSecret))
	{
		protected.POST("/user", h.AddUser)
	}
}

// RegisterTenant godoc
// @Summary Register a new tenant
// @Description Register a new tenant organization with an admin user
// @Tags Registration
// @Accept json
// @Produce json
// @Param tenantDTO body frameworkdto.TenantRegistrationDTO true "Tenant registration details"
// @Success 201 {object} frameworkdto.CreatedResponseDTO "Tenant registered successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "Invalid request body"
// @Failure 409 {object} frameworkdto.ErrorResponseDTO "Tenant already exists"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /registration/tenant [post]
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

	frameworkutils.CreatedResponse(c, nil, "Tenant registered successfully")
}

// AddUser godoc
// @Summary Add a new user to tenant
// @Description Add a new user to the authenticated tenant (requires authentication)
// @Tags Registration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userDTO body frameworkdto.UserRegistrationDTO true "User registration details"
// @Success 201 {object} frameworkdto.CreatedResponseDTO "User added successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "Invalid request body"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /registration/user [post]
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

	frameworkutils.CreatedResponse(c, nil, "User added successfully")
}
