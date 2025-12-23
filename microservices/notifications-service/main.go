package main

import (
	"context"
	"fintech-notifications-service/api"
	app "fintech-notifications-service/application"
	"fintech-notifications-service/infra/persistence"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Banco de dados próprio do serviço (autonomia de dados)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Repositório usando banco próprio
	notificationRepo := persistence.NewPgNotificationRepository(pool)

	// Use case
	createUC := app.NewCreateNotificationUseCase(notificationRepo)

	handler := api.NewNotificationsHandler(createUC, notificationRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"notifications","type":"microservice"}`))
	})
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Notifications Service listening on :" + port)
	log.Println("NOTE: This is a microservice with its own database")
	log.Fatal(srv.ListenAndServe())
}
