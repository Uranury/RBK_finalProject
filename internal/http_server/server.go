package http_server

import (
	"github.com/Uranury/RBK_finalProject/config"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	router      *gin.Engine
	cfg         *config.Config
	db          *sqlx.DB
	redisClient *redis.Client
	asynqClient *asynq.Client
}

func NewServer(cfg *config.Config, redisClient *redis.Client, asynqClient *asynq.Client) *Server {
}
