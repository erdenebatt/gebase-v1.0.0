package repository

import (
	"context"
	"time"

	"gebase/internal/domain"

	"gorm.io/gorm"
)

type SessionRepository struct {
	*BaseRepository[domain.Session]
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		BaseRepository: NewBaseRepository[domain.Session](db),
	}
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session
	err := r.DB.WithContext(ctx).
		Preload("User").
		Preload("Device").
		Preload("CurrentSystem").
		Where("session_token = ?", token).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) FindByUserID(ctx context.Context, userID int64) ([]domain.Session, error) {
	var sessions []domain.Session
	err := r.DB.WithContext(ctx).
		Preload("Device").
		Preload("CurrentSystem").
		Where("user_id = ? AND is_active = true", userID).
		Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) FindByDeviceID(ctx context.Context, deviceID int64) ([]domain.Session, error) {
	var sessions []domain.Session
	err := r.DB.WithContext(ctx).
		Preload("User").
		Where("device_id = ? AND is_active = true", deviceID).
		Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) FindActiveSessions(ctx context.Context, params PaginationParams) (*PaginatedResult[domain.Session], error) {
	var sessions []domain.Session
	var total int64

	query := r.DB.WithContext(ctx).Model(&domain.Session{}).
		Where("is_active = true AND expires_at > ?", time.Now())

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Preload("User").
		Preload("Device").
		Preload("CurrentSystem").
		Offset(params.GetOffset()).
		Limit(params.GetLimit()).
		Order("created_date DESC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / params.GetLimit()
	if int(total)%params.GetLimit() > 0 {
		totalPages++
	}

	return &PaginatedResult[domain.Session]{
		Data:       sessions,
		Page:       params.Page,
		PageSize:   params.GetLimit(),
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *SessionRepository) UpdateActivity(ctx context.Context, sessionID int64) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&domain.Session{}).
		Where("id = ?", sessionID).
		Update("last_activity", &now).Error
}

func (r *SessionRepository) UpdateCurrentSystem(ctx context.Context, sessionID int64, systemID int) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&domain.Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"current_system_id":  systemID,
			"last_system_switch": &now,
		}).Error
}

func (r *SessionRepository) Logout(ctx context.Context, sessionID int64, reason string) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&domain.Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"is_active":     false,
			"logout_at":     &now,
			"logout_reason": reason,
		}).Error
}

func (r *SessionRepository) LogoutByUserID(ctx context.Context, userID int64, reason string) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&domain.Session{}).
		Where("user_id = ? AND is_active = true", userID).
		Updates(map[string]interface{}{
			"is_active":     false,
			"logout_at":     &now,
			"logout_reason": reason,
		}).Error
}

func (r *SessionRepository) LogoutByDeviceID(ctx context.Context, deviceID int64, reason string) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&domain.Session{}).
		Where("device_id = ? AND is_active = true", deviceID).
		Updates(map[string]interface{}{
			"is_active":     false,
			"logout_at":     &now,
			"logout_reason": reason,
		}).Error
}

func (r *SessionRepository) CountActiveSessions(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&domain.Session{}).
		Where("is_active = true AND expires_at > ?", time.Now()).
		Count(&count).Error
	return count, err
}

func (r *SessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&domain.Session{}).
		Where("expires_at < ? AND is_active = true", now).
		Updates(map[string]interface{}{
			"is_active":     false,
			"logout_at":     &now,
			"logout_reason": "expired",
		}).Error
}

type SessionSystemHistoryRepository struct {
	*BaseRepository[domain.SessionSystemHistory]
}

func NewSessionSystemHistoryRepository(db *gorm.DB) *SessionSystemHistoryRepository {
	return &SessionSystemHistoryRepository{
		BaseRepository: NewBaseRepository[domain.SessionSystemHistory](db),
	}
}

func (r *SessionSystemHistoryRepository) FindBySessionID(ctx context.Context, sessionID int64) ([]domain.SessionSystemHistory, error) {
	var history []domain.SessionSystemHistory
	err := r.DB.WithContext(ctx).
		Preload("System").
		Where("session_id = ?", sessionID).
		Order("switched_at DESC").
		Find(&history).Error
	return history, err
}

func (r *SessionSystemHistoryRepository) RecordSwitch(ctx context.Context, sessionID int64, systemID int, ipAddress string) error {
	history := domain.SessionSystemHistory{
		SessionID:  sessionID,
		SystemID:   systemID,
		SwitchedAt: time.Now(),
		IPAddress:  ipAddress,
	}
	return r.DB.WithContext(ctx).Create(&history).Error
}
