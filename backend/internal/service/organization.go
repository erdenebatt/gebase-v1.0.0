package service

import (
	"context"
	"errors"

	"gebase/internal/domain"
	"gebase/internal/repository"
)

var (
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrOrgRegNoExists       = errors.New("organization registration number already exists")
)

type OrganizationService struct {
	orgRepo       *repository.OrganizationRepository
	orgTypeRepo   *repository.OrganizationTypeRepository
	orgSystemRepo *repository.OrganizationSystemRepository
}

func NewOrganizationService(
	orgRepo *repository.OrganizationRepository,
	orgTypeRepo *repository.OrganizationTypeRepository,
	orgSystemRepo *repository.OrganizationSystemRepository,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:       orgRepo,
		orgTypeRepo:   orgTypeRepo,
		orgSystemRepo: orgSystemRepo,
	}
}

type CreateOrganizationRequest struct {
	RegNo         string  `json:"reg_no" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	ShortName     string  `json:"short_name"`
	TypeID        int     `json:"type_id" binding:"required"`
	PhoneNo       string  `json:"phone_no"`
	Email         string  `json:"email"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	AimagID       int     `json:"aimag_id"`
	SumID         int     `json:"sum_id"`
	BagID         int     `json:"bag_id"`
	AddressDetail string  `json:"address_detail"`
	ParentID      *int64  `json:"parent_id"`
	Sequence      int     `json:"sequence"`
}

