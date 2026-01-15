# Gateway Domain Models

All Gateway models are located in `internal/domain/gateway/` directory.

## API Client Model (gateway/client.go)

```go
package gateway

import (
    "time"
    "gebase/internal/domain"
)

type ClientStatus string

const (
    ClientStatusActive    ClientStatus = "active"
    ClientStatusInactive  ClientStatus = "inactive"
    ClientStatusRevoked   ClientStatus = "revoked"
    ClientStatusSuspended ClientStatus = "suspended"
)

type ClientType string

const (
    ClientTypeConfidential ClientType = "confidential"  // Server-side apps
    ClientTypePublic       ClientType = "public"        // Mobile/SPA apps
)

// Client represents an OAuth/API client (3rd party apps accessing our APIs)
type Client struct {
    ID             int64        `json:"id" gorm:"primaryKey;autoIncrement"`
    ClientID       string       `json:"client_id" gorm:"unique;type:varchar(100)"`         // Public identifier
    ClientSecret   string       `json:"-" gorm:"type:varchar(255)"`                        // Hashed secret
    Name           string       `json:"name" gorm:"type:varchar(255)"`
    Description    string       `json:"description" gorm:"type:text"`
    ClientType     ClientType   `json:"client_type" gorm:"type:varchar(50)"`
    Status         ClientStatus `json:"status" gorm:"type:varchar(50);default:'active'"`
    
    // URLs
    RedirectURIs   string       `json:"redirect_uris" gorm:"type:jsonb"`                   // ["https://app.example.com/callback"]
    WebhookURL     string       `json:"webhook_url" gorm:"type:varchar(500)"`
    
    // Scopes & Permissions
    AllowedScopes  string       `json:"allowed_scopes" gorm:"type:jsonb"`                  // ["read", "write", "admin"]
    AllowedIPs     string       `json:"allowed_ips" gorm:"type:jsonb"`                     // IP whitelist
    
    // Rate Limiting
    RateLimitID    *int         `json:"rate_limit_id"`
    RateLimit      *RateLimit   `json:"rate_limit,omitempty" gorm:"foreignKey:RateLimitID"`
    
    // Organization ownership
    OrganizationID int64        `json:"organization_id"`
    Organization   *domain.Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
    
    // Audit
    LastUsedAt     *time.Time   `json:"last_used_at"`
    RevokedAt      *time.Time   `json:"revoked_at"`
    RevokedBy      *int64       `json:"revoked_by"`
    RevokedReason  string       `json:"revoked_reason" gorm:"type:varchar(255)"`
    
    domain.ExtraFields
}

func (Client) TableName() string {
    return "gateway_clients"
}
```

## Client Token Model (gateway/client_token.go)

```go
package gateway

import (
    "time"
    "gebase/internal/domain"
)

// ClientToken represents tokens issued to clients
type ClientToken struct {
    ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
    ClientID     int64     `json:"client_id"`
    Client       *Client   `json:"client,omitempty" gorm:"foreignKey:ClientID"`
    AccessToken  string    `json:"access_token" gorm:"unique;type:varchar(500)"`
    RefreshToken string    `json:"refresh_token" gorm:"unique;type:varchar(500)"`
    Scopes       string    `json:"scopes" gorm:"type:jsonb"`
    ExpiresAt    time.Time `json:"expires_at"`
    IssuedAt     time.Time `json:"issued_at"`
    RevokedAt    *time.Time `json:"revoked_at"`
    IPAddress    string    `json:"ip_address" gorm:"type:varchar(45)"`
    UserAgent    string    `json:"user_agent" gorm:"type:varchar(500)"`
    domain.ExtraFields
}

func (ClientToken) TableName() string {
    return "gateway_client_tokens"
}
```

## API Endpoint Model (gateway/endpoint.go)

