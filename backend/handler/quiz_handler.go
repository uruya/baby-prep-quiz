package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"baby-prep-quiz/domain"
	"baby-prep-quiz/usecase"
)

type QuizHandler struct {
	quizUC *usecase.QuizUsecase
	authUC *usecase.AuthUsecase
}

func NewQuizHandler(quizUC *usecase.QuizUsecase, authUC *usecase.AuthUsecase) *QuizHandler {
	return &QuizHandler{quizUC: quizUC, authUC: authUC}
}

func (h *QuizHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	category := strings.TrimPrefix(r.URL.Path, "/api/quiz/")

	// freeカテゴリ以外はプレミアム判定
	if !domain.FreeCategories[category] {
		cookie, err := r.Cookie("session")
		if err != nil {
			writePremiumRequired(w)
			return
		}
		claims, err := h.authUC.ParseToken(cookie.Value)
		if err != nil {
			writePremiumRequired(w)
			return
		}
		user, err := h.authUC.GetUserByID(claims.UserID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "ユーザー情報の取得に失敗しました")
			return
		}
		sub := &domain.Subscription{
			UserID:    user.ID,
			Tier:      user.SubscriptionTier,
			ExpiresAt: user.SubscriptionExpiresAt,
		}
		if !sub.IsActive() {
			writePremiumRequired(w)
			return
		}
	}

	questions, err := h.quizUC.GetQuestions(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func (h *QuizHandler) SaveResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "未ログインです")
		return
	}
	var req struct {
		Category string `json:"category"`
		Score    int    `json:"score"`
		Total    int    `json:"total"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストが不正です")
		return
	}
	if err := h.quizUC.SaveResult(userID, req.Category, req.Score, req.Total); err != nil {
		writeError(w, http.StatusInternalServerError, "保存に失敗しました")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *QuizHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := getUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "未ログインです")
		return
	}
	stats, err := h.quizUC.GetStats(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "取得に失敗しました")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func writePremiumRequired(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{"error": "premium_required"})
}
