package repository

import (
	"context"

	"gebase/internal/domain"

	"gorm.io/gorm"
)

type RoleRepository struct {
	*BaseRepository[domain.Role]
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		BaseRepository: NewBaseRepository[domain.Role](db),
	}
}

func (r *RoleRepository) FindByCode(ctx context.Context, code string) (*domain.Role, error) {
	var role domain.Role
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindBySystemID(ctx context.Context, systemID int) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.DB.WithContext(ctx).Where("system_id = ?", systemID).Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) FindPlatformRoles(ctx context.Context) ([]domain.Role, error) {
	var roles []domain.Role
	err := r.DB.WithContext(ctx).Where("system_id IS NULL").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) FindWithPermissions(ctx context.Context, id int) (*domain.Role, error) {
	var role domain.Role
	err := r.DB.WithContext(ctx).
		Preload("Permissions").
		Preload("Permissions.Permission").
		First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindWithMenus(ctx context.Context, id int) (*domain.Role, error) {
	var role domain.Role
	err := r.DB.WithContext(ctx).
		Preload("Menus").
		Preload("Menus.Menu").
		First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

type RolePermissionRepository struct {
	*BaseRepository[domain.RolePermission]
}

func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{
		BaseRepository: NewBaseRepository[domain.RolePermission](db),
	}
}

func (r *RolePermissionRepository) DeleteByRoleID(ctx context.Context, roleID int) error {
	return r.DB.WithContext(ctx).Where("role_id = ?", roleID).Delete(&domain.RolePermission{}).Error
}

func (r *RolePermissionRepository) AssignPermissions(ctx context.Context, roleID int, permissionIDs []int, createdBy int64) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing
		if err := tx.Where("role_id = ?", roleID).Delete(&domain.RolePermission{}).Error; err != nil {
			return err
		}

		// Create new
		for _, permID := range permissionIDs {
			rp := domain.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
			}
			rp.CreatedBy = &createdBy
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

type RoleMenuRepository struct {
	*BaseRepository[domain.RoleMenu]
}

func NewRoleMenuRepository(db *gorm.DB) *RoleMenuRepository {
	return &RoleMenuRepository{
		BaseRepository: NewBaseRepository[domain.RoleMenu](db),
	}
}

func (r *RoleMenuRepository) DeleteByRoleID(ctx context.Context, roleID int) error {
	return r.DB.WithContext(ctx).Where("role_id = ?", roleID).Delete(&domain.RoleMenu{}).Error
}

func (r *RoleMenuRepository) AssignMenus(ctx context.Context, roleID int, menuIDs []int, createdBy int64) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&domain.RoleMenu{}).Error; err != nil {
			return err
		}

		for _, menuID := range menuIDs {
			rm := domain.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			}
			rm.CreatedBy = &createdBy
			if err := tx.Create(&rm).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

type UserSystemRoleRepository struct {
	*BaseRepository[domain.UserSystemRole]
}

func NewUserSystemRoleRepository(db *gorm.DB) *UserSystemRoleRepository {
	return &UserSystemRoleRepository{
		BaseRepository: NewBaseRepository[domain.UserSystemRole](db),
	}
}

func (r *UserSystemRoleRepository) FindByUserID(ctx context.Context, userID int64) ([]domain.UserSystemRole, error) {
	var roles []domain.UserSystemRole
	err := r.DB.WithContext(ctx).
		Preload("Role").
		Preload("System").
		Preload("Organization").
		Where("user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *UserSystemRoleRepository) FindByUserAndSystem(ctx context.Context, userID int64, systemID *int) ([]domain.UserSystemRole, error) {
	var roles []domain.UserSystemRole
	query := r.DB.WithContext(ctx).
		Preload("Role").
		Where("user_id = ?", userID)

	if systemID != nil {
		query = query.Where("system_id = ?", *systemID)
	} else {
		query = query.Where("system_id IS NULL")
	}

	err := query.Find(&roles).Error
	return roles, err
}

func (r *UserSystemRoleRepository) AssignRoles(ctx context.Context, userID int64, systemID *int, roleIDs []int, orgID *int64, createdBy int64) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		query := tx.Where("user_id = ?", userID)
		if systemID != nil {
			query = query.Where("system_id = ?", *systemID)
		} else {
			query = query.Where("system_id IS NULL")
		}
		if err := query.Delete(&domain.UserSystemRole{}).Error; err != nil {
			return err
		}

		for _, roleID := range roleIDs {
			usr := domain.UserSystemRole{
				UserID:         userID,
				SystemID:       systemID,
				RoleID:         roleID,
				OrganizationID: orgID,
				IsActive:       domain.Ptr(true),
			}
			usr.CreatedBy = &createdBy
			if err := tx.Create(&usr).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
