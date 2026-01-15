# Seed Data - Platform + Systems

## Helper

```go
func ptr[T any](v T) *T {
    return &v
}
```

## Languages

```go
languages := []domain.Language{
    {ID: 1, Code: "mn", Name: "Mongolian", NativeName: "Монгол", IsActive: ptr(true), IsDefault: ptr(true)},
    {ID: 2, Code: "en", Name: "English", NativeName: "English", IsActive: ptr(true), IsDefault: ptr(false)},
}
```

## Actions

```go
actions := []domain.Action{
    {ID: 1, Code: "view", Name: "Харах"},
    {ID: 2, Code: "create", Name: "Үүсгэх"},
    {ID: 3, Code: "update", Name: "Засах"},
    {ID: 4, Code: "delete", Name: "Устгах"},
    {ID: 5, Code: "export", Name: "Экспорт"},
    {ID: 6, Code: "import", Name: "Импорт"},
    {ID: 7, Code: "approve", Name: "Батлах"},
    {ID: 8, Code: "reject", Name: "Буцаах"},
    {ID: 9, Code: "execute", Name: "Гүйцэтгэх"},
    {ID: 10, Code: "publish", Name: "Нийтлэх"},
    {ID: 11, Code: "manage", Name: "Удирдах"},
    {ID: 12, Code: "terminate", Name: "Зогсоох"},
    {ID: 13, Code: "switch", Name: "Сэлгэх"},
    {ID: 14, Code: "revoke", Name: "Хүчингүй болгох"},
    {ID: 15, Code: "test", Name: "Тест"},
}
```

## Systems

```go
systems := []domain.System{
    {
        ID:          1,
        Code:        "dsl",
        Name:        "DSL систем",
        Description: "Domain Specific Language - Динамик схем, дүрэм, workflow",
        IconName:    "Code",
        Color:       "#6366f1",
        Sequence:    1,
        IsActive:    ptr(true),
    },
    {
        ID:          2,
        Code:        "gateway",
        Name:        "Gateway систем",
        Description: "API Management & Integration Hub",
        IconName:    "Network",
        Color:       "#22c55e",
        Sequence:    2,
        IsActive:    ptr(true),
    },
}
```

## Platform Modules (system_id = NULL)

```go
platformModules := []domain.Module{
    {ID: 1, Code: "user", Name: "Хэрэглэгч", SystemID: nil},
    {ID: 2, Code: "organization", Name: "Байгууллага", SystemID: nil},
    {ID: 3, Code: "system", Name: "Систем", SystemID: nil},
    {ID: 4, Code: "module", Name: "Модуль", SystemID: nil},
    {ID: 5, Code: "action", Name: "Үйлдэл", SystemID: nil},
    {ID: 6, Code: "role", Name: "Эрх", SystemID: nil},
    {ID: 7, Code: "permission", Name: "Зөвшөөрөл", SystemID: nil},
    {ID: 8, Code: "menu", Name: "Цэс", SystemID: nil},
    {ID: 9, Code: "device", Name: "Төхөөрөмж", SystemID: nil},
    {ID: 10, Code: "session", Name: "Сешн", SystemID: nil},
    {ID: 11, Code: "monitoring", Name: "Мониторинг", SystemID: nil},
    {ID: 12, Code: "language", Name: "Хэл", SystemID: nil},
    {ID: 13, Code: "translation", Name: "Орчуулга", SystemID: nil},
}
```

## DSL System Modules (system_id = 1)

```go
dslModules := []domain.Module{
    {ID: 14, Code: "schema", Name: "Схем", SystemID: ptr(1)},
    {ID: 15, Code: "field", Name: "Талбар", SystemID: ptr(1)},
    {ID: 16, Code: "rule", Name: "Дүрэм", SystemID: ptr(1)},
    {ID: 17, Code: "workflow", Name: "Ажлын урсгал", SystemID: ptr(1)},
    {ID: 18, Code: "template", Name: "Загвар", SystemID: ptr(1)},
    {ID: 19, Code: "function", Name: "Функц", SystemID: ptr(1)},
    {ID: 20, Code: "variable", Name: "Хувьсагч", SystemID: ptr(1)},
    {ID: 21, Code: "executor", Name: "Гүйцэтгэгч", SystemID: ptr(1)},
    {ID: 22, Code: "log", Name: "Лог", SystemID: ptr(1)},
}
```

