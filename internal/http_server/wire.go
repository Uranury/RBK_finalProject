package http_server

import (
	"github.com/Uranury/RBK_finalProject/internal/auth"
	"github.com/Uranury/RBK_finalProject/internal/handlers"
	"github.com/Uranury/RBK_finalProject/internal/repositories/user"
	"github.com/Uranury/RBK_finalProject/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (s *Server) initDependencies() error {
	// Initialize repositories
	userRepo := user.NewRepository(s.db)

	// Initialize services
	s.authService = auth.NewService(s.cfg.JWTKey)
	userService := services.NewUser(userRepo, s.authService)

	// Initialize handlers
	s.userHandler = handlers.NewUserHandler(userService)

	return nil
}

func (s *Server) initHTTPServer() {
	s.router = gin.Default()

	s.httpServer = &http.Server{
		Addr:         s.cfg.ListenAddr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
