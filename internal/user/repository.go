package user

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

func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT 
			id,
			username,
			password_hash,
			role,
			created_at
		FROM users
		WHERE username = $1
		`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *Repository) Create(
	ctx context.Context,
	user *User,
) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO users (
			id,
			username,
			password_hash
			role
		)
		VALUES ($1, $2, $3)
		`,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.Role,
	)

	return err
}

func (r *Repository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*User, error) {
	user := &User{}

	err := r.db.QueryRow(
		ctx,
		`
  SELECT
		id,
		username,
		password_hash,
		role,
		created_at
		FROM users
		WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}