## Gateway System Modules (system_id = 2)

```go
gatewayModules := []domain.Module{
    {ID: 23, Code: "client", Name: "API Client", SystemID: ptr(2)},
    {ID: 24, Code: "endpoint", Name: "Endpoint", SystemID: ptr(2)},
    {ID: 25, Code: "integration", Name: "Интеграц", SystemID: ptr(2)},
    {ID: 26, Code: "credential", Name: "Credential", SystemID: ptr(2)},
    {ID: 27, Code: "webhook", Name: "Webhook", SystemID: ptr(2)},
    {ID: 28, Code: "ratelimit", Name: "Rate Limit", SystemID: ptr(2)},
    {ID: 29, Code: "apilog", Name: "API Лог", SystemID: ptr(2)},
    {ID: 30, Code: "monitor", Name: "Мониторинг", SystemID: ptr(2)},
}
```

## Platform Roles (system_id = NULL)

```go
platformRoles := []domain.Role{
    {
        ID:          1,
        Code:        "platform_superadmin",
        Name:        "Platform Супер Админ",
        Description: "Платформын бүх эрхтэй",
        SystemID:    nil,
        IsSystem:    ptr(true),
        IsActive:    ptr(true),
    },
    {
        ID:          2,
        Code:        "platform_admin",
        Name:        "Platform Админ",
        Description: "Платформын админ",
        SystemID:    nil,
        IsSystem:    ptr(true),
        IsActive:    ptr(true),
    },
    {
        ID:          3,
        Code:        "platform_user",
        Name:        "Platform Хэрэглэгч",
        Description: "Энгийн хэрэглэгч, зөвхөн системд нэвтрэх",
        SystemID:    nil,
        IsSystem:    ptr(true),
        IsActive:    ptr(true),
    },
}
```

## DSL System Roles (system_id = 1)

```go
dslRoles := []domain.Role{
    {
        ID:          4,
        Code:        "dsl_admin",
        Name:        "DSL Админ",
        Description: "DSL системийн бүрэн эрхтэй",
        SystemID:    ptr(1),
        IsSystem:    ptr(true),
        IsActive:    ptr(true),
    },
    {
        ID:          5,
        Code:        "dsl_developer",
        Name:        "DSL Хөгжүүлэгч",
        Description: "Схем, дүрэм, workflow үүсгэх, засах",
        SystemID:    ptr(1),
        IsSystem:    ptr(true),
        IsActive:    ptr(true),
    },
    {
        ID:          6,
        Code:        "dsl_viewer",
        Name:        "DSL Үзэгч",
        Description: "Зөвхөн харах эрхтэй",
        SystemID:    ptr(1),
        IsSystem:    ptr(true),
        IsActive:    ptr(true),
    },
}
```

## Gateway System Roles (system_id = 2)

```go
gatewayRoles := []domain.Role{
    {
        ID:          7,
        Code:        "gateway_admin",
        Name:        "Gateway Админ",
        Description: "Gateway системийн бүрэн эрхтэй",
        SystemID:    ptr(2),
        IsSystem:    ptr(true),
        IsActive:    ptr(true),
    },
    {
        ID:          8,
        Code:        "gateway_developer",
        Name:        "Gateway Хөгжүүлэгч",
        Description: "API, Integration удирдах",
        SystemID:    ptr(2),
        IsSystem:    ptr(true),
        IsActive:    ptr(true),
    },
    {
        ID:          9,
        Code:        "gateway_viewer",
        Name:        "Gateway Үзэгч",
        Description: "Зөвхөн харах, мониторинг",
        SystemID:    ptr(2),
        IsSystem:    ptr(true),
        IsActive:    ptr(true),
    },
}
```

## Platform Menus (system_id = NULL)

