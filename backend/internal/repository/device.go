package repository

import (
	"context"
	"time"

	"gebase/internal/domain"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	*BaseRepository[domain.Device]
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{
		BaseRepository: NewBaseRepository[domain.Device](db),
	}
}

func (r *DeviceRepository) FindByUID(ctx context.Context, uid string) (*domain.Device, error) {
	var device domain.Device
	err := r.DB.WithContext(ctx).Where("device_uid = ?", uid).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) FindByOrganization(ctx context.Context, orgID int64, params PaginationParams) (*PaginatedResult[domain.Device], error) {
	var devices []domain.Device
	var total int64

	query := r.DB.WithContext(ctx).Model(&domain.Device{}).Where("organization_id = ?", orgID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset(params.GetOffset()).Limit(params.GetLimit()).Find(&devices).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / params.GetLimit()
	if int(total)%params.GetLimit() > 0 {
		totalPages++
	}

	return &PaginatedResult[domain.Device]{
		Data:       devices,
		Page:       params.Page,
		PageSize:   params.GetLimit(),
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *DeviceRepository) FindByPlatform(ctx context.Context, platform domain.DevicePlatform) ([]domain.Device, error) {
	var devices []domain.Device
	err := r.DB.WithContext(ctx).Where("platform = ?", platform).Find(&devices).Error
	return devices, err
}

func (r *DeviceRepository) UpdateHeartbeat(ctx context.Context, uid string) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&domain.Device{}).
		Where("device_uid = ?", uid).
		Update("last_heartbeat", &now).Error
}

func (r *DeviceRepository) UpdateConfig(ctx context.Context, id int64, configJSON string, updatedBy int64) error {
	return r.DB.WithContext(ctx).Model(&domain.Device{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"config_json": configJSON,
			"updated_by":  updatedBy,
		}).Error
}

func (r *DeviceRepository) Register(ctx context.Context, uid string, updatedBy int64) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&domain.Device{}).
		Where("device_uid = ?", uid).
		Updates(map[string]interface{}{
			"is_registered": true,
			"registered_at": &now,
			"updated_by":    updatedBy,
		}).Error
}

func (r *DeviceRepository) FindOnlineDevices(ctx context.Context, threshold time.Duration) ([]domain.Device, error) {
	var devices []domain.Device
	cutoff := time.Now().Add(-threshold)
	err := r.DB.WithContext(ctx).
		Where("last_heartbeat > ? AND is_active = true", cutoff).
		Find(&devices).Error
	return devices, err
}
