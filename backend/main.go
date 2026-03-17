package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type Question struct {
	ID            int      `json:"id"`
	Category      string   `json:"category"`
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectAnswer int      `json:"correctAnswer"`
	Explanation   string   `json:"explanation"`
}

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type claims struct {
	UserID int    `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found, using environment variables")
	}

	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.dbname")
	dbSSLMode := viper.GetString("database.sslmode")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	jwtSecret := viper.GetString("jwt.secret")
	frontendURL := strings.TrimRight(viper.GetString("app.frontend_url"), "/")

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}
	log.Println("Database successfully connected")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("could not run migrate up: %v", err)
	}
	log.Println("Database migration completed!")

	// CORS設定（cookieを使うためcredentials対応）
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", frontendURL)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			h.ServeHTTP(w, r)
		})
	}

	mux := http.NewServeMux()

	// POST /api/auth/signup
	mux.HandleFunc("/api/auth/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req signupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Message: "リクエストが不正です"})
			return
		}
		if req.Name == "" || req.Email == "" || req.Password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Message: "全ての項目を入力してください"})
			return
		}
		if len(req.Password) < 6 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Message: "パスワードは6文字以上である必要があります"})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{Message: "サーバーエラーが発生しました"})
			return
		}
		var user User
		err = db.QueryRow(
			`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)
			 RETURNING id, name, email, created_at`,
			req.Name, req.Email, string(hash),
		).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			if strings.Contains(err.Error(), "unique") {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(errorResponse{Message: "このメールアドレスは既に登録されています"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(errorResponse{Message: "登録に失敗しました"})
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})

	// POST /api/auth/login
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Message: "リクエストが不正です"})
			return
		}
		var user User
		var passwordHash string
		err := db.QueryRow(
			`SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`,
			req.Email,
		).Scan(&user.ID, &user.Name, &user.Email, &passwordHash, &user.CreatedAt)
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorResponse{Message: "メールアドレスまたはパスワードが正しくありません"})
			return
		}
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{Message: "サーバーエラーが発生しました"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorResponse{Message: "メールアドレスまたはパスワードが正しくありません"})
			return
		}

		// JWT生成
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
			UserID: user.ID,
			Name:   user.Name,
			Email:  user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		})
		tokenStr, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{Message: "サーバーエラーが発生しました"})
			return
		}

		// httpOnly cookieにセット
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    tokenStr,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode, // クロスオリジンでcookieを送信するために必要
			Path:     "/",
			MaxAge:   60 * 60 * 24, // 24時間
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	// GET /api/auth/me
	mux.HandleFunc("/api/auth/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		cookie, err := r.Cookie("session")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorResponse{Message: "未ログインです"})
			return
		}
		c := &claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, c, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorResponse{Message: "セッションが無効です"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{
			ID:    c.UserID,
			Name:  c.Name,
			Email: c.Email,
		})
	})

	// POST /api/auth/logout
	mux.HandleFunc("/api/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Path:     "/",
			MaxAge:   -1,
		})
		w.WriteHeader(http.StatusNoContent)
	})

	mux.Handle("/api/quiz/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
