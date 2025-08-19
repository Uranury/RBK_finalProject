package main

import (
	"github.com/hibiken/asynq"
	"log/slog"
	"os"
)

func main() {
	logger := slog.Default()

	deps, err := InitWorkerDeps(logger)
	if err != nil {
		logger.Error("failed to init worker", "err", err)
		os.Exit(1)
	}

	mux := asynq.NewServeMux()

	// Register your task handlers here
	// e.g. mux.HandleFunc(tasks.TypeGeneratePDF, handlers.HandleGeneratePDF(deps.DB, deps.Logger))

	if err := deps.Server.Run(mux); err != nil {
		logger.Error("could not run asynq server", "err", err)
		os.Exit(1)
	}
}
