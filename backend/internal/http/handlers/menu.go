package handlers

import (
	"strconv"

	"gebase/internal/domain"
	"gebase/internal/http/response"
	"gebase/internal/middleware"
	"gebase/internal/service"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuService   *service.MenuService
	systemService *service.SystemService
}

func NewMenuHandler(menuService *service.MenuService, systemService *service.SystemService) *MenuHandler {
	return &MenuHandler{
		menuService:   menuService,
		systemService: systemService,
	}
}

// List godoc
// @Summary List menus
// @Description Get list of menus with optional system filter
// @Tags Menus
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param system_id query int false "Filter by system ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /menus [get]
func (h *MenuHandler) List(c *gin.Context) {
	systemIDStr := c.Query("system_id")

	var menus []domain.Menu
	var err error

	if systemIDStr != "" {
		systemID, parseErr := strconv.Atoi(systemIDStr)
		if parseErr != nil {
			response.BadRequest(c, "Invalid system ID")
			return
		}
		menus, err = h.menuService.GetAllMenus(c.Request.Context(), systemID)
	} else {
		// Get all menus from all systems - need to get systems first
		systems, sysErr := h.systemService.ListSystems(c.Request.Context())
		if sysErr != nil {
			response.InternalError(c, "Failed to list systems")
			return
		}

		for _, sys := range systems {
			sysMenus, _ := h.menuService.GetAllMenus(c.Request.Context(), sys.ID)
			menus = append(menus, sysMenus...)
		}
	}

	if err != nil {
		response.InternalError(c, "Failed to list menus")
		return
	}

	response.Success(c, menus)
}

// GetTree godoc
// @Summary Get menu tree
// @Description Get menu tree structure for a system
// @Tags Menus
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param system_id query int true "System ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /menus/tree [get]
func (h *MenuHandler) GetTree(c *gin.Context) {
	systemIDStr := c.Query("system_id")
	if systemIDStr == "" {
		response.BadRequest(c, "System ID is required")
		return
	}

	systemID, err := strconv.Atoi(systemIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid system ID")
		return
	}

	tree, err := h.menuService.GetMenuTree(c.Request.Context(), systemID)
	if err != nil {
		response.InternalError(c, "Failed to get menu tree")
		return
	}

	response.Success(c, tree)
}

// Create godoc
// @Summary Create menu
// @Description Create a new menu
// @Tags Menus
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body CreateMenuRequest true "Menu info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /menus [post]
func (h *MenuHandler) Create(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)

	menu := &domain.Menu{
		Code:      req.Code,
		Name:      req.Name,
		SystemID:  req.SystemID,
		ParentID:  req.ParentID,
		Path:      req.Path,
		Icon:      req.Icon,
		Component: req.Component,
		Sequence:  req.Sequence,
		IsVisible: req.IsVisible,
		IsActive:  req.IsActive,
	}
	menu.CreatedBy = &userID

	if err := h.menuService.CreateMenu(c.Request.Context(), menu); err != nil {
		response.InternalError(c, "Failed to create menu")
		return
	}

	response.Created(c, menu)
}

// Get godoc
// @Summary Get menu
// @Description Get menu by ID
// @Tags Menus
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Menu ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /menus/{id} [get]
func (h *MenuHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	menu, err := h.menuService.GetMenuByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Menu not found")
		return
	}

	response.Success(c, menu)
}

// Update godoc
// @Summary Update menu
// @Description Update menu info
// @Tags Menus
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Menu ID"
// @Param request body UpdateMenuRequest true "Menu info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /menus/{id} [put]
func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	var req UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	menu, err := h.menuService.GetMenuByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Menu not found")
		return
	}

	userID := middleware.GetUserID(c)

	// Update fields
	if req.Code != "" {
		menu.Code = req.Code
	}
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.SystemID != nil {
		menu.SystemID = req.SystemID
	}
	menu.ParentID = req.ParentID
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Sequence != 0 {
		menu.Sequence = req.Sequence
	}
	if req.IsVisible != nil {
		menu.IsVisible = req.IsVisible
	}
	if req.IsActive != nil {
		menu.IsActive = req.IsActive
	}
	menu.UpdatedBy = &userID

	if err := h.menuService.UpdateMenu(c.Request.Context(), menu); err != nil {
		response.InternalError(c, "Failed to update menu")
		return
	}

	response.Success(c, menu)
}

// Delete godoc
// @Summary Delete menu
// @Description Soft delete menu
// @Tags Menus
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Menu ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /menus/{id} [delete]
func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid menu ID")
		return
	}

	if err := h.menuService.DeleteMenu(c.Request.Context(), id); err != nil {
		response.InternalError(c, "Failed to delete menu")
		return
	}

	response.Success(c, gin.H{"message": "Menu deleted"})
}

// Request DTOs
type CreateMenuRequest struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	SystemID  *int   `json:"system_id" binding:"required"`
	ParentID  *int   `json:"parent_id"`
	Path      string `json:"path"`
	Icon      string `json:"icon"`
	Component string `json:"component"`
	Sequence  int    `json:"sequence"`
	IsVisible *bool  `json:"is_visible"`
	IsActive  *bool  `json:"is_active"`
}

type UpdateMenuRequest struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	SystemID  *int   `json:"system_id"`
	ParentID  *int   `json:"parent_id"`
	Path      string `json:"path"`
	Icon      string `json:"icon"`
	Component string `json:"component"`
	Sequence  int    `json:"sequence"`
	IsVisible *bool  `json:"is_visible"`
	IsActive  *bool  `json:"is_active"`
}
