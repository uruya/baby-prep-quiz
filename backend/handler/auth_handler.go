package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"baby-prep-quiz/usecase"
)

type AuthHandler struct {
	authUC *usecase.AuthUsecase
}

func NewAuthHandler(authUC *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストが不正です")
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "全ての項目を入力してください")
		return
	}
	if len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "パスワードは6文字以上である必要があります")
		return
	}
	user, err := h.authUC.SignUp(req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrDuplicateEmail) {
			writeError(w, http.StatusConflict, "このメールアドレスは既に登録されています")
		} else {
			writeError(w, http.StatusInternalServerError, "登録に失敗しました")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストが不正です")
		return
	}
	user, token, err := h.authUC.Login(req.Email, req.Password)
	if errors.Is(err, usecase.ErrInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, "メールアドレスまたはパスワードが正しくありません")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "サーバーエラーが発生しました")
		return
	}
	setSessionCookie(w, token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "未ログインです")
		return
	}
	claims, err := h.authUC.ParseToken(cookie.Value)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "セッションが無効です")
		return
	}
	user, err := h.authUC.GetUserByID(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ユーザー情報の取得に失敗しました")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}
