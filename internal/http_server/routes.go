package http_server

import (
	"github.com/Uranury/RBK_finalProject/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) setupRoutes() {
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	protected := s.router.Group("/", middleware.JWTAuthMiddleware(s.authService))
	protected.GET("/lol", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "kek"})
	})
}
