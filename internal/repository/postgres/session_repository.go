package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shrin00/pen/internal/domain"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (s *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO sessions (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, revoked_at, created_at
	`

	return s.db.QueryRow(ctx, query, session.UserID, session.TokenHash, session.ExpiresAt).Scan(
		&session.ID,
		&session.RevokedAt,
		&session.CreatedAt,
	)
}

func (s *SessionRepository) FindByTokenHash(ctx context.Context, token_hash string) (*domain.Session, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM sessions
		WHERE token_hash = $1
	`
	var session domain.Session
	err := s.db.QueryRow(ctx, query, token_hash).Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

func (s *SessionRepository) FindValidByTokenHash(ctx context.Context, token_hash string) (*domain.Session, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM sessions
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > now()
	`
	var session domain.Session
	err := s.db.QueryRow(ctx, query, token_hash).Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

func (s *SessionRepository) RevokeByTokenHash(ctx context.Context, token_hash string) (*domain.Session, error) {
	query := `
		UPDATE sessions
		SET revoked_at = $1
		WHERE token_hash = $2 AND revoked_at IS NULL
		RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at
	`
	revoked_time := time.Now()

	var session domain.Session
	err := s.db.QueryRow(ctx, query, revoked_time, token_hash).Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}
