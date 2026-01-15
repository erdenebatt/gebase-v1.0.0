package service

import (
	"context"
	"errors"

	"gebase/internal/domain"
	"gebase/internal/repository"
)

var (
	ErrRoleNotFound   = errors.New("role not found")
	ErrRoleCodeExists = errors.New("role code already exists")
)

type RoleService struct {
	roleRepo           *repository.RoleRepository
	rolePermissionRepo *repository.RolePermissionRepository
	roleMenuRepo       *repository.RoleMenuRepository
}

func NewRoleService(
	roleRepo *repository.RoleRepository,
	rolePermissionRepo *repository.RolePermissionRepository,
	roleMenuRepo *repository.RoleMenuRepository,
) *RoleService {
	return &RoleService{
		roleRepo:           roleRepo,
		rolePermissionRepo: rolePermissionRepo,
		roleMenuRepo:       roleMenuRepo,
	}
}

type CreateRoleRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SystemID    *int   `json:"system_id"`
	IsSystem    bool   `json:"is_system"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// ListRoles returns all roles
func (s *RoleService) ListRoles(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[domain.Role], error) {
	return s.roleRepo.FindWithPagination(ctx, repository.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	})
}

// ListRolesBySystem returns roles for a specific system
func (s *RoleService) ListRolesBySystem(ctx context.Context, systemID int) ([]domain.Role, error) {
	return s.roleRepo.FindBySystemID(ctx, systemID)
}

// ListPlatformRoles returns platform-level roles
func (s *RoleService) ListPlatformRoles(ctx context.Context) ([]domain.Role, error) {
	return s.roleRepo.FindPlatformRoles(ctx)
}

// GetRole returns role by ID
func (s *RoleService) GetRole(ctx context.Context, id int) (*domain.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

// GetRoleByCode returns role by code
func (s *RoleService) GetRoleByCode(ctx context.Context, code string) (*domain.Role, error) {
	return s.roleRepo.FindByCode(ctx, code)
}

// GetRoleWithPermissions returns role with its permissions
func (s *RoleService) GetRoleWithPermissions(ctx context.Context, id int) (*domain.Role, error) {
	return s.roleRepo.FindWithPermissions(ctx, id)
}

// GetRoleWithMenus returns role with its menus
func (s *RoleService) GetRoleWithMenus(ctx context.Context, id int) (*domain.Role, error) {
	return s.roleRepo.FindWithMenus(ctx, id)
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(ctx context.Context, req *CreateRoleRequest, createdBy int64) (*domain.Role, error) {
	// Check code uniqueness
	existing, _ := s.roleRepo.FindByCode(ctx, req.Code)
	if existing != nil {
		return nil, ErrRoleCodeExists
	}

	role := &domain.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		SystemID:    req.SystemID,
		IsSystem:    domain.Ptr(req.IsSystem),
		IsActive:    domain.Ptr(true),
	}
	role.CreatedBy = &createdBy

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// UpdateRole updates a role
func (s *RoleService) UpdateRole(ctx context.Context, id int, req *UpdateRoleRequest, updatedBy int64) (*domain.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.IsActive != nil {
		role.IsActive = req.IsActive
	}

	role.UpdatedBy = &updatedBy

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// DeleteRole soft deletes a role
func (s *RoleService) DeleteRole(ctx context.Context, id int, deletedBy int64) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return ErrRoleNotFound
	}

	// Don't allow deletion of system roles
	if role.IsSystem != nil && *role.IsSystem {
		return errors.New("cannot delete system role")
	}

	role.DeletedBy = &deletedBy
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return err
	}

	return s.roleRepo.Delete(ctx, id)
}

// AssignPermissions assigns permissions to a role
func (s *RoleService) AssignPermissions(ctx context.Context, roleID int, permissionIDs []int, assignedBy int64) error {
	return s.rolePermissionRepo.AssignPermissions(ctx, roleID, permissionIDs, assignedBy)
}

// AssignMenus assigns menus to a role
func (s *RoleService) AssignMenus(ctx context.Context, roleID int, menuIDs []int, assignedBy int64) error {
	return s.roleMenuRepo.AssignMenus(ctx, roleID, menuIDs, assignedBy)
}

// GetRolePermissionIDs returns permission IDs for a role
func (s *RoleService) GetRolePermissionIDs(ctx context.Context, roleID int) ([]int, error) {
	role, err := s.roleRepo.FindWithPermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(role.Permissions))
	for i, rp := range role.Permissions {
		ids[i] = rp.PermissionID
	}
	return ids, nil
}

// GetRoleMenuIDs returns menu IDs for a role
func (s *RoleService) GetRoleMenuIDs(ctx context.Context, roleID int) ([]int, error) {
	role, err := s.roleRepo.FindWithMenus(ctx, roleID)
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(role.Menus))
	for i, rm := range role.Menus {
		ids[i] = rm.MenuID
	}
	return ids, nil
}
