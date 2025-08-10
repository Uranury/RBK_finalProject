package models

import (
	"github.com/google/uuid"
	"time"
)

type Skin struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	OwnerID   *uuid.UUID `json:"owner_id" db:"owner_id"`
	Name      string     `json:"name" db:"name"`
	Rarity    string     `json:"rarity" db:"rarity"`
	Condition float64    `json:"condition" db:"condition"`
	Price     float64    `json:"price" db:"price"`
	Image     string     `json:"image" db:"image"`
	Available bool       `json:"available" db:"available"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}
