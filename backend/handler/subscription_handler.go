package handler

import (
	"encoding/json"
	"net/http"

	"baby-prep-quiz/usecase"
)

type SubscriptionHandler struct {
	subUC *usecase.SubscriptionUsecase
}

func NewSubscriptionHandler(subUC *usecase.SubscriptionUsecase) *SubscriptionHandler {
	return &SubscriptionHandler{subUC: subUC}
}

// Status は現在のサブスクリプションプランを返す
func (h *SubscriptionHandler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "未ログインです")
		return
	}
	sub, err := h.subUC.GetStatus(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "サブスク情報の取得に失敗しました")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

// Upgrade はプレミアムプランへのアップグレードを行う（仮実装）
func (h *SubscriptionHandler) Upgrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "未ログインです")
		return
	}
	sub, err := h.subUC.Upgrade(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "アップグレードに失敗しました")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}
