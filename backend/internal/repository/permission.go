package repository

import (
	"context"

	"gebase/internal/domain"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	*BaseRepository[domain.Permission]
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		BaseRepository: NewBaseRepository[domain.Permission](db),
	}
}

func (r *PermissionRepository) FindByCode(ctx context.Context, code string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) FindBySystemID(ctx context.Context, systemID int) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.DB.WithContext(ctx).Where("system_id = ?", systemID).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) FindPlatformPermissions(ctx context.Context) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.DB.WithContext(ctx).Where("system_id IS NULL").Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) FindByRoleID(ctx context.Context, roleID int) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.DB.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) FindUserPermissions(ctx context.Context, userID int64, systemID *int) ([]domain.Permission, error) {
	var permissions []domain.Permission

	query := r.DB.WithContext(ctx).
		Distinct("permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_system_roles ON user_system_roles.role_id = role_permissions.role_id").
		Where("user_system_roles.user_id = ? AND user_system_roles.is_active = true", userID)

	if systemID != nil {
		query = query.Where("(user_system_roles.system_id = ? OR user_system_roles.system_id IS NULL)", *systemID)
	} else {
		query = query.Where("user_system_roles.system_id IS NULL")
	}

	err := query.Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) CheckUserPermission(ctx context.Context, userID int64, systemID *int, permissionCode string) (bool, error) {
	var count int64

	query := r.DB.WithContext(ctx).
		Model(&domain.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_system_roles ON user_system_roles.role_id = role_permissions.role_id").
		Where("user_system_roles.user_id = ? AND user_system_roles.is_active = true", userID).
		Where("permissions.code = ?", permissionCode)

	if systemID != nil {
		query = query.Where("(user_system_roles.system_id = ? OR user_system_roles.system_id IS NULL)", *systemID)
	} else {
		query = query.Where("user_system_roles.system_id IS NULL")
	}

	err := query.Count(&count).Error
	return count > 0, err
}
