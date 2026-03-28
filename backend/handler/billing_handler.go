package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/stripe/stripe-go/v76"
	portalsession "github.com/stripe/stripe-go/v76/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"

	"baby-prep-quiz/usecase"
)

type BillingHandler struct {
	subUC               *usecase.SubscriptionUsecase
	authUC              *usecase.AuthUsecase
	stripePriceID       string
	stripeWebhookSecret string
	frontendURL         string
}

func NewBillingHandler(subUC *usecase.SubscriptionUsecase, authUC *usecase.AuthUsecase, priceID, webhookSecret, frontendURL string) *BillingHandler {
	return &BillingHandler{
		subUC:               subUC,
		authUC:              authUC,
		stripePriceID:       priceID,
		stripeWebhookSecret: webhookSecret,
		frontendURL:         frontendURL,
	}
}

// Checkout は Stripe Checkout セッションを作成してURLを返す
func (h *BillingHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "未ログインです")
		return
	}

	user, err := h.authUC.GetUserByID(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ユーザー情報の取得に失敗しました")
		return
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(h.stripePriceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL:        stripe.String(h.frontendURL + "/profile?upgraded=true"),
		CancelURL:         stripe.String(h.frontendURL + "/pricing"),
		CustomerEmail:     stripe.String(user.Email),
		ClientReferenceID: stripe.String(strconv.Itoa(userID)),
	}

	s, err := checkoutsession.New(params)
	if err != nil {
		log.Printf("Stripe Checkout error: %v", err)
		writeError(w, http.StatusInternalServerError, "Checkoutセッションの作成に失敗しました")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": s.URL})
}

// Portal は Stripe Customer Portal セッションを作成してURLを返す
func (h *BillingHandler) Portal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "未ログインです")
		return
	}

	sub, err := h.subUC.GetStatus(userID)
	if err != nil || sub.StripeCustomerID == "" {
		writeError(w, http.StatusBadRequest, "Stripeカスタマー情報が見つかりません")
		return
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(sub.StripeCustomerID),
		ReturnURL: stripe.String(h.frontendURL + "/profile"),
	}

	s, err := portalsession.New(params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Portalセッションの作成に失敗しました")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": s.URL})
}

// Webhook は Stripe からのイベントを処理する
func (h *BillingHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), h.stripeWebhookSecret)
	if err != nil {
		http.Error(w, "Webhook signature verification failed", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var cs stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &cs); err != nil {
			http.Error(w, "Failed to parse event", http.StatusBadRequest)
			return
		}
		userID, err := strconv.Atoi(cs.ClientReferenceID)
		if err != nil || userID == 0 {
			break
		}
		if cs.Customer != nil {
			h.subUC.ActivatePremium(userID, cs.Customer.ID)
		}
	case "customer.subscription.deleted":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			http.Error(w, "Failed to parse event", http.StatusBadRequest)
			return
		}
		if sub.Customer != nil {
			h.subUC.DeactivatePremiumByCustomerID(sub.Customer.ID)
		}
	}

	w.WriteHeader(http.StatusOK)
}
