package http_server

import (
	"net/http"

	"github.com/Uranury/RBK_finalProject/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (s *Server) setupRoutes() {
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	s.router.POST("/signup", s.userHandler.Signup)
	s.router.POST("/login", s.userHandler.Login)

	protected := s.router.Group("/", middleware.JWTAuthMiddleware(s.authService))
	// Marketplace
	protected.GET("/marketplace/skins", s.marketplaceHandler.ListAvailable)
	protected.GET("/marketplace/skins/mine", s.marketplaceHandler.ListMine)
	protected.GET("/marketplace/orders/:order_id", s.marketplaceHandler.GetOrder)
	protected.POST("/marketplace/purchase", s.marketplaceHandler.Purchase)
	// Skin creation (protected)
	protected.POST("/skins", s.skinHandler.Create)
}