```go
package gateway

import "gebase/internal/domain"

type HTTPMethod string

const (
    MethodGET    HTTPMethod = "GET"
    MethodPOST   HTTPMethod = "POST"
    MethodPUT    HTTPMethod = "PUT"
    MethodPATCH  HTTPMethod = "PATCH"
    MethodDELETE HTTPMethod = "DELETE"
)

type EndpointStatus string

const (
    EndpointStatusActive     EndpointStatus = "active"
    EndpointStatusInactive   EndpointStatus = "inactive"
    EndpointStatusDeprecated EndpointStatus = "deprecated"
)

// Endpoint represents an API endpoint we expose
type Endpoint struct {
    ID             int            `json:"id" gorm:"primaryKey"`
    Code           string         `json:"code" gorm:"unique;type:varchar(100)"`
    Name           string         `json:"name" gorm:"type:varchar(255)"`
    Description    string         `json:"description" gorm:"type:text"`
    Path           string         `json:"path" gorm:"type:varchar(255)"`                    // /api/v1/users
    Method         HTTPMethod     `json:"method" gorm:"type:varchar(10)"`
    Status         EndpointStatus `json:"status" gorm:"type:varchar(50);default:'active'"`
    
    // Backend target
    TargetURL      string         `json:"target_url" gorm:"type:varchar(500)"`              // Internal service URL
    Timeout        int            `json:"timeout" gorm:"default:30"`                        // seconds
    
    // Request/Response transformation
    RequestSchema  string         `json:"request_schema" gorm:"type:jsonb"`                 // JSON Schema
    ResponseSchema string         `json:"response_schema" gorm:"type:jsonb"`
    RequestTransform  string      `json:"request_transform" gorm:"type:text"`               // Transformation template
    ResponseTransform string      `json:"response_transform" gorm:"type:text"`
    
    // Auth & Rate Limiting
    RequiredScopes string         `json:"required_scopes" gorm:"type:jsonb"`                // ["read:users"]
    RateLimitID    *int           `json:"rate_limit_id"`
    RateLimit      *RateLimit     `json:"rate_limit,omitempty" gorm:"foreignKey:RateLimitID"`
    
    // Documentation
    Tags           string         `json:"tags" gorm:"type:jsonb"`
    Examples       string         `json:"examples" gorm:"type:jsonb"`
    
    // Versioning
    Version        string         `json:"version" gorm:"type:varchar(20)"`                  // v1, v2
    DeprecatedAt   *string        `json:"deprecated_at"`
    SunsetAt       *string        `json:"sunset_at"`
    
    OrganizationID *int64         `json:"organization_id,omitempty"`
    IsPublic       *bool          `json:"is_public" gorm:"default:false"`                   // Public API (no auth)
    IsActive       *bool          `json:"is_active" gorm:"default:true"`
    
    domain.ExtraFields
}

func (Endpoint) TableName() string {
    return "gateway_endpoints"
}
```

## Integration Model (gateway/integration.go)

```go
package gateway

import (
    "time"
    "gebase/internal/domain"
)

type IntegrationType string

const (
    IntegrationTypeREST    IntegrationType = "rest"
    IntegrationTypeSOAP    IntegrationType = "soap"
    IntegrationTypeGraphQL IntegrationType = "graphql"
    IntegrationTypeGRPC    IntegrationType = "grpc"
)

type IntegrationStatus string

const (
    IntegrationStatusActive   IntegrationStatus = "active"
    IntegrationStatusInactive IntegrationStatus = "inactive"
    IntegrationStatusError    IntegrationStatus = "error"
)

type AuthType string

const (
    AuthTypeNone        AuthType = "none"
    AuthTypeBasic       AuthType = "basic"
    AuthTypeAPIKey      AuthType = "api_key"
    AuthTypeOAuth2      AuthType = "oauth2"
    AuthTypeBearerToken AuthType = "bearer"
    AuthTypeCustom      AuthType = "custom"
)

// Integration represents a 3rd party API connection
type Integration struct {
    ID              int               `json:"id" gorm:"primaryKey"`
    Code            string            `json:"code" gorm:"unique;type:varchar(100)"`
    Name            string            `json:"name" gorm:"type:varchar(255)"`
    Description     string            `json:"description" gorm:"type:text"`
    Type            IntegrationType   `json:"type" gorm:"type:varchar(50)"`
    Status          IntegrationStatus `json:"status" gorm:"type:varchar(50);default:'active'"`
    
    // Connection
    BaseURL         string            `json:"base_url" gorm:"type:varchar(500)"`
    Timeout         int               `json:"timeout" gorm:"default:30"`
    RetryCount      int               `json:"retry_count" gorm:"default:3"`
    RetryDelay      int               `json:"retry_delay" gorm:"default:1000"`               // ms
    
    // Authentication
    AuthType        AuthType          `json:"auth_type" gorm:"type:varchar(50)"`
    CredentialID    *int              `json:"credential_id"`
    Credential      *Credential       `json:"credential,omitempty" gorm:"foreignKey:CredentialID"`
    
    // Headers & Config
    DefaultHeaders  string            `json:"default_headers" gorm:"type:jsonb"`             // {"X-API-Key": "..."}
    Config          string            `json:"config" gorm:"type:jsonb"`                      // Custom config
    
    // Health Check
    HealthCheckURL  string            `json:"health_check_url" gorm:"type:varchar(500)"`
    HealthCheckInterval int           `json:"health_check_interval" gorm:"default:60"`       // seconds
    LastHealthCheck *time.Time        `json:"last_health_check"`
    LastHealthStatus string           `json:"last_health_status" gorm:"type:varchar(50)"`
    
    // Rate Limiting (respect 3rd party limits)
    RateLimitPerSecond int            `json:"rate_limit_per_second" gorm:"default:10"`
    RateLimitPerMinute int            `json:"rate_limit_per_minute" gorm:"default:100"`
    
    OrganizationID  int64             `json:"organization_id"`
    Organization    *domain.Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
    
    domain.ExtraFields
}

func (Integration) TableName() string {
    return "gateway_integrations"
}
```

