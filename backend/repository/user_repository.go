package repository

import (
	"database/sql"
	"strings"
	"time"

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
		 RETURNING id, name, email, created_at, subscription_tier, subscription_expires_at`,
		name, email, passwordHash,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt,
		&user.SubscriptionTier, &user.SubscriptionExpiresAt)
	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, string, error) {
	var user domain.User
	var passwordHash string
	err := r.db.QueryRow(
		`SELECT id, name, email, password_hash, created_at, subscription_tier, subscription_expires_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &passwordHash, &user.CreatedAt,
		&user.SubscriptionTier, &user.SubscriptionExpiresAt)
	if err != nil {
		return nil, "", err
	}
	return &user, passwordHash, nil
}

func (r *UserRepository) FindByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(
		`SELECT id, name, email, created_at, subscription_tier, subscription_expires_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt,
		&user.SubscriptionTier, &user.SubscriptionExpiresAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateSubscription(userID int, tier string, expiresAt *time.Time) error {
	_, err := r.db.Exec(
		`UPDATE users SET subscription_tier = $1, subscription_expires_at = $2 WHERE id = $3`,
		tier, expiresAt, userID,
	)
	return err
}

func IsUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "unique")
}
