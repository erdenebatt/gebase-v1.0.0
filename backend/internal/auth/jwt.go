package auth

import (
	"errors"
	"time"

	"gebase/internal/config"
	"gebase/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenType string

const (
	TokenTypePlatform TokenType = "platform"
	TokenTypeSystem   TokenType = "system"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID         int64     `json:"user_id"`
	Email          string    `json:"email"`
	OrganizationID *int64    `json:"organization_id,omitempty"`
	SessionID      int64     `json:"session_id"`
	DeviceID       int64     `json:"device_id"`
	TokenType      TokenType `json:"token_type"`
	SystemID       *int      `json:"system_id,omitempty"`
	SystemCode     string    `json:"system_code,omitempty"`
	RoleIDs        []int     `json:"role_ids,omitempty"`
}

type JWTService struct {
	config *config.Config
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{config: cfg}
}

// GeneratePlatformToken creates a platform-level token (24h expiry)
func (s *JWTService) GeneratePlatformToken(user *domain.User, session *domain.Session) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWT.PlatformExpiry)),
			Issuer:    "gebase",
		},
		UserID:         user.ID,
		Email:          user.Email,
		OrganizationID: user.OrganizationID,
		SessionID:      session.ID,
		DeviceID:       session.DeviceID,
		TokenType:      TokenTypePlatform,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

// GenerateSystemToken creates a system-level token (8h expiry)
func (s *JWTService) GenerateSystemToken(user *domain.User, session *domain.Session, system *domain.System, roleIDs []int) (string, error) {
	now := time.Now()
	systemExpiry := 8 * time.Hour // System token expires in 8 hours

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(systemExpiry)),
			Issuer:    "gebase",
		},
		UserID:         user.ID,
		Email:          user.Email,
		OrganizationID: user.OrganizationID,
		SessionID:      session.ID,
		DeviceID:       session.DeviceID,
		TokenType:      TokenTypeSystem,
		SystemID:       &system.ID,
		SystemCode:     system.Code,
		RoleIDs:        roleIDs,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

// GenerateRefreshToken creates a refresh token (7 days expiry)
func (s *JWTService) GenerateRefreshToken(user *domain.User, session *domain.Session) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWT.RefreshExpiry)),
			Issuer:    "gebase",
		},
		UserID:    user.ID,
		Email:     user.Email,
		SessionID: session.ID,
		DeviceID:  session.DeviceID,
		TokenType: TokenTypePlatform,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

// ValidateToken validates and parses a JWT token
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// IsPlatformToken checks if token is a platform token
func (c *Claims) IsPlatformToken() bool {
	return c.TokenType == TokenTypePlatform
}

// IsSystemToken checks if token is a system token
func (c *Claims) IsSystemToken() bool {
	return c.TokenType == TokenTypeSystem
}

// HasSystem checks if token has system context
func (c *Claims) HasSystem() bool {
	return c.SystemID != nil
}
