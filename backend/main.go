package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Question struct {
	ID int   `json:"id"`
	Category string `json:"category"`
	Question string `json:"question"`
	Options []string `json:"options"`
	CorrectAnswer string `json:"correctAnswer"`
	Explanation string `json:"explanation"`
}

func main() {
	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" +os.Getenv("DB_NAME") + "?sslmode=disable"
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/quiz/", func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Path[len("/quiz/"):]
		rows, err := db.Query(`SELECT id, category, question, options, correct_answer, explanation
			FROM questions WHERE category = $1`,
			category)
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
			qs = append(qs, q)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(qs)
	})
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