## Credential Model (gateway/credential.go)

```go
package gateway

import "gebase/internal/domain"

type CredentialType string

const (
    CredentialTypeAPIKey       CredentialType = "api_key"
    CredentialTypeBasicAuth    CredentialType = "basic_auth"
    CredentialTypeOAuth2       CredentialType = "oauth2"
    CredentialTypeBearerToken  CredentialType = "bearer_token"
    CredentialTypeCertificate  CredentialType = "certificate"
    CredentialTypeCustom       CredentialType = "custom"
)

// Credential stores API keys, secrets, tokens for integrations
type Credential struct {
    ID             int            `json:"id" gorm:"primaryKey"`
    Code           string         `json:"code" gorm:"unique;type:varchar(100)"`
    Name           string         `json:"name" gorm:"type:varchar(255)"`
    Description    string         `json:"description" gorm:"type:text"`
    Type           CredentialType `json:"type" gorm:"type:varchar(50)"`
    
    // Encrypted values (stored encrypted in DB)
    APIKey         string         `json:"-" gorm:"type:text"`                               // Encrypted
    APISecret      string         `json:"-" gorm:"type:text"`                               // Encrypted
    Username       string         `json:"-" gorm:"type:varchar(255)"`                       // Encrypted
    Password       string         `json:"-" gorm:"type:text"`                               // Encrypted
    AccessToken    string         `json:"-" gorm:"type:text"`                               // Encrypted
    RefreshToken   string         `json:"-" gorm:"type:text"`                               // Encrypted
    Certificate    string         `json:"-" gorm:"type:text"`                               // Encrypted PEM
    PrivateKey     string         `json:"-" gorm:"type:text"`                               // Encrypted
    CustomData     string         `json:"-" gorm:"type:text"`                               // Encrypted JSON
    
    // OAuth2 specific
    TokenURL       string         `json:"token_url" gorm:"type:varchar(500)"`
    Scopes         string         `json:"scopes" gorm:"type:jsonb"`
    ExpiresAt      *string        `json:"expires_at"`
    
    // Metadata
    Environment    string         `json:"environment" gorm:"type:varchar(50)"`              // dev, staging, prod
    
    OrganizationID int64          `json:"organization_id"`
    Organization   *domain.Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
    
    IsActive       *bool          `json:"is_active" gorm:"default:true"`
    domain.ExtraFields
}

func (Credential) TableName() string {
    return "gateway_credentials"
}
```

## Webhook Model (gateway/webhook.go)

