# API Routes - Platform + Systems

## Route Structure

```
/api/v1/
├── auth/                          # Authentication (Public + Protected)
├── platform/                      # Platform layer APIs
│   ├── users/
│   ├── organizations/
│   ├── roles/
│   ├── permissions/
│   ├── devices/
│   ├── sessions/
│   └── settings/
└── systems/
    ├── {system_code}/             # System-specific APIs
    │   ├── dsl/...
    │   └── gateway/...
    └── switch                     # System switching
```

---

## Authentication Routes

### Public Routes
```
POST   /api/v1/auth/login                    # User login
POST   /api/v1/auth/refresh                  # Refresh platform token
POST   /api/v1/auth/forgot-password          # Password reset request
POST   /api/v1/auth/reset-password           # Password reset
POST   /api/v1/devices/register              # Register new device
POST   /api/v1/devices/heartbeat             # Device heartbeat (semi-public)
```

### Protected Routes (Platform Token Required)
```
POST   /api/v1/auth/logout                   # Logout from platform
GET    /api/v1/auth/me                       # Get current user info
PUT    /api/v1/auth/me                       # Update profile
PUT    /api/v1/auth/change-password          # Change password
GET    /api/v1/auth/available-systems        # Get user's available systems
POST   /api/v1/auth/switch-system            # Switch to a system
POST   /api/v1/auth/exit-system              # Exit system, return to platform
GET    /api/v1/auth/current-context          # Get current system context
```

---

## System Switching Endpoints

### GET /api/v1/auth/available-systems

Returns list of systems user has access to.

**Request Headers:**
```
Authorization: Bearer {platform_token}
X-Device-UID: device-uuid
```

**Response:**
```json
{
    "success": true,
    "data": {
        "systems": [
            {
                "id": 1,
                "code": "dsl",
                "name": "DSL System",
                "description": "Domain Specific Language",
                "icon_name": "Code",
                "color": "#6366f1",
                "roles": [
                    {
                        "id": 5,
                        "code": "dsl_developer",
                        "name": "DSL Developer",
                        "organization_id": 789,
                        "organization_name": "Gerege LLC"
                    }
                ]
            },
            {
                "id": 2,
                "code": "gateway",
                "name": "Gateway System",
                "description": "API Management",
                "icon_name": "Network",
                "color": "#22c55e",
                "roles": [
                    {
                        "id": 8,
                        "code": "gateway_admin",
                        "name": "Gateway Admin",
                        "organization_id": 789,
                        "organization_name": "Gerege LLC"
                    }
                ]
            }
        ],
        "default_system_code": "dsl"
    }
}
```

### POST /api/v1/auth/switch-system

Switch to a specific system and get system token.

**Request Headers:**
```
Authorization: Bearer {platform_token}
X-Device-UID: device-uuid
```

**Request Body:**
```json
{
    "system_code": "dsl",
    "organization_id": 789,    // Optional, if user has multiple orgs for this system
    "role_id": 5               // Optional, if user has multiple roles
}
```

**Response:**
```json
{
    "success": true,
    "data": {
        "system_token": "eyJhbGciOiJIUzI1NiIs...",
        "expires_in": 28800,
        "current_system": {
            "id": 1,
            "code": "dsl",
            "name": "DSL System"
        },
        "current_role": {
            "id": 5,
            "code": "dsl_developer",
            "name": "DSL Developer"
        },
        "current_organization": {
            "id": 789,
            "name": "Gerege LLC"
        },
        "permissions": [
            "dsl.schema.view",
            "dsl.schema.create",
            "dsl.workflow.view",
            "dsl.workflow.execute"
        ],
        "menus": [
            {
                "id": 1,
                "code": "dsl_dashboard",
                "name": "Dashboard",
                "path": "/dashboard",
                "icon": "LayoutDashboard",
                "children": []
            },
            {
                "id": 2,
                "code": "schemas",
                "name": "Schemas",
                "path": "/schemas",
                "icon": "Table",
                "children": []
            }
        ]
    }
}
```

### POST /api/v1/auth/exit-system

Exit current system and return to platform context.

**Request Headers:**
```
Authorization: Bearer {system_token}
X-Device-UID: device-uuid
```

