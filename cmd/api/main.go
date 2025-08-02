package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Uranury/RBK_finalProject/internal/http_server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize dependencies
	appDeps, err := InitDeps(logger)
	if err != nil {
		logger.Error("Failed to initialize dependencies", "error", err)
		os.Exit(1)
	}

	// Create server
	server, err := http_server.NewServer(
		appDeps.cfg,
		appDeps.db,
		appDeps.redisClient,
		appDeps.asynqClient,
		logger,
	)
	if err != nil {
		logger.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	// Graceful shutdown setup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		logger.Info("Starting server", "port", appDeps.cfg.ListenAddr)
		if err := server.Run(); err != nil {
			logger.Error("Server failed to start", "error", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	select {
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", "signal", sig)
	case <-ctx.Done():
		logger.Info("Context cancelled")
	}

	// Graceful shutdown
	logger.Info("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Server stopped gracefully")
}
