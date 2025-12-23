package notifications

type NotificationRepository interface {
	Save(notification *Notification) (*Notification, error)
	FindByID(id int64) (*Notification, error)
}