```go
platformMenus := []domain.Menu{
    // Dashboard
    {ID: 1, Code: "platform_dashboard", Name: "Хянах самбар", SystemID: nil, Path: "/platform", Icon: "LayoutDashboard", Sequence: 1},
    
    // User Management
    {ID: 2, Code: "user_management", Name: "Хэрэглэгч", SystemID: nil, Path: "", Icon: "Users", Sequence: 2},
    {ID: 3, Code: "users", Name: "Хэрэглэгчид", SystemID: nil, ParentID: ptr(2), Path: "/platform/users", Icon: "User", Sequence: 1},
    {ID: 4, Code: "roles", Name: "Эрхүүд", SystemID: nil, ParentID: ptr(2), Path: "/platform/roles", Icon: "Shield", Sequence: 2},
    {ID: 5, Code: "permissions", Name: "Зөвшөөрлүүд", SystemID: nil, ParentID: ptr(2), Path: "/platform/permissions", Icon: "Key", Sequence: 3},
    
    // Organization
    {ID: 6, Code: "org_management", Name: "Байгууллага", SystemID: nil, Path: "", Icon: "Building", Sequence: 3},
    {ID: 7, Code: "organizations", Name: "Байгууллагууд", SystemID: nil, ParentID: ptr(6), Path: "/platform/organizations", Icon: "Building2", Sequence: 1},
    {ID: 8, Code: "org_types", Name: "Төрлүүд", SystemID: nil, ParentID: ptr(6), Path: "/platform/organization-types", Icon: "Tag", Sequence: 2},
    
    // System Config
    {ID: 9, Code: "system_config", Name: "Тохиргоо", SystemID: nil, Path: "", Icon: "Settings", Sequence: 4},
    {ID: 10, Code: "systems", Name: "Системүүд", SystemID: nil, ParentID: ptr(9), Path: "/platform/systems", Icon: "Server", Sequence: 1},
    {ID: 11, Code: "modules", Name: "Модулиуд", SystemID: nil, ParentID: ptr(9), Path: "/platform/modules", Icon: "Package", Sequence: 2},
    {ID: 12, Code: "menus", Name: "Цэсүүд", SystemID: nil, ParentID: ptr(9), Path: "/platform/menus", Icon: "Menu", Sequence: 3},
    
    // Device & Session
    {ID: 13, Code: "device_session", Name: "Төхөөрөмж", SystemID: nil, Path: "", Icon: "Smartphone", Sequence: 5},
    {ID: 14, Code: "devices", Name: "Төхөөрөмжүүд", SystemID: nil, ParentID: ptr(13), Path: "/platform/devices", Icon: "Monitor", Sequence: 1},
    {ID: 15, Code: "sessions", Name: "Сешнүүд", SystemID: nil, ParentID: ptr(13), Path: "/platform/sessions", Icon: "Activity", Sequence: 2},
    
    // Monitoring
    {ID: 16, Code: "monitoring", Name: "Мониторинг", SystemID: nil, Path: "/platform/monitoring", Icon: "BarChart", Sequence: 6},
    
    // Localization
    {ID: 17, Code: "localization", Name: "Орчуулга", SystemID: nil, Path: "", Icon: "Globe", Sequence: 7},
    {ID: 18, Code: "languages", Name: "Хэлүүд", SystemID: nil, ParentID: ptr(17), Path: "/platform/languages", Icon: "Languages", Sequence: 1},
    {ID: 19, Code: "translations", Name: "Орчуулгууд", SystemID: nil, ParentID: ptr(17), Path: "/platform/translations", Icon: "FileText", Sequence: 2},
}
```

## DSL System Menus (system_id = 1)

