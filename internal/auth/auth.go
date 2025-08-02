package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type Role string
type Claims = jwt.MapClaims

const (
	Admin Role = "admin"
	User  Role = "user"
)

type Service struct {
	jwtKey []byte
}

func NewService(secret string) *Service {
	return &Service{
		jwtKey: []byte(secret),
	}
}

func (s *Service) GenerateJWT(userID uuid.UUID, role Role) (string, error) {
	if role == "" {
		role = User
	}

	claims := jwt.MapClaims{
		"user_id": userID.String(), // Explicitly convert to string
		"role":    string(role),    // Be explicit about string conversion
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}

func (s *Service) VerifyJWT(tokenString string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtKey, nil
	})
	if err != nil || !token.Valid {
		if err != nil {
			return Claims{}, fmt.Errorf("failed to parse token: %w", err)
		}
		return Claims{}, fmt.Errorf("token is not valid")
	}
	return claims, nil
}
