# Updated Domain Models - Platform + Systems

## Base Model (internal/domain/base.go)

```go
package domain

import (
    "time"
    "gorm.io/gorm"
)

type ExtraFields struct {
    CreatedDate *time.Time     `json:"created_date,omitempty" gorm:"autoCreateTime"`
    UpdatedDate *time.Time     `json:"updated_date,omitempty" gorm:"autoUpdateTime"`
    CreatedBy   *int64         `json:"created_by,omitempty"`
    UpdatedBy   *int64         `json:"updated_by,omitempty"`
    DeletedDate gorm.DeletedAt `json:"deleted_date,omitempty" gorm:"index"`
    DeletedBy   *int64         `json:"deleted_by,omitempty"`
}
```

## System Model (internal/domain/system.go)

```go
package domain

type System struct {
    ID          int      `json:"id" gorm:"primaryKey"`
    Code        string   `json:"code" gorm:"unique;type:varchar(50)"` // dsl, gateway
    Name        string   `json:"name" gorm:"type:varchar(100)"`
    Description string   `json:"description" gorm:"type:varchar(255)"`
    IconURL     string   `json:"icon_url" gorm:"type:varchar(255)"`
    IconName    string   `json:"icon_name" gorm:"type:varchar(50)"`   // Lucide icon name
    BaseURL     string   `json:"base_url" gorm:"type:varchar(255)"`
    Color       string   `json:"color" gorm:"type:varchar(20)"`       // Brand color
    IsActive    *bool    `json:"is_active" gorm:"default:true"`
    Sequence    int      `json:"sequence" gorm:"default:0"`
    Modules     []Module `json:"modules,omitempty" gorm:"foreignKey:SystemID"`
    Menus       []Menu   `json:"menus,omitempty" gorm:"foreignKey:SystemID"`
    ExtraFields
}

func (System) TableName() string {
    return "systems"
}
```

## Module Model (internal/domain/module.go)

```go
package domain

type Module struct {
    ID            int            `json:"id" gorm:"primaryKey"`
    Code          string         `json:"code" gorm:"type:varchar(50)"`
    Name          string         `json:"name" gorm:"type:varchar(100)"`
    Description   string         `json:"description" gorm:"type:varchar(255)"`
    SystemID      *int           `json:"system_id"`                              // NULL = Platform module
    System        *System        `json:"system,omitempty" gorm:"foreignKey:SystemID"`
    IsActive      *bool          `json:"is_active" gorm:"default:true"`
    ModuleActions []ModuleAction `json:"module_actions,omitempty" gorm:"foreignKey:ModuleID"`
    ExtraFields
}

func (Module) TableName() string {
    return "modules"
}

// IsPlatformModule returns true if module belongs to platform layer
func (m *Module) IsPlatformModule() bool {
    return m.SystemID == nil
}

// Unique constraint: system_id + code (NULL system_id treated as unique value)
```

## Permission Model (internal/domain/permission.go)

```go
package domain

// Permission code format:
// Platform: {module.action} - e.g., "user.view", "organization.create"
// System:   {system.module.action} - e.g., "dsl.schema.view", "gateway.client.create"
type Permission struct {
    ID          int     `json:"id" gorm:"primaryKey"`
    Code        string  `json:"code" gorm:"unique;not null;type:varchar(255)"`
    Name        string  `json:"name" gorm:"type:varchar(255)"`
    Description string  `json:"description" gorm:"type:varchar(255)"`
    SystemID    *int    `json:"system_id"`                              // NULL = Platform permission
    System      *System `json:"system,omitempty" gorm:"foreignKey:SystemID"`
    ModuleID    int     `json:"module_id"`
    Module      *Module `json:"module,omitempty" gorm:"foreignKey:ModuleID"`
    ActionID    *int64  `json:"action_id"`
    Action      *Action `json:"action,omitempty" gorm:"foreignKey:ActionID"`
    IsActive    *bool   `json:"is_active" gorm:"default:true"`
    ExtraFields
}

func (Permission) TableName() string {
    return "permissions"
}

// IsPlatformPermission returns true if permission belongs to platform layer
func (p *Permission) IsPlatformPermission() bool {
    return p.SystemID == nil
}
```

## Role Model (internal/domain/role.go)