```go
dslMenus := []domain.Menu{
    // Dashboard
    {ID: 20, Code: "dsl_dashboard", Name: "Хянах самбар", SystemID: ptr(1), Path: "/dsl", Icon: "LayoutDashboard", Sequence: 1},
    
    // Data Modeling
    {ID: 21, Code: "data_modeling", Name: "Өгөгдөл", SystemID: ptr(1), Path: "", Icon: "Database", Sequence: 2},
    {ID: 22, Code: "schemas", Name: "Схемүүд", SystemID: ptr(1), ParentID: ptr(21), Path: "/dsl/schemas", Icon: "Table", Sequence: 1},
    {ID: 23, Code: "fields", Name: "Талбарууд", SystemID: ptr(1), ParentID: ptr(21), Path: "/dsl/fields", Icon: "Columns", Sequence: 2},
    
    // Business Logic
    {ID: 24, Code: "business_logic", Name: "Бизнес логик", SystemID: ptr(1), Path: "", Icon: "GitBranch", Sequence: 3},
    {ID: 25, Code: "rules", Name: "Дүрмүүд", SystemID: ptr(1), ParentID: ptr(24), Path: "/dsl/rules", Icon: "CheckSquare", Sequence: 1},
    {ID: 26, Code: "workflows", Name: "Workflow", SystemID: ptr(1), ParentID: ptr(24), Path: "/dsl/workflows", Icon: "Workflow", Sequence: 2},
    {ID: 27, Code: "functions", Name: "Функцүүд", SystemID: ptr(1), ParentID: ptr(24), Path: "/dsl/functions", Icon: "Code", Sequence: 3},
    
    // Templates
    {ID: 28, Code: "templates", Name: "Загварууд", SystemID: ptr(1), Path: "/dsl/templates", Icon: "FileCode", Sequence: 4},
    
    // Variables
    {ID: 29, Code: "variables", Name: "Хувьсагчид", SystemID: ptr(1), Path: "/dsl/variables", Icon: "Variable", Sequence: 5},
    
    // Execution
    {ID: 30, Code: "execution", Name: "Гүйцэтгэл", SystemID: ptr(1), Path: "", Icon: "Play", Sequence: 6},
    {ID: 31, Code: "executor", Name: "Executor", SystemID: ptr(1), ParentID: ptr(30), Path: "/dsl/executor", Icon: "Terminal", Sequence: 1},
    {ID: 32, Code: "dsl_logs", Name: "Логууд", SystemID: ptr(1), ParentID: ptr(30), Path: "/dsl/logs", Icon: "ScrollText", Sequence: 2},
}
```

## Gateway System Menus (system_id = 2)

```go
gatewayMenus := []domain.Menu{
    // Dashboard
    {ID: 33, Code: "gateway_dashboard", Name: "Хянах самбар", SystemID: ptr(2), Path: "/gateway", Icon: "LayoutDashboard", Sequence: 1},
    
    // API Management
    {ID: 34, Code: "api_management", Name: "API удирдлага", SystemID: ptr(2), Path: "", Icon: "Plug", Sequence: 2},
    {ID: 35, Code: "clients", Name: "API Clients", SystemID: ptr(2), ParentID: ptr(34), Path: "/gateway/clients", Icon: "Key", Sequence: 1},
    {ID: 36, Code: "endpoints", Name: "Endpoints", SystemID: ptr(2), ParentID: ptr(34), Path: "/gateway/endpoints", Icon: "Link", Sequence: 2},
    {ID: 37, Code: "ratelimits", Name: "Rate Limits", SystemID: ptr(2), ParentID: ptr(34), Path: "/gateway/rate-limits", Icon: "Gauge", Sequence: 3},
    
    // Integrations
    {ID: 38, Code: "integrations_menu", Name: "Интеграц", SystemID: ptr(2), Path: "", Icon: "Network", Sequence: 3},
    {ID: 39, Code: "integrations", Name: "Интеграцууд", SystemID: ptr(2), ParentID: ptr(38), Path: "/gateway/integrations", Icon: "Unplug", Sequence: 1},
    {ID: 40, Code: "credentials", Name: "Credentials", SystemID: ptr(2), ParentID: ptr(38), Path: "/gateway/credentials", Icon: "Lock", Sequence: 2},
    
    // Webhooks
    {ID: 41, Code: "webhooks", Name: "Webhooks", SystemID: ptr(2), Path: "/gateway/webhooks", Icon: "Webhook", Sequence: 4},
    
    // Monitoring
    {ID: 42, Code: "gateway_monitoring", Name: "Мониторинг", SystemID: ptr(2), Path: "", Icon: "BarChart", Sequence: 5},
    {ID: 43, Code: "api_logs", Name: "API Logs", SystemID: ptr(2), ParentID: ptr(42), Path: "/gateway/logs", Icon: "ScrollText", Sequence: 1},
    {ID: 44, Code: "health_check", Name: "Health Check", SystemID: ptr(2), ParentID: ptr(42), Path: "/gateway/health", Icon: "HeartPulse", Sequence: 2},
    {ID: 45, Code: "analytics", Name: "Analytics", SystemID: ptr(2), ParentID: ptr(42), Path: "/gateway/analytics", Icon: "TrendingUp", Sequence: 3},
}
```

