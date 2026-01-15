package repository

import (
	"context"

	"gebase/internal/domain"

	"gorm.io/gorm"
)

type ModuleRepository struct {
	*BaseRepository[domain.Module]
}

func NewModuleRepository(db *gorm.DB) *ModuleRepository {
	return &ModuleRepository{
		BaseRepository: NewBaseRepository[domain.Module](db),
	}
}

func (r *ModuleRepository) FindByCode(ctx context.Context, code string) (*domain.Module, error) {
	var module domain.Module
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *ModuleRepository) FindBySystemID(ctx context.Context, systemID int) ([]domain.Module, error) {
	var modules []domain.Module
	err := r.DB.WithContext(ctx).Where("system_id = ?", systemID).Find(&modules).Error
	return modules, err
}

func (r *ModuleRepository) FindBySystemAndCode(ctx context.Context, systemID int, code string) (*domain.Module, error) {
	var module domain.Module
	err := r.DB.WithContext(ctx).
		Where("system_id = ? AND code = ?", systemID, code).
		First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *ModuleRepository) FindPlatformModules(ctx context.Context) ([]domain.Module, error) {
	var modules []domain.Module
	err := r.DB.WithContext(ctx).Where("system_id IS NULL").Find(&modules).Error
	return modules, err
}

func (r *ModuleRepository) FindWithActions(ctx context.Context, id int) (*domain.Module, error) {
	var module domain.Module
	err := r.DB.WithContext(ctx).
		Preload("ModuleActions").
		Preload("ModuleActions.Action").
		First(&module, id).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

type ActionRepository struct {
	*BaseRepository[domain.Action]
}

func NewActionRepository(db *gorm.DB) *ActionRepository {
	return &ActionRepository{
		BaseRepository: NewBaseRepository[domain.Action](db),
	}
}

func (r *ActionRepository) FindByCode(ctx context.Context, code string) (*domain.Action, error) {
	var action domain.Action
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&action).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func (r *ActionRepository) FindActive(ctx context.Context) ([]domain.Action, error) {
	var actions []domain.Action
	err := r.DB.WithContext(ctx).Where("is_active = ?", true).Find(&actions).Error
	return actions, err
}

type ModuleActionRepository struct {
	*BaseRepository[domain.ModuleAction]
}

func NewModuleActionRepository(db *gorm.DB) *ModuleActionRepository {
	return &ModuleActionRepository{
		BaseRepository: NewBaseRepository[domain.ModuleAction](db),
	}
}

func (r *ModuleActionRepository) FindByModule(ctx context.Context, moduleID int) ([]domain.ModuleAction, error) {
	var moduleActions []domain.ModuleAction
	err := r.DB.WithContext(ctx).
		Preload("Action").
		Where("module_id = ?", moduleID).
		Find(&moduleActions).Error
	return moduleActions, err
}
