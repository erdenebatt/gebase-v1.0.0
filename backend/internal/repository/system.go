package repository

import (
	"context"

	"gebase/internal/domain"

	"gorm.io/gorm"
)

type SystemRepository struct {
	*BaseRepository[domain.System]
}

func NewSystemRepository(db *gorm.DB) *SystemRepository {
	return &SystemRepository{
		BaseRepository: NewBaseRepository[domain.System](db),
	}
}

func (r *SystemRepository) FindByCode(ctx context.Context, code string) (*domain.System, error) {
	var system domain.System
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&system).Error
	if err != nil {
		return nil, err
	}
	return &system, nil
}

func (r *SystemRepository) FindWithModules(ctx context.Context, id int) (*domain.System, error) {
	var system domain.System
	err := r.DB.WithContext(ctx).
		Preload("Modules").
		Preload("Modules.ModuleActions").
		Preload("Modules.ModuleActions.Action").
		First(&system, id).Error
	if err != nil {
		return nil, err
	}
	return &system, nil
}

func (r *SystemRepository) FindWithMenus(ctx context.Context, id int) (*domain.System, error) {
	var system domain.System
	err := r.DB.WithContext(ctx).
		Preload("Menus", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		First(&system, id).Error
	if err != nil {
		return nil, err
	}
	return &system, nil
}

func (r *SystemRepository) FindActive(ctx context.Context) ([]domain.System, error) {
	var systems []domain.System
	err := r.DB.WithContext(ctx).Where("is_active = ?", true).Order("sequence").Find(&systems).Error
	return systems, err
}