## Organization Types

```go
orgTypes := []domain.OrganizationType{
    {ID: 1, Code: "government", Name: "Төрийн байгууллага"},
    {ID: 2, Code: "private", Name: "Хувийн хэвшил"},
    {ID: 3, Code: "ngo", Name: "ТББ"},
    {ID: 4, Code: "education", Name: "Боловсролын байгууллага"},
    {ID: 5, Code: "healthcare", Name: "Эрүүл мэндийн байгууллага"},
}
```

## Sample Platform Permissions

```go
platformPermissions := []domain.Permission{
    // User module
    {Code: "user.view", Name: "Хэрэглэгч харах", SystemID: nil, ModuleID: 1, ActionID: ptr(int64(1))},
    {Code: "user.create", Name: "Хэрэглэгч үүсгэх", SystemID: nil, ModuleID: 1, ActionID: ptr(int64(2))},
    {Code: "user.update", Name: "Хэрэглэгч засах", SystemID: nil, ModuleID: 1, ActionID: ptr(int64(3))},
    {Code: "user.delete", Name: "Хэрэглэгч устгах", SystemID: nil, ModuleID: 1, ActionID: ptr(int64(4))},
    
    // Organization module
    {Code: "organization.view", Name: "Байгууллага харах", SystemID: nil, ModuleID: 2, ActionID: ptr(int64(1))},
    {Code: "organization.create", Name: "Байгууллага үүсгэх", SystemID: nil, ModuleID: 2, ActionID: ptr(int64(2))},
    {Code: "organization.update", Name: "Байгууллага засах", SystemID: nil, ModuleID: 2, ActionID: ptr(int64(3))},
    {Code: "organization.delete", Name: "Байгууллага устгах", SystemID: nil, ModuleID: 2, ActionID: ptr(int64(4))},
    
    // Role module
    {Code: "role.view", Name: "Эрх харах", SystemID: nil, ModuleID: 6, ActionID: ptr(int64(1))},
    {Code: "role.create", Name: "Эрх үүсгэх", SystemID: nil, ModuleID: 6, ActionID: ptr(int64(2))},
    {Code: "role.manage", Name: "Эрх удирдах", SystemID: nil, ModuleID: 6, ActionID: ptr(int64(11))},
    
    // Monitoring
    {Code: "monitoring.view", Name: "Мониторинг харах", SystemID: nil, ModuleID: 11, ActionID: ptr(int64(1))},
    
    // Session
    {Code: "session.view", Name: "Сешн харах", SystemID: nil, ModuleID: 10, ActionID: ptr(int64(1))},
    {Code: "session.terminate", Name: "Сешн зогсоох", SystemID: nil, ModuleID: 10, ActionID: ptr(int64(12))},
    
    // System switch permission
    {Code: "system.switch", Name: "Систем сэлгэх", SystemID: nil, ModuleID: 3, ActionID: ptr(int64(13))},
}
```

## Sample DSL System Permissions

```go
dslPermissions := []domain.Permission{
    // Schema
    {Code: "dsl.schema.view", Name: "Схем харах", SystemID: ptr(1), ModuleID: 14, ActionID: ptr(int64(1))},
    {Code: "dsl.schema.create", Name: "Схем үүсгэх", SystemID: ptr(1), ModuleID: 14, ActionID: ptr(int64(2))},
    {Code: "dsl.schema.update", Name: "Схем засах", SystemID: ptr(1), ModuleID: 14, ActionID: ptr(int64(3))},
    {Code: "dsl.schema.delete", Name: "Схем устгах", SystemID: ptr(1), ModuleID: 14, ActionID: ptr(int64(4))},
    {Code: "dsl.schema.execute", Name: "Схем deploy", SystemID: ptr(1), ModuleID: 14, ActionID: ptr(int64(9))},
    
    // Rule
    {Code: "dsl.rule.view", Name: "Дүрэм харах", SystemID: ptr(1), ModuleID: 16, ActionID: ptr(int64(1))},
    {Code: "dsl.rule.create", Name: "Дүрэм үүсгэх", SystemID: ptr(1), ModuleID: 16, ActionID: ptr(int64(2))},
    {Code: "dsl.rule.execute", Name: "Дүрэм тест", SystemID: ptr(1), ModuleID: 16, ActionID: ptr(int64(9))},
    
    // Workflow
    {Code: "dsl.workflow.view", Name: "Workflow харах", SystemID: ptr(1), ModuleID: 17, ActionID: ptr(int64(1))},
    {Code: "dsl.workflow.create", Name: "Workflow үүсгэх", SystemID: ptr(1), ModuleID: 17, ActionID: ptr(int64(2))},
    {Code: "dsl.workflow.execute", Name: "Workflow гүйцэтгэх", SystemID: ptr(1), ModuleID: 17, ActionID: ptr(int64(9))},
    {Code: "dsl.workflow.publish", Name: "Workflow нийтлэх", SystemID: ptr(1), ModuleID: 17, ActionID: ptr(int64(10))},
    
    // Template, Function, Variable, Log...
}
```

