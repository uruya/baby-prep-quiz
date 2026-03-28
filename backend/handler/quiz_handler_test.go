package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"baby-prep-quiz/domain"
	"baby-prep-quiz/usecase"
)

// mockQuizRepoForHandler は handler テスト用のモッククイズリポジトリ
type mockQuizRepoForHandler struct {
	findByCategoryFn func(category string) ([]domain.Question, error)
	saveResultFn     func(userID int, category string, score, total int) error
	getStatsFn       func(userID int) (*domain.QuizStats, error)
}

func (m *mockQuizRepoForHandler) FindByCategory(category string) ([]domain.Question, error) {
	return m.findByCategoryFn(category)
}

func (m *mockQuizRepoForHandler) SaveResult(userID int, category string, score, total int) error {
	return m.saveResultFn(userID, category, score, total)
}

func (m *mockQuizRepoForHandler) GetStats(userID int) (*domain.QuizStats, error) {
	return m.getStatsFn(userID)
}

func newTestQuizHandler(repo domain.QuizRepository) *QuizHandler {
	uc := usecase.NewQuizUsecase(repo)
	return NewQuizHandler(uc, nil)
}

func TestGetByCategoryHandler_Success(t *testing.T) {
	questions := []domain.Question{
		{ID: 1, Category: "pregnancy", Question: "テスト問題", Options: []string{"A", "B", "C", "D"}, CorrectAnswer: 0},
	}
	repo := &mockQuizRepoForHandler{
		findByCategoryFn: func(category string) ([]domain.Question, error) {
			return questions, nil
		},
	}
	h := newTestQuizHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/quiz/pregnancy", nil)
	w := httptest.NewRecorder()

	h.GetByCategory(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	var got []domain.Question
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 question, got %d", len(got))
	}
}

func TestSaveResultHandler_Unauthorized(t *testing.T) {
	h := newTestQuizHandler(&mockQuizRepoForHandler{})

	body := `{"category":"pregnancy","score":8,"total":10}`
	req := httptest.NewRequest(http.MethodPost, "/api/quiz/results", bytes.NewBufferString(body))
	// contextにuserIDをセットしない（未認証状態）
	w := httptest.NewRecorder()

	h.SaveResult(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestSaveResultHandler_Success(t *testing.T) {
	saved := false
	repo := &mockQuizRepoForHandler{
		saveResultFn: func(userID int, category string, score, total int) error {
			saved = true
			if userID != 1 {
				t.Errorf("expected userID 1, got %d", userID)
			}
			return nil
		},
	}
	h := newTestQuizHandler(repo)

	body := `{"category":"pregnancy","score":8,"total":10}`
	req := httptest.NewRequest(http.MethodPost, "/api/quiz/results", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	// AuthMiddlewareが行うcontextへのuserIDセットを再現
	ctx := context.WithValue(req.Context(), userIDKey, 1)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.SaveResult(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
	if !saved {
		t.Error("expected SaveResult to be called")
	}
}

func TestSaveResultHandler_MethodNotAllowed(t *testing.T) {
	h := newTestQuizHandler(&mockQuizRepoForHandler{})

	req := httptest.NewRequest(http.MethodGet, "/api/quiz/results", nil)
	ctx := context.WithValue(req.Context(), userIDKey, 1)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.SaveResult(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestGetStatsHandler_Success(t *testing.T) {
	stats := &domain.QuizStats{
		CompletedQuizzes: 3,
		TotalScore:       25,
		TotalPossible:    30,
		Categories: map[string]domain.CategoryStat{
			"pregnancy": {BestScore: 9, Total: 10},
		},
	}
	repo := &mockQuizRepoForHandler{
		getStatsFn: func(userID int) (*domain.QuizStats, error) {
			return stats, nil
		},
	}
	h := newTestQuizHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/quiz/stats", nil)
	ctx := context.WithValue(req.Context(), userIDKey, 1)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.GetStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	var got domain.QuizStats
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.CompletedQuizzes != 3 {
		t.Errorf("expected completedQuizzes 3, got %d", got.CompletedQuizzes)
	}
}

func TestGetStatsHandler_Unauthorized(t *testing.T) {
	h := newTestQuizHandler(&mockQuizRepoForHandler{})

	req := httptest.NewRequest(http.MethodGet, "/api/quiz/stats", nil)
	// contextにuserIDなし
	w := httptest.NewRecorder()

	h.GetStats(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}