```go
package gateway

import (
    "time"
    "gebase/internal/domain"
)

type WebhookDirection string

const (
    WebhookIncoming WebhookDirection = "incoming"  // Receiving webhooks
    WebhookOutgoing WebhookDirection = "outgoing"  // Sending webhooks
)

type WebhookStatus string

const (
    WebhookStatusActive   WebhookStatus = "active"
    WebhookStatusInactive WebhookStatus = "inactive"
    WebhookStatusFailed   WebhookStatus = "failed"   // Too many failures
)

// Webhook represents webhook configuration
type Webhook struct {
    ID             int              `json:"id" gorm:"primaryKey"`
    Code           string           `json:"code" gorm:"unique;type:varchar(100)"`
    Name           string           `json:"name" gorm:"type:varchar(255)"`
    Description    string           `json:"description" gorm:"type:text"`
    Direction      WebhookDirection `json:"direction" gorm:"type:varchar(50)"`
    Status         WebhookStatus    `json:"status" gorm:"type:varchar(50);default:'active'"`
    
    // For outgoing webhooks
    TargetURL      string           `json:"target_url" gorm:"type:varchar(500)"`
    Events         string           `json:"events" gorm:"type:jsonb"`                        // ["user.created", "order.completed"]
    Headers        string           `json:"headers" gorm:"type:jsonb"`                       // Custom headers
    
    // For incoming webhooks
    EndpointPath   string           `json:"endpoint_path" gorm:"type:varchar(255)"`          // /webhooks/stripe
    
    // Security
    Secret         string           `json:"-" gorm:"type:varchar(255)"`                      // Signing secret
    SignatureHeader string          `json:"signature_header" gorm:"type:varchar(100)"`       // X-Webhook-Signature
    
    // Retry configuration
    RetryCount     int              `json:"retry_count" gorm:"default:3"`
    RetryDelay     int              `json:"retry_delay" gorm:"default:60"`                   // seconds
    
    // Stats
    LastTriggeredAt *time.Time      `json:"last_triggered_at"`
    SuccessCount   int64            `json:"success_count" gorm:"default:0"`
    FailureCount   int64            `json:"failure_count" gorm:"default:0"`
    
    OrganizationID int64            `json:"organization_id"`
    Organization   *domain.Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
    
    IsActive       *bool            `json:"is_active" gorm:"default:true"`
    domain.ExtraFields
}

func (Webhook) TableName() string {
    return "gateway_webhooks"
}
```

## Webhook Delivery Log (gateway/webhook_delivery.go)

```go
package gateway

import (
    "time"
    "gebase/internal/domain"
)

type DeliveryStatus string

const (
    DeliveryStatusPending   DeliveryStatus = "pending"
    DeliveryStatusSuccess   DeliveryStatus = "success"
    DeliveryStatusFailed    DeliveryStatus = "failed"
    DeliveryStatusRetrying  DeliveryStatus = "retrying"
)

type WebhookDelivery struct {
    ID             int64          `json:"id" gorm:"primaryKey;autoIncrement"`
    WebhookID      int            `json:"webhook_id"`
    Webhook        *Webhook       `json:"webhook,omitempty" gorm:"foreignKey:WebhookID"`
    Event          string         `json:"event" gorm:"type:varchar(100)"`
    Status         DeliveryStatus `json:"status" gorm:"type:varchar(50)"`
    
    // Request
    RequestURL     string         `json:"request_url" gorm:"type:varchar(500)"`
    RequestHeaders string         `json:"request_headers" gorm:"type:jsonb"`
    RequestBody    string         `json:"request_body" gorm:"type:text"`
    
    // Response
    ResponseStatus int            `json:"response_status"`
    ResponseHeaders string        `json:"response_headers" gorm:"type:jsonb"`
    ResponseBody   string         `json:"response_body" gorm:"type:text"`
    
    // Timing
    AttemptCount   int            `json:"attempt_count" gorm:"default:1"`
    Duration       int64          `json:"duration"`                                         // ms
    TriggeredAt    time.Time      `json:"triggered_at"`
    DeliveredAt    *time.Time     `json:"delivered_at"`
    NextRetryAt    *time.Time     `json:"next_retry_at"`
    
    ErrorMessage   string         `json:"error_message" gorm:"type:text"`
    
    domain.ExtraFields
}

func (WebhookDelivery) TableName() string {
    return "gateway_webhook_deliveries"
}
```

## Rate Limit Model (gateway/rate_limit.go)

