package auth

import (
	"context"
	"time"

	"gebase/internal/config"
	"gebase/internal/domain"
	"gebase/internal/repository"

	"github.com/google/uuid"
)

type SessionService struct {
	config      *config.Config
	sessionRepo *repository.SessionRepository
	historyRepo *repository.SessionSystemHistoryRepository
}

func NewSessionService(
	cfg *config.Config,
	sessionRepo *repository.SessionRepository,
	historyRepo *repository.SessionSystemHistoryRepository,
) *SessionService {
	return &SessionService{
		config:      cfg,
		sessionRepo: sessionRepo,
		historyRepo: historyRepo,
	}
}

// CreateSession creates a new session for user + device
func (s *SessionService) CreateSession(ctx context.Context, userID int64, deviceID int64, ipAddress, userAgent string, orgID *int64) (*domain.Session, error) {
	session := &domain.Session{
		SessionToken:   uuid.New().String(),
		UserID:         userID,
		DeviceID:       deviceID,
		OrganizationID: orgID,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		IsActive:       domain.Ptr(true),
		ExpiresAt:      time.Now().Add(s.config.JWT.PlatformExpiry),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession retrieves session by token
func (s *SessionService) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	return s.sessionRepo.FindByToken(ctx, token)
}

// GetSessionByID retrieves session by ID
func (s *SessionService) GetSessionByID(ctx context.Context, id int64) (*domain.Session, error) {
	session, err := s.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// UpdateActivity updates last activity timestamp
func (s *SessionService) UpdateActivity(ctx context.Context, sessionID int64) error {
	return s.sessionRepo.UpdateActivity(ctx, sessionID)
}

// SwitchSystem switches the current system for a session
func (s *SessionService) SwitchSystem(ctx context.Context, sessionID int64, systemID int, ipAddress string) error {
	// Update session's current system
	if err := s.sessionRepo.UpdateCurrentSystem(ctx, sessionID, systemID); err != nil {
		return err
	}

	// Record system switch history
	if err := s.historyRepo.RecordSwitch(ctx, sessionID, systemID, ipAddress); err != nil {
		return err
	}

	return nil
}

// Logout terminates a session
func (s *SessionService) Logout(ctx context.Context, sessionID int64, reason string) error {
	return s.sessionRepo.Logout(ctx, sessionID, reason)
}

// LogoutUser terminates all sessions for a user
func (s *SessionService) LogoutUser(ctx context.Context, userID int64, reason string) error {
	return s.sessionRepo.LogoutByUserID(ctx, userID, reason)
}

// LogoutDevice terminates all sessions for a device
func (s *SessionService) LogoutDevice(ctx context.Context, deviceID int64, reason string) error {
	return s.sessionRepo.LogoutByDeviceID(ctx, deviceID, reason)
}

// GetUserSessions retrieves all active sessions for a user
func (s *SessionService) GetUserSessions(ctx context.Context, userID int64) ([]domain.Session, error) {
	return s.sessionRepo.FindByUserID(ctx, userID)
}

// GetDeviceSessions retrieves all active sessions for a device
func (s *SessionService) GetDeviceSessions(ctx context.Context, deviceID int64) ([]domain.Session, error) {
	return s.sessionRepo.FindByDeviceID(ctx, deviceID)
}

// GetActiveSessions retrieves all active sessions with pagination
func (s *SessionService) GetActiveSessions(ctx context.Context, page, pageSize int) (*repository.PaginatedResult[domain.Session], error) {
	return s.sessionRepo.FindActiveSessions(ctx, repository.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	})
}

// CountActiveSessions returns count of active sessions
func (s *SessionService) CountActiveSessions(ctx context.Context) (int64, error) {
	return s.sessionRepo.CountActiveSessions(ctx)
}

// CleanupExpiredSessions marks expired sessions as inactive
func (s *SessionService) CleanupExpiredSessions(ctx context.Context) error {
	return s.sessionRepo.CleanupExpiredSessions(ctx)
}

// GetSystemSwitchHistory retrieves system switch history for a session
func (s *SessionService) GetSystemSwitchHistory(ctx context.Context, sessionID int64) ([]domain.SessionSystemHistory, error) {
	return s.historyRepo.FindBySessionID(ctx, sessionID)
}

// IsSessionValid checks if session is valid (active and not expired)
func (s *SessionService) IsSessionValid(ctx context.Context, session *domain.Session) bool {
	if session == nil {
		return false
	}
	if session.IsActive == nil || !*session.IsActive {
		return false
	}
	if time.Now().After(session.ExpiresAt) {
		return false
	}
	return true
}
