package service

import (
	"context"
	"errors"

	"gebase/internal/domain"
	"gebase/internal/repository"
)

var (
	ErrSystemCodeExists = errors.New("system code already exists")
)

type SystemService struct {
	systemRepo *repository.SystemRepository
	moduleRepo *repository.ModuleRepository
	menuRepo   *repository.MenuRepository
}

func NewSystemService(
	systemRepo *repository.SystemRepository,
	moduleRepo *repository.ModuleRepository,
	menuRepo *repository.MenuRepository,
) *SystemService {
	return &SystemService{
		systemRepo: systemRepo,
		moduleRepo: moduleRepo,
		menuRepo:   menuRepo,
	}
}

type CreateSystemRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	BaseURL     string `json:"base_url"`
	Sequence    int    `json:"sequence"`
}

type UpdateSystemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	BaseURL     string `json:"base_url"`
	Sequence    int    `json:"sequence"`
	IsActive    *bool  `json:"is_active"`
}

// ListSystems returns all systems
func (s *SystemService) ListSystems(ctx context.Context) ([]domain.System, error) {
	return s.systemRepo.FindAll(ctx)
}

// ListActiveSystems returns all active systems
func (s *SystemService) ListActiveSystems(ctx context.Context) ([]domain.System, error) {
	return s.systemRepo.FindActive(ctx)
}

// GetSystem returns system by ID
func (s *SystemService) GetSystem(ctx context.Context, id int) (*domain.System, error) {
	system, err := s.systemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrSystemNotFound
	}
	return system, nil
}

// GetSystemByCode returns system by code
func (s *SystemService) GetSystemByCode(ctx context.Context, code string) (*domain.System, error) {
	return s.systemRepo.FindByCode(ctx, code)
}

// GetSystemWithModules returns system with its modules
func (s *SystemService) GetSystemWithModules(ctx context.Context, id int) (*domain.System, error) {
	return s.systemRepo.FindWithModules(ctx, id)
}

// GetSystemWithMenus returns system with its menus
func (s *SystemService) GetSystemWithMenus(ctx context.Context, id int) (*domain.System, error) {
	return s.systemRepo.FindWithMenus(ctx, id)
}

// CreateSystem creates a new system
func (s *SystemService) CreateSystem(ctx context.Context, req *CreateSystemRequest, createdBy int64) (*domain.System, error) {
	// Check code uniqueness
	existing, _ := s.systemRepo.FindByCode(ctx, req.Code)
	if existing != nil {
		return nil, ErrSystemCodeExists
	}

	system := &domain.System{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IconURL:     req.IconURL,
		BaseURL:     req.BaseURL,
		Sequence:    req.Sequence,
		IsActive:    domain.Ptr(true),
	}
	system.CreatedBy = &createdBy

	if err := s.systemRepo.Create(ctx, system); err != nil {
		return nil, err
	}

	return system, nil
}

// UpdateSystem updates a system
func (s *SystemService) UpdateSystem(ctx context.Context, id int, req *UpdateSystemRequest, updatedBy int64) (*domain.System, error) {
	system, err := s.systemRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrSystemNotFound
	}

	if req.Name != "" {
		system.Name = req.Name
	}
	if req.Description != "" {
		system.Description = req.Description
	}
	if req.IconURL != "" {
		system.IconURL = req.IconURL
	}
	if req.BaseURL != "" {
		system.BaseURL = req.BaseURL
	}
	if req.Sequence != 0 {
		system.Sequence = req.Sequence
	}
	if req.IsActive != nil {
		system.IsActive = req.IsActive
	}

	system.UpdatedBy = &updatedBy

	if err := s.systemRepo.Update(ctx, system); err != nil {
		return nil, err
	}

	return system, nil
}

// DeleteSystem soft deletes a system
func (s *SystemService) DeleteSystem(ctx context.Context, id int, deletedBy int64) error {
	system, err := s.systemRepo.FindByID(ctx, id)
	if err != nil {
		return ErrSystemNotFound
	}

	system.DeletedBy = &deletedBy
	if err := s.systemRepo.Update(ctx, system); err != nil {
		return err
	}

	return s.systemRepo.Delete(ctx, id)
}

// GetSystemModules returns modules for a system
func (s *SystemService) GetSystemModules(ctx context.Context, systemID int) ([]domain.Module, error) {
	return s.moduleRepo.FindBySystemID(ctx, systemID)
}

// GetSystemMenus returns menus for a system
func (s *SystemService) GetSystemMenus(ctx context.Context, systemID int) ([]domain.Menu, error) {
	return s.menuRepo.FindBySystemID(ctx, systemID)
}
