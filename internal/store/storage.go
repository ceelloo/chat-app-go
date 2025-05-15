package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Users interface {
		Create (ctx context.Context, user User) error
		GetByEmail (ctx context.Context, email string) (User, error)
		GetById (ctx context.Context, id string) (User, error)
	}
	Sessions interface {
		Create (ctx context.Context, session Session) (Session, error)
		Get (ctx context.Context, token string) (Session, error)
	}
}

func NewStorage (db *sql.DB) Storage {
	return Storage{
		Users: &UserStore{db},
		Sessions: &SessionStore{db},
	}
}