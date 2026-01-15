package middleware

import (
	"net/http"
	"strings"

	"gebase/internal/auth"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService     *auth.JWTService
	sessionService *auth.SessionService
}

func NewAuthMiddleware(jwtService *auth.JWTService, sessionService *auth.SessionService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:     jwtService,
		sessionService: sessionService,
	}
}

// Auth validates JWT token and sets user context
func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authorization header required",
				},
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Invalid authorization header format",
				},
			})
			return
		}

		claims, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			if err == auth.ErrExpiredToken {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "TOKEN_EXPIRED",
						"message": "Token has expired",
					},
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
				},
			})
			return
		}

		// Validate session
		session, err := m.sessionService.GetSessionByID(c.Request.Context(), claims.SessionID)
		if err != nil || !m.sessionService.IsSessionValid(c.Request.Context(), session) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "SESSION_INVALID",
					"message": "Session is invalid or expired",
				},
			})
			return
		}

		// Update session activity
		go m.sessionService.UpdateActivity(c.Request.Context(), claims.SessionID)

		// Set context values
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("organization_id", claims.OrganizationID)
		c.Set("session_id", claims.SessionID)
		c.Set("device_id", claims.DeviceID)
		c.Set("token_type", string(claims.TokenType))
		c.Set("claims", claims)

		if claims.SystemID != nil {
			c.Set("system_id", *claims.SystemID)
			c.Set("system_code", claims.SystemCode)
		}
		if claims.RoleIDs != nil {
			c.Set("role_ids", claims.RoleIDs)
		}

		c.Next()
	}
}

// RequirePlatformToken ensures the token is a platform token
func (m *AuthMiddleware) RequirePlatformToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authentication required",
				},
			})
			return
		}

		authClaims := claims.(*auth.Claims)
		if !authClaims.IsPlatformToken() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "PLATFORM_TOKEN_REQUIRED",
					"message": "Platform token required for this operation",
				},
			})
			return
		}

		c.Next()
	}
}

// RequireSystemToken ensures the token is a system token
func (m *AuthMiddleware) RequireSystemToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authentication required",
				},
			})
			return
		}

		authClaims := claims.(*auth.Claims)
		if !authClaims.IsSystemToken() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "SYSTEM_TOKEN_REQUIRED",
					"message": "System token required. Please switch to a system first.",
				},
			})
			return
		}

		c.Next()
	}
}

// GetUserID helper to get user ID from context
func GetUserID(c *gin.Context) int64 {
	userID, _ := c.Get("user_id")
	return userID.(int64)
}

// GetSessionID helper to get session ID from context
func GetSessionID(c *gin.Context) int64 {
	sessionID, _ := c.Get("session_id")
	return sessionID.(int64)
}

// GetSystemID helper to get system ID from context (may be nil)
func GetSystemID(c *gin.Context) *int {
	systemID, exists := c.Get("system_id")
	if !exists {
		return nil
	}
	id := systemID.(int)
	return &id
}

// GetClaims helper to get claims from context
func GetClaims(c *gin.Context) *auth.Claims {
	claims, _ := c.Get("claims")
	return claims.(*auth.Claims)
}
