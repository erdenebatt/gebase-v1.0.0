package middleware

import (
	"net/http"

	"gebase/internal/service"

	"github.com/gin-gonic/gin"
)

type RBACMiddleware struct {
	permissionService *service.PermissionService
}

func NewRBACMiddleware(permissionService *service.PermissionService) *RBACMiddleware {
	return &RBACMiddleware{
		permissionService: permissionService,
	}
}

// RequirePermission checks if user has the specified permission
func (m *RBACMiddleware) RequirePermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		systemID := GetSystemID(c)

		hasPermission, err := m.permissionService.CheckPermission(
			c.Request.Context(),
			userID,
			systemID,
			permissionCode,
		)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "PERMISSION_CHECK_FAILED",
					"message": "Failed to check permission",
				},
			})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":       "FORBIDDEN",
					"message":    "You don't have permission to perform this action",
					"permission": permissionCode,
				},
			})
			return
		}

		c.Next()
	}
}

// RequireAnyPermission checks if user has any of the specified permissions
func (m *RBACMiddleware) RequireAnyPermission(permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		systemID := GetSystemID(c)

		for _, code := range permissionCodes {
			hasPermission, err := m.permissionService.CheckPermission(
				c.Request.Context(),
				userID,
				systemID,
				code,
			)

			if err == nil && hasPermission {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":        "FORBIDDEN",
				"message":     "You don't have permission to perform this action",
				"permissions": permissionCodes,
			},
		})
	}
}

// RequireAllPermissions checks if user has all of the specified permissions
func (m *RBACMiddleware) RequireAllPermissions(permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		systemID := GetSystemID(c)

		for _, code := range permissionCodes {
			hasPermission, err := m.permissionService.CheckPermission(
				c.Request.Context(),
				userID,
				systemID,
				code,
			)

			if err != nil || !hasPermission {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"success": false,
					"error": gin.H{
						"code":              "FORBIDDEN",
						"message":           "You don't have all required permissions",
						"missing_permission": code,
					},
				})
				return
			}
		}

		c.Next()
	}
}

// RequireSystemContext ensures request has system context
func (m *RBACMiddleware) RequireSystemContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		systemID := GetSystemID(c)
		if systemID == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "SYSTEM_CONTEXT_REQUIRED",
					"message": "System context required. Please switch to a system first.",
				},
			})
			return
		}

		c.Next()
	}
}
