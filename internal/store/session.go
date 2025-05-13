package store

import (
	"context"
	"database/sql"
)

type Session struct {
	Id string `json:"_id"`
	UserId string `json:"user_id"`
	Token string `json:"token"`
	CsrfToken string `json:"csrf_token"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

type SessionStore struct {
	db *sql.DB
}

func (s *SessionStore) Create(ctx context.Context, session Session) (Session, error) {
	query := `
		INSERT INTO session (id, user_id, token, csrf_token, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING token, csrf_token
	`

	var returnedSession Session
	err := s.db.QueryRowContext(ctx, query,
		session.Id,
		session.UserId,
		session.Token,
		session.CsrfToken,
		session.ExpiresAt,
		session.CreatedAt,
	).Scan(&returnedSession.Token, &returnedSession.CsrfToken)

	if err != nil {
		return Session{}, err
	}

	returnedSession.Id = session.Id
	returnedSession.UserId = session.UserId
	returnedSession.ExpiresAt = session.ExpiresAt
	returnedSession.CreatedAt = session.CreatedAt

	return returnedSession, nil
}
