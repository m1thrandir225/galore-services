package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/m1thrandir225/galore-services/internal/cache"
	db "github.com/m1thrandir225/galore-services/internal/db/sqlc"
)

type SessionService interface {
	CreateSession(
		ctx context.Context,
		id uuid.UUID,
		email,
		refreshToken,
		userAgent,
		clientIp string,
		expiresAt time.Time,
		isBlocked bool,
	) (db.Session, error)
	InvalidateSession(ctx context.Context, sessionID uuid.UUID) (db.Session, error)
	GetAllUserSessions(ctx context.Context, email string) ([]db.Session, error)
	GetSession(ctx context.Context, sessionID uuid.UUID) (db.Session, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
}

type sessionService struct {
	store db.Store
	cache cache.Store
}

func NewSessionService(store db.Store, cache cache.Store) SessionService {
	return &sessionService{
		store: store,
		cache: cache,
	}
}

func (s *sessionService) CreateSession(ctx context.Context, id uuid.UUID, email, refreshToken, userAgent, clientIp string, expiresAt time.Time) (db.Session, error) {
	args := db.CreateSessionParams{
		ID:           id,
		Email:        email,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIp,
		ExpiresAt:    expiresAt,
	}
	return s.store.CreateSession(ctx, args)
}

func (s *sessionService) InvalidateSession(ctx context.Context, sessionID uuid.UUID) (db.Session, error) {
	return s.store.InvalidateSession(ctx, sessionID)
}

func (s *sessionService) GetAllUserSessions(ctx context.Context, email string) ([]db.Session, error) {
	return s.store.GetAllUserSessions(ctx, email)
}

func (s *sessionService) GetSession(ctx context.Context, sessionID uuid.UUID) (db.Session, error) {
	return s.store.GetSession(ctx, sessionID)
}

func (s *sessionService) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.store.DeleteSession(ctx, sessionID)
}
