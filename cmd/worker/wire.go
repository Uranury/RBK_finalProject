package main

import (
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/Uranury/RBK_finalProject/pkg/config"
	"github.com/hibiken/asynq"
	"log/slog"
)

type WorkerDeps struct {
	Cfg    *config.Config
	Server *asynq.Server
	Logger *slog.Logger
}

func InitWorkerDeps(logger *slog.Logger) (*WorkerDeps, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, apperrors.NewInternalError("couldn't load config", err)
	}

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &WorkerDeps{
		Cfg:    cfg,
		Server: server,
		Logger: logger,
	}, nil
}
