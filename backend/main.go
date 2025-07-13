package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Question struct {
	ID int   `json:"id"`
	Category string `json:"category"`
	Question string `json:"question"`
	Options []string `json:"options"`
	CorrectAnswer int `json:"correctAnswer"`
	Explanation string `json:"explanation"`
}

func main() {
	dbURL := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@localhost:5432/" + os.Getenv("POSTGRES_DB") + "?sslmode=disable"

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// 疎通確認
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}
	log.Println("Database successfully connected")

	// CORS設定
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				return
			}
			h.ServeHTTP(w, r)
		})
	}

	mux := http.NewServeMux()
	mux.Handle("/api/quiz/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// URLから/api/quiz/を取り除いてカテゴリー名を取得
		category := strings.TrimPrefix(r.URL.Path, "/api/quiz/")

		rows, err := db.Query(`SELECT id, category, question, options, correct_answer, explanation
			FROM questions WHERE category = $1`, category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var qs []Question
		for rows.Next() {
			var q Question
			var opts []byte
			if err := rows.Scan(&q.ID, &q.Category, &q.Question, &opts, &q.CorrectAnswer, &q.Explanation); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := json.Unmarshal(opts, &q.Options); err != nil {
				http.Error(w, "Failed to parse options", http.StatusInternalServerError)
				return
			}
			qs = append(qs, q)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(qs)
	}))

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", corsHandler(mux)); err != nil {
		log.Fatal(err)
	}

}
