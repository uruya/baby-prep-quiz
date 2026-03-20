package repository

import (
	"database/sql"
	"strings"

	"baby-prep-quiz/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(name, email, passwordHash string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(
		`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)
		 RETURNING id, name, email, created_at`,
		name, email, passwordHash,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, string, error) {
	var user domain.User
	var passwordHash string
	err := r.db.QueryRow(
		`SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &passwordHash, &user.CreatedAt)
	if err != nil {
		return nil, "", err
	}
	return &user, passwordHash, nil
}

func IsUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "unique")
}
