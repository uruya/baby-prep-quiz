package usecase

import "baby-prep-quiz/domain"

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

// ActivatePremium は Stripe Checkout 完了後にプレミアムを有効化する
func (u *SubscriptionUsecase) ActivatePremium(userID int, stripeCustomerID string) error {
	return u.subRepo.ActivatePremium(userID, stripeCustomerID)
}

// ActivatePremiumByCustomerID は既知のStripeカスタマーのプレミアム権限を有効化する
func (u *SubscriptionUsecase) ActivatePremiumByCustomerID(stripeCustomerID string) error {
	return u.subRepo.ActivatePremiumByCustomerID(stripeCustomerID)
}

// DeactivatePremiumByCustomerID はStripe上の契約が無効になったときにプレミアムを無効化する
func (u *SubscriptionUsecase) DeactivatePremiumByCustomerID(stripeCustomerID string) error {
	return u.subRepo.DeactivatePremiumByCustomerID(stripeCustomerID)
}
