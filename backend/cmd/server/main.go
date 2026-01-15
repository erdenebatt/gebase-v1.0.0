package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gebase/internal/app"
	"gebase/internal/config"
	"gebase/internal/db"

	"gorm.io/gorm"
)

// @title Gebase API
// @version 1.0
// @description Enterprise Auth & RBAC Platform API
// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	database, err := db.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Check for seed command
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		log.Println("Running seed data...")
		if err := runSeed(database); err != nil {
			log.Fatalf("Failed to seed data: %v", err)
		}
		log.Println("Seed completed successfully")
		return
	}

	// Create application container
	container := app.NewContainer(cfg, database)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      container.Router.GetEngine(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func runSeed(database interface{}) error {
	gormDB, ok := database.(*gorm.DB)
	if !ok {
		return fmt.Errorf("database is not *gorm.DB")
	}
	return db.RunSeed(gormDB)
}
