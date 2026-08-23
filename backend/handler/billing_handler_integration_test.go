package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stripe/stripe-go/v76"

	"baby-prep-quiz/domain"
	"baby-prep-quiz/usecase"
)

type webhookSubscriptionRepo struct {
	subscriptions map[int]*domain.Subscription
	updateErr     error
}

func (r *webhookSubscriptionRepo) GetByUserID(userID int) (*domain.Subscription, error) {
	sub, ok := r.subscriptions[userID]
	if !ok {
		return nil, errors.New("subscription not found")
	}
	copy := *sub
	return &copy, nil
}

func (r *webhookSubscriptionRepo) ActivatePremium(userID int, customerID string) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	sub, ok := r.subscriptions[userID]
	if !ok {
		sub = &domain.Subscription{UserID: userID}
		r.subscriptions[userID] = sub
	}
	sub.Tier = domain.TierPremium
	sub.ExpiresAt = nil
	sub.StripeCustomerID = customerID
	return nil
}

func (r *webhookSubscriptionRepo) ActivatePremiumByCustomerID(customerID string) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	for _, sub := range r.subscriptions {
		if sub.StripeCustomerID == customerID {
			sub.Tier = domain.TierPremium
			sub.ExpiresAt = nil
		}
	}
	return nil
}

func (r *webhookSubscriptionRepo) DeactivatePremiumByCustomerID(customerID string) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	for _, sub := range r.subscriptions {
		if sub.StripeCustomerID == customerID {
			sub.Tier = domain.TierFree
			sub.ExpiresAt = nil
		}
	}
	return nil
}

// signedWebhookRequest はStripeが送る t=...,v1=... 形式の署名を生成する。
// これにより、テストでも本番と同じConstructEventの署名検証を通す。
func signedWebhookRequest(t *testing.T, secret, payload string) *http.Request {
	t.Helper()
	timestamp := time.Now().Unix()
	signedPayload := fmt.Sprintf("%d.%s", timestamp, payload)
	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(signedPayload)); err != nil {
		t.Fatalf("failed to sign webhook payload: %v", err)
	}
	signature := hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest(http.MethodPost, "/api/billing/webhook", bytes.NewBufferString(payload))
	req.Header.Set("Stripe-Signature", fmt.Sprintf("t=%d,v1=%s", timestamp, signature))
	return req
}

func stripeEventPayload(eventID, eventType, object string) string {
	return fmt.Sprintf(
		`{"id":%q,"object":"event","api_version":%q,"type":%q,"data":{"object":%s}}`,
		eventID,
		stripe.APIVersion,
		eventType,
		object,
	)
}

func TestBillingWebhookIntegration_CheckoutActivatesPremium(t *testing.T) {
	const secret = "whsec_test"
	repo := &webhookSubscriptionRepo{subscriptions: map[int]*domain.Subscription{
		42: {UserID: 42, Tier: domain.TierFree},
	}}
	h := NewBillingHandler(usecase.NewSubscriptionUsecase(repo), nil, "", secret, "")

	payload := stripeEventPayload(
		"evt_checkout",
		"checkout.session.completed",
		`{"id":"cs_test","object":"checkout.session","mode":"subscription","payment_status":"paid","client_reference_id":"42","customer":"cus_123"}`,
	)
	recorder := httptest.NewRecorder()
	h.Webhook(recorder, signedWebhookRequest(t, secret, payload))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	got := repo.subscriptions[42]
	if got.Tier != domain.TierPremium || got.StripeCustomerID != "cus_123" {
		t.Fatalf("expected premium subscription for cus_123, got %+v", got)
	}
}

func TestBillingWebhookIntegration_SubscriptionStateChanges(t *testing.T) {
	const secret = "whsec_test"
	tests := []struct {
		name        string
		eventType   string
		status      string
		initialTier string
		wantTier    string
	}{
		{
			name:        "updated active restores access",
			eventType:   "customer.subscription.updated",
			status:      "active",
			initialTier: domain.TierFree,
			wantTier:    domain.TierPremium,
		},
		{
			name:        "updated unpaid revokes access",
			eventType:   "customer.subscription.updated",
			status:      "unpaid",
			initialTier: domain.TierPremium,
			wantTier:    domain.TierFree,
		},
		{
			name:        "deleted revokes access",
			eventType:   "customer.subscription.deleted",
			status:      "canceled",
			initialTier: domain.TierPremium,
			wantTier:    domain.TierFree,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &webhookSubscriptionRepo{subscriptions: map[int]*domain.Subscription{
				42: {
					UserID:           42,
					Tier:             tt.initialTier,
					StripeCustomerID: "cus_123",
				},
			}}
			h := NewBillingHandler(usecase.NewSubscriptionUsecase(repo), nil, "", secret, "")
			object := fmt.Sprintf(
				`{"id":"sub_123","object":"subscription","customer":"cus_123","status":%q}`,
				tt.status,
			)
			payload := stripeEventPayload("evt_subscription", tt.eventType, object)
			recorder := httptest.NewRecorder()

			h.Webhook(recorder, signedWebhookRequest(t, secret, payload))

			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
			}
			if got := repo.subscriptions[42].Tier; got != tt.wantTier {
				t.Fatalf("expected tier %q, got %q", tt.wantTier, got)
			}
		})
	}
}

func TestBillingWebhookIntegration_RejectsInvalidSignature(t *testing.T) {
	const secret = "whsec_test"
	repo := &webhookSubscriptionRepo{subscriptions: map[int]*domain.Subscription{}}
	h := NewBillingHandler(usecase.NewSubscriptionUsecase(repo), nil, "", secret, "")
	payload := stripeEventPayload("evt_invalid", "checkout.session.completed", `{}`)
	req := signedWebhookRequest(t, "wrong_secret", payload)
	recorder := httptest.NewRecorder()

	h.Webhook(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestBillingWebhookIntegration_ReturnsServerErrorWhenPersistenceFails(t *testing.T) {
	const secret = "whsec_test"
	repo := &webhookSubscriptionRepo{
		subscriptions: map[int]*domain.Subscription{
			42: {UserID: 42, Tier: domain.TierFree},
		},
		updateErr: errors.New("database unavailable"),
	}
	h := NewBillingHandler(usecase.NewSubscriptionUsecase(repo), nil, "", secret, "")
	payload := stripeEventPayload(
		"evt_retry",
		"checkout.session.completed",
		`{"id":"cs_test","object":"checkout.session","mode":"subscription","payment_status":"paid","client_reference_id":"42","customer":"cus_123"}`,
	)
	recorder := httptest.NewRecorder()

	h.Webhook(recorder, signedWebhookRequest(t, secret, payload))

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 so Stripe retries the event, got %d", recorder.Code)
	}
	if got := repo.subscriptions[42].Tier; got != domain.TierFree {
		t.Fatalf("expected failed update to leave tier free, got %q", got)
	}
}
