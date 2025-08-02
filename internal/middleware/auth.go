package middleware

import (
	"github.com/Uranury/RBK_finalProject/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := authService.VerifyJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		return uuid.UUID{}, false
	}

	claims, ok := claimsVal.(auth.Claims)
	if !ok {
		return uuid.UUID{}, false
	}

	userIDRaw, ok := claims["user_id"]
	if !ok {
		return uuid.UUID{}, false
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		return uuid.UUID{}, false
	}

	id, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, false
	}

	return id, true
}

func GetUserRole(c *gin.Context) (auth.Role, bool) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		return "", false
	}

	claims, ok := claimsVal.(auth.Claims)
	if !ok {
		return "", false
	}

	roleRaw, ok := claims["role"]
	if !ok {
		return "", false
	}

	roleStr, ok := roleRaw.(string)
	if !ok {
		return "", false
	}

	return auth.Role(roleStr), true
}
