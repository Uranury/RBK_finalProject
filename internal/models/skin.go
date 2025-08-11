package models

import (
	"time"

	"github.com/google/uuid"
)

type Gun string

const (
	// Pistols
	AK47        Gun = "AK-47"
	M4A4        Gun = "M4A4"
	M4A1S       Gun = "M4A1-S"
	DesertEagle Gun = "Desert Eagle"
	USPS        Gun = "USP-S"
	Glock18     Gun = "Glock-18"
	P250        Gun = "P250"
	Tec9        Gun = "Tec-9"
	CZ75        Gun = "CZ75-Auto"

	// Rifles
	AWP    Gun = "AWP"
	SSG08  Gun = "SSG 08"
	SCAR20 Gun = "SCAR-20"
	G3SG1  Gun = "G3SG1"

	// SMGs
	MP9     Gun = "MP9"
	MAC10   Gun = "MAC-10"
	MP7     Gun = "MP7"
	P90     Gun = "P90"
	UMP45   Gun = "UMP-45"
	PPBizon Gun = "PP-Bizon"

	// Shotguns
	Nova     Gun = "Nova"
	XM1014   Gun = "XM1014"
	MAG7     Gun = "MAG-7"
	SawedOff Gun = "Sawed-Off"

	// Machine Guns
	M249  Gun = "M249"
	Negev Gun = "Negev"

	// Knives
	Karambit      Gun = "Karambit"
	Butterfly     Gun = "Butterfly Knife"
	M9Bayonet     Gun = "M9 Bayonet"
	Bayonet       Gun = "Bayonet"
	FlipKnife     Gun = "Flip Knife"
	GutKnife      Gun = "Gut Knife"
	Huntsman      Gun = "Huntsman Knife"
	ShadowDaggers Gun = "Shadow Daggers"

	// Other
	Falchion Gun = "Falchion Knife"
	Bowie    Gun = "Bowie Knife"
	Navaja   Gun = "Navaja Knife"
	Stiletto Gun = "Stiletto Knife"
	Ursus    Gun = "Ursus Knife"
	Nomad    Gun = "Nomad Knife"
	Paracord Gun = "Paracord Knife"
	Survival Gun = "Survival Knife"
	Classic  Gun = "Classic Knife"
)

type Skin struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	OwnerID   *uuid.UUID `json:"owner_id" db:"owner_id"`
	Name      string     `json:"name" db:"name"`
	Gun       Gun        `json:"gun" db:"gun"`
	Rarity    string     `json:"rarity" db:"rarity"`
	Condition float64    `json:"condition" db:"condition"`
	Price     float64    `json:"price" db:"price"`
	Image     string     `json:"image" db:"image"`
	Available bool       `json:"available" db:"available"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}
