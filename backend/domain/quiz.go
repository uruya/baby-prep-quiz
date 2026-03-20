package domain

type Question struct {
	ID            int      `json:"id"`
	Category      string   `json:"category"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectAnswer int      `json:"correctAnswer"`
	Explanation   string   `json:"explanation"`
}

type CategoryStat struct {
	BestScore int `json:"bestScore"`
	Total     int `json:"total"`
}

type QuizStats struct {
	CompletedQuizzes int                     `json:"completedQuizzes"`
	TotalScore       int                     `json:"totalScore"`
	TotalPossible    int                     `json:"totalPossible"`
	Categories       map[string]CategoryStat `json:"categories"`
}

type QuizRepository interface {
	FindByCategory(category string) ([]Question, error)
	SaveResult(userID int, category string, score, total int) error
	GetStats(userID int) (*QuizStats, error)
}