**Response:**
```json
{
    "success": true,
    "data": {
        "platform_token": "eyJhbGciOiJIUzI1NiIs...",
        "message": "Exited from DSL System"
    }
}
```

### GET /api/v1/auth/current-context

Get current session context (platform or system).

**Response (Platform Context):**
```json
{
    "success": true,
    "data": {
        "context_type": "platform",
        "user": {
            "id": 123,
            "email": "user@example.com",
            "first_name": "Батболд"
        },
        "current_system": null,
        "available_systems": [...]
    }
}
```

**Response (System Context):**
```json
{
    "success": true,
    "data": {
        "context_type": "system",
        "user": {
            "id": 123,
            "email": "user@example.com",
            "first_name": "Батболд"
        },
        "current_system": {
            "id": 1,
            "code": "dsl",
            "name": "DSL System"
        },
        "current_role": {
            "id": 5,
            "code": "dsl_developer"
        },
        "current_organization": {
            "id": 789,
            "name": "Gerege LLC"
        },
        "permissions": [...]
    }
}
```

---

## Platform Routes (Platform Token Required)

### User Management
```
GET    /api/v1/platform/users                     # List users
POST   /api/v1/platform/users                     # Create user
GET    /api/v1/platform/users/:id                 # Get user
PUT    /api/v1/platform/users/:id                 # Update user
DELETE /api/v1/platform/users/:id                 # Soft delete user
GET    /api/v1/platform/users/:id/roles           # Get user's all system roles
PUT    /api/v1/platform/users/:id/roles           # Assign system roles
GET    /api/v1/platform/users/:id/sessions        # Get user's active sessions
```

### Organization Management
```
GET    /api/v1/platform/organizations             # List organizations
POST   /api/v1/platform/organizations             # Create organization
GET    /api/v1/platform/organizations/:id         # Get organization
PUT    /api/v1/platform/organizations/:id         # Update organization
DELETE /api/v1/platform/organizations/:id         # Soft delete
GET    /api/v1/platform/organizations/:id/children # Get children
GET    /api/v1/platform/organizations/:id/systems  # Get enabled systems
PUT    /api/v1/platform/organizations/:id/systems  # Enable/disable systems
GET    /api/v1/platform/organization-types         # Get org types
```

### System Management
```
GET    /api/v1/platform/systems                   # List all systems
POST   /api/v1/platform/systems                   # Create system
GET    /api/v1/platform/systems/:id               # Get system
PUT    /api/v1/platform/systems/:id               # Update system
GET    /api/v1/platform/systems/:id/modules       # Get system modules
GET    /api/v1/platform/systems/:id/roles         # Get system roles
GET    /api/v1/platform/systems/:id/menus         # Get system menus
```

### Role Management
```
GET    /api/v1/platform/roles                     # List roles (filter: system_id)
POST   /api/v1/platform/roles                     # Create role
GET    /api/v1/platform/roles/:id                 # Get role
PUT    /api/v1/platform/roles/:id                 # Update role
DELETE /api/v1/platform/roles/:id                 # Soft delete
GET    /api/v1/platform/roles/:id/permissions     # Get role permissions
PUT    /api/v1/platform/roles/:id/permissions     # Assign permissions
GET    /api/v1/platform/roles/:id/menus           # Get role menus
PUT    /api/v1/platform/roles/:id/menus           # Assign menus
```

### Permission Management
```
GET    /api/v1/platform/permissions               # List permissions
POST   /api/v1/platform/permissions               # Create permission
GET    /api/v1/platform/permissions/:id           # Get permission
PUT    /api/v1/platform/permissions/:id           # Update permission
POST   /api/v1/platform/permissions/sync          # Auto-sync from modules
```

### Module & Action Management
```
GET    /api/v1/platform/modules                   # List modules
POST   /api/v1/platform/modules                   # Create module
GET    /api/v1/platform/modules/:id               # Get module
PUT    /api/v1/platform/modules/:id               # Update module
GET    /api/v1/platform/actions                   # List actions
POST   /api/v1/platform/actions                   # Create action
```

