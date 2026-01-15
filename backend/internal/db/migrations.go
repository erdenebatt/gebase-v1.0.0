package db

import (
	"log"

	"gebase/internal/domain"
	"gebase/internal/domain/dsl"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Core domain models
	err := db.AutoMigrate(
		// Base entities
		&domain.Language{},
		&domain.Translation{},
		&domain.Action{},
		&domain.System{},
		&domain.Module{},
		&domain.ModuleAction{},
		&domain.Permission{},
		&domain.Role{},
		&domain.RolePermission{},
		&domain.Menu{},
		&domain.RoleMenu{},

		// Organization entities
		&domain.OrganizationType{},
		&domain.Organization{},
		&domain.OrganizationSystem{},

		// User entities
		&domain.User{},
		&domain.UserSystemRole{},

		// Device & Session entities
		&domain.Device{},
		&domain.Session{},
		&domain.SessionSystemHistory{},
	)
	if err != nil {
		return err
	}

	// DSL domain models
	err = db.AutoMigrate(
		&dsl.Schema{},
		&dsl.Field{},
		&dsl.Rule{},
		&dsl.Workflow{},
		&dsl.WorkflowStep{},
		&dsl.WorkflowInstance{},
		&dsl.Template{},
		&dsl.Function{},
		&dsl.Variable{},
		&dsl.ExecutionLog{},
	)
	if err != nil {
		return err
	}

	// Create unique constraints
	if err := createUniqueConstraints(db); err != nil {
		log.Printf("Warning: some constraints may already exist: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func createUniqueConstraints(db *gorm.DB) error {
	constraints := []string{
		// Module: system_id + code
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_modules_system_code ON modules(system_id, code) WHERE deleted_date IS NULL`,

		// ModuleAction: module_id + action_id
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_module_actions_module_action ON module_actions(module_id, action_id) WHERE deleted_date IS NULL`,

		// RolePermission: role_id + permission_id
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_role_permissions_role_permission ON role_permissions(role_id, permission_id) WHERE deleted_date IS NULL`,

		// RoleMenu: role_id + menu_id
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_role_menus_role_menu ON role_menus(role_id, menu_id) WHERE deleted_date IS NULL`,

		// Menu: system_id + code
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_menus_system_code ON menus(system_id, code) WHERE deleted_date IS NULL`,

		// OrganizationSystem: organization_id + system_id
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_org_systems_org_system ON organization_systems(organization_id, system_id) WHERE deleted_date IS NULL`,

		// UserSystemRole: user_id + system_id + role_id + organization_id
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_user_system_roles_unique ON user_system_roles(user_id, COALESCE(system_id, 0), role_id, COALESCE(organization_id, 0)) WHERE deleted_date IS NULL`,

		// Translation: language_code + key
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_translations_lang_key ON translations(language_code, key) WHERE deleted_date IS NULL`,

		// DSL Field: schema_id + code
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_dsl_fields_schema_code ON dsl_fields(schema_id, code) WHERE deleted_date IS NULL`,

		// DSL Variable: code + scope + organization_id + user_id
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_dsl_variables_unique ON dsl_variables(code, scope, COALESCE(organization_id, 0), COALESCE(user_id, 0)) WHERE deleted_date IS NULL`,
	}

	for _, constraint := range constraints {
		if err := db.Exec(constraint).Error; err != nil {
			log.Printf("Warning: constraint creation issue: %v", err)
		}
	}

	return nil
}
