package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPNotificationClient implementa comunicação síncrona via HTTP
// com o serviço de notificações
type HTTPNotificationClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPNotificationClient(baseURL string) *HTTPNotificationClient {
	return &HTTPNotificationClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type notificationRequest struct {
	PaymentID int64   `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
}

func (c *HTTPNotificationClient) SendPaymentCreatedNotification(paymentID int64, amount float64) error {
	return c.sendNotification(paymentID, amount, "PAYMENT_CREATED")
}

func (c *HTTPNotificationClient) SendPaymentAuthorizedNotification(paymentID int64, amount float64) error {
	return c.sendNotification(paymentID, amount, "PAYMENT_AUTHORIZED")
}

func (c *HTTPNotificationClient) SendPaymentSettledNotification(paymentID int64, amount float64) error {
	return c.sendNotification(paymentID, amount, "PAYMENT_SETTLED")
}

func (c *HTTPNotificationClient) sendNotification(paymentID int64, amount float64, notificationType string) error {
	reqBody := notificationRequest{
		PaymentID: paymentID,
		Amount:    amount,
		Type:      notificationType,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/notifications", c.baseURL)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// Em produção, isso deveria ser logado e possivelmente retentado
		// Por enquanto, apenas retornamos o erro (eventual consistency)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("notification service returned status %d", resp.StatusCode)
	}

	return nil
}
