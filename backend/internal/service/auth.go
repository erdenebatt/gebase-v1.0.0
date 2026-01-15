package service

import (
	"context"
	"errors"

	"gebase/internal/auth"
	"gebase/internal/domain"
	"gebase/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user is inactive")
	ErrDeviceNotFound     = errors.New("device not found")
	ErrDeviceNotActive    = errors.New("device is not active")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session has expired")
	ErrSystemNotFound     = errors.New("system not found")
	ErrNoSystemAccess     = errors.New("user does not have access to this system")
)

type AuthService struct {
	userRepo          *repository.UserRepository
	deviceRepo        *repository.DeviceRepository
	systemRepo        *repository.SystemRepository
	userSystemRoleRepo *repository.UserSystemRoleRepository
	roleRepo          *repository.RoleRepository
	permissionRepo    *repository.PermissionRepository
	menuRepo          *repository.MenuRepository
	jwtService        *auth.JWTService
	sessionService    *auth.SessionService
}

func NewAuthService(
	userRepo *repository.UserRepository,
	deviceRepo *repository.DeviceRepository,
	systemRepo *repository.SystemRepository,
	userSystemRoleRepo *repository.UserSystemRoleRepository,
	roleRepo *repository.RoleRepository,
	permissionRepo *repository.PermissionRepository,
	menuRepo *repository.MenuRepository,
	jwtService *auth.JWTService,
	sessionService *auth.SessionService,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		deviceRepo:        deviceRepo,
		systemRepo:        systemRepo,
		userSystemRoleRepo: userSystemRoleRepo,
		roleRepo:          roleRepo,
		permissionRepo:    permissionRepo,
		menuRepo:          menuRepo,
		jwtService:        jwtService,
		sessionService:    sessionService,
	}
}

type LoginRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	DeviceUID string `json:"device_uid" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	TokenType    string          `json:"token_type"`
	ExpiresIn    int64           `json:"expires_in"`
	User         *domain.User    `json:"user"`
	Systems      []domain.System `json:"available_systems"`
}

// Login authenticates user and returns platform token
func (s *AuthService) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if user.IsActive == nil || !*user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify device
	device, err := s.deviceRepo.FindByUID(ctx, req.DeviceUID)
	if err != nil {
		return nil, ErrDeviceNotFound
	}

	if device.IsActive == nil || !*device.IsActive {
		return nil, ErrDeviceNotActive
	}

	// Create session
	session, err := s.sessionService.CreateSession(ctx, user.ID, device.ID, ipAddress, userAgent, user.OrganizationID)
	if err != nil {
		return nil, err
	}

	// Generate platform token
	accessToken, err := s.jwtService.GeneratePlatformToken(user, session)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.jwtService.GenerateRefreshToken(user, session)
	if err != nil {
		return nil, err
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Get available systems
	userWithRoles, err := s.userRepo.FindWithRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	systems := userWithRoles.GetAvailableSystems()

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400, // 24 hours in seconds
		User:         user,
		Systems:      systems,
	}, nil
}

type SwitchSystemRequest struct {
	SystemCode string `json:"system_code" binding:"required"`
}

type SwitchSystemResponse struct {
	SystemToken         string         `json:"system_token"`
	TokenType           string         `json:"token_type"`
	ExpiresIn           int64          `json:"expires_in"`
	CurrentSystem       *domain.System `json:"current_system"`
	CurrentRole         *domain.Role   `json:"current_role"`
	CurrentOrganization *domain.Organization `json:"current_organization,omitempty"`
	Permissions         []string       `json:"permissions"`
	Menus               []domain.Menu  `json:"menus"`
}

// SwitchSystem switches to a specific system and returns system token
func (s *AuthService) SwitchSystem(ctx context.Context, claims *auth.Claims, req *SwitchSystemRequest, ipAddress string) (*SwitchSystemResponse, error) {
	// Get system
	system, err := s.systemRepo.FindByCode(ctx, req.SystemCode)
	if err != nil {
		return nil, ErrSystemNotFound
	}

	// Check user access to system
	userSystemRoles, err := s.userSystemRoleRepo.FindByUserAndSystem(ctx, claims.UserID, &system.ID)
	if err != nil || len(userSystemRoles) == 0 {
		return nil, ErrNoSystemAccess
	}

	// Get role IDs and first role
	roleIDs := make([]int, len(userSystemRoles))
	for i, ur := range userSystemRoles {
		roleIDs[i] = ur.RoleID
	}

	// Get the first role details
	firstRole, err := s.roleRepo.FindByID(ctx, roleIDs[0])
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Get session
	session, err := s.sessionService.GetSessionByID(ctx, claims.SessionID)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Switch system in session
	if err := s.sessionService.SwitchSystem(ctx, session.ID, system.ID, ipAddress); err != nil {
		return nil, err
	}

	// Generate system token
	accessToken, err := s.jwtService.GenerateSystemToken(user, session, system, roleIDs)
	if err != nil {
		return nil, err
	}

	// Get user permissions for this system
	permissions, err := s.permissionRepo.FindUserPermissions(ctx, claims.UserID, &system.ID)
	if err != nil {
		permissions = []domain.Permission{}
	}

	// Convert permissions to string array
	permissionCodes := make([]string, len(permissions))
	for i, p := range permissions {
		permissionCodes[i] = p.Code
	}

	// Get user menus for this system
	menus, err := s.menuRepo.FindUserMenus(ctx, claims.UserID, system.ID)
	if err != nil {
		menus = []domain.Menu{}
	}

	// Build menu tree
	menuTree := s.menuRepo.BuildMenuTree(menus)

	return &SwitchSystemResponse{
		SystemToken:   accessToken,
		TokenType:     "Bearer",
		ExpiresIn:     28800, // 8 hours in seconds
		CurrentSystem: system,
		CurrentRole:   firstRole,
		Permissions:   permissionCodes,
		Menus:         menuTree,
	}, nil
}

// RefreshToken refreshes the access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Get session
	session, err := s.sessionService.GetSessionByID(ctx, claims.SessionID)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Check session validity
	if !s.sessionService.IsSessionValid(ctx, session) {
		return nil, ErrSessionExpired
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new tokens
	accessToken, err := s.jwtService.GeneratePlatformToken(user, session)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user, session)
	if err != nil {
		return nil, err
	}

	// Get available systems
	userWithRoles, err := s.userRepo.FindWithRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	systems := userWithRoles.GetAvailableSystems()

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400,
		User:         user,
		Systems:      systems,
	}, nil
}

// Logout terminates the current session
func (s *AuthService) Logout(ctx context.Context, sessionID int64) error {
	return s.sessionService.Logout(ctx, sessionID, "user")
}

// RemoteLogout terminates sessions remotely
func (s *AuthService) RemoteLogout(ctx context.Context, userID int64, deviceID *int64, reason string) error {
	if deviceID != nil {
		return s.sessionService.LogoutDevice(ctx, *deviceID, reason)
	}
	return s.sessionService.LogoutUser(ctx, userID, reason)
}

// ValidateToken validates a token and returns claims
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
	return s.jwtService.ValidateToken(token)
}

// GetCurrentUser returns the current user from claims
func (s *AuthService) GetCurrentUser(ctx context.Context, userID int64) (*domain.User, error) {
	return s.userRepo.FindWithRoles(ctx, userID)
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
