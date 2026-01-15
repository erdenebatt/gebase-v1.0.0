package service

import (
	"context"

	"gebase/internal/domain"
	"gebase/internal/repository"
)

type MenuService struct {
	menuRepo *repository.MenuRepository
}

func NewMenuService(menuRepo *repository.MenuRepository) *MenuService {
	return &MenuService{menuRepo: menuRepo}
}

type MenuTree struct {
	ID        int        `json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Icon      string     `json:"icon"`
	Component string     `json:"component"`
	Sequence  int        `json:"sequence"`
	IsVisible *bool      `json:"is_visible"`
	Children  []MenuTree `json:"children,omitempty"`
}

// GetMenuTree returns menu tree for a system
func (s *MenuService) GetMenuTree(ctx context.Context, systemID int) ([]MenuTree, error) {
	menus, err := s.menuRepo.FindBySystemID(ctx, systemID)
	if err != nil {
		return nil, err
	}

	return s.buildMenuTree(menus, nil), nil
}

// GetUserMenus returns menu tree for user based on roles
func (s *MenuService) GetUserMenus(ctx context.Context, userID int64, systemID int) ([]MenuTree, error) {
	menus, err := s.menuRepo.FindUserMenus(ctx, userID, systemID)
	if err != nil {
		return nil, err
	}

	return s.buildMenuTree(menus, nil), nil
}

// GetAllMenus returns all menus for a system (flat list)
func (s *MenuService) GetAllMenus(ctx context.Context, systemID int) ([]domain.Menu, error) {
	return s.menuRepo.FindBySystemID(ctx, systemID)
}

// GetMenuByID returns menu by ID
func (s *MenuService) GetMenuByID(ctx context.Context, id int) (*domain.Menu, error) {
	menu, err := s.menuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return menu, nil
}

// GetMenuByCode returns menu by code
func (s *MenuService) GetMenuByCode(ctx context.Context, code string) (*domain.Menu, error) {
	return s.menuRepo.FindByCode(ctx, code)
}

// CreateMenu creates a new menu
func (s *MenuService) CreateMenu(ctx context.Context, menu *domain.Menu) error {
	return s.menuRepo.Create(ctx, menu)
}

// UpdateMenu updates an existing menu
func (s *MenuService) UpdateMenu(ctx context.Context, menu *domain.Menu) error {
	return s.menuRepo.Update(ctx, menu)
}

// DeleteMenu soft deletes a menu
func (s *MenuService) DeleteMenu(ctx context.Context, id int) error {
	return s.menuRepo.Delete(ctx, id)
}

// buildMenuTree builds a tree structure from flat menu list
func (s *MenuService) buildMenuTree(menus []domain.Menu, parentID *int) []MenuTree {
	var tree []MenuTree

	for _, menu := range menus {
		// Check if menu belongs to current parent
		if (parentID == nil && menu.ParentID == nil) || (parentID != nil && menu.ParentID != nil && *parentID == *menu.ParentID) {
			node := MenuTree{
				ID:        menu.ID,
				Code:      menu.Code,
				Name:      menu.Name,
				Path:      menu.Path,
				Icon:      menu.Icon,
				Component: menu.Component,
				Sequence:  menu.Sequence,
				IsVisible: menu.IsVisible,
				Children:  s.buildMenuTree(menus, &menu.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}

// ConvertToMenuTree converts domain menus to MenuTree
func ConvertToMenuTree(menus []domain.Menu) []MenuTree {
	menuMap := make(map[int]*MenuTree)
	var roots []MenuTree

	// First pass: create all nodes
	for _, m := range menus {
		node := &MenuTree{
			ID:        m.ID,
			Code:      m.Code,
			Name:      m.Name,
			Path:      m.Path,
			Icon:      m.Icon,
			Component: m.Component,
			Sequence:  m.Sequence,
			IsVisible: m.IsVisible,
			Children:  []MenuTree{},
		}
		menuMap[m.ID] = node
	}

	// Second pass: build tree structure
	for _, m := range menus {
		node := menuMap[m.ID]
		if m.ParentID == nil {
			roots = append(roots, *node)
		} else if parent, ok := menuMap[*m.ParentID]; ok {
			parent.Children = append(parent.Children, *node)
		}
	}

	return roots
}
