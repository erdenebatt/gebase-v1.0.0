package router

import (
	"gebase/internal/config"
	"gebase/internal/http/handlers"
	"gebase/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	cfg    *config.Config
}

func NewRouter(cfg *config.Config) *Router {
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	return &Router{
		engine: engine,
		cfg:    cfg,
	}
}

func (r *Router) Setup(
	authMiddleware *middleware.AuthMiddleware,
	rbacMiddleware *middleware.RBACMiddleware,
	deviceMiddleware *middleware.DeviceMiddleware,
	authHandler *handlers.AuthHandler,
	deviceHandler *handlers.DeviceHandler,
	userHandler *handlers.UserHandler,
	orgHandler *handlers.OrganizationHandler,
	systemHandler *handlers.SystemHandler,
	roleHandler *handlers.RoleHandler,
	menuHandler *handlers.MenuHandler,
) *gin.Engine {
	// Global middleware
	r.engine.Use(middleware.CORS(r.cfg))
	r.engine.Use(middleware.Logger())

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	api := r.engine.Group("/api/v1")

	// Public routes
	r.setupPublicRoutes(api, authHandler, deviceHandler)

	// Protected routes
	r.setupProtectedRoutes(api, authMiddleware, rbacMiddleware, deviceMiddleware,
		authHandler, deviceHandler, userHandler, orgHandler, systemHandler, roleHandler, menuHandler)

	return r.engine
}

func (r *Router) setupPublicRoutes(
	api *gin.RouterGroup,
	authHandler *handlers.AuthHandler,
	deviceHandler *handlers.DeviceHandler,
) {
	// Auth
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Device registration
	devices := api.Group("/devices")
	{
		devices.POST("/register", deviceHandler.Register)
		devices.POST("/heartbeat", deviceHandler.Heartbeat)
	}
}

