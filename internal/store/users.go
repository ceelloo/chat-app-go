package store

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string   `json:"_id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

type UserStore struct {
	db *sql.DB
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

func (s *UserStore) Create(ctx context.Context, user User) error {
	query := `
		INSERT INTO users (id, name, email, password, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.ExecContext(ctx, query, user.Id, user.Name, user.Email, user.Password.hash, user.CreatedAt)

	return err
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (User, error) {
	query := `
		SELECT id, name, email, password, created_at
		FROM users
		WHERE email = ?
	`

	var user User
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)

	return user, err
}

func (s *UserStore) GetById(ctx context.Context, id string) (User, error) {
	query := `
		SELECT id, name, email, password, created_at
		FROM users
		WHERE id = ?
	`

	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)

	return user, err
}