type UpdateOrganizationRequest struct {
	Name          string  `json:"name"`
	ShortName     string  `json:"short_name"`
	TypeID        int     `json:"type_id"`
	PhoneNo       string  `json:"phone_no"`
	Email         string  `json:"email"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	AimagID       int     `json:"aimag_id"`
	SumID         int     `json:"sum_id"`
	BagID         int     `json:"bag_id"`
	AddressDetail string  `json:"address_detail"`
	AimagName     string  `json:"aimag_name"`
	SumName       string  `json:"sum_name"`
	BagName       string  `json:"bag_name"`
	ParentID      *int64  `json:"parent_id"`
	Sequence      int     `json:"sequence"`
	IsActive      *bool   `json:"is_active"`
}

// ListOrganizations returns paginated list of organizations
func (s *OrganizationService) ListOrganizations(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[domain.Organization], error) {
	return s.orgRepo.FindWithPagination(ctx, repository.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	})
}

// GetOrganization returns organization by ID
func (s *OrganizationService) GetOrganization(ctx context.Context, id int64) (*domain.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrOrganizationNotFound
	}
	return org, nil
}

// GetOrganizationWithChildren returns organization with children
func (s *OrganizationService) GetOrganizationWithChildren(ctx context.Context, id int64) (*domain.Organization, error) {
	return s.orgRepo.FindWithChildren(ctx, id)
}

// GetRootOrganizations returns root-level organizations
func (s *OrganizationService) GetRootOrganizations(ctx context.Context) ([]domain.Organization, error) {
	return s.orgRepo.FindRootOrganizations(ctx)
}

// GetChildOrganizations returns child organizations
func (s *OrganizationService) GetChildOrganizations(ctx context.Context, parentID int64) ([]domain.Organization, error) {
	return s.orgRepo.FindChildren(ctx, parentID)
}

// CreateOrganization creates a new organization
func (s *OrganizationService) CreateOrganization(ctx context.Context, req *CreateOrganizationRequest, createdBy int64) (*domain.Organization, error) {
	// Check reg_no uniqueness
	existing, _ := s.orgRepo.FindByRegNo(ctx, req.RegNo)
	if existing != nil {
		return nil, ErrOrgRegNoExists
	}

	org := &domain.Organization{
		RegNo:         req.RegNo,
		Name:          req.Name,
		ShortName:     req.ShortName,
		TypeID:        req.TypeID,
		PhoneNo:       req.PhoneNo,
		Email:         req.Email,
		Longitude:     req.Longitude,
		Latitude:      req.Latitude,
		AimagID:       req.AimagID,
		SumID:         req.SumID,
		BagID:         req.BagID,
		AddressDetail: req.AddressDetail,
		ParentID:      req.ParentID,
		Sequence:      req.Sequence,
		IsActive:      domain.Ptr(true),
	}
	org.CreatedBy = &createdBy

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// UpdateOrganization updates an organization
func (s *OrganizationService) UpdateOrganization(ctx context.Context, id int64, req *UpdateOrganizationRequest, updatedBy int64) (*domain.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrOrganizationNotFound
	}

	if req.Name != "" {
		org.Name = req.Name
	}
	if req.ShortName != "" {
		org.ShortName = req.ShortName
	}
	if req.TypeID != 0 {
		org.TypeID = req.TypeID
	}
	if req.PhoneNo != "" {
		org.PhoneNo = req.PhoneNo
	}
	if req.Email != "" {
		org.Email = req.Email
	}
	if req.Longitude != 0 {
		org.Longitude = req.Longitude
	}
	if req.Latitude != 0 {
		org.Latitude = req.Latitude
	}
	if req.AimagID != 0 {
		org.AimagID = req.AimagID
	}
	if req.SumID != 0 {
		org.SumID = req.SumID
	}
	if req.BagID != 0 {
		org.BagID = req.BagID
	}
	if req.AddressDetail != "" {
		org.AddressDetail = req.AddressDetail
	}
	if req.AimagName != "" {
		org.AimagName = req.AimagName
	}
	if req.SumName != "" {
		org.SumName = req.SumName
	}
	if req.BagName != "" {
		org.BagName = req.BagName
	}
	if req.ParentID != nil {
		org.ParentID = req.ParentID
	}
	if req.Sequence != 0 {
		org.Sequence = req.Sequence
	}
	if req.IsActive != nil {
		org.IsActive = req.IsActive
	}

	org.UpdatedBy = &updatedBy

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// DeleteOrganization soft deletes an organization
func (s *OrganizationService) DeleteOrganization(ctx context.Context, id int64, deletedBy int64) error {
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return ErrOrganizationNotFound
	}

	org.DeletedBy = &deletedBy
	if err := s.orgRepo.Update(ctx, org); err != nil {
		return err
	}

	return s.orgRepo.Delete(ctx, id)
}

// GetOrganizationTypes returns all organization types
func (s *OrganizationService) GetOrganizationTypes(ctx context.Context) ([]domain.OrganizationType, error) {
	return s.orgTypeRepo.FindAll(ctx)
}

// GetEnabledSystems returns enabled systems for an organization
func (s *OrganizationService) GetEnabledSystems(ctx context.Context, orgID int64) ([]domain.System, error) {
	return s.orgRepo.FindEnabledSystems(ctx, orgID)
}

// EnableSystem enables a system for an organization
func (s *OrganizationService) EnableSystem(ctx context.Context, orgID int64, systemID int, createdBy int64) error {
	// Check if already enabled
	existing, _ := s.orgSystemRepo.FindByOrgAndSystem(ctx, orgID, systemID)
	if existing != nil {
		if existing.IsActive != nil && *existing.IsActive {
			return nil // Already enabled
		}
		// Reactivate
		existing.IsActive = domain.Ptr(true)
		existing.UpdatedBy = &createdBy
		return s.orgSystemRepo.Update(ctx, existing)
	}

	orgSystem := &domain.OrganizationSystem{
		OrganizationID: orgID,
		SystemID:       systemID,
		IsActive:       domain.Ptr(true),
	}
	orgSystem.CreatedBy = &createdBy

	return s.orgSystemRepo.Create(ctx, orgSystem)
}

// DisableSystem disables a system for an organization
func (s *OrganizationService) DisableSystem(ctx context.Context, orgID int64, systemID int, updatedBy int64) error {
	existing, err := s.orgSystemRepo.FindByOrgAndSystem(ctx, orgID, systemID)
	if err != nil {
		return nil // Not enabled, nothing to disable
	}

	existing.IsActive = domain.Ptr(false)
	existing.UpdatedBy = &updatedBy
	return s.orgSystemRepo.Update(ctx, existing)
}