### Menu Management
```
GET    /api/v1/platform/menus                     # List menus
POST   /api/v1/platform/menus                     # Create menu
GET    /api/v1/platform/menus/:id                 # Get menu
PUT    /api/v1/platform/menus/:id                 # Update menu
DELETE /api/v1/platform/menus/:id                 # Soft delete
GET    /api/v1/platform/menus/tree                # Get menu tree
PUT    /api/v1/platform/menus/reorder             # Reorder menus
```

### Device Management
```
GET    /api/v1/platform/devices                   # List devices
GET    /api/v1/platform/devices/:id               # Get device
PUT    /api/v1/platform/devices/:id               # Update device
DELETE /api/v1/platform/devices/:id               # Deactivate device
PUT    /api/v1/platform/devices/:id/config        # Update remote config
```

### Session Management
```
GET    /api/v1/platform/sessions                  # List active sessions
GET    /api/v1/platform/sessions/:id              # Get session
DELETE /api/v1/platform/sessions/:id              # Terminate session
POST   /api/v1/platform/sessions/terminate-user   # Terminate all user sessions
```

### Monitoring (Platform Admin)
```
GET    /api/v1/platform/monitoring/dashboard      # Dashboard stats
GET    /api/v1/platform/monitoring/sessions       # Active sessions
GET    /api/v1/platform/monitoring/devices        # Online devices
GET    /api/v1/platform/monitoring/health         # System health
GET    /api/v1/platform/monitoring/login-history  # Login audit
POST   /api/v1/platform/monitoring/remote-logout  # Force logout
```

### Language & Translation
```
GET    /api/v1/platform/languages                 # List languages
POST   /api/v1/platform/languages                 # Create language
GET    /api/v1/platform/translations              # List translations
POST   /api/v1/platform/translations              # Create translation
GET    /api/v1/platform/translations/export       # Export
POST   /api/v1/platform/translations/import       # Import
```

---

## System Routes (System Token Required)

### DSL System - /api/v1/systems/dsl/*

```
# Schema
GET    /api/v1/systems/dsl/schemas                # List schemas
POST   /api/v1/systems/dsl/schemas                # Create schema
GET    /api/v1/systems/dsl/schemas/:id            # Get schema
PUT    /api/v1/systems/dsl/schemas/:id            # Update schema
DELETE /api/v1/systems/dsl/schemas/:id            # Delete schema
POST   /api/v1/systems/dsl/schemas/:id/deploy     # Deploy schema
POST   /api/v1/systems/dsl/schemas/:id/clone      # Clone schema
GET    /api/v1/systems/dsl/schemas/:id/fields     # Get fields

# Field
GET    /api/v1/systems/dsl/fields                 # List fields
POST   /api/v1/systems/dsl/fields                 # Create field
PUT    /api/v1/systems/dsl/fields/:id             # Update field
DELETE /api/v1/systems/dsl/fields/:id             # Delete field

# Rule
GET    /api/v1/systems/dsl/rules                  # List rules
POST   /api/v1/systems/dsl/rules                  # Create rule
GET    /api/v1/systems/dsl/rules/:id              # Get rule
PUT    /api/v1/systems/dsl/rules/:id              # Update rule
DELETE /api/v1/systems/dsl/rules/:id              # Delete rule
POST   /api/v1/systems/dsl/rules/:id/test         # Test rule

# Workflow
GET    /api/v1/systems/dsl/workflows              # List workflows
POST   /api/v1/systems/dsl/workflows              # Create workflow
GET    /api/v1/systems/dsl/workflows/:id          # Get workflow
PUT    /api/v1/systems/dsl/workflows/:id          # Update workflow
DELETE /api/v1/systems/dsl/workflows/:id          # Delete workflow
POST   /api/v1/systems/dsl/workflows/:id/publish  # Publish
POST   /api/v1/systems/dsl/workflows/:id/execute  # Execute

# Template
GET    /api/v1/systems/dsl/templates              # List templates
POST   /api/v1/systems/dsl/templates              # Create template
GET    /api/v1/systems/dsl/templates/:id          # Get template
PUT    /api/v1/systems/dsl/templates/:id          # Update template
POST   /api/v1/systems/dsl/templates/:id/render   # Render template

# Function
GET    /api/v1/systems/dsl/functions              # List functions
POST   /api/v1/systems/dsl/functions              # Create function
GET    /api/v1/systems/dsl/functions/:id          # Get function
POST   /api/v1/systems/dsl/functions/:id/test     # Test function
GET    /api/v1/systems/dsl/functions/builtin      # Built-in functions

# Variable
GET    /api/v1/systems/dsl/variables              # List variables
POST   /api/v1/systems/dsl/variables              # Create variable
PUT    /api/v1/systems/dsl/variables/:id          # Update variable
DELETE /api/v1/systems/dsl/variables/:id          # Delete variable

# Logs
GET    /api/v1/systems/dsl/logs                   # Execution logs
GET    /api/v1/systems/dsl/logs/stats             # Statistics
```