## Sample Gateway System Permissions

```go
gatewayPermissions := []domain.Permission{
    // Client
    {Code: "gateway.client.view", Name: "Client харах", SystemID: ptr(2), ModuleID: 23, ActionID: ptr(int64(1))},
    {Code: "gateway.client.create", Name: "Client үүсгэх", SystemID: ptr(2), ModuleID: 23, ActionID: ptr(int64(2))},
    {Code: "gateway.client.update", Name: "Client засах", SystemID: ptr(2), ModuleID: 23, ActionID: ptr(int64(3))},
    {Code: "gateway.client.delete", Name: "Client устгах", SystemID: ptr(2), ModuleID: 23, ActionID: ptr(int64(4))},
    {Code: "gateway.client.revoke", Name: "Client хүчингүй болгох", SystemID: ptr(2), ModuleID: 23, ActionID: ptr(int64(14))},
    
    // Integration
    {Code: "gateway.integration.view", Name: "Интеграц харах", SystemID: ptr(2), ModuleID: 25, ActionID: ptr(int64(1))},
    {Code: "gateway.integration.create", Name: "Интеграц үүсгэх", SystemID: ptr(2), ModuleID: 25, ActionID: ptr(int64(2))},
    {Code: "gateway.integration.update", Name: "Интеграц засах", SystemID: ptr(2), ModuleID: 25, ActionID: ptr(int64(3))},
    {Code: "gateway.integration.execute", Name: "Интеграц тест", SystemID: ptr(2), ModuleID: 25, ActionID: ptr(int64(9))},
    
    // Webhook
    {Code: "gateway.webhook.view", Name: "Webhook харах", SystemID: ptr(2), ModuleID: 27, ActionID: ptr(int64(1))},
    {Code: "gateway.webhook.create", Name: "Webhook үүсгэх", SystemID: ptr(2), ModuleID: 27, ActionID: ptr(int64(2))},
    
    // Monitor
    {Code: "gateway.monitor.view", Name: "Мониторинг харах", SystemID: ptr(2), ModuleID: 30, ActionID: ptr(int64(1))},
}
```

## Default Super Admin User

```go
superAdminUser := domain.User{
    ID:              1,
    RegNo:           "АА00000000",
    FamilyName:      "Систем",
    LastName:        "Админ",
    FirstName:       "Супер",
    Gender:          1,
    BirthDate:       "1990-01-01",
    PhoneNo:         "99999999",
    Email:           "admin@gebase.mn",
    PasswordHash:    hashPassword("Admin@123"),
    IsActive:        ptr(true),
    LanguageCode:    "mn",
    DefaultSystemID: ptr(1), // Default to DSL system
}

// Super admin has all roles
superAdminRoles := []domain.UserSystemRole{
    // Platform superadmin
    {UserID: 1, SystemID: nil, RoleID: 1, IsActive: ptr(true), IsDefault: ptr(true)},
    // DSL admin
    {UserID: 1, SystemID: ptr(1), RoleID: 4, IsActive: ptr(true), IsDefault: ptr(true)},
    // Gateway admin
    {UserID: 1, SystemID: ptr(2), RoleID: 7, IsActive: ptr(true), IsDefault: ptr(false)},
}
```
