package main

import (
	"context"
	"github.com/Uranury/RBK_finalProject/internal/auth"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Uranury/RBK_finalProject/internal/http_server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	appDeps, err := InitDeps(logger)
	if err != nil {
		logger.Error("Failed to initialize dependencies", "error", err)
		os.Exit(1)
	}
	auth.NewService(appDeps.cfg.JWTKey)

	server, err := http_server.NewServer(
		appDeps.cfg,
		appDeps.db,
		appDeps.redisClient,
		appDeps.asynqClient,
		logger,
	)
	if err != nil {
		logger.Error("Failed to create http_server", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start http_server in goroutine
	go func() {
		logger.Info("Starting http_server", "port", appDeps.cfg.ListenAddr)
		if err := server.Run(); err != nil {
			logger.Error("Server failed to start", "error", err)
			cancel()
		}
	}()

	select {
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", "signal", sig)
	case <-ctx.Done():
		logger.Info("Context cancelled")
	}

	logger.Info("Shutting down http_server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Server stopped gracefully")
}
