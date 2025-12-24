package handlers

import (
	"net/http"
	"strconv"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	frameworkutils "github.com/geekible-ltd/serviceframework/framework-utils"
	"github.com/geekible-ltd/serviceframework/internal/middleware"
	"github.com/geekible-ltd/serviceframework/internal/services"
	"github.com/gin-gonic/gin"
)

type LicenceTypeHandler struct {
	jwtSecret          string
	licenceTypeService *services.LicenceTypeService
}

func NewLicenceTypeHandler(jwtSecret string, licenceTypeService *services.LicenceTypeService) *LicenceTypeHandler {
	return &LicenceTypeHandler{jwtSecret: jwtSecret, licenceTypeService: licenceTypeService}
}

func (h *LicenceTypeHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/licence-type")
	protected := api.Use(middleware.BearerAuthMiddleware(h.jwtSecret))
	{
		protected.GET("/get-all", h.GetAll)
		protected.GET("/get-by-id", h.GetById)
		protected.POST("/create", h.Create)
		protected.PUT("/update", h.Update)
		protected.DELETE("/delete", h.Delete)
	}
}

// GetAll godoc
// @Summary Get all licence types
// @Description Get list of all licence types (requires authentication, super admin only)
// @Tags Licence Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} frameworkdto.SuccessResponseDTO{data=[]frameworkdto.GetLicenceTypeDTO} "Licence types fetched successfully"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Forbidden"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /licence-type/get-all [get]
func (h *LicenceTypeHandler) GetAll(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role == string(frameworkconstants.UserRoleTenantAdmin) || tokenDto.Role == string(frameworkconstants.UserRoleTenantUser) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("Forbidden"))
		return
	}

	licenceTypes, err := h.licenceTypeService.GetAll()
	if err != nil {
		frameworkutils.ErrorResponse(c, err)
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, licenceTypes, "Licence types fetched successfully")
}

// GetById godoc
// @Summary Get licence type by ID
// @Description Get a specific licence type by ID (requires authentication, super admin only)
// @Tags Licence Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Licence Type ID"
// @Success 200 {object} frameworkdto.SuccessResponseDTO{data=frameworkdto.GetLicenceTypeDTO} "Licence type fetched successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "ID is required or Invalid ID"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Forbidden"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /licence-type/get-by-id [get]
func (h *LicenceTypeHandler) GetById(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role == string(frameworkconstants.UserRoleTenantAdmin) || tokenDto.Role == string(frameworkconstants.UserRoleTenantUser) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("Forbidden"))
		return
	}

	id := c.Query("id")
	if id == "" {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("ID is required"))
		return
	}

	licenceTypeId, err := strconv.Atoi(id)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid ID"))
		return
	}

	licenceType, err := h.licenceTypeService.GetByID(uint(licenceTypeId))
	if err != nil {
		frameworkutils.ErrorResponse(c, err)
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, licenceType, "Licence type fetched successfully")
}

// Create godoc
// @Summary Create a new licence type
// @Description Create a new licence type (requires authentication, super admin only)
// @Tags Licence Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param createLicenceTypeDTO body frameworkdto.LicenceTypeCreateRequestDTO true "Licence type creation details"
// @Success 201 {object} frameworkdto.CreatedResponseDTO "Licence type created successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "Invalid request body"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Forbidden"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /licence-type/create [post]
func (h *LicenceTypeHandler) Create(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role != string(frameworkconstants.UserRoleSuperAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("Forbidden"))
		return
	}

	var createLicenceTypeDTO frameworkdto.LicenceTypeCreateRequestDTO
	if err := c.ShouldBindJSON(&createLicenceTypeDTO); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	err = h.licenceTypeService.Create(createLicenceTypeDTO)
	if err != nil {
		frameworkutils.ErrorResponse(c, err)
		return
	}

	frameworkutils.CreatedResponse(c, nil, "Licence type created successfully")
}

// Update godoc
// @Summary Update a licence type
// @Description Update an existing licence type (requires authentication, super admin only)
// @Tags Licence Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dto body frameworkdto.LicenceTypeUpdateRequestDTO true "Licence type update details"
// @Success 202 {object} frameworkdto.SuccessResponseDTO "Licence type updated successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "Invalid request body"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Forbidden"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /licence-type/update [put]
func (h *LicenceTypeHandler) Update(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role != string(frameworkconstants.UserRoleSuperAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("Forbidden"))
		return
	}

	var dto frameworkdto.LicenceTypeUpdateRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	err = h.licenceTypeService.Update(dto)
	if err != nil {
		frameworkutils.ErrorResponse(c, err)
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "Licence type updated successfully")
}

// Delete godoc
// @Summary Delete a licence type
// @Description Delete a licence type by ID (requires authentication, super admin only)
// @Tags Licence Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Licence Type ID"
// @Success 202 {object} frameworkdto.SuccessResponseDTO "Licence type deleted successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "ID is required or Invalid ID"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Forbidden"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /licence-type/delete [delete]
func (h *LicenceTypeHandler) Delete(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role != string(frameworkconstants.UserRoleSuperAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("Forbidden"))
		return
	}

	id := c.Query("id")
	if id == "" {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("ID is required"))
		return
	}

	licenceTypeId, err := strconv.Atoi(id)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid ID"))
		return
	}

	err = h.licenceTypeService.Delete(uint(licenceTypeId))
	if err != nil {
		frameworkutils.ErrorResponse(c, err)
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "Licence type deleted successfully")
}