func (r *Router) setupProtectedRoutes(
	api *gin.RouterGroup,
	authMiddleware *middleware.AuthMiddleware,
	rbacMiddleware *middleware.RBACMiddleware,
	deviceMiddleware *middleware.DeviceMiddleware,
	authHandler *handlers.AuthHandler,
	deviceHandler *handlers.DeviceHandler,
	userHandler *handlers.UserHandler,
	orgHandler *handlers.OrganizationHandler,
	systemHandler *handlers.SystemHandler,
	roleHandler *handlers.RoleHandler,
	menuHandler *handlers.MenuHandler,
) {
	// Protected routes require auth and device verification
	protected := api.Group("")
	protected.Use(authMiddleware.Auth())
	protected.Use(deviceMiddleware.Device())

	// Auth routes (authenticated)
	auth := protected.Group("/auth")
	{
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", authHandler.Me)
		auth.POST("/switch-system", authHandler.SwitchSystem)
		auth.POST("/exit-system", authHandler.ExitSystem)
		auth.GET("/systems", authHandler.GetAvailableSystems)
		auth.GET("/permissions", authHandler.GetPermissions)
		auth.GET("/menus", authHandler.GetMenus)
	}

	// Users
	users := protected.Group("/users")
	{
		users.GET("", rbacMiddleware.RequirePermission("admin.user.view"), userHandler.List)
		users.POST("", rbacMiddleware.RequirePermission("admin.user.create"), userHandler.Create)
		users.GET("/:id", rbacMiddleware.RequirePermission("admin.user.view"), userHandler.Get)
		users.PUT("/:id", rbacMiddleware.RequirePermission("admin.user.update"), userHandler.Update)
		users.DELETE("/:id", rbacMiddleware.RequirePermission("admin.user.delete"), userHandler.Delete)
		users.GET("/:id/roles", rbacMiddleware.RequirePermission("admin.user.view"), userHandler.GetRoles)
		users.PUT("/:id/roles", rbacMiddleware.RequirePermission("admin.user.update"), userHandler.AssignRoles)
		users.POST("/:id/reset-password", rbacMiddleware.RequirePermission("admin.user.update"), userHandler.ResetPassword)
	}

	// Organizations
	orgs := protected.Group("/organizations")
	{
		orgs.GET("", rbacMiddleware.RequirePermission("admin.organization.view"), orgHandler.List)
		orgs.POST("", rbacMiddleware.RequirePermission("admin.organization.create"), orgHandler.Create)
		orgs.GET("/types", rbacMiddleware.RequirePermission("admin.organization.view"), orgHandler.GetTypes)
		orgs.GET("/:id", rbacMiddleware.RequirePermission("admin.organization.view"), orgHandler.Get)
		orgs.PUT("/:id", rbacMiddleware.RequirePermission("admin.organization.update"), orgHandler.Update)
		orgs.DELETE("/:id", rbacMiddleware.RequirePermission("admin.organization.delete"), orgHandler.Delete)
		orgs.GET("/:id/children", rbacMiddleware.RequirePermission("admin.organization.view"), orgHandler.GetChildren)
		orgs.GET("/:id/systems", rbacMiddleware.RequirePermission("admin.organization.view"), orgHandler.GetSystems)
		orgs.POST("/:id/systems", rbacMiddleware.RequirePermission("admin.organization.update"), orgHandler.EnableSystem)
		orgs.DELETE("/:id/systems/:system_id", rbacMiddleware.RequirePermission("admin.organization.update"), orgHandler.DisableSystem)
	}

	// Systems
	systems := protected.Group("/systems")
	{
		systems.GET("", rbacMiddleware.RequirePermission("admin.system.view"), systemHandler.List)
		systems.POST("", rbacMiddleware.RequirePermission("admin.system.create"), systemHandler.Create)
		systems.GET("/:id", rbacMiddleware.RequirePermission("admin.system.view"), systemHandler.Get)
		systems.PUT("/:id", rbacMiddleware.RequirePermission("admin.system.update"), systemHandler.Update)
		systems.DELETE("/:id", rbacMiddleware.RequirePermission("admin.system.delete"), systemHandler.Delete)
		systems.GET("/:id/modules", rbacMiddleware.RequirePermission("admin.system.view"), systemHandler.GetModules)
		systems.GET("/:id/menus", rbacMiddleware.RequirePermission("admin.system.view"), systemHandler.GetMenus)
	}

	// Roles
	roles := protected.Group("/roles")
	{
		roles.GET("", rbacMiddleware.RequirePermission("admin.role.view"), roleHandler.List)
		roles.POST("", rbacMiddleware.RequirePermission("admin.role.create"), roleHandler.Create)
		roles.GET("/:id", rbacMiddleware.RequirePermission("admin.role.view"), roleHandler.Get)
		roles.PUT("/:id", rbacMiddleware.RequirePermission("admin.role.update"), roleHandler.Update)
		roles.DELETE("/:id", rbacMiddleware.RequirePermission("admin.role.delete"), roleHandler.Delete)
		roles.GET("/:id/permissions", rbacMiddleware.RequirePermission("admin.role.view"), roleHandler.GetPermissions)
		roles.PUT("/:id/permissions", rbacMiddleware.RequirePermission("admin.role.update"), roleHandler.AssignPermissions)
		roles.GET("/:id/menus", rbacMiddleware.RequirePermission("admin.role.view"), roleHandler.GetMenus)
		roles.PUT("/:id/menus", rbacMiddleware.RequirePermission("admin.role.update"), roleHandler.AssignMenus)
	}

	// Menus
	menus := protected.Group("/menus")
	{
		menus.GET("", rbacMiddleware.RequirePermission("admin.menu.view"), menuHandler.List)
		menus.GET("/tree", rbacMiddleware.RequirePermission("admin.menu.view"), menuHandler.GetTree)
		menus.POST("", rbacMiddleware.RequirePermission("admin.menu.create"), menuHandler.Create)
		menus.GET("/:id", rbacMiddleware.RequirePermission("admin.menu.view"), menuHandler.Get)
		menus.PUT("/:id", rbacMiddleware.RequirePermission("admin.menu.update"), menuHandler.Update)
		menus.DELETE("/:id", rbacMiddleware.RequirePermission("admin.menu.delete"), menuHandler.Delete)
	}

	// Devices (management)
	devices := protected.Group("/devices")
	{
		devices.GET("", rbacMiddleware.RequirePermission("admin.device.view"), deviceHandler.List)
		devices.GET("/:id", rbacMiddleware.RequirePermission("admin.device.view"), deviceHandler.Get)
		devices.PUT("/:id", rbacMiddleware.RequirePermission("admin.device.update"), deviceHandler.Update)
		devices.DELETE("/:id", rbacMiddleware.RequirePermission("admin.device.delete"), deviceHandler.Deactivate)
		devices.PUT("/:id/config", rbacMiddleware.RequirePermission("admin.device.update"), deviceHandler.UpdateConfig)
		devices.GET("/:id/sessions", rbacMiddleware.RequirePermission("admin.device.view"), deviceHandler.GetSessions)
	}
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
