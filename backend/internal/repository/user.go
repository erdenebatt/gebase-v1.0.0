package repository

import (
	"context"

	"gebase/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository[domain.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository[domain.User](db),
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByRegNo(ctx context.Context, regNo string) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).Where("reg_no = ?", regNo).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindBySsoUserID(ctx context.Context, ssoUserID int64) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).Where("sso_user_id = ?", ssoUserID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindWithRoles(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).
		Preload("UserSystemRoles").
		Preload("UserSystemRoles.Role").
		Preload("UserSystemRoles.System").
		Preload("Organization").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByOrganization(ctx context.Context, orgID int64, params PaginationParams) (*PaginatedResult[domain.User], error) {
	var users []domain.User
	var total int64

	query := r.DB.WithContext(ctx).Model(&domain.User{}).Where("organization_id = ?", orgID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset(params.GetOffset()).Limit(params.GetLimit()).Find(&users).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / params.GetLimit()
	if int(total)%params.GetLimit() > 0 {
		totalPages++
	}

	return &PaginatedResult[domain.User]{
		Data:       users,
		Page:       params.Page,
		PageSize:   params.GetLimit(),
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	return r.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userID).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}
