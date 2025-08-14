package http_server

import (
	"net/http"

	"github.com/Uranury/RBK_finalProject/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// TODO: Add query params for RemoveFromListing endpoint instead of JSON
// TODO: Consider doing the same for Purchase endpoint

func (s *Server) setupRoutes() {
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	s.router.POST("/signup", s.userHandler.Signup)
	s.router.POST("/login", s.userHandler.Login)

	// Swagger documentation
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public endpoints
	s.router.GET("/guns", s.skinHandler.GetGuns)
	s.router.GET("/wears", s.skinHandler.GetWears)

	protected := s.router.Group("/", middleware.JWTAuthMiddleware(s.authService))
	// Marketplace
	protected.GET("/marketplace/skins", s.marketplaceHandler.ListAvailable)
	protected.GET("/marketplace/skins/mine", s.marketplaceHandler.ListMine)
	protected.GET("/marketplace/orders/:order_id", s.marketplaceHandler.GetOrder)
	protected.POST("/marketplace/purchase", s.marketplaceHandler.Purchase)
	protected.DELETE("/marketplace/skins", s.marketplaceHandler.RemoveFromListing)
	// Skin creation (protected)
	protected.POST("/skins", s.skinHandler.Create)
	// Transactions
	protected.POST("/transactions/withdraw", s.transactionHandler.Withdraw)
	protected.POST("/transactions/deposit", s.transactionHandler.Deposit)
	protected.GET("/transactions/history", s.transactionHandler.GetHistory)
}