```go
package domain

type Role struct {
    ID          int              `json:"id" gorm:"primaryKey"`
    Code        string           `json:"code" gorm:"type:varchar(50)"`
    Name        string           `json:"name" gorm:"type:varchar(100)"`
    Description string           `json:"description" gorm:"type:varchar(255)"`
    SystemID    *int             `json:"system_id"`                              // NULL = Platform role
    System      *System          `json:"system,omitempty" gorm:"foreignKey:SystemID"`
    IsSystem    *bool            `json:"is_system" gorm:"default:false"`         // Built-in role
    IsActive    *bool            `json:"is_active" gorm:"default:true"`
    Permissions []RolePermission `json:"permissions,omitempty" gorm:"foreignKey:RoleID"`
    Menus       []RoleMenu       `json:"menus,omitempty" gorm:"foreignKey:RoleID"`
    ExtraFields
}

func (Role) TableName() string {
    return "roles"
}

// IsPlatformRole returns true if role belongs to platform layer
func (r *Role) IsPlatformRole() bool {
    return r.SystemID == nil
}

// Unique constraint: system_id + code
```

## Menu Model (internal/domain/menu.go)

```go
package domain

type Menu struct {
    ID        int     `json:"id" gorm:"primaryKey"`
    Code      string  `json:"code" gorm:"type:varchar(50)"`
    Name      string  `json:"name" gorm:"type:varchar(100)"`
    SystemID  *int    `json:"system_id"`                              // NULL = Platform menu
    System    *System `json:"system,omitempty" gorm:"foreignKey:SystemID"`
    ParentID  *int    `json:"parent_id"`
    Parent    *Menu   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
    Children  []Menu  `json:"children,omitempty" gorm:"foreignKey:ParentID"`
    Path      string  `json:"path" gorm:"type:varchar(255)"`
    Icon      string  `json:"icon" gorm:"type:varchar(100)"`
    Component string  `json:"component" gorm:"type:varchar(255)"`
    Sequence  int     `json:"sequence" gorm:"default:0"`
    IsVisible *bool   `json:"is_visible" gorm:"default:true"`
    IsActive  *bool   `json:"is_active" gorm:"default:true"`
    ExtraFields
}

func (Menu) TableName() string {
    return "menus"
}

// Unique constraint: system_id + code
```

## User System Role Model (internal/domain/user_system_role.go)

```go
package domain

type UserSystemRole struct {
    ID             int           `json:"id" gorm:"primaryKey"`
    UserID         int64         `json:"user_id"`
    User           *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
    SystemID       *int          `json:"system_id"`                              // NULL = Platform role
    System         *System       `json:"system,omitempty" gorm:"foreignKey:SystemID"`
    RoleID         int           `json:"role_id"`
    Role           *Role         `json:"role,omitempty" gorm:"foreignKey:RoleID"`
    OrganizationID *int64        `json:"organization_id,omitempty"`              // Org-scoped role
    Organization   *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
    IsActive       *bool         `json:"is_active" gorm:"default:true"`
    IsDefault      *bool         `json:"is_default" gorm:"default:false"`        // Default system when login
    ExtraFields
}

func (UserSystemRole) TableName() string {
    return "user_system_roles"
}

// Unique constraint: user_id + system_id + role_id + organization_id
```

## Session Model - Updated (internal/domain/session.go)

```go
package domain

import "time"

type TokenType string

const (
    TokenTypePlatform TokenType = "platform"
    TokenTypeSystem   TokenType = "system"
)

type Session struct {
    ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
    SessionToken    string     `json:"session_token" gorm:"unique;type:varchar(255)"`
    UserID          int64      `json:"user_id"`
    User            *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
    DeviceID        int64      `json:"device_id"`
    Device          *Device    `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
    
    // Current active system (NULL = platform context only)
    CurrentSystemID *int       `json:"current_system_id"`
    CurrentSystem   *System    `json:"current_system,omitempty" gorm:"foreignKey:CurrentSystemID"`
    
    // Organization context for current system
    OrganizationID  *int64     `json:"organization_id,omitempty"`
    Organization    *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
    
    IPAddress       string     `json:"ip_address" gorm:"type:varchar(45)"`
    UserAgent       string     `json:"user_agent" gorm:"type:varchar(500)"`
    IsActive        *bool      `json:"is_active" gorm:"default:true"`
    ExpiresAt       time.Time  `json:"expires_at"`
    LastActivity    *time.Time `json:"last_activity"`
    LastSystemSwitch *time.Time `json:"last_system_switch"`
    LogoutAt        *time.Time `json:"logout_at"`
    LogoutReason    string     `json:"logout_reason" gorm:"type:varchar(100)"`
    ExtraFields
}

func (Session) TableName() string {
    return "sessions"
}

// IsInSystemContext returns true if session is currently in a system context
func (s *Session) IsInSystemContext() bool {
    return s.CurrentSystemID != nil
}
```

## Session System History Model (internal/domain/session_system_history.go)

```go
package domain

import "time"

