package http_server

import (
	"context"
	"github.com/Uranury/RBK_finalProject/internal/auth"
	"github.com/Uranury/RBK_finalProject/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	router      *gin.Engine
	httpServer  *http.Server
	cfg         *config.Config
	db          *sqlx.DB
	asynqClient *asynq.Client
	authService *auth.Service
	redisClient *redis.Client
	logger      *slog.Logger
}

func NewServer(
	cfg *config.Config,
	db *sqlx.DB,
	redisClient *redis.Client,
	asynqClient *asynq.Client,
	logger *slog.Logger) (*Server, error) {

	router := gin.Default()
	httpServer := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	authService := auth.NewService(cfg.JWTKey)

	s := &Server{
		router:      router,
		httpServer:  httpServer,
		cfg:         cfg,
		db:          db,
		asynqClient: asynqClient,
		authService: authService,
		redisClient: redisClient,
		logger:      logger,
	}
	s.setupRoutes()
	return s, nil
}

func (s *Server) Run() error {
	s.logger.Info("Server starting", "address", s.cfg.ListenAddr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Starting graceful shutdown")
	done := make(chan error, 1)

	go func() {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Error("HTTP http_server shutdown failed", "error", err)
			done <- err
			return
		}
		if s.db != nil {
			s.logger.Info("Closing database connections...")
			if err := s.db.Close(); err != nil {
				s.logger.Error("Failed to close database", "error", err)
			}
		}

		if s.redisClient != nil {
			s.logger.Info("Closing redis client...")
			if err := s.redisClient.Close(); err != nil {
				s.logger.Error("Failed to close redis client", "error", err)
			}
		}

		if s.asynqClient != nil {
			s.logger.Info("Closing Asynq client...")
			if err := s.asynqClient.Close(); err != nil {
				s.logger.Error("Failed to close Asynq client", "error", err)
			}
		}

		s.logger.Info("Graceful shutdown completed")
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		s.logger.Info("Graceful shutdown timed out")
		return ctx.Err()
	}
}