### Gateway System - /api/v1/systems/gateway/*

```
# API Clients (OAuth clients)
GET    /api/v1/systems/gateway/clients            # List clients
POST   /api/v1/systems/gateway/clients            # Create client
GET    /api/v1/systems/gateway/clients/:id        # Get client
PUT    /api/v1/systems/gateway/clients/:id        # Update client
DELETE /api/v1/systems/gateway/clients/:id        # Delete client
POST   /api/v1/systems/gateway/clients/:id/regenerate-secret  # Regenerate secret
POST   /api/v1/systems/gateway/clients/:id/revoke             # Revoke client

# API Endpoints (Own APIs to expose)
GET    /api/v1/systems/gateway/endpoints          # List endpoints
POST   /api/v1/systems/gateway/endpoints          # Create endpoint
GET    /api/v1/systems/gateway/endpoints/:id      # Get endpoint
PUT    /api/v1/systems/gateway/endpoints/:id      # Update endpoint
DELETE /api/v1/systems/gateway/endpoints/:id      # Delete endpoint
POST   /api/v1/systems/gateway/endpoints/:id/test # Test endpoint

# Integrations (3rd party APIs)
GET    /api/v1/systems/gateway/integrations       # List integrations
POST   /api/v1/systems/gateway/integrations       # Create integration
GET    /api/v1/systems/gateway/integrations/:id   # Get integration
PUT    /api/v1/systems/gateway/integrations/:id   # Update integration
DELETE /api/v1/systems/gateway/integrations/:id   # Delete integration
POST   /api/v1/systems/gateway/integrations/:id/test    # Test connection
GET    /api/v1/systems/gateway/integrations/:id/health  # Health check

# Credentials (API keys, secrets)
GET    /api/v1/systems/gateway/credentials        # List credentials
POST   /api/v1/systems/gateway/credentials        # Create credential
GET    /api/v1/systems/gateway/credentials/:id    # Get credential
PUT    /api/v1/systems/gateway/credentials/:id    # Update credential
DELETE /api/v1/systems/gateway/credentials/:id    # Delete credential

# Webhooks
GET    /api/v1/systems/gateway/webhooks           # List webhooks
POST   /api/v1/systems/gateway/webhooks           # Create webhook
GET    /api/v1/systems/gateway/webhooks/:id       # Get webhook
PUT    /api/v1/systems/gateway/webhooks/:id       # Update webhook
DELETE /api/v1/systems/gateway/webhooks/:id       # Delete webhook
GET    /api/v1/systems/gateway/webhooks/:id/logs  # Webhook delivery logs

# Rate Limiting
GET    /api/v1/systems/gateway/rate-limits        # List rate limit rules
POST   /api/v1/systems/gateway/rate-limits        # Create rule
PUT    /api/v1/systems/gateway/rate-limits/:id    # Update rule
DELETE /api/v1/systems/gateway/rate-limits/:id    # Delete rule

# Monitoring
GET    /api/v1/systems/gateway/monitoring/overview       # Overview
GET    /api/v1/systems/gateway/monitoring/api-calls      # API call logs
GET    /api/v1/systems/gateway/monitoring/errors         # Error logs
GET    /api/v1/systems/gateway/monitoring/latency        # Latency metrics
GET    /api/v1/systems/gateway/monitoring/integrations   # Integration health
```

---

## Router Implementation

