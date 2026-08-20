package audit

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID
	UserID    *uuid.UUID
	SessionID *uuid.UUID
	Event     string
	IPAddress *string
	UserAgent *string
	CreatedAt time.Time
}
