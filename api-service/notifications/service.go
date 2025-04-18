package notifications

// Description:
// Interface for a notification service used to for sending notifications to
// user's mobile devices
type NotificationService interface {
	//Description:
	//Method for sending the notification to the users device(s)
	//
	//Parameters:
	//title: string
	//body: string
	//deviceTokens: []string
	//
	//Return:
	//error
	SendNotification(title, body string, deviceTokens []string) error
}
