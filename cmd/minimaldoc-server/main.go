// Package main provides the entrypoint for minimaldoc-server.
// This is a standalone backend service that adds optional dynamic features
// to MinimalDoc static sites.
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

	"github.com/studiowebux/minimaldoc/internal/server/api"
	"github.com/studiowebux/minimaldoc/internal/server/config"
	"github.com/studiowebux/minimaldoc/internal/server/email"
	"github.com/studiowebux/minimaldoc/internal/server/scheduler"
	"github.com/studiowebux/minimaldoc/internal/server/storage"
	"github.com/studiowebux/minimaldoc/internal/server/store"
)

const version = "0.1.0"

func main() {
	// Print banner
	fmt.Printf(`
╔══════════════════════════════════════╗
║     MinimalDoc Server v%s        ║
║       Backend API Service            ║
╚══════════════════════════════════════╝
`, version)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := store.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize email sender
	emailSender, err := email.NewSender(cfg.Email)
	if err != nil {
		log.Fatalf("Failed to initialize email sender: %v", err)
	}
	log.Printf("Email provider: %s", cfg.Email.Provider)

	// Initialize storage
	fileStorage, err := storage.New(cfg.Storage)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	log.Printf("Storage provider: %s", cfg.Storage.Provider)

	// Initialize routers
	publicRouter := api.NewPublicRouter(cfg, db, emailSender)
	adminRouter := api.NewAdminRouter(cfg, db, emailSender, fileStorage)

	// Initialize and start background scheduler
	sched := scheduler.New(db, time.Minute)
	sched.Start()

	// Configure public HTTP server
	publicAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	publicSrv := &http.Server{
		Addr:         publicAddr,
		Handler:      publicRouter,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Configure admin HTTP server
	adminAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.AdminPort)
	adminSrv := &http.Server{
		Addr:         adminAddr,
		Handler:      adminRouter,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start public server
	go func() {
		log.Printf("Public API:  http://%s (tracking, feedback, newsletter)", publicAddr)
		if err := publicSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Public server failed: %v", err)
		}
	}()

	// Start admin server
	go func() {
		log.Printf("Admin API:   http://%s%s", adminAddr, cfg.Server.APIPath)
		log.Printf("Admin UI:    http://%s%s", adminAddr, cfg.Server.AdminPath)
		if err := adminSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Admin server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Stop scheduler first
	sched.Stop()

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown both servers
	if err := publicSrv.Shutdown(ctx); err != nil {
		log.Printf("Public server forced to shutdown: %v", err)
	}
	if err := adminSrv.Shutdown(ctx); err != nil {
		log.Printf("Admin server forced to shutdown: %v", err)
	}

	log.Println("Servers stopped")
}
