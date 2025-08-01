package models

import "github.com/Uranury/RBK_finalProject/internal/auth"

type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Role     auth.Role `json:"role"`
}
