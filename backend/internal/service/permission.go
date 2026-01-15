package service

import (
	"context"
	"fmt"

	"gebase/internal/domain"
	"gebase/internal/repository"
)

type PermissionService struct {
	permissionRepo *repository.PermissionRepository
	moduleRepo     *repository.ModuleRepository
	actionRepo     *repository.ActionRepository
	systemRepo     *repository.SystemRepository
}

func NewPermissionService(
	permissionRepo *repository.PermissionRepository,
	moduleRepo *repository.ModuleRepository,
	actionRepo *repository.ActionRepository,
	systemRepo *repository.SystemRepository,
) *PermissionService {
	return &PermissionService{
		permissionRepo: permissionRepo,
		moduleRepo:     moduleRepo,
		actionRepo:     actionRepo,
		systemRepo:     systemRepo,
	}
}

// CheckPermission checks if user has a specific permission
func (s *PermissionService) CheckPermission(ctx context.Context, userID int64, systemID *int, permissionCode string) (bool, error) {
	return s.permissionRepo.CheckUserPermission(ctx, userID, systemID, permissionCode)
}

// GetUserPermissions returns all permissions for a user in a system
func (s *PermissionService) GetUserPermissions(ctx context.Context, userID int64, systemID *int) ([]domain.Permission, error) {
	return s.permissionRepo.FindUserPermissions(ctx, userID, systemID)
}

// GetUserPermissionCodes returns permission codes for a user in a system
func (s *PermissionService) GetUserPermissionCodes(ctx context.Context, userID int64, systemID *int) ([]string, error) {
	permissions, err := s.permissionRepo.FindUserPermissions(ctx, userID, systemID)
	if err != nil {
		return nil, err
	}

	codes := make([]string, len(permissions))
	for i, p := range permissions {
		codes[i] = p.Code
	}
	return codes, nil
}

// GetRolePermissions returns all permissions for a role
func (s *PermissionService) GetRolePermissions(ctx context.Context, roleID int) ([]domain.Permission, error) {
	return s.permissionRepo.FindByRoleID(ctx, roleID)
}

// GetSystemPermissions returns all permissions for a system
func (s *PermissionService) GetSystemPermissions(ctx context.Context, systemID int) ([]domain.Permission, error) {
	return s.permissionRepo.FindBySystemID(ctx, systemID)
}

// GetPlatformPermissions returns all platform-level permissions
func (s *PermissionService) GetPlatformPermissions(ctx context.Context) ([]domain.Permission, error) {
	return s.permissionRepo.FindPlatformPermissions(ctx)
}

// SyncPermissions generates permissions from system.module.action combinations
func (s *PermissionService) SyncPermissions(ctx context.Context, systemID *int, createdBy int64) error {
	var modules []domain.Module
	var err error

	if systemID != nil {
		// Get modules for specific system
		modules, err = s.moduleRepo.FindBySystemID(ctx, *systemID)
	} else {
		// Get platform modules (system_id is NULL)
		modules, err = s.moduleRepo.FindPlatformModules(ctx)
	}
	if err != nil {
		return err
	}

	// Get system code if applicable
	var systemCode string
	if systemID != nil {
		system, err := s.systemRepo.FindByID(ctx, *systemID)
		if err != nil {
			return err
		}
		systemCode = system.Code
	}

	// Generate permissions
	for _, module := range modules {
		moduleWithActions, err := s.moduleRepo.FindWithActions(ctx, module.ID)
		if err != nil {
			continue
		}

		for _, ma := range moduleWithActions.ModuleActions {
			if ma.Action == nil {
				continue
			}

			var code string
			var name string
			if systemCode != "" {
				// System permission: {system.module.action}
				code = fmt.Sprintf("%s.%s.%s", systemCode, module.Code, ma.Action.Code)
				name = fmt.Sprintf("%s - %s - %s", systemCode, module.Name, ma.Action.Name)
			} else {
				// Platform permission: {module.action}
				code = fmt.Sprintf("%s.%s", module.Code, ma.Action.Code)
				name = fmt.Sprintf("%s - %s", module.Name, ma.Action.Name)
			}

			// Check if permission already exists
			existing, err := s.permissionRepo.FindByCode(ctx, code)
			if err == nil && existing != nil {
				continue // Skip existing permissions
			}

			// Create new permission
			permission := &domain.Permission{
				Code:     code,
				Name:     name,
				ModuleID: module.ID,
				IsActive: domain.Ptr(true),
			}

			if systemID != nil {
				permission.SystemID = systemID
			}

			actionID := int64(ma.ActionID)
			permission.ActionID = &actionID
			permission.CreatedBy = &createdBy

			if err := s.permissionRepo.Create(ctx, permission); err != nil {
				// Log error but continue
				continue
			}
		}
	}

	return nil
}

// GeneratePermissionCode generates permission code from components
func GeneratePermissionCode(systemCode, moduleCode, actionCode string) string {
	if systemCode == "" {
		return fmt.Sprintf("%s.%s", moduleCode, actionCode)
	}
	return fmt.Sprintf("%s.%s.%s", systemCode, moduleCode, actionCode)
}

// ParsePermissionCode parses permission code into components
func ParsePermissionCode(code string) (systemCode, moduleCode, actionCode string) {
	parts := splitPermissionCode(code)
	switch len(parts) {
	case 2:
		return "", parts[0], parts[1]
	case 3:
		return parts[0], parts[1], parts[2]
	default:
		return "", "", ""
	}
}

func splitPermissionCode(code string) []string {
	var parts []string
	current := ""
	for _, c := range code {
		if c == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
