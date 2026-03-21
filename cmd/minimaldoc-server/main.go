// Package main provides the entrypoint for minimaldoc-server.
// This is a standalone backend service that adds optional dynamic features
// to MinimalDoc static sites.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/studiowebux/minimaldoc/internal/server/api"
	"github.com/studiowebux/minimaldoc/internal/server/auth"
	"github.com/studiowebux/minimaldoc/internal/server/config"
	"github.com/studiowebux/minimaldoc/internal/server/email"
	"github.com/studiowebux/minimaldoc/internal/server/scheduler"
	"github.com/studiowebux/minimaldoc/internal/server/storage"
	"github.com/studiowebux/minimaldoc/internal/server/store"
	"github.com/studiowebux/minimaldoc/internal/server/telemetry"
	"github.com/studiowebux/minimaldoc/internal/version"
)

func main() {
	// Print banner
	fmt.Printf(`
╔══════════════════════════════════════╗
║     MinimalDoc Server v%s        ║
║       Backend API Service            ║
╚══════════════════════════════════════╝
`, version.Version)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// slog not yet initialized, use fmt+os.Exit
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Validate OAuth provider URLs (HTTPS required)
	if cfg.Auth.EnableOAuth {
		if err := auth.ValidateOAuthProviderURLs(cfg.OAuth.Providers); err != nil {
			fmt.Fprintf(os.Stderr, "OAuth configuration error: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize structured logging
	var logHandler slog.Handler
	if cfg.Server.Environment == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(logHandler))

	// Initialize OpenTelemetry
	otelShutdown, err := telemetry.Init(context.Background(), telemetry.Config{
		Enabled:     cfg.Telemetry.Enabled,
		Endpoint:    cfg.Telemetry.Endpoint,
		ServiceName: cfg.Telemetry.ServiceName,
		Environment: cfg.Server.Environment,
	})
	if err != nil {
		slog.Error("failed to initialize opentelemetry", "error", err)
		os.Exit(1)
	}
	defer otelShutdown(context.Background())

	// Initialize database
	db, err := store.New(cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize email sender
	emailSender, err := email.NewSender(cfg.Email)
	if err != nil {
		slog.Error("failed to initialize email sender", "error", err)
		os.Exit(1)
	}
	slog.Info("email provider initialized", "provider", cfg.Email.Provider)

	// Initialize storage
	fileStorage, err := storage.New(cfg.Storage)
	if err != nil {
		slog.Error("failed to initialize storage", "error", err)
		os.Exit(1)
	}
	slog.Info("storage provider initialized", "provider", cfg.Storage.Provider)

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

	// Start servers — send errors to channel instead of os.Exit to allow cleanup
	errCh := make(chan error, 2)

	go func() {
		slog.Info("public API started", "addr", publicAddr)
		if err := publicSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("public server failed", "error", err)
			errCh <- err
		}
	}()

	go func() {
		slog.Info("admin API started", "addr", adminAddr, "api_path", cfg.Server.APIPath, "admin_path", cfg.Server.AdminPath)
		if err := adminSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("admin server failed", "error", err)
			errCh <- err
		}
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case err := <-errCh:
		slog.Error("server startup failed, shutting down", "error", err)
	}

	slog.Info("shutting down servers")

	// Stop scheduler first
	sched.Stop()

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown both servers
	if err := publicSrv.Shutdown(ctx); err != nil {
		slog.Warn("public server forced to shutdown", "error", err)
	}
	if err := adminSrv.Shutdown(ctx); err != nil {
		slog.Warn("admin server forced to shutdown", "error", err)
	}

	slog.Info("servers stopped")
}
