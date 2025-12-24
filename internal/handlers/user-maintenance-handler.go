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

type UserMaintenanceHandler struct {
	jwtSecret              string
	userMaintenanceService *services.UserMaintenanceService
}

func NewUserMaintenanceHandler(jwtSecret string, userMaintenanceService *services.UserMaintenanceService) *UserMaintenanceHandler {
	return &UserMaintenanceHandler{jwtSecret: jwtSecret, userMaintenanceService: userMaintenanceService}
}

func (h *UserMaintenanceHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/user-maintenance")

	api.POST("/reset-password-request", h.ResetPasswordRequest)
	api.POST("/reset-password", h.ResetPassword)
	api.POST("/verify-email", h.VerifyEmail)

	protected := api.Use(middleware.BearerAuthMiddleware(h.jwtSecret))
	{
		protected.DELETE("/user", h.DeleteUser)
		protected.PUT("/user", h.UpdateUser)
		protected.GET("/users/get-all", h.GetAllUsers)
		protected.GET("/users/get-roles", h.GetUserRoles)
	}
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user from the tenant (requires authentication, tenant admin only)
// @Tags User Maintenance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId query int true "User ID to delete"
// @Success 202 {object} frameworkdto.SuccessResponseDTO "User deleted successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "Invalid User ID format or User ID is required"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "You cannot delete yourself or not authorized"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /user-maintenance/user [delete]
func (h *UserMaintenanceHandler) DeleteUser(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	userID, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid User ID format"))
		return
	}

	if userID <= 0 {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("User ID is required"))
		return
	}

	currentUserID, err := strconv.Atoi(tokenDto.Sub)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid User ID format"))
		return
	}

	if userID == currentUserID {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You cannot delete yourself"))
		return
	}

	if tokenDto.Role != string(frameworkconstants.UserRoleTenantUser) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to delete this user"))
		return
	}

	err = h.userMaintenanceService.DeleteUser(tokenDto.TenantID, uint(userID))
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "User deleted successfully")
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update user details (requires authentication, tenant admin or self)
// @Tags User Maintenance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param updateUserDTO body frameworkdto.UserUpdateRequestDTO true "User update details"
// @Success 202 {object} frameworkdto.SuccessResponseDTO "User updated successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "Invalid request body or User ID format"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Not authorized to update this user"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /user-maintenance/user [put]
func (h *UserMaintenanceHandler) UpdateUser(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	var updateUserDTO frameworkdto.UserUpdateRequestDTO
	if err := c.ShouldBindJSON(&updateUserDTO); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	currentUserID, err := strconv.Atoi(tokenDto.Sub)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid User ID format"))
		return
	}

	if updateUserDTO.UserID != uint(currentUserID) || tokenDto.Role != string(frameworkconstants.UserRoleTenantAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to update this user"))
		return
	}

	err = h.userMaintenanceService.UpdateUser(tokenDto.TenantID, updateUserDTO.UserID, updateUserDTO)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "User updated successfully")
}

// ResetPasswordRequest godoc
// @Summary Request password reset
// @Description Request a password reset token to be sent to the user's email
// @Tags User Maintenance
// @Accept json
// @Produce json
// @Param resetPasswordRequestDTO body frameworkdto.ResetPasswordRequestDTO true "Email address"
// @Success 202 {object} frameworkdto.SuccessResponseDTO "Reset password request sent successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "Invalid request body"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /user-maintenance/reset-password-request [post]
func (h *UserMaintenanceHandler) ResetPasswordRequest(c *gin.Context) {
	var resetPasswordRequestDTO frameworkdto.ResetPasswordRequestDTO
	if err := c.ShouldBindJSON(&resetPasswordRequestDTO); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	err := h.userMaintenanceService.SetResetPasswordToken(resetPasswordRequestDTO.Email)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "Reset password request sent successfully")
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user password using the reset token received via email
// @Tags User Maintenance
// @Accept json
// @Produce json
// @Param resetPasswordDTO body frameworkdto.ResetPasswordDTO true "Reset token and new password"
// @Success 202 {object} frameworkdto.SuccessResponseDTO "Password reset successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "Invalid request body"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /user-maintenance/reset-password [post]
func (h *UserMaintenanceHandler) ResetPassword(c *gin.Context) {
	var resetPasswordDTO frameworkdto.ResetPasswordDTO
	if err := c.ShouldBindJSON(&resetPasswordDTO); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	err := h.userMaintenanceService.UpdateUserPassword(resetPasswordDTO.ResetToken, resetPasswordDTO.NewPassword)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "Password reset successfully")
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify user email address using the verification token sent via email
// @Tags User Maintenance
// @Accept json
// @Produce json
// @Param verifyEmailDTO body frameworkdto.VerifyEmailDTO true "Email verification details"
// @Success 202 {object} frameworkdto.SuccessResponseDTO "Email verified successfully"
// @Failure 400 {object} frameworkdto.ErrorResponseDTO "Invalid request body"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /user-maintenance/verify-email [post]
func (h *UserMaintenanceHandler) VerifyEmail(c *gin.Context) {
	var verifyEmailDTO frameworkdto.VerifyEmailDTO
	if err := c.ShouldBindJSON(&verifyEmailDTO); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	err := h.userMaintenanceService.VerifyEmail(verifyEmailDTO.TenantID, verifyEmailDTO.UserID, verifyEmailDTO.Token)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "Email verified successfully")
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users for the authenticated tenant (requires authentication, tenant admin only)
// @Tags User Maintenance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} frameworkdto.SuccessResponseDTO{data=[]frameworkdto.GetUsersResponseDTO} "Users fetched successfully"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Not authorized to get all users"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /user-maintenance/users/get-all [get]
func (h *UserMaintenanceHandler) GetAllUsers(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role != string(frameworkconstants.UserRoleTenantAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to get all users"))
		return
	}

	users, err := h.userMaintenanceService.GetAllUsersByTenantID(tokenDto.TenantID)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, users, "Users fetched successfully")
}

// GetUserRoles godoc
// @Summary Get all user roles
// @Description Get list of all available user roles (requires authentication, tenant admin only)
// @Tags User Maintenance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} frameworkdto.SuccessResponseDTO{data=[]frameworkdto.GetUserRoles} "User roles fetched successfully"
// @Failure 401 {object} frameworkdto.ErrorResponseDTO "Unauthorized"
// @Failure 403 {object} frameworkdto.ErrorResponseDTO "Not authorized to get this resource"
// @Failure 500 {object} frameworkdto.ErrorResponseDTO "Internal server error"
// @Router /user-maintenance/users/get-roles [get]
func (h *UserMaintenanceHandler) GetUserRoles(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role != string(frameworkconstants.UserRoleTenantAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to get this resource"))
		return
	}

	roles, err := h.userMaintenanceService.GetUserRoles()
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, roles, "User roles fetched successfully")
}