type SessionSystemHistory struct {
    ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
    SessionID  int64     `json:"session_id"`
    Session    *Session  `json:"session,omitempty" gorm:"foreignKey:SessionID"`
    SystemID   int       `json:"system_id"`
    System     *System   `json:"system,omitempty" gorm:"foreignKey:SystemID"`
    SwitchedAt time.Time `json:"switched_at"`
    IPAddress  string    `json:"ip_address" gorm:"type:varchar(45)"`
}

func (SessionSystemHistory) TableName() string {
    return "session_system_history"
}
```

## Organization System Model (internal/domain/organization_system.go)

```go
package domain

import "time"

type OrganizationSystem struct {
    ID             int           `json:"id" gorm:"primaryKey"`
    OrganizationID int64         `json:"organization_id"`
    Organization   *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
    SystemID       int           `json:"system_id"`
    System         *System       `json:"system,omitempty" gorm:"foreignKey:SystemID"`
    IsActive       *bool         `json:"is_active" gorm:"default:true"`
    ActivatedAt    *time.Time    `json:"activated_at"`
    ExpiresAt      *time.Time    `json:"expires_at"`
    MaxUsers       *int          `json:"max_users"`                              // License limit
    Config         string        `json:"config" gorm:"type:jsonb"`               // System-specific config
    ExtraFields
}

func (OrganizationSystem) TableName() string {
    return "organization_systems"
}

// Unique constraint: organization_id + system_id
```

## User Model - Updated (internal/domain/user.go)

```go
package domain

import "time"

type User struct {
    ID            int64  `json:"id" gorm:"primaryKey;autoIncrement"`
    SsoUserID     int64  `json:"sso_user_id" gorm:"uniqueIndex"`
    CivilID       int64  `json:"civil_id"`
    RegNo         string `json:"reg_no" gorm:"type:varchar(10);uniqueIndex"`
    FamilyName    string `json:"family_name" gorm:"type:varchar(80)"`
    LastName      string `json:"last_name" gorm:"type:varchar(150)"`
    FirstName     string `json:"first_name" gorm:"type:varchar(150)"`
    Gender        int    `json:"gender"`
    BirthDate     string `json:"birth_date" gorm:"type:varchar(10)"`
    PhoneNo       string `json:"phone_no" gorm:"type:varchar(8)"`
    Email         string `json:"email" gorm:"type:varchar(80);uniqueIndex"`
    PasswordHash  string `json:"-" gorm:"type:varchar(255)"`
    AvatarURL     *string    `json:"avatar_url,omitempty" gorm:"type:varchar(500)"`
    IsActive      *bool      `json:"is_active" gorm:"default:true"`
    LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
    LanguageCode  string     `json:"language_code" gorm:"type:varchar(5);default:'mn'"`
    
    // Primary organization (home org)
    OrganizationID  *int64           `json:"organization_id,omitempty"`
    Organization    *Organization    `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
    
    // Default system to redirect after login
    DefaultSystemID *int             `json:"default_system_id,omitempty"`
    DefaultSystem   *System          `json:"default_system,omitempty" gorm:"foreignKey:DefaultSystemID"`
    
    // All system roles
    UserSystemRoles []UserSystemRole `json:"user_system_roles,omitempty" gorm:"foreignKey:UserID"`
    
    ExtraFields
}

func (User) TableName() string {
    return "users"
}

// GetAvailableSystems returns list of systems user has access to
func (u *User) GetAvailableSystems() []System {
    systemMap := make(map[int]System)
    for _, usr := range u.UserSystemRoles {
        if usr.SystemID != nil && usr.System != nil && usr.IsActive != nil && *usr.IsActive {
            systemMap[*usr.SystemID] = *usr.System
        }
    }
    
    systems := make([]System, 0, len(systemMap))
    for _, s := range systemMap {
        systems = append(systems, s)
    }
    return systems
}
```

## JWT Claims Structure (internal/auth/claims.go)

```go
package auth

import "github.com/golang-jwt/jwt/v5"

// PlatformClaims for platform-level token
type PlatformClaims struct {
    jwt.RegisteredClaims
    TokenType string `json:"type"`      // "platform"
    UserID    int64  `json:"user_id"`
    DeviceID  int64  `json:"device_id"`
    SessionID int64  `json:"session_id"`
}

// SystemClaims for system-level token
type SystemClaims struct {
    jwt.RegisteredClaims
    TokenType      string   `json:"type"`           // "system"
    UserID         int64    `json:"user_id"`
    DeviceID       int64    `json:"device_id"`
    SessionID      int64    `json:"session_id"`
    SystemID       int      `json:"system_id"`
    SystemCode     string   `json:"system_code"`
    OrganizationID *int64   `json:"organization_id,omitempty"`
    RoleID         int      `json:"role_id"`
    RoleCode       string   `json:"role_code"`
    Permissions    []string `json:"permissions"`    // Cached permissions
}
```
