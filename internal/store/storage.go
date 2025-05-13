package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Users interface {
		Create (ctx context.Context, user User) error
	}
	Sessions interface {
		Create (ctx context.Context, session Session) (Session, error)
	}
}

func NewStorage (db *sql.DB) Storage {
	return Storage{
		Users: &UserStore{db},
		Sessions: &SessionStore{db},
	}
}