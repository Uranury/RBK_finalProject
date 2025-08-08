package main

import (
	"github.com/Uranury/RBK_finalProject/pkg/config"
	"github.com/Uranury/RBK_finalProject/pkg/db"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type AppDeps struct {
	cfg         *config.Config
	db          *sqlx.DB
	redisClient *redis.Client
	asynqClient *asynq.Client
	logger      *slog.Logger
}

func InitDeps(logger *slog.Logger) (*AppDeps, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	database, err := db.InitDB("postgres", cfg.DbURL, cfg.MigrationsPath, logger)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.RedisAddr,
	})

	return &AppDeps{
		cfg:         cfg,
		db:          database,
		redisClient: redisClient,
		asynqClient: asynqClient,
		logger:      logger,
	}, nil
}
