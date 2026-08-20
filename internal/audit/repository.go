package audit

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	auditLog *AuditLog,
) error {
	const query = `
		INSERT INTO audit_logs (
			id,
			user_id,
			session_id,
			event,
			ip_address,
			user_agent
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		auditLog.ID,
		auditLog.UserID,
		auditLog.SessionID,
		auditLog.Event,
		auditLog.IPAddress,
		auditLog.UserAgent,
	)

	return err
}
