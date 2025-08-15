package http_server

import (
	"github.com/Uranury/RBK_finalProject/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupRoutes() {
	protected := s.router.Group("/", middleware.JWTAuthMiddleware(s.authService))
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.router.POST("/signup", s.userHandler.Signup)
	s.router.POST("/login", s.userHandler.Login)
	protected.GET("/profile", s.userHandler.Profile)

	// Public endpoints
	s.router.GET("/guns", s.skinHandler.GetGuns)
	s.router.GET("/wears", s.skinHandler.GetWears)

	// Marketplace
	s.router.GET("/marketplace/skins", s.marketplaceHandler.ListAvailable)
	protected.GET("/marketplace/skins/mine", s.marketplaceHandler.ListMine)
	protected.GET("/marketplace/orders/:order_id", s.marketplaceHandler.GetOrder)
	protected.POST("/marketplace/purchase", s.marketplaceHandler.Purchase)
	protected.DELETE("/marketplace/skins/:skin_id", s.marketplaceHandler.RemoveFromListing)
	protected.POST("/marketplace/sell", s.marketplaceHandler.Sell)
	// Skin creation (protected)
	protected.POST("/skins", s.skinHandler.Create)
	// Transactions
	protected.POST("/transactions/withdraw", s.transactionHandler.Withdraw)
	protected.POST("/transactions/deposit", s.transactionHandler.Deposit)
	protected.GET("/transactions/history", s.transactionHandler.GetHistory)
}
