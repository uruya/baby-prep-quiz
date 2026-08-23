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

const maxWebhookBodyBytes = 64 << 10

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
	sub, err := h.subUC.GetStatus(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "サブスク情報の取得に失敗しました")
		return
	}
	// 契約中に別のSubscriptionを作ると二重請求になるため、Checkoutを開始させない。
	if sub.IsActive() {
		writeError(w, http.StatusConflict, "すでにプレミアムプランを利用中です")
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
		ClientReferenceID: stripe.String(strconv.Itoa(userID)),
	}
	// 再契約時も同じCustomerを使い、Portalから全履歴を管理できるようにする。
	if sub.StripeCustomerID != "" {
		params.Customer = stripe.String(sub.StripeCustomerID)
	} else {
		params.CustomerEmail = stripe.String(user.Email)
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

	// 公開エンドポイントなので、署名検証より前に読み込む本文サイズを制限する。
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxWebhookBodyBytes))
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), h.stripeWebhookSecret)
	if err != nil {
		http.Error(w, "Webhook signature verification failed", http.StatusBadRequest)
		return
	}

	// 永続化エラーは5xxにして、Stripeに同じイベントを再送してもらう。
	if err := h.processWebhookEvent(event); err != nil {
		log.Printf("Stripe webhook processing error for event %s: %v", event.ID, err)
		http.Error(w, "Webhook processing failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BillingHandler) processWebhookEvent(event stripe.Event) error {
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return err
		}
		// 別用途のCheckoutイベントで有料権限を付与しない。
		if string(session.Mode) != string(stripe.CheckoutSessionModeSubscription) {
			return nil
		}
		paymentStatus := string(session.PaymentStatus)
		if paymentStatus != "paid" && paymentStatus != "no_payment_required" {
			return nil
		}
		userID, err := strconv.Atoi(session.ClientReferenceID)
		if err != nil || userID <= 0 || session.Customer == nil || session.Customer.ID == "" {
			return nil
		}
		return h.subUC.ActivatePremium(userID, session.Customer.ID)

	case "customer.subscription.updated":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return err
		}
		if sub.Customer == nil || sub.Customer.ID == "" {
			return nil
		}
		// 一時的なpast_due等はStripeの再試行に任せ、終端状態だけ権限を外す。
		switch string(sub.Status) {
		case "active", "trialing":
			return h.subUC.ActivatePremiumByCustomerID(sub.Customer.ID)
		case "canceled", "unpaid", "incomplete_expired", "paused":
			return h.subUC.DeactivatePremiumByCustomerID(sub.Customer.ID)
		}

	case "customer.subscription.deleted":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return err
		}
		if sub.Customer == nil || sub.Customer.ID == "" {
			return nil
		}
		return h.subUC.DeactivatePremiumByCustomerID(sub.Customer.ID)
	}
	return nil
}
