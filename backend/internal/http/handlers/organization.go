package handlers

import (
	"strconv"

	"gebase/internal/http/response"
	"gebase/internal/middleware"
	"gebase/internal/service"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	orgService *service.OrganizationService
}

func NewOrganizationHandler(orgService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

// List godoc
// @Summary List organizations
// @Description Get paginated list of organizations
// @Tags Organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /organizations [get]
func (h *OrganizationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.orgService.ListOrganizations(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, "Failed to list organizations")
		return
	}

	response.SuccessWithMeta(c, result.Data, response.FromPagination(
		result.Page, result.PageSize, result.Total, result.TotalPages,
	))
}

// Create godoc
// @Summary Create organization
// @Description Create a new organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body service.CreateOrganizationRequest true "Organization info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /organizations [post]
func (h *OrganizationHandler) Create(c *gin.Context) {
	var req service.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)
	org, err := h.orgService.CreateOrganization(c.Request.Context(), &req, userID)
	if err != nil {
		if err == service.ErrOrgRegNoExists {
			response.Conflict(c, "Organization registration number already exists")
			return
		}
		response.InternalError(c, "Failed to create organization")
		return
	}

	response.Created(c, org)
}

// Get godoc
// @Summary Get organization
// @Description Get organization by ID
// @Tags Organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Organization ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /organizations/{id} [get]
func (h *OrganizationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	org, err := h.orgService.GetOrganizationWithChildren(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Organization not found")
		return
	}

	response.Success(c, org)
}

// Update godoc
// @Summary Update organization
// @Description Update organization info
// @Tags Organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Organization ID"
// @Param request body service.UpdateOrganizationRequest true "Organization info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /organizations/{id} [put]
func (h *OrganizationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	var req service.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)
	org, err := h.orgService.UpdateOrganization(c.Request.Context(), id, &req, userID)
	if err != nil {
		if err == service.ErrOrganizationNotFound {
			response.NotFound(c, "Organization not found")
			return
		}
		response.InternalError(c, "Failed to update organization")
		return
	}

	response.Success(c, org)
}

// Delete godoc
// @Summary Delete organization
// @Description Soft delete organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Organization ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /organizations/{id} [delete]
func (h *OrganizationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.orgService.DeleteOrganization(c.Request.Context(), id, userID); err != nil {
		if err == service.ErrOrganizationNotFound {
			response.NotFound(c, "Organization not found")
			return
		}
		response.InternalError(c, "Failed to delete organization")
		return
	}

	response.Success(c, gin.H{"message": "Organization deleted"})
}

// GetChildren godoc
// @Summary Get child organizations
// @Description Get child organizations of a parent
// @Tags Organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Parent Organization ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /organizations/{id}/children [get]
func (h *OrganizationHandler) GetChildren(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	children, err := h.orgService.GetChildOrganizations(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "Failed to get child organizations")
		return
	}

	response.Success(c, gin.H{"children": children})
}

// GetTypes godoc
// @Summary Get organization types
// @Description Get all organization types
// @Tags Organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /organizations/types [get]
func (h *OrganizationHandler) GetTypes(c *gin.Context) {
	types, err := h.orgService.GetOrganizationTypes(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Failed to get organization types")
		return
	}

	response.Success(c, gin.H{"types": types})
}

// GetSystems godoc
// @Summary Get enabled systems
// @Description Get systems enabled for an organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Organization ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /organizations/{id}/systems [get]
func (h *OrganizationHandler) GetSystems(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	systems, err := h.orgService.GetEnabledSystems(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "Failed to get enabled systems")
		return
	}

	response.Success(c, gin.H{"systems": systems})
}

// EnableSystem godoc
// @Summary Enable system for organization
// @Description Enable a system for an organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Organization ID"
// @Param request body map[string]int true "System ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /organizations/{id}/systems [post]
func (h *OrganizationHandler) EnableSystem(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	var req struct {
		SystemID int `json:"system_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "System ID is required")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.orgService.EnableSystem(c.Request.Context(), id, req.SystemID, userID); err != nil {
		response.InternalError(c, "Failed to enable system")
		return
	}

	response.Success(c, gin.H{"message": "System enabled"})
}

// DisableSystem godoc
// @Summary Disable system for organization
// @Description Disable a system for an organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Organization ID"
// @Param system_id path int true "System ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /organizations/{id}/systems/{system_id} [delete]
func (h *OrganizationHandler) DisableSystem(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid organization ID")
		return
	}

	systemID, err := strconv.Atoi(c.Param("system_id"))
	if err != nil {
		response.BadRequest(c, "Invalid system ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.orgService.DisableSystem(c.Request.Context(), id, systemID, userID); err != nil {
		response.InternalError(c, "Failed to disable system")
		return
	}

	response.Success(c, gin.H{"message": "System disabled"})
}
