# Gerege Gebase Platform MVP - Claude Code Prompt

Enterprise-grade Auth & RBAC platform with multi-platform device support, organization hierarchy, granular permission management, and system switching capability.

## Tech Stack
- **Backend:** Go 1.22+, Gin framework, GORM
- **Database:** PostgreSQL 15+ (schema: gerege_base)
- **Frontend:** Next.js 14+ (App Router)
- **Docs:** Swagger/OpenAPI 3.0
- **Containerization:** Docker & Docker Compose

## Architecture Overview

```
┌────────────────────────────────────────────────────────────┐
│                    GEBASE PLATFORM                          │
│ ┌────────────────────────────────────────────────────────┐ │
│ │              Platform Layer (Core)                     │ │
│ │                 system_id = NULL                       │ │
│ │   User, Organization, Role, Permission, Device, etc.  │ │
│ └────────────────────────────────────────────────────────┘ │
│                                                            │
│                     System Switcher                        │
│                  ┌────────┴────────┐                       │
│                  ▼                 ▼                       │
│  ┌──────────────────────┐ ┌──────────────────────┐        │
│  │     DSL System       │ │   Gateway System     │        │
│  │    system_id = 1     │ │    system_id = 2     │        │
│  │  Schema, Workflow,   │ │  API Client, Integ., │        │
│  │  Rule, Template      │ │  Webhook, RateLimit  │        │
│  └──────────────────────┘ └──────────────────────┘        │
└────────────────────────────────────────────────────────────┘
```

## Systems Overview

| Layer/System | Code     | Description                              |
|--------------|----------|------------------------------------------|
| Platform     | (core)   | User, Org, Role, Permission, Device      |
| DSL          | dsl      | Domain Specific Language Engine          |
| Gateway      | gateway  | API Management & Integration Hub         |

## Multi-Platform Device Support

| Platform        | Technology | Device Code     |
|-----------------|------------|-----------------|
| Web Browser     | Next.js    | web             |
| iOS App         | Swift      | ios             |
| Android App     | Kotlin     | android         |
| iPad            | Swift      | tablet_ios      |
| Android Tablet  | Kotlin     | tablet_android  |
| Windows Desktop | C# WPF     | windows_desktop |
| Mac Desktop     | Swift      | mac_desktop     |
| Gerege Kiosk    | C# WPF     | kiosk           |
| Android POS     | Kotlin     | pos_android     |
| Linux POS       | Electron/Go| pos_linux       |

## Core Features

- **System Switching** - Two-token strategy (platform + system tokens)
- Төхөөрөмж бүртгэл (Device registration)
- Session tracking (Хэрэглэгч + Төхөөрөмж + System)
- Remote configuration
- Heartbeat monitoring
- Remote logout
- Multi-language support via Translations table
- Role-based dynamic menu rendering per system
- Permission format:
  - Platform: `{module.action}`
  - System: `{system.module.action}`

## Project Structure

