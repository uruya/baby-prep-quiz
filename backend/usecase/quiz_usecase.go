package usecase

import "baby-prep-quiz/domain"

type QuizUsecase struct {
	quizRepo domain.QuizRepository
}

func NewQuizUsecase(quizRepo domain.QuizRepository) *QuizUsecase {
	return &QuizUsecase{quizRepo: quizRepo}
}

func (u *QuizUsecase) GetQuestions(category string) ([]domain.Question, error) {
	return u.quizRepo.FindByCategory(category)
}

func (u *QuizUsecase) SaveResult(userID int, category string, score, total int) error {
	return u.quizRepo.SaveResult(userID, category, score, total)
}

func (u *QuizUsecase) GetStats(userID int) (*domain.QuizStats, error) {
	return u.quizRepo.GetStats(userID)
}
