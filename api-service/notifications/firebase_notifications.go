package notifications

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"log"
)

type FirebaseNotifications struct {
	fcm *messaging.Client
}

func getDecodedKey(key string) ([]byte, error) {
	decodedSecretKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		panic(err)
	}

	return decodedSecretKey, nil
}

func NewFirebaseNotifications(encodedSecretKey string) (*FirebaseNotifications, error) {
	key, err := getDecodedKey(encodedSecretKey)
	if err != nil {
		return nil, err
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(key)}

	app, err := firebase.NewApp(context.Background(), nil, opts...)
	if err != nil {
		return nil, err
	}

	fcmClient, err := app.Messaging(context.Background())

	log.Println(app)

	if err != nil {
		return nil, err
	}

	return &FirebaseNotifications{
		fcm: fcmClient,
	}, nil
}

func (client *FirebaseNotifications) SendNotification(title, body string, deviceTokens []string) error {
	_, err := client.fcm.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: deviceTokens, // it's an array of device tokens
	})

	if err != nil {
		return err
	}
	return nil
}
