package application

import "fintech-notifications-service/domain"

type CreateNotificationUseCase struct {
	repo domain.NotificationRepository
}

func NewCreateNotificationUseCase(repo domain.NotificationRepository) *CreateNotificationUseCase {
	return &CreateNotificationUseCase{repo: repo}
}

func (uc *CreateNotificationUseCase) Execute(paymentID int64, notificationType, recipient, message string) (*domain.Notification, error) {
	notification := domain.NewNotification(paymentID, notificationType, recipient, message)

	// Salvar no banco próprio do serviço
	saved, err := uc.repo.Save(notification)
	if err != nil {
		return nil, err
	}

	// Simular envio (em produção, isso seria um worker assíncrono)
	saved.MarkAsSent()
	_, _ = uc.repo.Save(saved)

	return saved, nil
}
