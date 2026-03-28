package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stripe/stripe-go/v76"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"baby-prep-quiz/config"
	"baby-prep-quiz/handler"
	"baby-prep-quiz/repository"
	"baby-prep-quiz/usecase"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
	}
	log.Println("Database successfully connected")

	// マイグレーション
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not create postgres driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("could not run migrate up: %v", err)
	}
	log.Println("Database migration completed!")

	// Stripe初期化
	stripe.Key = cfg.StripeSecretKey

	// DI
	userRepo := repository.NewUserRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret)
	quizUC := usecase.NewQuizUsecase(quizRepo)
	subUC := usecase.NewSubscriptionUsecase(subRepo)
	authH := handler.NewAuthHandler(authUC)
	quizH := handler.NewQuizHandler(quizUC, authUC)
	subH := handler.NewSubscriptionHandler(subUC)
	billingH := handler.NewBillingHandler(subUC, authUC, cfg.StripePriceID, cfg.StripeWebhookSecret, cfg.FrontendURL)

	// CORS
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", cfg.FrontendURL)
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

	// ルーティング
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/signup", authH.SignUp)
	mux.HandleFunc("/api/auth/login", authH.Login)
	mux.HandleFunc("/api/auth/me", authH.Me)
	mux.HandleFunc("/api/auth/logout", authH.Logout)
	mux.HandleFunc("/api/quiz/results", handler.AuthMiddleware(authUC, quizH.SaveResult))
	mux.HandleFunc("/api/quiz/stats", handler.AuthMiddleware(authUC, quizH.GetStats))
	mux.Handle("/api/quiz/", http.HandlerFunc(quizH.GetByCategory))
	mux.HandleFunc("/api/subscription/status", handler.AuthMiddleware(authUC, subH.Status))
	mux.HandleFunc("/api/subscription/upgrade", handler.AuthMiddleware(authUC, subH.Upgrade))
	mux.HandleFunc("/api/billing/checkout", handler.AuthMiddleware(authUC, billingH.Checkout))
	mux.HandleFunc("/api/billing/portal", handler.AuthMiddleware(authUC, billingH.Portal))
	mux.HandleFunc("/api/billing/webhook", billingH.Webhook)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", corsHandler(mux)); err != nil {
		log.Fatal(err)
	}
}
