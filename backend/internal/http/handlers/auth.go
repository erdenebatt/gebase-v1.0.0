package handlers

import (
	"net/http"

	"gebase/internal/auth"
	"gebase/internal/http/response"
	"gebase/internal/middleware"
	"gebase/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService       *service.AuthService
	permissionService *service.PermissionService
	menuService       *service.MenuService
}

func NewAuthHandler(
	authService *service.AuthService,
	permissionService *service.PermissionService,
	menuService *service.MenuService,
) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		permissionService: permissionService,
		menuService:       menuService,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return platform token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=service.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	ipAddress := middleware.GetClientIP(c)
	userAgent := c.Request.UserAgent()

	result, err := h.authService.Login(c.Request.Context(), &req, ipAddress, userAgent)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			response.Unauthorized(c, "Invalid email or password")
		case service.ErrUserInactive:
			response.Forbidden(c, "User account is inactive")
		case service.ErrDeviceNotFound:
			response.Error(c, http.StatusBadRequest, "DEVICE_NOT_REGISTERED", "Device is not registered")
		case service.ErrDeviceNotActive:
			response.Forbidden(c, "Device is deactivated")
		default:
			response.InternalError(c, "Login failed")
		}
		return
	}

	response.Success(c, result)
}

// SwitchSystem godoc
// @Summary Switch to a system
// @Description Switch to a specific system and get system token
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body service.SwitchSystemRequest true "System to switch to"
// @Success 200 {object} response.Response{data=service.SwitchSystemResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /auth/switch-system [post]
func (h *AuthHandler) SwitchSystem(c *gin.Context) {
	var req service.SwitchSystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	claims := middleware.GetClaims(c)
	ipAddress := middleware.GetClientIP(c)

	result, err := h.authService.SwitchSystem(c.Request.Context(), claims, &req, ipAddress)
	if err != nil {
		switch err {
		case service.ErrSystemNotFound:
			response.NotFound(c, "System not found")
		case service.ErrNoSystemAccess:
			response.Forbidden(c, "You don't have access to this system")
		default:
			response.InternalError(c, "Failed to switch system")
		}
		return
	}

	response.Success(c, result)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} response.Response{data=service.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Refresh token is required")
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == auth.ErrExpiredToken {
			response.Unauthorized(c, "Refresh token has expired")
			return
		}
		response.Unauthorized(c, "Invalid refresh token")
		return
	}

	response.Success(c, result)
}

// Logout godoc
// @Summary User logout
// @Description Terminate current session
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID := middleware.GetSessionID(c)

	if err := h.authService.Logout(c.Request.Context(), sessionID); err != nil {
		response.InternalError(c, "Failed to logout")
		return
	}

	response.Success(c, gin.H{"message": "Logged out successfully"})
}

// Me godoc
// @Summary Get current user
// @Description Get current authenticated user info
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.authService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, "Failed to get user info")
		return
	}

	response.Success(c, user)
}

// GetPermissions godoc
// @Summary Get user permissions
// @Description Get permissions for current user in current system
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/permissions [get]
func (h *AuthHandler) GetPermissions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	systemID := middleware.GetSystemID(c)

	permissions, err := h.permissionService.GetUserPermissionCodes(c.Request.Context(), userID, systemID)
	if err != nil {
		response.InternalError(c, "Failed to get permissions")
		return
	}

	response.Success(c, gin.H{"permissions": permissions})
}

// GetMenus godoc
// @Summary Get user menus
// @Description Get menu tree for current user in current system
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/menus [get]
func (h *AuthHandler) GetMenus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	systemID := middleware.GetSystemID(c)

	if systemID == nil {
		response.BadRequest(c, "System context required. Please switch to a system first.")
		return
	}

	menus, err := h.menuService.GetUserMenus(c.Request.Context(), userID, *systemID)
	if err != nil {
		response.InternalError(c, "Failed to get menus")
		return
	}

	response.Success(c, gin.H{"menus": menus})
}

// ExitSystem godoc
// @Summary Exit current system
// @Description Exit current system and return to platform level
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/exit-system [post]
func (h *AuthHandler) ExitSystem(c *gin.Context) {
	// Exit system simply acknowledges the request
	// The frontend handles clearing the system context
	// Session tracking can be updated here if needed
	response.Success(c, gin.H{"message": "Exited system successfully"})
}

// GetAvailableSystems godoc
// @Summary Get available systems
// @Description Get systems available to current user
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/systems [get]
func (h *AuthHandler) GetAvailableSystems(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.authService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, "Failed to get user info")
		return
	}

	systems := user.GetAvailableSystems()
	response.Success(c, gin.H{"systems": systems})
}
