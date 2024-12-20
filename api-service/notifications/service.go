package notifications

type NotificationService interface {
	SendNotification(title, body string, deviceTokens []string) error
}
