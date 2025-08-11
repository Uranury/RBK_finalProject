package http_server

import (
	"net/http"
	"time"

	"github.com/Uranury/RBK_finalProject/internal/auth"
	"github.com/Uranury/RBK_finalProject/internal/handlers"
	orderRepoPkg "github.com/Uranury/RBK_finalProject/internal/repositories/order"
	skinRepoPkg "github.com/Uranury/RBK_finalProject/internal/repositories/skin"
	"github.com/Uranury/RBK_finalProject/internal/repositories/user"
	"github.com/Uranury/RBK_finalProject/internal/services"
	"github.com/gin-gonic/gin"
)

func (s *Server) initDependencies() error {
	// Initialize repositories
	userRepo := user.NewRepository(s.db)
	skinRepo := skinRepoPkg.NewRepository(s.db)
	ordRepo := orderRepoPkg.NewRepository(s.db)

	// Initialize services
	s.authService = auth.NewService(s.cfg.JWTKey)
	userService := services.NewUser(userRepo, s.authService, s.logger)
	skinService := services.NewSkin(skinRepo, s.logger)
	marketplaceService := services.NewMarketplaceService(skinRepo, ordRepo, userRepo, s.asynqClient, s.db, s.logger)

	// Initialize handlers
	s.userHandler = handlers.NewUserHandler(userService)
	s.skinHandler = handlers.NewSkinHandler(skinService)
	s.marketplaceHandler = handlers.NewMarketplaceHandler(marketplaceService)

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
