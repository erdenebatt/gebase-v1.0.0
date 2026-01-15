package repository

import (
	"context"

	"gebase/internal/domain"

	"gorm.io/gorm"
)

type MenuRepository struct {
	*BaseRepository[domain.Menu]
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{
		BaseRepository: NewBaseRepository[domain.Menu](db),
	}
}

func (r *MenuRepository) FindByCode(ctx context.Context, code string) (*domain.Menu, error) {
	var menu domain.Menu
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepository) FindBySystemID(ctx context.Context, systemID int) ([]domain.Menu, error) {
	var menus []domain.Menu
	err := r.DB.WithContext(ctx).
		Where("system_id = ?", systemID).
		Order("sequence ASC").
		Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) FindBySystemAndCode(ctx context.Context, systemID int, code string) (*domain.Menu, error) {
	var menu domain.Menu
	err := r.DB.WithContext(ctx).
		Where("system_id = ? AND code = ?", systemID, code).
		First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepository) FindRootMenus(ctx context.Context, systemID int) ([]domain.Menu, error) {
	var menus []domain.Menu
	err := r.DB.WithContext(ctx).
		Where("system_id = ? AND parent_id IS NULL", systemID).
		Order("sequence ASC").
		Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) FindChildren(ctx context.Context, parentID int) ([]domain.Menu, error) {
	var menus []domain.Menu
	err := r.DB.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("sequence ASC").
		Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) FindMenuTree(ctx context.Context, systemID int) ([]domain.Menu, error) {
	var menus []domain.Menu
	err := r.DB.WithContext(ctx).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		Where("system_id = ? AND parent_id IS NULL", systemID).
		Order("sequence ASC").
		Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) FindUserMenus(ctx context.Context, userID int64, systemID int) ([]domain.Menu, error) {
	var menus []domain.Menu
	err := r.DB.WithContext(ctx).
		Distinct("menus.*").
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Joins("JOIN user_system_roles ON user_system_roles.role_id = role_menus.role_id").
		Where("user_system_roles.user_id = ? AND user_system_roles.is_active = true", userID).
		Where("menus.system_id = ? AND menus.is_visible = true AND menus.is_active = true", systemID).
		Order("menus.sequence ASC").
		Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) BuildMenuTree(menus []domain.Menu) []domain.Menu {
	if len(menus) == 0 {
		return []domain.Menu{}
	}

	// Create map with copies (to avoid modifying original)
	menuMap := make(map[int]*domain.Menu)
	for i := range menus {
		menu := menus[i]
		menu.Children = nil // Reset children
		menuMap[menu.ID] = &menu
	}

	// Build parent-child relationships
	var rootMenus []*domain.Menu
	for _, menu := range menuMap {
		if menu.ParentID == nil {
			rootMenus = append(rootMenus, menu)
		} else if parent, ok := menuMap[*menu.ParentID]; ok {
			parent.Children = append(parent.Children, *menu)
		}
	}

	// Sort root menus by sequence
	sortMenusBySequence(rootMenus)

	// Convert to slice of values and sort children
	result := make([]domain.Menu, len(rootMenus))
	for i, menu := range rootMenus {
		sortChildrenBySequence(menu)
		result[i] = *menu
	}

	return result
}

func sortMenusBySequence(menus []*domain.Menu) {
	for i := 0; i < len(menus)-1; i++ {
		for j := i + 1; j < len(menus); j++ {
			if menus[i].Sequence > menus[j].Sequence {
				menus[i], menus[j] = menus[j], menus[i]
			}
		}
	}
}

func sortChildrenBySequence(menu *domain.Menu) {
	if len(menu.Children) == 0 {
		return
	}
	for i := 0; i < len(menu.Children)-1; i++ {
		for j := i + 1; j < len(menu.Children); j++ {
			if menu.Children[i].Sequence > menu.Children[j].Sequence {
				menu.Children[i], menu.Children[j] = menu.Children[j], menu.Children[i]
			}
		}
	}
}
