package repository

import (
	"context"

	"gebase/internal/domain"

	"gorm.io/gorm"
)

type OrganizationRepository struct {
	*BaseRepository[domain.Organization]
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{
		BaseRepository: NewBaseRepository[domain.Organization](db),
	}
}

func (r *OrganizationRepository) FindByRegNo(ctx context.Context, regNo string) (*domain.Organization, error) {
	var org domain.Organization
	err := r.DB.WithContext(ctx).Where("reg_no = ?", regNo).First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) FindBySsoOrgID(ctx context.Context, ssoOrgID int64) (*domain.Organization, error) {
	var org domain.Organization
	err := r.DB.WithContext(ctx).Where("sso_org_id = ?", ssoOrgID).First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) FindWithChildren(ctx context.Context, id int64) (*domain.Organization, error) {
	var org domain.Organization
	err := r.DB.WithContext(ctx).
		Preload("Children").
		Preload("Type").
		Preload("EnabledSystems").
		Preload("EnabledSystems.System").
		First(&org, id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) FindChildren(ctx context.Context, parentID int64) ([]domain.Organization, error) {
	var orgs []domain.Organization
	err := r.DB.WithContext(ctx).Where("parent_id = ?", parentID).Order("sequence").Find(&orgs).Error
	return orgs, err
}

func (r *OrganizationRepository) FindRootOrganizations(ctx context.Context) ([]domain.Organization, error) {
	var orgs []domain.Organization
	err := r.DB.WithContext(ctx).Where("parent_id IS NULL").Order("sequence").Find(&orgs).Error
	return orgs, err
}

func (r *OrganizationRepository) FindEnabledSystems(ctx context.Context, orgID int64) ([]domain.System, error) {
	var systems []domain.System
	err := r.DB.WithContext(ctx).
		Joins("JOIN organization_systems ON organization_systems.system_id = systems.id").
		Where("organization_systems.organization_id = ? AND organization_systems.is_active = true", orgID).
		Find(&systems).Error
	return systems, err
}

type OrganizationTypeRepository struct {
	*BaseRepository[domain.OrganizationType]
}

func NewOrganizationTypeRepository(db *gorm.DB) *OrganizationTypeRepository {
	return &OrganizationTypeRepository{
		BaseRepository: NewBaseRepository[domain.OrganizationType](db),
	}
}

func (r *OrganizationTypeRepository) FindByCode(ctx context.Context, code string) (*domain.OrganizationType, error) {
	var orgType domain.OrganizationType
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&orgType).Error
	if err != nil {
		return nil, err
	}
	return &orgType, nil
}

type OrganizationSystemRepository struct {
	*BaseRepository[domain.OrganizationSystem]
}

func NewOrganizationSystemRepository(db *gorm.DB) *OrganizationSystemRepository {
	return &OrganizationSystemRepository{
		BaseRepository: NewBaseRepository[domain.OrganizationSystem](db),
	}
}

func (r *OrganizationSystemRepository) FindByOrgAndSystem(ctx context.Context, orgID int64, systemID int) (*domain.OrganizationSystem, error) {
	var orgSystem domain.OrganizationSystem
	err := r.DB.WithContext(ctx).
		Where("organization_id = ? AND system_id = ?", orgID, systemID).
		First(&orgSystem).Error
	if err != nil {
		return nil, err
	}
	return &orgSystem, nil
}
