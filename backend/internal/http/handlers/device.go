package handlers

import (
	"strconv"

	"gebase/internal/http/response"
	"gebase/internal/middleware"
	"gebase/internal/service"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	deviceService *service.DeviceService
}

func NewDeviceHandler(deviceService *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// Register godoc
// @Summary Register device
// @Description Register a new device or update existing
// @Tags Devices
// @Accept json
// @Produce json
// @Param request body service.RegisterDeviceRequest true "Device info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /devices/register [post]
func (h *DeviceHandler) Register(c *gin.Context) {
	var req service.RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	device, err := h.deviceService.RegisterDevice(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, "Failed to register device")
		return
	}

	response.Success(c, device)
}

// Heartbeat godoc
// @Summary Device heartbeat
// @Description Update device heartbeat timestamp
// @Tags Devices
// @Accept json
// @Produce json
// @Param X-Device-UID header string true "Device UID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /devices/heartbeat [post]
func (h *DeviceHandler) Heartbeat(c *gin.Context) {
	deviceUID := c.GetHeader("X-Device-UID")
	if deviceUID == "" {
		response.BadRequest(c, "X-Device-UID header is required")
		return
	}

	if err := h.deviceService.Heartbeat(c.Request.Context(), deviceUID); err != nil {
		response.InternalError(c, "Failed to update heartbeat")
		return
	}

	response.Success(c, gin.H{"message": "Heartbeat recorded"})
}

// List godoc
// @Summary List devices
// @Description Get paginated list of devices
// @Tags Devices
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /devices [get]
func (h *DeviceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.deviceService.ListDevices(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, "Failed to list devices")
		return
	}

	response.SuccessWithMeta(c, result.Data, response.FromPagination(
		result.Page, result.PageSize, result.Total, result.TotalPages,
	))
}

// Get godoc
// @Summary Get device
// @Description Get device by ID
// @Tags Devices
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Device ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /devices/{id} [get]
func (h *DeviceHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	device, err := h.deviceService.GetDevice(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "Device not found")
		return
	}

	response.Success(c, device)
}

// Update godoc
// @Summary Update device
// @Description Update device info
// @Tags Devices
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Device ID"
// @Param request body service.UpdateDeviceRequest true "Device info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /devices/{id} [put]
func (h *DeviceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	var req service.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	userID := middleware.GetUserID(c)
	device, err := h.deviceService.UpdateDevice(c.Request.Context(), id, &req, userID)
	if err != nil {
		if err == service.ErrDeviceNotFound {
			response.NotFound(c, "Device not found")
			return
		}
		response.InternalError(c, "Failed to update device")
		return
	}

	response.Success(c, device)
}

// Deactivate godoc
// @Summary Deactivate device
// @Description Deactivate device and terminate its sessions
// @Tags Devices
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Device ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /devices/{id} [delete]
func (h *DeviceHandler) Deactivate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.deviceService.DeactivateDevice(c.Request.Context(), id, "admin", userID); err != nil {
		if err == service.ErrDeviceNotFound {
			response.NotFound(c, "Device not found")
			return
		}
		response.InternalError(c, "Failed to deactivate device")
		return
	}

	response.Success(c, gin.H{"message": "Device deactivated"})
}

// UpdateConfig godoc
// @Summary Update device config
// @Description Update device remote configuration
// @Tags Devices
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Device ID"
// @Param request body map[string]interface{} true "Config JSON"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /devices/{id}/config [put]
func (h *DeviceHandler) UpdateConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	var req struct {
		Config string `json:"config" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Config is required")
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.deviceService.UpdateDeviceConfig(c.Request.Context(), id, req.Config, userID); err != nil {
		response.InternalError(c, "Failed to update config")
		return
	}

	response.Success(c, gin.H{"message": "Config updated"})
}

// GetSessions godoc
// @Summary Get device sessions
// @Description Get active sessions for a device
// @Tags Devices
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "Device ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /devices/{id}/sessions [get]
func (h *DeviceHandler) GetSessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid device ID")
		return
	}

	sessions, err := h.deviceService.GetDeviceSessions(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "Failed to get sessions")
		return
	}

	response.Success(c, gin.H{"sessions": sessions})
}
