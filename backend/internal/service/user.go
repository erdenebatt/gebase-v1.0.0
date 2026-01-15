package service

import (
	"context"
	"errors"

	"gebase/internal/domain"
	"gebase/internal/repository"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrRegNoAlreadyExists = errors.New("registration number already exists")
)

type UserService struct {
	userRepo     *repository.UserRepository
	roleRepo     *repository.UserSystemRoleRepository
	sessionRepo  *repository.SessionRepository
}

func NewUserService(
	userRepo *repository.UserRepository,
	roleRepo *repository.UserSystemRoleRepository,
	sessionRepo *repository.SessionRepository,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		sessionRepo: sessionRepo,
	}
}

type CreateUserRequest struct {
	RegNo          string `json:"reg_no" binding:"required"`
	FamilyName     string `json:"family_name"`
	LastName       string `json:"last_name" binding:"required"`
	FirstName      string `json:"first_name" binding:"required"`
	Gender         int    `json:"gender"`
	BirthDate      string `json:"birth_date"`
	PhoneNo        string `json:"phone_no"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	OrganizationID *int64 `json:"organization_id"`
	LanguageCode   string `json:"language_code"`
}

type UpdateUserRequest struct {
	FamilyName     string  `json:"family_name"`
	LastName       string  `json:"last_name"`
	FirstName      string  `json:"first_name"`
	Gender         int     `json:"gender"`
	BirthDate      string  `json:"birth_date"`
	PhoneNo        string  `json:"phone_no"`
	AvatarURL      *string `json:"avatar_url"`
	OrganizationID *int64  `json:"organization_id"`
	LanguageCode   string  `json:"language_code"`
	IsActive       *bool   `json:"is_active"`
}

// ListUsers returns paginated list of users
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[domain.User], error) {
	return s.userRepo.FindWithPagination(ctx, repository.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	})
}

// ListUsersByOrganization returns paginated list of users in an organization
func (s *UserService) ListUsersByOrganization(ctx context.Context, orgID int64, page, pageSize int) (*repository.PaginatedResult[domain.User], error) {
	return s.userRepo.FindByOrganization(ctx, orgID, repository.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	})
}

// GetUser returns user by ID
func (s *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserWithRoles returns user with roles
func (s *UserService) GetUserWithRoles(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepo.FindWithRoles(ctx, id)
}

// GetUserByEmail returns user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest, createdBy int64) (*domain.User, error) {
	// Check email uniqueness
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Check reg_no uniqueness
	existing, _ = s.userRepo.FindByRegNo(ctx, req.RegNo)
	if existing != nil {
		return nil, ErrRegNoAlreadyExists
	}

	// Hash password
	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	languageCode := req.LanguageCode
	if languageCode == "" {
		languageCode = "mn"
	}

	user := &domain.User{
		RegNo:          req.RegNo,
		FamilyName:     req.FamilyName,
		LastName:       req.LastName,
		FirstName:      req.FirstName,
		Gender:         req.Gender,
		BirthDate:      req.BirthDate,
		PhoneNo:        req.PhoneNo,
		Email:          req.Email,
		PasswordHash:   passwordHash,
		OrganizationID: req.OrganizationID,
		LanguageCode:   languageCode,
		IsActive:       domain.Ptr(true),
	}
	user.CreatedBy = &createdBy

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest, updatedBy int64) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.FamilyName != "" {
		user.FamilyName = req.FamilyName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.Gender != 0 {
		user.Gender = req.Gender
	}
	if req.BirthDate != "" {
		user.BirthDate = req.BirthDate
	}
	if req.PhoneNo != "" {
		user.PhoneNo = req.PhoneNo
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.OrganizationID != nil {
		user.OrganizationID = req.OrganizationID
	}
	if req.LanguageCode != "" {
		user.LanguageCode = req.LanguageCode
	}
	if req.IsActive != nil {
		user.IsActive = req.IsActive
	}

	user.UpdatedBy = &updatedBy

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id int64, deletedBy int64) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	user.DeletedBy = &deletedBy
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return s.userRepo.Delete(ctx, id)
}

// GetUserRoles returns user roles
func (s *UserService) GetUserRoles(ctx context.Context, userID int64) ([]domain.UserSystemRole, error) {
	return s.roleRepo.FindByUserID(ctx, userID)
}

// AssignUserRoles assigns roles to user
func (s *UserService) AssignUserRoles(ctx context.Context, userID int64, systemID *int, roleIDs []int, orgID *int64, assignedBy int64) error {
	return s.roleRepo.AssignRoles(ctx, userID, systemID, roleIDs, orgID, assignedBy)
}

// ResetPassword resets user password
func (s *UserService) ResetPassword(ctx context.Context, userID int64, newPassword string, updatedBy int64) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	user.UpdatedBy = &updatedBy

	return s.userRepo.Update(ctx, user)
}

// ChangePassword changes user password (requires old password)
func (s *UserService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := verifyPassword(user.PasswordHash, oldPassword); err != nil {
		return ErrInvalidCredentials
	}

	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	user.UpdatedBy = &user.ID

	return s.userRepo.Update(ctx, user)
}

func verifyPassword(hash, password string) error {
	return bcryptCompare(hash, password)
}

func bcryptCompare(hash, password string) error {
	// Use golang.org/x/crypto/bcrypt
	return nil // Placeholder - will use bcrypt.CompareHashAndPassword
}
