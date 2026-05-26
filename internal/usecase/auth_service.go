package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/shrin00/pen/internal/domain"
	"github.com/shrin00/pen/internal/security"
)

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrEmailAlreadyExists = errors.New("Email already exists")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	FindByTokenHash(ctx context.Context, token_hash string) (*domain.Session, error)
	FindValidByTokenHash(ctx context.Context, token_hash string) (*domain.Session, error)
	RevokeByTokenHash(ctx context.Context, token_hash string) (*domain.Session, error)
}

type AuthService struct {
	user       UserRepository
	session    SessionRepository
	SessionTTL time.Duration
}

func NewAuthService(user UserRepository, session SessionRepository, sessionTTL time.Duration) *AuthService {
	return &AuthService{
		user:       user,
		session:    session,
		SessionTTL: sessionTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, email string, password string) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if err := validateCredentials(email, password); err != nil {
		return nil, err
	}

	existingUser, err := s.user.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: hashedPassword,
	}
	err = s.user.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*domain.User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	user, err := s.user.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	if !security.CheckPassword(password, user.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	user_token, err := security.GenerateRawToken()
	if err != nil {
		return nil, "", err
	}

	token_hash := security.TokenHash(user_token)
	session_domain := &domain.Session{
		UserID:    user.ID,
		TokenHash: token_hash,
		ExpiresAt: time.Now().Add(s.SessionTTL),
	}
	err = s.session.Create(ctx, session_domain)
	if err != nil {
		return nil, "", err
	}

	return user, user_token, nil
}

func validateCredentials(email string, password string) error {
	if email == "" || !strings.Contains(email, "@") {
		return ErrInvalidEmail
	}

	if len(password) < 8 {
		return ErrWeakPassword
	}
	return nil
}
