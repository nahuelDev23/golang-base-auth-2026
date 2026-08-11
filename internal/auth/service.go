package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"test.nahueldev23.com/internal/session"
	"test.nahueldev23.com/internal/user"
)

type Service struct {
	users    *user.Repository
	sessions *session.Repository

	jwtSecret []byte
}

func NewService(
	users *user.Repository,
	sessions *session.Repository,
	jwtSecret string,
) *Service {
	return &Service{
		users:     users,
		sessions:  sessions,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *Service) Login(
	ctx context.Context,
	username string,
	password string,
) (string, string, error) {

	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	sessionID := uuid.New()

	token, err := s.GenerateAccessToken(
		user.ID,
		sessionID,
	)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	session := &session.Session{
		ID:               sessionID,
		UserId:           user.ID,
		RefreshTokenHash: HashToken(refreshToken),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		Revoked:          false,
	}

	err = s.sessions.Create(
		ctx,
		session,
	)

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil

}

func (s *Service) ValidateToken(
	ctx context.Context,
	token string,
) (*Claims, error) {

	claims, err := s.ParseToken(token)
	if err != nil {
		return nil, err
	}

	session, err := s.sessions.FindByID(
		ctx,
		claims.SessionID,
	)

	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, ErrInvalidSession
	}

	if session.Revoked {
		return nil, ErrInvalidSession
	}

	return claims, nil
}

func (s *Service) Logout(
	ctx context.Context,
	sessionID uuid.UUID,
) error {

	return s.sessions.Revoke(ctx, sessionID)
}

func (s *Service) LogoutAll(
	ctx context.Context,
	userID uuid.UUID,
) error {
	return s.sessions.RevokeAll(ctx, userID)
}

func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
) (string, error) {

	hash := HashToken(refreshToken)

	session, err := s.sessions.FindByRefreshTokenHash(
		ctx,
		hash,
	)

	if err != nil {
		return "", err
	}

	if session == nil {
		return "", ErrInvalidSession
	}

	if session.Revoked {
		return "", ErrInvalidSession
	}

	if time.Now().After(session.ExpiresAt) {
		return "", ErrInvalidSession
	}

	accessToken, err := s.GenerateAccessToken(
		session.UserId,
		session.ID,
	)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}
