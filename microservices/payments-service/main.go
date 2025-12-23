package main

import (
	"context"
	"fintech-payments-service/api"
	app "fintech-payments-service/application"
	"fintech-payments-service/infra/messaging/pix"
	"fintech-payments-service/infra/notifications"
	"fintech-payments-service/infra/persistence"
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

	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationServiceURL == "" {
		notificationServiceURL = "http://notifications-service:8080"
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
	paymentRepo := persistence.NewPgPixPaymentRepository(pool)

	// Cliente HTTP para comunicação com serviço de notificações
	notificationClient := notifications.NewHTTPNotificationClient(notificationServiceURL)

	// Gateway do BACEN (simulação)
	gateway := pix.NewBacenPixGateway()

	// Event broadcaster para observabilidade em tempo real
	eventBroadcaster := api.GetBroadcaster()

	// Use case que usa o cliente de notificações, gateway e event broadcaster
	createUC := app.NewCreatePixPaymentUseCase(paymentRepo, notificationClient, gateway, eventBroadcaster)

	handler := api.NewPaymentsHandler(createUC, paymentRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"payments","type":"microservice"}`))
	})
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Payments Service listening on :" + port)
	log.Println("Notification Service URL:", notificationServiceURL)
	log.Println("NOTE: This is a microservice with its own database")
	log.Fatal(srv.ListenAndServe())
}
