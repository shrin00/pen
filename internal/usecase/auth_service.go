package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/shrin00/pen/internal/domain"
	"github.com/shrin00/pen/internal/security"
)

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrEmailAlreadyExists = errors.New("Email already exists")
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type AuthService struct {
	user UserRepository
}

func NewAuthService(user UserRepository) *AuthService {
	return &AuthService{
		user: user,
	}
}

func (s *AuthService) Register(ctx context.Context, email string, password string) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

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

func (s *AuthService) Login(ctx context.Context, email string, password string) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	user, err := s.user.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !security.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
