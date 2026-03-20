package usecase_test

import (
	"database/sql"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"baby-prep-quiz/domain"
	"baby-prep-quiz/usecase"
)

// mockUserRepo は domain.UserRepository のモック実装
type mockUserRepo struct {
	createFn      func(name, email, passwordHash string) (*domain.User, error)
	findByEmailFn func(email string) (*domain.User, string, error)
}

func (m *mockUserRepo) Create(name, email, passwordHash string) (*domain.User, error) {
	return m.createFn(name, email, passwordHash)
}

func (m *mockUserRepo) FindByEmail(email string) (*domain.User, string, error) {
	return m.findByEmailFn(email)
}

func TestSignUp_Success(t *testing.T) {
	repo := &mockUserRepo{
		createFn: func(name, email, passwordHash string) (*domain.User, error) {
			return &domain.User{ID: 1, Name: name, Email: email}, nil
		},
	}
	uc := usecase.NewAuthUsecase(repo, "test-secret")

	user, err := uc.SignUp("田中太郎", "taro@example.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "田中太郎" {
		t.Errorf("expected name 田中太郎, got %s", user.Name)
	}
	if user.Email != "taro@example.com" {
		t.Errorf("expected email taro@example.com, got %s", user.Email)
	}
}

func TestSignUp_DuplicateEmail(t *testing.T) {
	repo := &mockUserRepo{
		createFn: func(name, email, passwordHash string) (*domain.User, error) {
			return nil, errors.New("unique constraint violation")
		},
	}
	uc := usecase.NewAuthUsecase(repo, "test-secret")

	_, err := uc.SignUp("田中太郎", "taro@example.com", "password123")
	if !errors.Is(err, usecase.ErrDuplicateEmail) {
		t.Errorf("expected ErrDuplicateEmail, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, string, error) {
			return &domain.User{ID: 1, Name: "田中太郎", Email: email}, string(hash), nil
		},
	}
	uc := usecase.NewAuthUsecase(repo, "test-secret")

	user, token, err := uc.Login("taro@example.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if token == "" {
		t.Error("expected JWT token, got empty string")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	repo := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, string, error) {
			return &domain.User{ID: 1, Name: "田中太郎", Email: email}, string(hash), nil
		},
	}
	uc := usecase.NewAuthUsecase(repo, "test-secret")

	_, _, err := uc.Login("taro@example.com", "wrongpassword")
	if !errors.Is(err, usecase.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, string, error) {
			return nil, "", sql.ErrNoRows
		},
	}
	uc := usecase.NewAuthUsecase(repo, "test-secret")

	_, _, err := uc.Login("notfound@example.com", "password123")
	if !errors.Is(err, usecase.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestParseToken_Valid(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, string, error) {
			return &domain.User{ID: 42, Name: "田中太郎", Email: email}, string(hash), nil
		},
	}
	uc := usecase.NewAuthUsecase(repo, "test-secret")

	_, token, _ := uc.Login("taro@example.com", "password123")

	claims, err := uc.ParseToken(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("expected userID 42, got %d", claims.UserID)
	}
	if claims.Email != "taro@example.com" {
		t.Errorf("expected email taro@example.com, got %s", claims.Email)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	uc := usecase.NewAuthUsecase(&mockUserRepo{}, "test-secret")

	_, err := uc.ParseToken("invalid.token.string")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockUserRepo{
		findByEmailFn: func(email string) (*domain.User, string, error) {
			return &domain.User{ID: 1, Name: "田中太郎", Email: email}, string(hash), nil
		},
	}
	ucSigner := usecase.NewAuthUsecase(repo, "secret-A")
	ucVerifier := usecase.NewAuthUsecase(repo, "secret-B")

	_, token, _ := ucSigner.Login("taro@example.com", "password123")

	_, err := ucVerifier.ParseToken(token)
	if err == nil {
		t.Error("expected error when verifying with wrong secret, got nil")
	}
}
