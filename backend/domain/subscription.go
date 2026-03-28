package domain

import "time"

const (
	TierFree    = "free"
	TierPremium = "premium"
)

// freeユーザーがアクセス可能なカテゴリ
var FreeCategories = map[string]bool{
	"pregnancy": true,
	"birth":     true,
}

type Subscription struct {
	UserID           int        `json:"userId"`
	Tier             string     `json:"tier"`
	ExpiresAt        *time.Time `json:"expiresAt"`
	StripeCustomerID string     `json:"stripeCustomerId,omitempty"`
}

func (s *Subscription) IsActive() bool {
	if s.Tier != TierPremium {
		return false
	}
	if s.ExpiresAt != nil && s.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

type SubscriptionRepository interface {
	GetByUserID(userID int) (*Subscription, error)
	Upsert(userID int, tier string, expiresAt *time.Time) error
	ActivatePremium(userID int, stripeCustomerID string) error
	DeactivatePremiumByCustomerID(stripeCustomerID string) error
}
