package db

import (
	"log"

	"gebase/internal/domain"
	"gebase/internal/domain/dsl"
	"gebase/internal/service"

	"gorm.io/gorm"
)

func ptr[T any](v T) *T {
	return &v
}

func RunSeed(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	if err := seedLanguages(db); err != nil {
		return err
	}

	if err := seedActions(db); err != nil {
		return err
	}

	if err := seedSystems(db); err != nil {
		return err
	}

	if err := seedModules(db); err != nil {
		return err
	}

	if err := seedOrganizationTypes(db); err != nil {
		return err
	}

	if err := seedRoles(db); err != nil {
		return err
	}

	if err := seedMenus(db); err != nil {
		return err
	}

	if err := seedTranslations(db); err != nil {
		return err
	}

	if err := seedDSLFunctions(db); err != nil {
		return err
	}

	if err := seedAdminUser(db); err != nil {
		return err
	}

	if err := seedRoleMenus(db); err != nil {
		return err
	}

	if err := seedPermissions(db); err != nil {
		return err
	}

	if err := seedRolePermissions(db); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

func seedLanguages(db *gorm.DB) error {
	languages := []domain.Language{
		{ID: 1, Code: "mn", Name: "Mongolian", NativeName: "Монгол", IsActive: ptr(true), IsDefault: ptr(true)},
		{ID: 2, Code: "en", Name: "English", NativeName: "English", IsActive: ptr(true), IsDefault: ptr(false)},
	}

	for _, lang := range languages {
		if err := db.Where("id = ?", lang.ID).FirstOrCreate(&lang).Error; err != nil {
			return err
		}
	}
	log.Println("Languages seeded")
	return nil
}

func seedActions(db *gorm.DB) error {
	actions := []domain.Action{
		{ID: 1, Code: "view", Name: "Харах", IsActive: ptr(true)},
		{ID: 2, Code: "create", Name: "Үүсгэх", IsActive: ptr(true)},
		{ID: 3, Code: "update", Name: "Засах", IsActive: ptr(true)},
		{ID: 4, Code: "delete", Name: "Устгах", IsActive: ptr(true)},
		{ID: 5, Code: "export", Name: "Экспорт", IsActive: ptr(true)},
		{ID: 6, Code: "import", Name: "Импорт", IsActive: ptr(true)},
		{ID: 7, Code: "approve", Name: "Батлах", IsActive: ptr(true)},
		{ID: 8, Code: "reject", Name: "Буцаах", IsActive: ptr(true)},
		{ID: 9, Code: "execute", Name: "Гүйцэтгэх", IsActive: ptr(true)},
		{ID: 10, Code: "publish", Name: "Нийтлэх", IsActive: ptr(true)},
	}

	for _, action := range actions {
		if err := db.Where("id = ?", action.ID).FirstOrCreate(&action).Error; err != nil {
			return err
		}
	}
	log.Println("Actions seeded")
	return nil
}

func seedSystems(db *gorm.DB) error {
	systems := []domain.System{
		{
			ID:          1,
			Code:        "admin",
			Name:        "Админ систем",
			Description: "Хэрэглэгч, байгууллага, эрх, тохиргооны удирдлага",
			IconName:    "Settings",
			Color:       "#6366f1",
			Sequence:    1,
			IsActive:    ptr(true),
		},
		{
			ID:          2,
			Code:        "dsl",
			Name:        "DSL систем",
			Description: "Domain Specific Language - Динамик схем, дүрэм, workflow",
			IconName:    "Code",
			Color:       "#10b981",
			Sequence:    2,
			IsActive:    ptr(true),
		},
	}

	for _, system := range systems {
		if err := db.Where("id = ?", system.ID).FirstOrCreate(&system).Error; err != nil {
			return err
		}
	}
	log.Println("Systems seeded")
	return nil
}

func seedModules(db *gorm.DB) error {
	// Admin modules
	adminModules := []domain.Module{
		{ID: 1, Code: "user", Name: "Хэрэглэгч", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 2, Code: "organization", Name: "Байгууллага", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 3, Code: "system", Name: "Систем", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 4, Code: "module", Name: "Модуль", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 5, Code: "action", Name: "Үйлдэл", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 6, Code: "role", Name: "Эрх", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 7, Code: "permission", Name: "Зөвшөөрөл", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 8, Code: "menu", Name: "Цэс", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 9, Code: "device", Name: "Төхөөрөмж", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 10, Code: "session", Name: "Сешн", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 11, Code: "monitoring", Name: "Мониторинг", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 12, Code: "language", Name: "Хэл", SystemID: ptr(1), IsActive: ptr(true)},
		{ID: 13, Code: "translation", Name: "Орчуулга", SystemID: ptr(1), IsActive: ptr(true)},
	}

	// DSL modules
	dslModules := []domain.Module{
		{ID: 14, Code: "schema", Name: "Схем", SystemID: ptr(2), IsActive: ptr(true)},
		{ID: 15, Code: "field", Name: "Талбар", SystemID: ptr(2), IsActive: ptr(true)},
		{ID: 16, Code: "rule", Name: "Дүрэм", SystemID: ptr(2), IsActive: ptr(true)},
		{ID: 17, Code: "workflow", Name: "Ажлын урсгал", SystemID: ptr(2), IsActive: ptr(true)},
		{ID: 18, Code: "template", Name: "Загвар", SystemID: ptr(2), IsActive: ptr(true)},
		{ID: 19, Code: "function", Name: "Функц", SystemID: ptr(2), IsActive: ptr(true)},
		{ID: 20, Code: "variable", Name: "Хувьсагч", SystemID: ptr(2), IsActive: ptr(true)},
		{ID: 21, Code: "executor", Name: "Гүйцэтгэгч", SystemID: ptr(2), IsActive: ptr(true)},
		{ID: 22, Code: "log", Name: "Лог", SystemID: ptr(2), IsActive: ptr(true)},
	}

	allModules := append(adminModules, dslModules...)
	for _, module := range allModules {
		if err := db.Where("id = ?", module.ID).FirstOrCreate(&module).Error; err != nil {
			return err
		}
	}
	log.Println("Modules seeded")
	return nil
}

func seedOrganizationTypes(db *gorm.DB) error {
	orgTypes := []domain.OrganizationType{
		{ID: 1, Code: "government", Name: "Төрийн байгууллага", IsActive: ptr(true)},
		{ID: 2, Code: "private", Name: "Хувийн хэвшил", IsActive: ptr(true)},
		{ID: 3, Code: "ngo", Name: "ТББ", IsActive: ptr(true)},
		{ID: 4, Code: "education", Name: "Боловсролын байгууллага", IsActive: ptr(true)},
		{ID: 5, Code: "healthcare", Name: "Эрүүл мэндийн байгууллага", IsActive: ptr(true)},
	}

	for _, orgType := range orgTypes {
		if err := db.Where("id = ?", orgType.ID).FirstOrCreate(&orgType).Error; err != nil {
			return err
		}
	}
	log.Println("Organization types seeded")
	return nil
}

func seedRoles(db *gorm.DB) error {
	roles := []domain.Role{
		// Admin system roles
		{ID: 1, Code: "super_admin", Name: "Супер админ", Description: "Бүх эрхтэй систем админ", SystemID: ptr(1), IsSystem: ptr(true), IsActive: ptr(true)},
		{ID: 2, Code: "admin", Name: "Админ", Description: "Байгууллагын админ", SystemID: ptr(1), IsSystem: ptr(true), IsActive: ptr(true)},
		{ID: 3, Code: "operator", Name: "Оператор", Description: "Энгийн оператор", SystemID: ptr(1), IsSystem: ptr(true), IsActive: ptr(true)},
		// DSL system roles
		{ID: 4, Code: "dsl_admin", Name: "DSL Админ", Description: "DSL системийн бүрэн эрхтэй", SystemID: ptr(2), IsSystem: ptr(true), IsActive: ptr(true)},
		{ID: 5, Code: "dsl_developer", Name: "DSL Хөгжүүлэгч", Description: "Схем, дүрэм, workflow үүсгэх", SystemID: ptr(2), IsSystem: ptr(true), IsActive: ptr(true)},
		{ID: 6, Code: "dsl_viewer", Name: "DSL Үзэгч", Description: "Зөвхөн харах эрхтэй", SystemID: ptr(2), IsSystem: ptr(true), IsActive: ptr(true)},
	}

	for _, role := range roles {
		if err := db.Where("id = ?", role.ID).FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}
	log.Println("Roles seeded")
	return nil
}

func seedMenus(db *gorm.DB) error {
	// Admin menus
	adminMenus := []domain.Menu{
		{ID: 1, Code: "dashboard", Name: "Хянах самбар", SystemID: ptr(1), Path: "/dashboard", Icon: "LayoutDashboard", Sequence: 1, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 2, Code: "user_management", Name: "Хэрэглэгчийн удирдлага", SystemID: ptr(1), Path: "", Icon: "Users", Sequence: 2, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 3, Code: "users", Name: "Хэрэглэгчид", SystemID: ptr(1), ParentID: ptr(2), Path: "/users", Icon: "User", Sequence: 1, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 4, Code: "roles", Name: "Эрхүүд", SystemID: ptr(1), ParentID: ptr(2), Path: "/roles", Icon: "Shield", Sequence: 2, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 5, Code: "permissions", Name: "Зөвшөөрлүүд", SystemID: ptr(1), ParentID: ptr(2), Path: "/permissions", Icon: "Key", Sequence: 3, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 6, Code: "org_management", Name: "Байгууллагын удирдлага", SystemID: ptr(1), Path: "", Icon: "Building", Sequence: 3, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 7, Code: "organizations", Name: "Байгууллагууд", SystemID: ptr(1), ParentID: ptr(6), Path: "/organizations", Icon: "Building2", Sequence: 1, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 8, Code: "org_types", Name: "Байгууллагын төрөл", SystemID: ptr(1), ParentID: ptr(6), Path: "/organization-types", Icon: "Tag", Sequence: 2, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 9, Code: "system_config", Name: "Системийн тохиргоо", SystemID: ptr(1), Path: "", Icon: "Settings", Sequence: 4, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 10, Code: "systems", Name: "Системүүд", SystemID: ptr(1), ParentID: ptr(9), Path: "/systems", Icon: "Server", Sequence: 1, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 11, Code: "modules", Name: "Модулиуд", SystemID: ptr(1), ParentID: ptr(9), Path: "/modules", Icon: "Package", Sequence: 2, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 12, Code: "actions", Name: "Үйлдлүүд", SystemID: ptr(1), ParentID: ptr(9), Path: "/actions", Icon: "Zap", Sequence: 3, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 13, Code: "menus", Name: "Цэсүүд", SystemID: ptr(1), ParentID: ptr(9), Path: "/menus", Icon: "Menu", Sequence: 4, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 14, Code: "device_management", Name: "Төхөөрөмж & Сешн", SystemID: ptr(1), Path: "", Icon: "Smartphone", Sequence: 5, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 15, Code: "devices", Name: "Төхөөрөмжүүд", SystemID: ptr(1), ParentID: ptr(14), Path: "/devices", Icon: "Monitor", Sequence: 1, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 16, Code: "sessions", Name: "Сешнүүд", SystemID: ptr(1), ParentID: ptr(14), Path: "/sessions", Icon: "Activity", Sequence: 2, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 17, Code: "monitoring", Name: "Мониторинг", SystemID: ptr(1), Path: "/monitoring", Icon: "BarChart", Sequence: 6, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 18, Code: "localization", Name: "Орчуулга", SystemID: ptr(1), Path: "", Icon: "Globe", Sequence: 7, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 19, Code: "languages", Name: "Хэлүүд", SystemID: ptr(1), ParentID: ptr(18), Path: "/languages", Icon: "Languages", Sequence: 1, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 20, Code: "translations", Name: "Орчуулгууд", SystemID: ptr(1), ParentID: ptr(18), Path: "/translations", Icon: "FileText", Sequence: 2, IsVisible: ptr(true), IsActive: ptr(true)},
	}

	// DSL menus
	dslMenus := []domain.Menu{
		{ID: 21, Code: "dsl_dashboard", Name: "Хянах самбар", SystemID: ptr(2), Path: "/dashboard", Icon: "LayoutDashboard", Sequence: 1, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 22, Code: "data_modeling", Name: "Өгөгдлийн загварчлал", SystemID: ptr(2), Path: "", Icon: "Database", Sequence: 2, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 23, Code: "schemas", Name: "Схемүүд", SystemID: ptr(2), ParentID: ptr(22), Path: "/schemas", Icon: "Table", Sequence: 1, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 24, Code: "fields", Name: "Талбарууд", SystemID: ptr(2), ParentID: ptr(22), Path: "/fields", Icon: "Columns", Sequence: 2, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 25, Code: "business_logic", Name: "Бизнес логик", SystemID: ptr(2), Path: "", Icon: "GitBranch", Sequence: 3, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 26, Code: "rules", Name: "Дүрмүүд", SystemID: ptr(2), ParentID: ptr(25), Path: "/rules", Icon: "CheckSquare", Sequence: 1, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 27, Code: "workflows", Name: "Workflow", SystemID: ptr(2), ParentID: ptr(25), Path: "/workflows", Icon: "Workflow", Sequence: 2, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 28, Code: "functions", Name: "Функцүүд", SystemID: ptr(2), ParentID: ptr(25), Path: "/functions", Icon: "Code", Sequence: 3, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 29, Code: "templates", Name: "Загварууд", SystemID: ptr(2), Path: "/templates", Icon: "FileCode", Sequence: 4, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 30, Code: "variables", Name: "Хувьсагчид", SystemID: ptr(2), Path: "/variables", Icon: "Variable", Sequence: 5, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 31, Code: "execution", Name: "Гүйцэтгэл", SystemID: ptr(2), Path: "", Icon: "Play", Sequence: 6, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 32, Code: "executor", Name: "Executor", SystemID: ptr(2), ParentID: ptr(31), Path: "/executor", Icon: "Terminal", Sequence: 1, IsVisible: ptr(true), IsActive: ptr(true)},
		{ID: 33, Code: "logs", Name: "Логууд", SystemID: ptr(2), ParentID: ptr(31), Path: "/logs", Icon: "ScrollText", Sequence: 2, IsVisible: ptr(true), IsActive: ptr(true)},
	}

	allMenus := append(adminMenus, dslMenus...)
	for _, menu := range allMenus {
		if err := db.Where("id = ?", menu.ID).FirstOrCreate(&menu).Error; err != nil {
			return err
		}
	}
	log.Println("Menus seeded")
	return nil
}

func seedTranslations(db *gorm.DB) error {
	translations := []domain.Translation{
		// Common
		{LanguageCode: "mn", Key: "common.save", Value: "Хадгалах", Module: "common"},
		{LanguageCode: "en", Key: "common.save", Value: "Save", Module: "common"},
		{LanguageCode: "mn", Key: "common.cancel", Value: "Цуцлах", Module: "common"},
		{LanguageCode: "en", Key: "common.cancel", Value: "Cancel", Module: "common"},
		{LanguageCode: "mn", Key: "common.delete", Value: "Устгах", Module: "common"},
		{LanguageCode: "en", Key: "common.delete", Value: "Delete", Module: "common"},
		{LanguageCode: "mn", Key: "common.edit", Value: "Засах", Module: "common"},
		{LanguageCode: "en", Key: "common.edit", Value: "Edit", Module: "common"},
		{LanguageCode: "mn", Key: "common.search", Value: "Хайх", Module: "common"},
		{LanguageCode: "en", Key: "common.search", Value: "Search", Module: "common"},
		// Error messages
		{LanguageCode: "mn", Key: "error.unauthorized", Value: "Нэвтрэх эрхгүй", Module: "error"},
		{LanguageCode: "en", Key: "error.unauthorized", Value: "Unauthorized", Module: "error"},
		{LanguageCode: "mn", Key: "error.forbidden", Value: "Хандах эрхгүй", Module: "error"},
		{LanguageCode: "en", Key: "error.forbidden", Value: "Forbidden", Module: "error"},
	}

	for _, t := range translations {
		if err := db.Where("language_code = ? AND key = ?", t.LanguageCode, t.Key).FirstOrCreate(&t).Error; err != nil {
			return err
		}
	}
	log.Println("Translations seeded")
	return nil
}

func seedDSLFunctions(db *gorm.DB) error {
	functions := []dsl.Function{
		{Code: "now", Name: "Now", Description: "Returns current datetime", FunctionType: dsl.FunctionTypeBuiltIn, ReturnType: "datetime", Parameters: "[]", Example: "now()", IsActive: ptr(true)},
		{Code: "today", Name: "Today", Description: "Returns current date", FunctionType: dsl.FunctionTypeBuiltIn, ReturnType: "date", Parameters: "[]", Example: "today()", IsActive: ptr(true)},
		{Code: "concat", Name: "Concatenate", Description: "Concatenates strings", FunctionType: dsl.FunctionTypeBuiltIn, ReturnType: "string", Parameters: `[{"name": "values", "type": "string[]", "required": true}]`, Example: `concat("Hello", " ", "World")`, IsActive: ptr(true)},
		{Code: "sum", Name: "Sum", Description: "Returns sum of numbers", FunctionType: dsl.FunctionTypeBuiltIn, ReturnType: "number", Parameters: `[{"name": "values", "type": "number[]", "required": true}]`, Example: "sum(1, 2, 3)", IsActive: ptr(true)},
		{Code: "avg", Name: "Average", Description: "Returns average of numbers", FunctionType: dsl.FunctionTypeBuiltIn, ReturnType: "number", Parameters: `[{"name": "values", "type": "number[]", "required": true}]`, Example: "avg(1, 2, 3)", IsActive: ptr(true)},
		{Code: "if", Name: "If", Description: "Conditional expression", FunctionType: dsl.FunctionTypeBuiltIn, ReturnType: "any", Parameters: `[{"name": "condition", "type": "boolean", "required": true}, {"name": "then", "type": "any", "required": true}, {"name": "else", "type": "any", "required": false}]`, Example: `if(age >= 18, "Adult", "Minor")`, IsActive: ptr(true)},
		{Code: "format_date", Name: "Format Date", Description: "Formats date to string", FunctionType: dsl.FunctionTypeBuiltIn, ReturnType: "string", Parameters: `[{"name": "date", "type": "date", "required": true}, {"name": "format", "type": "string", "required": false, "default": "YYYY-MM-DD"}]`, Example: `format_date(today(), "YYYY/MM/DD")`, IsActive: ptr(true)},
	}

	for _, f := range functions {
		if err := db.Where("code = ?", f.Code).FirstOrCreate(&f).Error; err != nil {
			return err
		}
	}
	log.Println("DSL functions seeded")
	return nil
}

func seedAdminUser(db *gorm.DB) error {
	// Hash password
	passwordHash, err := service.HashPassword("Admin@123")
	if err != nil {
		return err
	}

	adminUser := domain.User{
		ID:           1,
		RegNo:        "АА00000000",
		FamilyName:   "Систем",
		LastName:     "Админ",
		FirstName:    "Супер",
		Gender:       1,
		BirthDate:    "1990-01-01",
		PhoneNo:      "99999999",
		Email:        "admin@gerege.mn",
		PasswordHash: passwordHash,
		IsActive:     ptr(true),
		LanguageCode: "mn",
	}

	if err := db.Where("id = ?", adminUser.ID).FirstOrCreate(&adminUser).Error; err != nil {
		return err
	}

	// Assign super_admin role
	adminRole := domain.UserSystemRole{
		ID:       1,
		UserID:   1,
		SystemID: ptr(1),
		RoleID:   1,
		IsActive: ptr(true),
	}

	if err := db.Where("id = ?", adminRole.ID).FirstOrCreate(&adminRole).Error; err != nil {
		return err
	}

	// Assign DSL admin role
	dslRole := domain.UserSystemRole{
		ID:       2,
		UserID:   1,
		SystemID: ptr(2),
		RoleID:   4,
		IsActive: ptr(true),
	}

	if err := db.Where("id = ?", dslRole.ID).FirstOrCreate(&dslRole).Error; err != nil {
		return err
	}

	log.Println("Admin user seeded")
	return nil
}

func seedRoleMenus(db *gorm.DB) error {
	// Assign all admin menus (1-20) to super_admin role (1)
	for menuID := 1; menuID <= 20; menuID++ {
		roleMenu := domain.RoleMenu{
			RoleID: 1,
			MenuID: menuID,
		}
		if err := db.Where("role_id = ? AND menu_id = ?", roleMenu.RoleID, roleMenu.MenuID).FirstOrCreate(&roleMenu).Error; err != nil {
			return err
		}
	}

	// Assign all admin menus to admin role (2)
	for menuID := 1; menuID <= 20; menuID++ {
		roleMenu := domain.RoleMenu{
			RoleID: 2,
			MenuID: menuID,
		}
		if err := db.Where("role_id = ? AND menu_id = ?", roleMenu.RoleID, roleMenu.MenuID).FirstOrCreate(&roleMenu).Error; err != nil {
			return err
		}
	}

	// Assign all DSL menus (21-33) to dsl_admin role (4)
	for menuID := 21; menuID <= 33; menuID++ {
		roleMenu := domain.RoleMenu{
			RoleID: 4,
			MenuID: menuID,
		}
		if err := db.Where("role_id = ? AND menu_id = ?", roleMenu.RoleID, roleMenu.MenuID).FirstOrCreate(&roleMenu).Error; err != nil {
			return err
		}
	}

	// Assign all DSL menus to dsl_developer role (5)
	for menuID := 21; menuID <= 33; menuID++ {
		roleMenu := domain.RoleMenu{
			RoleID: 5,
			MenuID: menuID,
		}
		if err := db.Where("role_id = ? AND menu_id = ?", roleMenu.RoleID, roleMenu.MenuID).FirstOrCreate(&roleMenu).Error; err != nil {
			return err
		}
	}

	log.Println("Role menus seeded")
	return nil
}

func seedPermissions(db *gorm.DB) error {
	// Permission format: {system.module.action}
	// Admin system permissions (system_id = 1)
	adminModules := []string{"user", "organization", "system", "module", "action", "role", "permission", "menu", "device", "session", "monitoring", "language", "translation"}
	adminActions := []string{"view", "create", "update", "delete"}

	permID := 1
	var permissions []domain.Permission

	// Admin system permissions
	for moduleIdx, module := range adminModules {
		for actionIdx, action := range adminActions {
			code := "admin." + module + "." + action
			actionID := int64(actionIdx + 1)
			perm := domain.Permission{
				ID:       permID,
				Code:     code,
				Name:     code,
				SystemID: ptr(1),
				ModuleID: moduleIdx + 1,
				ActionID: &actionID,
				IsActive: ptr(true),
			}
			permissions = append(permissions, perm)
			permID++
		}
	}

	// DSL system permissions (system_id = 2)
	dslModules := []string{"schema", "field", "rule", "workflow", "template", "function", "variable", "executor", "log"}
	dslActions := []string{"view", "create", "update", "delete", "execute"}

	for moduleIdx, module := range dslModules {
		for actionIdx, action := range dslActions {
			code := "dsl." + module + "." + action
			actionID := int64(actionIdx + 1)
			if actionIdx == 4 {
				actionID = 9 // execute action ID
			}
			perm := domain.Permission{
				ID:       permID,
				Code:     code,
				Name:     code,
				SystemID: ptr(2),
				ModuleID: moduleIdx + 14, // DSL modules start at ID 14
				ActionID: &actionID,
				IsActive: ptr(true),
			}
			permissions = append(permissions, perm)
			permID++
		}
	}

	for _, perm := range permissions {
		if err := db.Where("id = ?", perm.ID).FirstOrCreate(&perm).Error; err != nil {
			return err
		}
	}

	log.Println("Permissions seeded:", len(permissions))
	return nil
}

func seedRolePermissions(db *gorm.DB) error {
	// Get all admin permissions (system_id = 1)
	var adminPermissions []domain.Permission
	if err := db.Where("system_id = ?", 1).Find(&adminPermissions).Error; err != nil {
		return err
	}

	// Assign all admin permissions to super_admin role (1)
	for _, perm := range adminPermissions {
		rp := domain.RolePermission{
			RoleID:       1, // super_admin
			PermissionID: perm.ID,
		}
		if err := db.Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	// Assign all admin permissions to admin role (2)
	for _, perm := range adminPermissions {
		rp := domain.RolePermission{
			RoleID:       2, // admin
			PermissionID: perm.ID,
		}
		if err := db.Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	// Get all DSL permissions (system_id = 2)
	var dslPermissions []domain.Permission
	if err := db.Where("system_id = ?", 2).Find(&dslPermissions).Error; err != nil {
		return err
	}

	// Assign all DSL permissions to dsl_admin role (4)
	for _, perm := range dslPermissions {
		rp := domain.RolePermission{
			RoleID:       4, // dsl_admin
			PermissionID: perm.ID,
		}
		if err := db.Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	// Assign all DSL permissions to dsl_developer role (5)
	for _, perm := range dslPermissions {
		rp := domain.RolePermission{
			RoleID:       5, // dsl_developer
			PermissionID: perm.ID,
		}
		if err := db.Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).FirstOrCreate(&rp).Error; err != nil {
			return err
		}
	}

	log.Println("Role permissions seeded")
	return nil
}
