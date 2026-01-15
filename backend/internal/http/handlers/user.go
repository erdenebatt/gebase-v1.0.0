package handlers

import (
	"strconv"

	"gebase/internal/http/response"
	"gebase/internal/middleware"
	"gebase/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// List godoc
// @Summary List users
// @Description Get paginated list of users
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, "Failed to list users")
		return
	}

	response.SuccessWithMeta(c, result.Data, response.FromPagination(
		result.Page, result.PageSize, result.Total, result.TotalPages,
	))
}

// Create godoc
// @Summary Create user
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body service.CreateUserRequest true "User info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)
	user, err := h.userService.CreateUser(c.Request.Context(), &req, userID)
	if err != nil {
		switch err {
		case service.ErrEmailAlreadyExists:
			response.Conflict(c, "Email already exists")
		case service.ErrRegNoAlreadyExists:
			response.Conflict(c, "Registration number already exists")
		default:
			response.InternalError(c, "Failed to create user")
		}
		return
	}

	response.Created(c, user)
}

// Get godoc
// @Summary Get user
// @Description Get user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserWithRoles(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, user)
}

// Update godoc
// @Summary Update user
// @Description Update user info
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Param request body service.UpdateUserRequest true "User info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	currentUserID := middleware.GetUserID(c)
	user, err := h.userService.UpdateUser(c.Request.Context(), id, &req, currentUserID)
	if err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalError(c, "Failed to update user")
		return
	}

	response.Success(c, user)
}

// Delete godoc
// @Summary Delete user
// @Description Soft delete user
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	currentUserID := middleware.GetUserID(c)
	if err := h.userService.DeleteUser(c.Request.Context(), id, currentUserID); err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalError(c, "Failed to delete user")
		return
	}

	response.Success(c, gin.H{"message": "User deleted"})
}

// GetRoles godoc
// @Summary Get user roles
// @Description Get roles assigned to user
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /users/{id}/roles [get]
func (h *UserHandler) GetRoles(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	roles, err := h.userService.GetUserRoles(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "Failed to get user roles")
		return
	}

	response.Success(c, gin.H{"roles": roles})
}

// AssignRoles godoc
// @Summary Assign roles to user
// @Description Assign roles to user for a specific system
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Param request body map[string]interface{} true "Role assignment"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /users/{id}/roles [put]
func (h *UserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req struct {
		SystemID       *int   `json:"system_id"`
		RoleIDs        []int  `json:"role_ids" binding:"required"`
		OrganizationID *int64 `json:"organization_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	currentUserID := middleware.GetUserID(c)
	if err := h.userService.AssignUserRoles(c.Request.Context(), id, req.SystemID, req.RoleIDs, req.OrganizationID, currentUserID); err != nil {
		response.InternalError(c, "Failed to assign roles")
		return
	}

	response.Success(c, gin.H{"message": "Roles assigned"})
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Admin reset user password
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Param request body map[string]string true "New password"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /users/{id}/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Password must be at least 6 characters")
		return
	}

	currentUserID := middleware.GetUserID(c)
	if err := h.userService.ResetPassword(c.Request.Context(), id, req.Password, currentUserID); err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalError(c, "Failed to reset password")
		return
	}

	response.Success(c, gin.H{"message": "Password reset successfully"})
}
