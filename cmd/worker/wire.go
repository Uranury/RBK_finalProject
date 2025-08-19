package main

import (
	"log/slog"

	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/Uranury/RBK_finalProject/pkg/config"
	"github.com/Uranury/RBK_finalProject/pkg/db"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
)

type WorkerDeps struct {
	Cfg    *config.Config
	Server *asynq.Server
	DB     *sqlx.DB
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

	database, err := db.InitDBWithoutMigrations("postgres", cfg.DbURL, logger)
	if err != nil {
		return nil, apperrors.NewInternalError("couldn't init database", err)
	}

	return &WorkerDeps{
		Cfg:    cfg,
		Server: server,
		DB:     database,
		Logger: logger,
	}, nil
}
