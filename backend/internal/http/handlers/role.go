package handlers

import (
	"strconv"

	"gebase/internal/http/response"
	"gebase/internal/middleware"
	"gebase/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// List godoc
// @Summary List roles
// @Description Get paginated list of roles
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param system_id query int false "Filter by system ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	systemIDStr := c.Query("system_id")
	if systemIDStr != "" {
		systemID, err := strconv.Atoi(systemIDStr)
		if err != nil {
			response.BadRequest(c, "Invalid system ID")
			return
		}
		roles, err := h.roleService.ListRolesBySystem(c.Request.Context(), systemID)
		if err != nil {
			response.InternalError(c, "Failed to list roles")
			return
		}
		response.Success(c, gin.H{"roles": roles})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.roleService.ListRoles(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, "Failed to list roles")
		return
	}

	response.SuccessWithMeta(c, result.Data, response.FromPagination(
		result.Page, result.PageSize, result.Total, result.TotalPages,
	))
}

// Create godoc
// @Summary Create role
// @Description Create a new role
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body service.CreateRoleRequest true "Role info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)
	role, err := h.roleService.CreateRole(c.Request.Context(), &req, userID)
	if err != nil {
		if err == service.ErrRoleCodeExists {
			response.Conflict(c, "Role code already exists")
			return
		}
		response.InternalError(c, "Failed to create role")
		return
	}

	response.Created(c, role)
}

// Get godoc
// @Summary Get role
// @Description Get role by ID
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Role ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /roles/{id} [get]
func (h *RoleHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	role, err := h.roleService.GetRoleWithPermissions(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Role not found")
		return
	}

	response.Success(c, role)
}

// Update godoc
// @Summary Update role
// @Description Update role info
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Role ID"
// @Param request body service.UpdateRoleRequest true "Role info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)
	role, err := h.roleService.UpdateRole(c.Request.Context(), id, &req, userID)
	if err != nil {
		if err == service.ErrRoleNotFound {
			response.NotFound(c, "Role not found")
			return
		}
		response.InternalError(c, "Failed to update role")
		return
	}

	response.Success(c, role)
}

// Delete godoc
// @Summary Delete role
// @Description Soft delete role
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Role ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.roleService.DeleteRole(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrRoleNotFound {
			response.NotFound(c, "Role not found")
			return
		}
		response.InternalError(c, "Failed to delete role")
		return
	}

	response.Success(c, gin.H{"message": "Role deleted"})
}

// GetPermissions godoc
// @Summary Get role permissions
// @Description Get permissions assigned to a role
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Role ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /roles/{id}/permissions [get]
func (h *RoleHandler) GetPermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	role, err := h.roleService.GetRoleWithPermissions(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Role not found")
		return
	}

	response.Success(c, gin.H{"permissions": role.Permissions})
}

// AssignPermissions godoc
// @Summary Assign permissions to role
// @Description Assign permissions to a role
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Role ID"
// @Param request body map[string][]int true "Permission IDs"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /roles/{id}/permissions [put]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	var req struct {
		PermissionIDs []int `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Permission IDs are required")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.roleService.AssignPermissions(c.Request.Context(), id, req.PermissionIDs, userID); err != nil {
		response.InternalError(c, "Failed to assign permissions")
		return
	}

	response.Success(c, gin.H{"message": "Permissions assigned"})
}

// GetMenus godoc
// @Summary Get role menus
// @Description Get menus assigned to a role
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Role ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /roles/{id}/menus [get]
func (h *RoleHandler) GetMenus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	role, err := h.roleService.GetRoleWithMenus(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Role not found")
		return
	}

	response.Success(c, gin.H{"menus": role.Menus})
}

// AssignMenus godoc
// @Summary Assign menus to role
// @Description Assign menus to a role
// @Tags Roles
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Role ID"
// @Param request body map[string][]int true "Menu IDs"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /roles/{id}/menus [put]
func (h *RoleHandler) AssignMenus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	var req struct {
		MenuIDs []int `json:"menu_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Menu IDs are required")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.roleService.AssignMenus(c.Request.Context(), id, req.MenuIDs, userID); err != nil {
		response.InternalError(c, "Failed to assign menus")
		return
	}

	response.Success(c, gin.H{"message": "Menus assigned"})
}
