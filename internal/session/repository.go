package session

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	session *Session,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO sessions (
			id,
			user_id,
			refresh_token_hash,
			expires_at,
			revoked
		)
		VALUES ($1, $2, $3, $4, $5)
		`,
		session.ID,
		session.UserId,
		session.RefreshTokenHash,
		session.ExpiresAt,
		session.Revoked,
	)

	return err
}

func (r *Repository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*Session, error) {

	session := &Session{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			user_id,
			refresh_token_hash,
			expires_at,
			revoked,
			created_at
		FROM sessions
		WHERE id = $1
		`,
		id,
	).Scan(
		&session.ID,
		&session.UserId,
		&session.RefreshTokenHash,
		&session.ExpiresAt,
		&session.Revoked,
		&session.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return session, nil
}

func (r *Repository) Revoke(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := r.db.Exec(
		ctx,
		`
			UPDATE sessions
			SET revoked = true
			WHERE id = $1
		`,
		id,
	)

	return err
}

func (r *Repository) RevokeAll(
	ctx context.Context,
	userID uuid.UUID,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE sessions
		SET revoked = true
		WHERE user_id = $1
		`,
		userID,
	)

	return err
}

func (r *Repository) FindByRefreshTokenHash(
	ctx context.Context,
	hash string,
) (*Session, error) {

	session := &Session{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			user_id,
			refresh_token_hash,
			expires_at,
			revoked,
			created_at
		FROM sessions
		WHERE refresh_token_hash = $1
		`,
		hash,
	).Scan(
		&session.ID,
		&session.UserId,
		&session.RefreshTokenHash,
		&session.ExpiresAt,
		&session.Revoked,
		&session.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return session, nil
}
