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
	protected.GET("/lol", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "kek"})
	})
}
