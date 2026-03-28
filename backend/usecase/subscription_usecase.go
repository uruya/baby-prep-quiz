package usecase

import (
	"time"

	"baby-prep-quiz/domain"
)

type SubscriptionUsecase struct {
	subRepo domain.SubscriptionRepository
}

func NewSubscriptionUsecase(subRepo domain.SubscriptionRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{subRepo: subRepo}
}

// GetStatus はユーザーの現在のサブスク状態を返す
func (u *SubscriptionUsecase) GetStatus(userID int) (*domain.Subscription, error) {
	return u.subRepo.GetByUserID(userID)
}

// Upgrade はユーザーをプレミアムプランに変更する（将来Stripe連携、今は仮実装）
func (u *SubscriptionUsecase) Upgrade(userID int) (*domain.Subscription, error) {
	// 仮実装: 30日間のプレミアムを付与
	expiresAt := time.Now().AddDate(0, 1, 0)
	if err := u.subRepo.Upsert(userID, domain.TierPremium, &expiresAt); err != nil {
		return nil, err
	}
	return u.subRepo.GetByUserID(userID)
}

// ActivatePremium は Stripe Checkout 完了後にプレミアムを有効化する
func (u *SubscriptionUsecase) ActivatePremium(userID int, stripeCustomerID string) error {
	return u.subRepo.ActivatePremium(userID, stripeCustomerID)
}

// DeactivatePremiumByCustomerID はサブスクキャンセル時にプレミアムを無効化する
func (u *SubscriptionUsecase) DeactivatePremiumByCustomerID(stripeCustomerID string) error {
	return u.subRepo.DeactivatePremiumByCustomerID(stripeCustomerID)
}
