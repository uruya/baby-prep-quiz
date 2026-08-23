package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"baby-prep-quiz/domain"
	"baby-prep-quiz/handler"
	"baby-prep-quiz/usecase"
)

// mockUserRepoForHandler は handler テスト用のモックリポジトリ
type mockUserRepoForHandler struct {
	createFn             func(name, email, passwordHash string) (*domain.User, error)
	findByEmailFn        func(email string) (*domain.User, string, error)
	findByIDFn           func(id int) (*domain.User, error)
	updateSubscriptionFn func(userID int, tier string, expiresAt *time.Time) error
}

func (m *mockUserRepoForHandler) Create(name, email, passwordHash string) (*domain.User, error) {
	return m.createFn(name, email, passwordHash)
}

func (m *mockUserRepoForHandler) FindByEmail(email string) (*domain.User, string, error) {
	return m.findByEmailFn(email)
}

func (m *mockUserRepoForHandler) FindByID(id int) (*domain.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(id)
	}
	return &domain.User{ID: id, SubscriptionTier: "free"}, nil
}

func (m *mockUserRepoForHandler) UpdateSubscription(userID int, tier string, expiresAt *time.Time) error {
	if m.updateSubscriptionFn != nil {
		return m.updateSubscriptionFn(userID, tier, expiresAt)
	}
	return nil
}

func newTestAuthHandler(repo domain.UserRepository) *handler.AuthHandler {
	uc := usecase.NewAuthUsecase(repo, "test-secret")
	return handler.NewAuthHandler(uc)
}

func TestSignUpHandler_Success(t *testing.T) {
	repo := &mockUserRepoForHandler{
		createFn: func(name, email, passwordHash string) (*domain.User, error) {
			return &domain.User{ID: 1, Name: name, Email: email}, nil
		},
	}
	h := newTestAuthHandler(repo)

	body := `{"name":"田中太郎","email":"taro@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.SignUp(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
	var user domain.User
	if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if user.Name != "田中太郎" {
		t.Errorf("expected name 田中太郎, got %s", user.Name)
	}
}

func TestSignUpHandler_MethodNotAllowed(t *testing.T) {
	h := newTestAuthHandler(&mockUserRepoForHandler{})

	req := httptest.NewRequest(http.MethodGet, "/api/auth/signup", nil)
	w := httptest.NewRecorder()

	h.SignUp(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestSignUpHandler_MissingFields(t *testing.T) {
	h := newTestAuthHandler(&mockUserRepoForHandler{})

	body := `{"name":"","email":"taro@example.com","password":"pass123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.SignUp(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSignUpHandler_ShortPassword(t *testing.T) {
	h := newTestAuthHandler(&mockUserRepoForHandler{})

	body := `{"name":"太郎","email":"taro@example.com","password":"abc"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.SignUp(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSignUpHandler_DuplicateEmail(t *testing.T) {
	repo := &mockUserRepoForHandler{
		createFn: func(name, email, passwordHash string) (*domain.User, error) {
			return nil, sql.ErrNoRows // unique制約エラーをシミュレート
		},
	}
	// unique制約エラーをusecase.ErrDuplicateEmailに変換するためにerror文字列を合わせる
	repo.createFn = func(name, email, passwordHash string) (*domain.User, error) {
		return nil, &duplicateError{}
	}
	h := newTestAuthHandler(repo)

	body := `{"name":"太郎","email":"taro@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.SignUp(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", w.Code)
	}
}

// duplicateError は unique制約違反エラーをシミュレートする
type duplicateError struct{}

func (e *duplicateError) Error() string { return "unique constraint violation" }

func TestLoginHandler_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockUserRepoForHandler{
		findByEmailFn: func(email string) (*domain.User, string, error) {
			return &domain.User{ID: 1, Name: "田中太郎", Email: email}, string(hash), nil
		},
	}
	h := newTestAuthHandler(repo)

	body := `{"email":"taro@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	// Cookieにsessionトークンがセットされているか確認
	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "session" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected session cookie to be set")
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	repo := &mockUserRepoForHandler{
		findByEmailFn: func(email string) (*domain.User, string, error) {
			return &domain.User{ID: 1, Name: "田中太郎", Email: email}, string(hash), nil
		},
	}
	h := newTestAuthHandler(repo)

	body := `{"email":"taro@example.com","password":"wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestMeHandler_NoCookie(t *testing.T) {
	h := newTestAuthHandler(&mockUserRepoForHandler{})

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	w := httptest.NewRecorder()

	h.Me(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestMeHandler_ValidSession(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockUserRepoForHandler{
		findByEmailFn: func(email string) (*domain.User, string, error) {
			return &domain.User{ID: 1, Name: "田中太郎", Email: email}, string(hash), nil
		},
		findByIDFn: func(id int) (*domain.User, error) {
			return &domain.User{ID: id, Name: "田中太郎", Email: "taro@example.com", SubscriptionTier: "free"}, nil
		},
	}
	uc := usecase.NewAuthUsecase(repo, "test-secret")
	h := handler.NewAuthHandler(uc)

	// ログインしてトークンを取得
	_, token, _ := uc.Login("taro@example.com", "password123")

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: token})
	w := httptest.NewRecorder()

	h.Me(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	var user domain.User
	if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if user.Email != "taro@example.com" {
		t.Errorf("expected email taro@example.com, got %s", user.Email)
	}
}

func TestLogoutHandler_Success(t *testing.T) {
	h := newTestAuthHandler(&mockUserRepoForHandler{})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
	// sessionCookieがMaxAge=-1でクリアされているか確認
	cookies := w.Result().Cookies()
	for _, c := range cookies {
		if c.Name == "session" && c.MaxAge != -1 {
			t.Errorf("expected session cookie MaxAge=-1, got %d", c.MaxAge)
		}
	}
}
