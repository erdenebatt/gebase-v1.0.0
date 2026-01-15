package handlers

import (
	"strconv"

	"gebase/internal/http/response"
	"gebase/internal/middleware"
	"gebase/internal/service"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	systemService *service.SystemService
}

func NewSystemHandler(systemService *service.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

// List godoc
// @Summary List systems
// @Description Get all systems
// @Tags Systems
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /systems [get]
func (h *SystemHandler) List(c *gin.Context) {
	systems, err := h.systemService.ListSystems(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Failed to list systems")
		return
	}

	response.Success(c, gin.H{"systems": systems})
}

// Create godoc
// @Summary Create system
// @Description Create a new system
// @Tags Systems
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body service.CreateSystemRequest true "System info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /systems [post]
func (h *SystemHandler) Create(c *gin.Context) {
	var req service.CreateSystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)
	system, err := h.systemService.CreateSystem(c.Request.Context(), &req, userID)
	if err != nil {
		if err == service.ErrSystemCodeExists {
			response.Conflict(c, "System code already exists")
			return
		}
		response.InternalError(c, "Failed to create system")
		return
	}

	response.Created(c, system)
}

// Get godoc
// @Summary Get system
// @Description Get system by ID
// @Tags Systems
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "System ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /systems/{id} [get]
func (h *SystemHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid system ID")
		return
	}

	system, err := h.systemService.GetSystemWithModules(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "System not found")
		return
	}

	response.Success(c, system)
}

// Update godoc
// @Summary Update system
// @Description Update system info
// @Tags Systems
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "System ID"
// @Param request body service.UpdateSystemRequest true "System info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /systems/{id} [put]
func (h *SystemHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid system ID")
		return
	}

	var req service.UpdateSystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)
	system, err := h.systemService.UpdateSystem(c.Request.Context(), id, &req, userID)
	if err != nil {
		if err == service.ErrSystemNotFound {
			response.NotFound(c, "System not found")
			return
		}
		response.InternalError(c, "Failed to update system")
		return
	}

	response.Success(c, system)
}

// Delete godoc
// @Summary Delete system
// @Description Soft delete system
// @Tags Systems
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "System ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /systems/{id} [delete]
func (h *SystemHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid system ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.systemService.DeleteSystem(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrSystemNotFound {
			response.NotFound(c, "System not found")
			return
		}
		response.InternalError(c, "Failed to delete system")
		return
	}

	response.Success(c, gin.H{"message": "System deleted"})
}

// GetModules godoc
// @Summary Get system modules
// @Description Get modules for a system
// @Tags Systems
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "System ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /systems/{id}/modules [get]
func (h *SystemHandler) GetModules(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid system ID")
		return
	}

	modules, err := h.systemService.GetSystemModules(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "Failed to get modules")
		return
	}

	response.Success(c, gin.H{"modules": modules})
}

// GetMenus godoc
// @Summary Get system menus
// @Description Get menus for a system
// @Tags Systems
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "System ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /systems/{id}/menus [get]
func (h *SystemHandler) GetMenus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid system ID")
		return
	}

	menus, err := h.systemService.GetSystemMenus(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "Failed to get menus")
		return
	}

	response.Success(c, gin.H{"menus": menus})
}
