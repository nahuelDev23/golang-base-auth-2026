package user

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	return s.repository.FindByUsername(ctx, username)
}

func (s *Service) Register(
	ctx context.Context,
	username,
	password string,
) error {
	existing, err := s.repository.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	if existing != nil {
		return ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user := &User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hash),
	}

	return s.repository.Create(ctx, user)

}

func (s *Service) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*User, error) {
	return s.repository.FindByID(ctx, id)
}
