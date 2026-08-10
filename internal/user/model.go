package user

import (
	"github.com/google/uuid"

	"time"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}
