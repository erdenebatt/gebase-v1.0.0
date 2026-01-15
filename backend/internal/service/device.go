package service

import (
	"context"
	"errors"
	"time"

	"gebase/internal/domain"
	"gebase/internal/repository"
)

var (
	ErrDeviceUIDExists = errors.New("device UID already exists")
)

type DeviceService struct {
	deviceRepo  *repository.DeviceRepository
	sessionRepo *repository.SessionRepository
}

func NewDeviceService(
	deviceRepo *repository.DeviceRepository,
	sessionRepo *repository.SessionRepository,
) *DeviceService {
	return &DeviceService{
		deviceRepo:  deviceRepo,
		sessionRepo: sessionRepo,
	}
}

type RegisterDeviceRequest struct {
	DeviceUID      string                `json:"device_uid" binding:"required"`
	Name           string                `json:"name" binding:"required"`
	Platform       domain.DevicePlatform `json:"platform" binding:"required"`
	OSVersion      string                `json:"os_version"`
	AppVersion     string                `json:"app_version"`
	PushToken      string                `json:"push_token"`
	OrganizationID *int64                `json:"organization_id"`
}

type UpdateDeviceRequest struct {
	Name       string `json:"name"`
	OSVersion  string `json:"os_version"`
	AppVersion string `json:"app_version"`
	PushToken  string `json:"push_token"`
	IsActive   *bool  `json:"is_active"`
}

// RegisterDevice registers a new device or returns existing one
func (s *DeviceService) RegisterDevice(ctx context.Context, req *RegisterDeviceRequest) (*domain.Device, error) {
	// Check if device already exists
	existing, err := s.deviceRepo.FindByUID(ctx, req.DeviceUID)
	if err == nil && existing != nil {
		// Update existing device info
		existing.Name = req.Name
		existing.Platform = req.Platform
		existing.OSVersion = req.OSVersion
		existing.AppVersion = req.AppVersion
		if req.PushToken != "" {
			existing.PushToken = req.PushToken
		}
		now := time.Now()
		existing.LastHeartbeat = &now

		if err := s.deviceRepo.Update(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// Create new device
	device := &domain.Device{
		DeviceUID:      req.DeviceUID,
		Name:           req.Name,
		Platform:       req.Platform,
		OSVersion:      req.OSVersion,
		AppVersion:     req.AppVersion,
		PushToken:      req.PushToken,
		OrganizationID: req.OrganizationID,
		IsRegistered:   domain.Ptr(false),
		IsActive:       domain.Ptr(true),
		ConfigJSON:     "{}",
	}

	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return nil, err
	}

	return device, nil
}

// Heartbeat updates device heartbeat timestamp
func (s *DeviceService) Heartbeat(ctx context.Context, deviceUID string) error {
	return s.deviceRepo.UpdateHeartbeat(ctx, deviceUID)
}

// VerifyDevice verifies if device exists and is active
func (s *DeviceService) VerifyDevice(ctx context.Context, deviceUID string) (*domain.Device, error) {
	device, err := s.deviceRepo.FindByUID(ctx, deviceUID)
	if err != nil {
		return nil, ErrDeviceNotFound
	}

	if device.IsActive == nil || !*device.IsActive {
		return nil, ErrDeviceNotActive
	}

	return device, nil
}

// ListDevices returns paginated list of devices
func (s *DeviceService) ListDevices(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[domain.Device], error) {
	return s.deviceRepo.FindWithPagination(ctx, repository.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	})
}

// ListDevicesByOrganization returns devices for an organization
func (s *DeviceService) ListDevicesByOrganization(ctx context.Context, orgID int64, page, pageSize int) (*repository.PaginatedResult[domain.Device], error) {
	return s.deviceRepo.FindByOrganization(ctx, orgID, repository.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	})
}

// GetDevice returns device by ID
func (s *DeviceService) GetDevice(ctx context.Context, id int64) (*domain.Device, error) {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrDeviceNotFound
	}
	return device, nil
}

// GetDeviceByUID returns device by UID
func (s *DeviceService) GetDeviceByUID(ctx context.Context, uid string) (*domain.Device, error) {
	return s.deviceRepo.FindByUID(ctx, uid)
}

// UpdateDevice updates device info
func (s *DeviceService) UpdateDevice(ctx context.Context, id int64, req *UpdateDeviceRequest, updatedBy int64) (*domain.Device, error) {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrDeviceNotFound
	}

	if req.Name != "" {
		device.Name = req.Name
	}
	if req.OSVersion != "" {
		device.OSVersion = req.OSVersion
	}
	if req.AppVersion != "" {
		device.AppVersion = req.AppVersion
	}
	if req.PushToken != "" {
		device.PushToken = req.PushToken
	}
	if req.IsActive != nil {
		device.IsActive = req.IsActive
	}

	device.UpdatedBy = &updatedBy

	if err := s.deviceRepo.Update(ctx, device); err != nil {
		return nil, err
	}

	return device, nil
}

// UpdateDeviceConfig updates device remote configuration
func (s *DeviceService) UpdateDeviceConfig(ctx context.Context, id int64, configJSON string, updatedBy int64) error {
	return s.deviceRepo.UpdateConfig(ctx, id, configJSON, updatedBy)
}

// DeactivateDevice deactivates a device and terminates its sessions
func (s *DeviceService) DeactivateDevice(ctx context.Context, id int64, reason string, deactivatedBy int64) error {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return ErrDeviceNotFound
	}

	device.IsActive = domain.Ptr(false)
	device.UpdatedBy = &deactivatedBy

	if err := s.deviceRepo.Update(ctx, device); err != nil {
		return err
	}

	// Terminate all sessions for this device
	return s.sessionRepo.LogoutByDeviceID(ctx, id, reason)
}

// GetOnlineDevices returns devices with recent heartbeat
func (s *DeviceService) GetOnlineDevices(ctx context.Context, threshold time.Duration) ([]domain.Device, error) {
	return s.deviceRepo.FindOnlineDevices(ctx, threshold)
}

// GetDeviceSessions returns active sessions for a device
func (s *DeviceService) GetDeviceSessions(ctx context.Context, deviceID int64) ([]domain.Session, error) {
	return s.sessionRepo.FindByDeviceID(ctx, deviceID)
}

// ApproveDevice marks device as registered/approved
func (s *DeviceService) ApproveDevice(ctx context.Context, deviceUID string, approvedBy int64) error {
	return s.deviceRepo.Register(ctx, deviceUID, approvedBy)
}