```go
package router

import (
    "github.com/gin-gonic/gin"
    "gebase/internal/http/handlers"
    "gebase/internal/middleware"
)

func SetupRouter(h *handlers.Handlers, m *middleware.Middleware) *gin.Engine {
    r := gin.Default()
    
    // Global middlewares
    r.Use(m.CORS())
    r.Use(m.Logger())
    r.Use(m.RequestID())
    
    api := r.Group("/api/v1")
    
    // ========== Public Routes ==========
    public := api.Group("")
    {
        public.POST("/auth/login", h.Auth.Login)
        public.POST("/auth/refresh", h.Auth.Refresh)
        public.POST("/auth/forgot-password", h.Auth.ForgotPassword)
        public.POST("/auth/reset-password", h.Auth.ResetPassword)
        public.POST("/devices/register", h.Device.Register)
        public.POST("/devices/heartbeat", h.Device.Heartbeat)
    }
    
    // ========== Platform Protected Routes ==========
    platform := api.Group("")
    platform.Use(m.PlatformAuth()) // Requires platform token
    platform.Use(m.Device())
    {
        // Auth & System Switching
        platform.POST("/auth/logout", h.Auth.Logout)
        platform.GET("/auth/me", h.Auth.Me)
        platform.PUT("/auth/me", h.Auth.UpdateProfile)
        platform.PUT("/auth/change-password", h.Auth.ChangePassword)
        platform.GET("/auth/available-systems", h.Auth.AvailableSystems)
        platform.POST("/auth/switch-system", h.Auth.SwitchSystem)
        platform.GET("/auth/current-context", h.Auth.CurrentContext)
        
        // Platform Management
        pf := platform.Group("/platform")
        {
            // Users
            pf.GET("/users", m.RequirePlatformPermission("user.view"), h.User.List)
            pf.POST("/users", m.RequirePlatformPermission("user.create"), h.User.Create)
            pf.GET("/users/:id", m.RequirePlatformPermission("user.view"), h.User.Get)
            pf.PUT("/users/:id", m.RequirePlatformPermission("user.update"), h.User.Update)
            pf.DELETE("/users/:id", m.RequirePlatformPermission("user.delete"), h.User.Delete)
            pf.GET("/users/:id/roles", m.RequirePlatformPermission("user.view"), h.User.GetRoles)
            pf.PUT("/users/:id/roles", m.RequirePlatformPermission("role.manage"), h.User.AssignRoles)
            
            // Organizations
            pf.GET("/organizations", m.RequirePlatformPermission("organization.view"), h.Organization.List)
            pf.POST("/organizations", m.RequirePlatformPermission("organization.create"), h.Organization.Create)
            pf.GET("/organizations/:id", m.RequirePlatformPermission("organization.view"), h.Organization.Get)
            pf.PUT("/organizations/:id", m.RequirePlatformPermission("organization.update"), h.Organization.Update)
            pf.DELETE("/organizations/:id", m.RequirePlatformPermission("organization.delete"), h.Organization.Delete)
            pf.GET("/organizations/:id/systems", m.RequirePlatformPermission("organization.view"), h.Organization.GetSystems)
            pf.PUT("/organizations/:id/systems", m.RequirePlatformPermission("organization.update"), h.Organization.UpdateSystems)
            
            // Systems, Roles, Permissions, Modules, Actions, Menus
            // ... similar pattern
            
            // Monitoring
            monitor := pf.Group("/monitoring")
            monitor.Use(m.RequirePlatformPermission("monitoring.view"))
            {
                monitor.GET("/dashboard", h.Monitoring.Dashboard)
                monitor.GET("/sessions", h.Monitoring.Sessions)
                monitor.GET("/devices", h.Monitoring.Devices)
                monitor.GET("/health", h.Monitoring.Health)
                monitor.POST("/remote-logout", m.RequirePlatformPermission("session.terminate"), h.Monitoring.RemoteLogout)
            }
        }
    }
    
    // ========== System Protected Routes ==========
    systems := api.Group("/systems")
    systems.Use(m.SystemAuth()) // Requires system token
    systems.Use(m.Device())
    {
        // Exit system (available for all systems)
        systems.POST("/exit", h.Auth.ExitSystem)
        
        // DSL System Routes
        dsl := systems.Group("/dsl")
        dsl.Use(m.RequireSystem("dsl")) // Verify token is for DSL system
        {
            dsl.GET("/schemas", m.RequireSystemPermission("dsl.schema.view"), h.DSL.ListSchemas)
            dsl.POST("/schemas", m.RequireSystemPermission("dsl.schema.create"), h.DSL.CreateSchema)
            dsl.GET("/schemas/:id", m.RequireSystemPermission("dsl.schema.view"), h.DSL.GetSchema)
            dsl.PUT("/schemas/:id", m.RequireSystemPermission("dsl.schema.update"), h.DSL.UpdateSchema)
            dsl.DELETE("/schemas/:id", m.RequireSystemPermission("dsl.schema.delete"), h.DSL.DeleteSchema)
            dsl.POST("/schemas/:id/deploy", m.RequireSystemPermission("dsl.schema.execute"), h.DSL.DeploySchema)
            
            // Rules
            dsl.GET("/rules", m.RequireSystemPermission("dsl.rule.view"), h.DSL.ListRules)
            dsl.POST("/rules", m.RequireSystemPermission("dsl.rule.create"), h.DSL.CreateRule)
            dsl.POST("/rules/:id/test", m.RequireSystemPermission("dsl.rule.execute"), h.DSL.TestRule)
            
            // Workflows
            dsl.GET("/workflows", m.RequireSystemPermission("dsl.workflow.view"), h.DSL.ListWorkflows)
            dsl.POST("/workflows", m.RequireSystemPermission("dsl.workflow.create"), h.DSL.CreateWorkflow)
            dsl.POST("/workflows/:id/publish", m.RequireSystemPermission("dsl.workflow.publish"), h.DSL.PublishWorkflow)
            dsl.POST("/workflows/:id/execute", m.RequireSystemPermission("dsl.workflow.execute"), h.DSL.ExecuteWorkflow)
            
            // Templates, Functions, Variables, Logs...
        }
        
        // Gateway System Routes
        gateway := systems.Group("/gateway")
        gateway.Use(m.RequireSystem("gateway"))
        {
            // Clients
            gateway.GET("/clients", m.RequireSystemPermission("gateway.client.view"), h.Gateway.ListClients)
            gateway.POST("/clients", m.RequireSystemPermission("gateway.client.create"), h.Gateway.CreateClient)
            gateway.GET("/clients/:id", m.RequireSystemPermission("gateway.client.view"), h.Gateway.GetClient)
            gateway.PUT("/clients/:id", m.RequireSystemPermission("gateway.client.update"), h.Gateway.UpdateClient)
            gateway.POST("/clients/:id/regenerate-secret", m.RequireSystemPermission("gateway.client.update"), h.Gateway.RegenerateSecret)
            gateway.POST("/clients/:id/revoke", m.RequireSystemPermission("gateway.client.delete"), h.Gateway.RevokeClient)
            
            // Integrations
            gateway.GET("/integrations", m.RequireSystemPermission("gateway.integration.view"), h.Gateway.ListIntegrations)
            gateway.POST("/integrations", m.RequireSystemPermission("gateway.integration.create"), h.Gateway.CreateIntegration)
            gateway.POST("/integrations/:id/test", m.RequireSystemPermission("gateway.integration.execute"), h.Gateway.TestIntegration)
            
            // Webhooks, Rate Limits, Monitoring...
        }
    }
    
    return r
}
```

---

## Middleware Types

```go
package middleware

// PlatformAuth validates platform token
func (m *Middleware) PlatformAuth() gin.HandlerFunc

// SystemAuth validates system token
func (m *Middleware) SystemAuth() gin.HandlerFunc

// RequireSystem verifies the token belongs to specific system
func (m *Middleware) RequireSystem(systemCode string) gin.HandlerFunc

// RequirePlatformPermission checks platform-level permission
func (m *Middleware) RequirePlatformPermission(permission string) gin.HandlerFunc

// RequireSystemPermission checks system-level permission (from token cache)
func (m *Middleware) RequireSystemPermission(permission string) gin.HandlerFunc
```
