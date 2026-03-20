package usecase_test

import (
	"errors"
	"testing"

	"baby-prep-quiz/domain"
	"baby-prep-quiz/usecase"
)

// mockQuizRepo は domain.QuizRepository のモック実装
type mockQuizRepo struct {
	findByCategoryFn func(category string) ([]domain.Question, error)
	saveResultFn     func(userID int, category string, score, total int) error
	getStatsFn       func(userID int) (*domain.QuizStats, error)
}

func (m *mockQuizRepo) FindByCategory(category string) ([]domain.Question, error) {
	return m.findByCategoryFn(category)
}

func (m *mockQuizRepo) SaveResult(userID int, category string, score, total int) error {
	return m.saveResultFn(userID, category, score, total)
}

func (m *mockQuizRepo) GetStats(userID int) (*domain.QuizStats, error) {
	return m.getStatsFn(userID)
}

func TestGetQuestions_Success(t *testing.T) {
	want := []domain.Question{
		{ID: 1, Category: "pregnancy", Question: "テスト問題", Options: []string{"A", "B"}, CorrectAnswer: 0},
	}
	repo := &mockQuizRepo{
		findByCategoryFn: func(category string) ([]domain.Question, error) {
			if category != "pregnancy" {
				t.Errorf("expected category pregnancy, got %s", category)
			}
			return want, nil
		},
	}
	uc := usecase.NewQuizUsecase(repo)

	got, err := uc.GetQuestions("pregnancy")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != len(want) {
		t.Errorf("expected %d questions, got %d", len(want), len(got))
	}
}

func TestGetQuestions_RepositoryError(t *testing.T) {
	repo := &mockQuizRepo{
		findByCategoryFn: func(category string) ([]domain.Question, error) {
			return nil, errors.New("db error")
		},
	}
	uc := usecase.NewQuizUsecase(repo)

	_, err := uc.GetQuestions("pregnancy")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSaveResult_Success(t *testing.T) {
	called := false
	repo := &mockQuizRepo{
		saveResultFn: func(userID int, category string, score, total int) error {
			called = true
			if userID != 1 || category != "birth" || score != 8 || total != 10 {
				t.Errorf("unexpected args: userID=%d, category=%s, score=%d, total=%d", userID, category, score, total)
			}
			return nil
		},
	}
	uc := usecase.NewQuizUsecase(repo)

	err := uc.SaveResult(1, "birth", 8, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Error("expected SaveResult to be called")
	}
}

func TestSaveResult_RepositoryError(t *testing.T) {
	repo := &mockQuizRepo{
		saveResultFn: func(userID int, category string, score, total int) error {
			return errors.New("db error")
		},
	}
	uc := usecase.NewQuizUsecase(repo)

	err := uc.SaveResult(1, "birth", 8, 10)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetStats_Success(t *testing.T) {
	want := &domain.QuizStats{
		CompletedQuizzes: 2,
		TotalScore:       15,
		TotalPossible:    20,
		Categories: map[string]domain.CategoryStat{
			"pregnancy": {BestScore: 8, Total: 10},
			"birth":     {BestScore: 7, Total: 10},
		},
	}
	repo := &mockQuizRepo{
		getStatsFn: func(userID int) (*domain.QuizStats, error) {
			return want, nil
		},
	}
	uc := usecase.NewQuizUsecase(repo)

	got, err := uc.GetStats(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.CompletedQuizzes != want.CompletedQuizzes {
		t.Errorf("expected completedQuizzes %d, got %d", want.CompletedQuizzes, got.CompletedQuizzes)
	}
	if got.TotalScore != want.TotalScore {
		t.Errorf("expected totalScore %d, got %d", want.TotalScore, got.TotalScore)
	}
}
