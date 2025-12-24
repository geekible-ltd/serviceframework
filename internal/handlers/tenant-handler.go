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

type TenantHandler struct {
	jwtSecret     string
	tenantService *services.TenantService
}

func NewTenantHandler(jwtSecret string, tenantService *services.TenantService) *TenantHandler {
	return &TenantHandler{jwtSecret: jwtSecret, tenantService: tenantService}
}

func (h *TenantHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/tenant")
	protected := api.Use(middleware.BearerAuthMiddleware(h.jwtSecret))
	{
		protected.GET("/get-by-id", h.GetTenantByID)
		protected.GET("/get-all", h.GetAllTenants)
		protected.PUT("/update", h.UpdateTenant)
		protected.DELETE("/delete", h.DeleteTenant)
	}
}

// GetTenantByID godoc
// @Summary Get tenant by ID
// @Description Get tenant details by ID (requires authentication)
// @Tags Tenant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} frameworkdto.SuccessResponseDTO{data=frameworkdto.GetTenantDTO} "Tenant fetched successfully"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Not authorized to get this resource"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /tenant/get-by-id [get]
func (h *TenantHandler) GetTenantByID(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role == string(frameworkconstants.UserRoleTenantAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to get this resource"))
		return
	}

	tenant, err := h.tenantService.GetTenantByID(tokenDto.TenantID)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, tenant, "Tenant fetched successfully")
}

// GetAllTenants godoc
// @Summary Get all tenants
// @Description Get list of all tenants (requires authentication, super admin only)
// @Tags Tenant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} frameworkdto.SuccessResponseDTO{data=[]frameworkdto.GetTenantDTO} "Tenants fetched successfully"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Not authorized to get this resource"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /tenant/get-all [get]
func (h *TenantHandler) GetAllTenants(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role == string(frameworkconstants.UserRoleSuperAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to get this resource"))
		return
	}

	tenants, err := h.tenantService.GetAllTenants()
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, tenants, "Tenants fetched successfully")
}

// UpdateTenant godoc
// @Summary Update tenant
// @Description Update tenant details (requires authentication, authorized roles only)
// @Tags Tenant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param updateTenantDTO body frameworkdto.UpdateTenantDTO true "Tenant update details"
// @Success 202 {object} frameworkdto.SuccessResponseDTO "Tenant updated successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "Invalid request body"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Not authorized to update tenant"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /tenant/update [put]
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role == string(frameworkconstants.UserRoleTenantAdmin) || tokenDto.Role == string(frameworkconstants.UserRoleSuperAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to get this resource"))
		return
	}

	var updateTenantDTO frameworkdto.UpdateTenantDTO
	if err := c.ShouldBindJSON(&updateTenantDTO); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	err = h.tenantService.UpdateTenant(tokenDto.TenantID, updateTenantDTO)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "Tenant updated successfully")
}

// DeleteTenant godoc
// @Summary Delete tenant
// @Description Delete a tenant (requires authentication, authorized roles only)
// @Tags Tenant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 202 {object} frameworkdto.SuccessResponseDTO "Tenant deleted successfully"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Not authorized to delete tenant"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /tenant/delete [delete]
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role == string(frameworkconstants.UserRoleTenantAdmin) || tokenDto.Role == string(frameworkconstants.UserRoleSuperAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to get this resource"))
		return
	}

	err = h.tenantService.DeleteTenant(tokenDto.TenantID)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "Tenant deleted successfully")
}
