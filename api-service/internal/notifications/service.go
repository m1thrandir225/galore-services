// Package notifications handles sending notifications to users devices
package notifications

// Service provides a way to send a notification to a users device
type Service interface {
	SendNotification(title, body string, deviceTokens []string) error
}
