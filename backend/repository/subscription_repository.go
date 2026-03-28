package repository

import (
	"database/sql"
	"time"

	"baby-prep-quiz/domain"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) GetByUserID(userID int) (*domain.Subscription, error) {
	var sub domain.Subscription
	var stripeCustomerID *string
	err := r.db.QueryRow(
		`SELECT id, subscription_tier, subscription_expires_at, stripe_customer_id FROM users WHERE id = $1`,
		userID,
	).Scan(&sub.UserID, &sub.Tier, &sub.ExpiresAt, &stripeCustomerID)
	if err != nil {
		return nil, err
	}
	if stripeCustomerID != nil {
		sub.StripeCustomerID = *stripeCustomerID
	}
	return &sub, nil
}

func (r *SubscriptionRepository) Upsert(userID int, tier string, expiresAt *time.Time) error {
	_, err := r.db.Exec(
		`UPDATE users SET subscription_tier = $1, subscription_expires_at = $2 WHERE id = $3`,
		tier, expiresAt, userID,
	)
	return err
}

func (r *SubscriptionRepository) ActivatePremium(userID int, stripeCustomerID string) error {
	_, err := r.db.Exec(
		`UPDATE users SET subscription_tier = 'premium', subscription_expires_at = NULL, stripe_customer_id = $1 WHERE id = $2`,
		stripeCustomerID, userID,
	)
	return err
}

func (r *SubscriptionRepository) DeactivatePremiumByCustomerID(stripeCustomerID string) error {
	_, err := r.db.Exec(
		`UPDATE users SET subscription_tier = 'free', subscription_expires_at = NULL WHERE stripe_customer_id = $1`,
		stripeCustomerID,
	)
	return err
}
