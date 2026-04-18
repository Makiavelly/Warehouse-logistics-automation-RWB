package models

import "time"

type TokenClaims struct {
	UserID string
	Role   string
	Exp    time.Time
}