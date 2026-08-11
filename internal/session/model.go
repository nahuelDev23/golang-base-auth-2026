package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID
	UserId           uuid.UUID
	RefreshTokenHash string
	ExpiresAt        time.Time
	Revoked          bool
	CreatedAt        time.Time
}