```
gebase/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── app/
│   │   │   └── container.go          # DI container
│   │   ├── auth/
│   │   │   ├── jwt.go
│   │   │   ├── claims.go             # Platform & System claims
│   │   │   └── session.go
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── db/
│   │   │   ├── postgres.go
│   │   │   └── migrations.go
│   │   ├── domain/
│   │   │   ├── base.go
│   │   │   ├── user.go
│   │   │   ├── organization.go
│   │   │   ├── system.go
│   │   │   ├── module.go
│   │   │   ├── action.go
│   │   │   ├── permission.go
│   │   │   ├── role.go
│   │   │   ├── menu.go
│   │   │   ├── device.go
│   │   │   ├── session.go
│   │   │   ├── language.go
│   │   │   ├── dsl/                  # DSL System models
│   │   │   │   ├── schema.go
│   │   │   │   ├── field.go
│   │   │   │   ├── rule.go
│   │   │   │   ├── workflow.go
│   │   │   │   ├── template.go
│   │   │   │   ├── function.go
│   │   │   │   └── variable.go
│   │   │   └── gateway/              # Gateway System models
│   │   │       ├── client.go
│   │   │       ├── endpoint.go
│   │   │       ├── integration.go
│   │   │       ├── credential.go
│   │   │       ├── webhook.go
│   │   │       ├── rate_limit.go
│   │   │       └── api_log.go
│   │   ├── http/
│   │   │   ├── dto/
│   │   │   │   ├── request/
│   │   │   │   └── response/
│   │   │   ├── handlers/
│   │   │   └── router/
│   │   │       └── router.go
│   │   ├── middleware/
│   │   │   ├── platform_auth.go      # Platform token auth
│   │   │   ├── system_auth.go        # System token auth
│   │   │   ├── rbac.go
│   │   │   ├── device.go
│   │   │   ├── cors.go
│   │   │   └── logger.go
│   │   ├── repository/
│   │   └── service/
│   ├── docs/
│   │   └── swagger.yaml
│   ├── docker/
│   │   └── Dockerfile
│   ├── .env.example
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── admin/                        # Port 3001 (Platform + All systems admin)
│   │   ├── src/
│   │   │   ├── app/
│   │   │   │   ├── (auth)/
│   │   │   │   ├── (platform)/       # Platform management pages
│   │   │   │   ├── (dsl)/            # DSL system pages
│   │   │   │   └── (gateway)/        # Gateway system pages
│   │   │   ├── components/
│   │   │   │   └── system-switcher/  # System switching UI
│   │   │   ├── lib/
│   │   │   ├── hooks/
│   │   │   └── types/
│   │   ├── package.json
│   │   └── next.config.js
│   └── portal/                       # Port 3000 (End-user portal)
│       ├── src/
│       ├── package.json
│       └── next.config.js
├── docker-compose.yml
└── README.md
```

## Environment Variables

```env
# Server
SERVER_PORT=8000
SERVER_MODE=development
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gerege_db
DB_SCHEMA=gerege_base
DB_SSLMODE=disable
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100

# JWT
JWT_SECRET=your-super-secret-key-change-in-production
JWT_ACCESS_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# SSO
SSO_BASE_URL=https://sso.gerege.mn
SSO_CLIENT_ID=your-client-id
SSO_CLIENT_SECRET=your-client-secret

# Redis (for session caching)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

## Implementation Order

1. Set up project structure, domain models
2. Configuration and environment
3. Database connection and migrations
4. Repository layer implementation
5. Service layer with business logic
6. HTTP handlers and routing
7. Middleware (auth, RBAC, device)
8. Admin frontend with all CRUD pages
9. Portal frontend
10. Docker setup
11. Seed data and testing

## Key Requirements Checklist

- [ ] Platform layer with system_id = NULL for core modules
- [ ] System switching with two-token strategy (platform token + system token)
- [ ] All code fields in lowercase with dot notation
- [ ] ExtraFields embedded in all domain models
- [ ] NO NameEn fields - use Translations table instead
- [ ] Multi-language via translations table with key-value pairs
- [ ] 2 systems: dsl, gateway
- [ ] Device registration required before login
- [ ] Session tracking with current_system_id for system context
- [ ] Role-based dynamic menu rendering per system
- [ ] Permission code format: 
  - Platform: `{module.action}` (e.g., user.view)
  - System: `{system.module.action}` (e.g., dsl.schema.view)
- [ ] DSL system: schema, field, rule, workflow, template, function, variable
- [ ] Gateway system: client, endpoint, integration, credential, webhook, ratelimit
- [ ] Clean architecture
- [ ] Swagger documentation
- [ ] Docker containerization

## Documentation Files

1. `00-architecture.md` - Platform + Systems architecture, system switching flow
2. `01-domain-models-updated.md` - Core domain models with platform layer
3. `02-dsl-models.md` - DSL system domain models
4. `03-api-routes-updated.md` - All API endpoints with system switching
5. `04-seed-data-updated.md` - Initial seed data for platform and systems
6. `05-docker-config.md` - Docker configuration
7. `07-gateway-models.md` - Gateway system domain models
8. `08-frontend-skeleton.md` - Frontend skeleton with system switching UI
