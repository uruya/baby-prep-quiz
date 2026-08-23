package domain

import "time"

type User struct {
	ID                    int        `json:"id"`
	Name                  string     `json:"name"`
	Email                 string     `json:"email"`
	CreatedAt             string     `json:"createdAt"`
	SubscriptionTier      string     `json:"subscriptionTier"`
	SubscriptionExpiresAt *time.Time `json:"subscriptionExpiresAt"`
}

type UserRepository interface {
	Create(name, email, passwordHash string) (*User, error)
	FindByEmail(email string) (*User, string, error) // User, passwordHash, error
	FindByID(id int) (*User, error)
	UpdateSubscription(userID int, tier string, expiresAt *time.Time) error
}
