package user

import (
	"github.com/google/uuid"

	"time"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}