```go
package gateway

import "gebase/internal/domain"

type RateLimitType string

const (
    RateLimitTypeFixed   RateLimitType = "fixed"    // Fixed window
    RateLimitTypeSliding RateLimitType = "sliding"  // Sliding window
    RateLimitTypeToken   RateLimitType = "token"    // Token bucket
)

// RateLimit defines rate limiting rules
type RateLimit struct {
    ID             int           `json:"id" gorm:"primaryKey"`
    Code           string        `json:"code" gorm:"unique;type:varchar(100)"`
    Name           string        `json:"name" gorm:"type:varchar(255)"`
    Description    string        `json:"description" gorm:"type:text"`
    Type           RateLimitType `json:"type" gorm:"type:varchar(50)"`
    
    // Limits
    RequestsPerSecond int        `json:"requests_per_second" gorm:"default:10"`
    RequestsPerMinute int        `json:"requests_per_minute" gorm:"default:100"`
    RequestsPerHour   int        `json:"requests_per_hour" gorm:"default:1000"`
    RequestsPerDay    int        `json:"requests_per_day" gorm:"default:10000"`
    
    // Burst
    BurstSize      int           `json:"burst_size" gorm:"default:50"`
    
    // Response when exceeded
    ExceededStatus int           `json:"exceeded_status" gorm:"default:429"`
    ExceededMessage string       `json:"exceeded_message" gorm:"type:varchar(255)"`
    
    IsActive       *bool         `json:"is_active" gorm:"default:true"`
    domain.ExtraFields
}

func (RateLimit) TableName() string {
    return "gateway_rate_limits"
}
```

## API Log Model (gateway/api_log.go)

```go
package gateway

import (
    "time"
    "gebase/internal/domain"
)

// APILog stores API call logs for monitoring
type APILog struct {
    ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
    
    // Request info
    ClientID       *int64     `json:"client_id"`
    EndpointID     *int       `json:"endpoint_id"`
    IntegrationID  *int       `json:"integration_id"`
    
    Method         string     `json:"method" gorm:"type:varchar(10)"`
    Path           string     `json:"path" gorm:"type:varchar(500)"`
    QueryParams    string     `json:"query_params" gorm:"type:jsonb"`
    RequestHeaders string     `json:"request_headers" gorm:"type:jsonb"`
    RequestBody    string     `json:"request_body" gorm:"type:text"`
    RequestSize    int64      `json:"request_size"`
    
    // Response info
    StatusCode     int        `json:"status_code"`
    ResponseHeaders string    `json:"response_headers" gorm:"type:jsonb"`
    ResponseBody   string     `json:"response_body" gorm:"type:text"`
    ResponseSize   int64      `json:"response_size"`
    
    // Timing & Performance
    Duration       int64      `json:"duration"`                                            // ms
    Timestamp      time.Time  `json:"timestamp" gorm:"index"`
    
    // Client info
    IPAddress      string     `json:"ip_address" gorm:"type:varchar(45)"`
    UserAgent      string     `json:"user_agent" gorm:"type:varchar(500)"`
    
    // Error tracking
    ErrorCode      string     `json:"error_code" gorm:"type:varchar(50)"`
    ErrorMessage   string     `json:"error_message" gorm:"type:text"`
    
    // Organization context
    OrganizationID *int64     `json:"organization_id"`
    
    domain.ExtraFields
}

func (APILog) TableName() string {
    return "gateway_api_logs"
}

// Note: Consider partitioning this table by timestamp for large volumes
// CREATE TABLE gateway_api_logs (...) PARTITION BY RANGE (timestamp);
```

## Gateway Module Summary

| Module       | Description                           | Key Features                              |
|--------------|---------------------------------------|-------------------------------------------|
| client       | OAuth/API clients                     | Client ID/Secret, scopes, revocation      |
| endpoint     | API endpoints to expose               | Path mapping, transformation, versioning  |
| integration  | 3rd party API connections             | Connection pooling, health check, retry   |
| credential   | Secrets & API keys storage            | Encrypted storage, rotation               |
| webhook      | Incoming/outgoing webhooks            | Signature verification, retry             |
| ratelimit    | Rate limiting rules                   | Fixed/sliding window, token bucket        |
| apilog       | API call audit trail                  | Request/response logging, analytics       |
| monitor      | Health & performance monitoring       | Latency, errors, uptime                   |
