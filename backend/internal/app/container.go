package app

import (
	"gebase/internal/auth"
	"gebase/internal/config"
	"gebase/internal/http/handlers"
	"gebase/internal/http/router"
	"gebase/internal/middleware"
	"gebase/internal/repository"
	"gebase/internal/service"

	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config
	DB     *gorm.DB

	// Repositories
	UserRepo              *repository.UserRepository
	OrganizationRepo      *repository.OrganizationRepository
	OrganizationTypeRepo  *repository.OrganizationTypeRepository
	OrganizationSystemRepo *repository.OrganizationSystemRepository
	SystemRepo            *repository.SystemRepository
	ModuleRepo            *repository.ModuleRepository
	ActionRepo            *repository.ActionRepository
	ModuleActionRepo      *repository.ModuleActionRepository
	PermissionRepo        *repository.PermissionRepository
	RoleRepo              *repository.RoleRepository
	RolePermissionRepo    *repository.RolePermissionRepository
	RoleMenuRepo          *repository.RoleMenuRepository
	UserSystemRoleRepo    *repository.UserSystemRoleRepository
	MenuRepo              *repository.MenuRepository
	DeviceRepo            *repository.DeviceRepository
	SessionRepo           *repository.SessionRepository
	SessionHistoryRepo    *repository.SessionSystemHistoryRepository
	LanguageRepo          *repository.LanguageRepository
	TranslationRepo       *repository.TranslationRepository

	// Auth
	JWTService     *auth.JWTService
	SessionService *auth.SessionService

	// Services
	AuthService       *service.AuthService
	UserService       *service.UserService
	OrganizationService *service.OrganizationService
	SystemService     *service.SystemService
	RoleService       *service.RoleService
	PermissionService *service.PermissionService
	MenuService       *service.MenuService
	DeviceService     *service.DeviceService

	// Middleware
	AuthMiddleware   *middleware.AuthMiddleware
	RBACMiddleware   *middleware.RBACMiddleware
	DeviceMiddleware *middleware.DeviceMiddleware

	// Handlers
	AuthHandler   *handlers.AuthHandler
	UserHandler   *handlers.UserHandler
	OrgHandler    *handlers.OrganizationHandler
	SystemHandler *handlers.SystemHandler
	RoleHandler   *handlers.RoleHandler
	MenuHandler   *handlers.MenuHandler
	DeviceHandler *handlers.DeviceHandler

	// Router
	Router *router.Router
}

// NewContainer creates a new dependency injection container
func NewContainer(cfg *config.Config, db *gorm.DB) *Container {
	c := &Container{
		Config: cfg,
		DB:     db,
	}

	c.initRepositories()
	c.initAuth()
	c.initServices()
	c.initMiddleware()
	c.initHandlers()
	c.initRouter()

	return c
}

func (c *Container) initRepositories() {
	c.UserRepo = repository.NewUserRepository(c.DB)
	c.OrganizationRepo = repository.NewOrganizationRepository(c.DB)
	c.OrganizationTypeRepo = repository.NewOrganizationTypeRepository(c.DB)
	c.OrganizationSystemRepo = repository.NewOrganizationSystemRepository(c.DB)
	c.SystemRepo = repository.NewSystemRepository(c.DB)
	c.ModuleRepo = repository.NewModuleRepository(c.DB)
	c.ActionRepo = repository.NewActionRepository(c.DB)
	c.ModuleActionRepo = repository.NewModuleActionRepository(c.DB)
	c.PermissionRepo = repository.NewPermissionRepository(c.DB)
	c.RoleRepo = repository.NewRoleRepository(c.DB)
	c.RolePermissionRepo = repository.NewRolePermissionRepository(c.DB)
	c.RoleMenuRepo = repository.NewRoleMenuRepository(c.DB)
	c.UserSystemRoleRepo = repository.NewUserSystemRoleRepository(c.DB)
	c.MenuRepo = repository.NewMenuRepository(c.DB)
	c.DeviceRepo = repository.NewDeviceRepository(c.DB)
	c.SessionRepo = repository.NewSessionRepository(c.DB)
	c.SessionHistoryRepo = repository.NewSessionSystemHistoryRepository(c.DB)
	c.LanguageRepo = repository.NewLanguageRepository(c.DB)
	c.TranslationRepo = repository.NewTranslationRepository(c.DB)
}

func (c *Container) initAuth() {
	c.JWTService = auth.NewJWTService(c.Config)
	c.SessionService = auth.NewSessionService(c.Config, c.SessionRepo, c.SessionHistoryRepo)
}

func (c *Container) initServices() {
	c.AuthService = service.NewAuthService(
		c.UserRepo,
		c.DeviceRepo,
		c.SystemRepo,
		c.UserSystemRoleRepo,
		c.RoleRepo,
		c.PermissionRepo,
		c.MenuRepo,
		c.JWTService,
		c.SessionService,
	)
	c.UserService = service.NewUserService(c.UserRepo, c.UserSystemRoleRepo, c.SessionRepo)
	c.OrganizationService = service.NewOrganizationService(c.OrganizationRepo, c.OrganizationTypeRepo, c.OrganizationSystemRepo)
	c.SystemService = service.NewSystemService(c.SystemRepo, c.ModuleRepo, c.MenuRepo)
	c.RoleService = service.NewRoleService(c.RoleRepo, c.RolePermissionRepo, c.RoleMenuRepo)
	c.PermissionService = service.NewPermissionService(c.PermissionRepo, c.ModuleRepo, c.ActionRepo, c.SystemRepo)
	c.MenuService = service.NewMenuService(c.MenuRepo)
	c.DeviceService = service.NewDeviceService(c.DeviceRepo, c.SessionRepo)
}

func (c *Container) initMiddleware() {
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.JWTService, c.SessionService)
	c.RBACMiddleware = middleware.NewRBACMiddleware(c.PermissionService)
	c.DeviceMiddleware = middleware.NewDeviceMiddleware(c.DeviceService)
}

func (c *Container) initHandlers() {
	c.AuthHandler = handlers.NewAuthHandler(c.AuthService, c.PermissionService, c.MenuService)
	c.UserHandler = handlers.NewUserHandler(c.UserService)
	c.OrgHandler = handlers.NewOrganizationHandler(c.OrganizationService)
	c.SystemHandler = handlers.NewSystemHandler(c.SystemService)
	c.RoleHandler = handlers.NewRoleHandler(c.RoleService)
	c.MenuHandler = handlers.NewMenuHandler(c.MenuService, c.SystemService)
	c.DeviceHandler = handlers.NewDeviceHandler(c.DeviceService)
}

func (c *Container) initRouter() {
	c.Router = router.NewRouter(c.Config)
	c.Router.Setup(
		c.AuthMiddleware,
		c.RBACMiddleware,
		c.DeviceMiddleware,
		c.AuthHandler,
		c.DeviceHandler,
		c.UserHandler,
		c.OrgHandler,
		c.SystemHandler,
		c.RoleHandler,
		c.MenuHandler,
	)
}
