package repository

import (
	"database/sql"
	"encoding/json"

	"baby-prep-quiz/domain"
)

type QuizRepository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) FindByCategory(category string) ([]domain.Question, error) {
	rows, err := r.db.Query(`
		SELECT id, category, question, options, correct_answer, explanation
		FROM questions WHERE category = $1`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var qs []domain.Question
	for rows.Next() {
		var q domain.Question
		var opts []byte
		if err := rows.Scan(&q.ID, &q.Category, &q.Question, &opts, &q.CorrectAnswer, &q.Explanation); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(opts, &q.Options); err != nil {
			return nil, err
		}
		qs = append(qs, q)
	}
	return qs, nil
}

func (r *QuizRepository) SaveResult(userID int, category string, score, total int) error {
	_, err := r.db.Exec(
		`INSERT INTO quiz_results (user_id, category, score, total) VALUES ($1, $2, $3, $4)`,
		userID, category, score, total,
	)
	return err
}

func (r *QuizRepository) GetStats(userID int) (*domain.QuizStats, error) {
	rows, err := r.db.Query(`
		SELECT category, MAX(score) as best_score, MAX(total) as total
		FROM quiz_results
		WHERE user_id = $1
		GROUP BY category`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := &domain.QuizStats{
		Categories: map[string]domain.CategoryStat{},
	}
	for rows.Next() {
		var category string
		var stat domain.CategoryStat
		if err := rows.Scan(&category, &stat.BestScore, &stat.Total); err != nil {
			continue
		}
		stats.Categories[category] = stat
		stats.CompletedQuizzes++
		stats.TotalScore += stat.BestScore
		stats.TotalPossible += stat.Total
	}
	return stats, nil
}
