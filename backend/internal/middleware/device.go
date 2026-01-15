package middleware

import (
	"net/http"

	"gebase/internal/service"

	"github.com/gin-gonic/gin"
)

type DeviceMiddleware struct {
	deviceService *service.DeviceService
}

func NewDeviceMiddleware(deviceService *service.DeviceService) *DeviceMiddleware {
	return &DeviceMiddleware{
		deviceService: deviceService,
	}
}

// Device validates device headers and verifies device is registered
func (m *DeviceMiddleware) Device() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceUID := c.GetHeader("X-Device-UID")
		if deviceUID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "DEVICE_UID_REQUIRED",
					"message": "X-Device-UID header is required",
				},
			})
			return
		}

		platform := c.GetHeader("X-Platform")
		if platform == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "PLATFORM_REQUIRED",
					"message": "X-Platform header is required",
				},
			})
			return
		}

		// Verify device
		device, err := m.deviceService.VerifyDevice(c.Request.Context(), deviceUID)
		if err != nil {
			if err == service.ErrDeviceNotFound {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "DEVICE_NOT_REGISTERED",
						"message": "Device is not registered. Please register the device first.",
					},
				})
				return
			}
			if err == service.ErrDeviceNotActive {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "DEVICE_INACTIVE",
						"message": "Device is deactivated. Please contact administrator.",
					},
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "DEVICE_VERIFICATION_FAILED",
					"message": "Failed to verify device",
				},
			})
			return
		}

		// Update heartbeat asynchronously
		go m.deviceService.Heartbeat(c.Request.Context(), deviceUID)

		// Set device info in context
		c.Set("device", device)
		c.Set("device_uid", deviceUID)
		c.Set("platform", platform)

		c.Next()
	}
}

// OptionalDevice validates device headers if present but doesn't require them
func (m *DeviceMiddleware) OptionalDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceUID := c.GetHeader("X-Device-UID")
		platform := c.GetHeader("X-Platform")

		if deviceUID != "" {
			device, err := m.deviceService.VerifyDevice(c.Request.Context(), deviceUID)
			if err == nil {
				c.Set("device", device)
				c.Set("device_uid", deviceUID)
				c.Set("platform", platform)

				// Update heartbeat asynchronously
				go m.deviceService.Heartbeat(c.Request.Context(), deviceUID)
			}
		}

		c.Next()
	}
}
