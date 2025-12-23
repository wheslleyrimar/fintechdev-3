package api

import (
	"encoding/json"
	"fintech-notifications-service/application"
	"fintech-notifications-service/domain"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type NotificationsHandler struct {
	createUC *application.CreateNotificationUseCase
	repo     domain.NotificationRepository
}

type createNotificationRequest struct {
	PaymentID int64   `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
}

func NewNotificationsHandler(createUC *application.CreateNotificationUseCase, repo domain.NotificationRepository) *NotificationsHandler {
	return &NotificationsHandler{createUC: createUC, repo: repo}
}

func (h *NotificationsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/notifications", h.handleNotifications)
	mux.HandleFunc("/notifications/", h.handleNotificationByID)
}

func (h *NotificationsHandler) handleNotifications(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *NotificationsHandler) listAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: Listing all notifications")

	notificationsList, err := h.repo.FindAll()
	if err != nil {
		log.Printf("ERROR: Failed to list notifications: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if notificationsList == nil {
		notificationsList = []*domain.Notification{} // Retorna array vazio ao invés de null
	}

	log.Printf("INFO: Found %d notifications", len(notificationsList))
	writeJSON(w, http.StatusOK, notificationsList)
}

func (h *NotificationsHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Criar mensagem baseada no tipo
	message := "Payment created successfully"
	switch req.Type {
	case "PAYMENT_CREATED":
		message = "Your payment has been created"
	case "PAYMENT_AUTHORIZED":
		message = "Your payment has been authorized"
	case "PAYMENT_SETTLED":
		message = "Your payment has been settled"
	}

	log.Printf("INFO: Creating notification - Type: %s, PaymentID: %d", req.Type, req.PaymentID)

	notification, err := h.createUC.Execute(req.PaymentID, req.Type, "user@example.com", message)
	if err != nil {
		log.Printf("ERROR: Failed to create notification: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Notification created successfully - ID: %d, Type: %s", notification.ID, notification.Type)
	writeJSON(w, http.StatusCreated, notification)
}

func (h *NotificationsHandler) handleNotificationByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL: /notifications/{id}
	path := strings.TrimPrefix(r.URL.Path, "/notifications/")
	if path == "" {
		http.Error(w, "notification ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		log.Printf("ERROR: Invalid notification ID: %s", path)
		http.Error(w, "invalid notification ID", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Fetching notification with ID: %d", id)

	notification, err := h.repo.FindByID(id)
	if err != nil {
		log.Printf("ERROR: Failed to find notification %d: %v", id, err)
		http.Error(w, "notification not found", http.StatusNotFound)
		return
	}

	log.Printf("INFO: Notification found - ID: %d, Type: %s", notification.ID, notification.Type)
	writeJSON(w, http.StatusOK, notification)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
