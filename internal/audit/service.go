package audit

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Log(
	ctx context.Context,
	userID *uuid.UUID,
	sessionID *uuid.UUID,
	event string,
	ipAddress *string,
	userAgent *string,
) error {
	auditLog := &AuditLog{
		ID:        uuid.New(),
		UserID:    userID,
		SessionID: sessionID,
		Event:     event,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	return s.repository.Create(ctx, auditLog)
}
