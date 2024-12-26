package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/mail"
	"golang.org/x/exp/rand"
)

func (server *Server) registerBackgroundHandlers() {
	server.scheduler.RegisterJob("send_mail", server.sendMailJob)
	/**
	Generate Daily Featured Cocktails Cron
	*/
	server.scheduler.RegisterCronJob("generate_daily_featured", "0 0 1 * *") //Cron for every day
	server.scheduler.RegisterJob("generate_daily_featured", server.generateDailyFeatured)

	/**
	Cocktail Migration Cron
	*/
	server.scheduler.RegisterCronJob("migrate_cocktails", "0 0 1 * *")
	server.scheduler.RegisterJob("migrate_cocktails", server.migrateCocktails)

	/**
	Send Notification Job
	*/
	server.scheduler.RegisterJob("send_notification", server.sendNotification)

	/**
	Generate Image Job
	*/
	//TODO: implement generate image background job

	/**
	Generate Cocktail Job
	*/
	//TODO: implement generate cocktail background job

}

func (server *Server) sendMailJob(args map[string]interface{}) error {
	log.Println("JOB STARTED: Sending mail")

	email, ok := args["email"].(string)

	if !ok {
		return fmt.Errorf("missing required arguments")
	}

	fmt.Println("Started: sending mail job")

	mailTemplate := mail.GenerateWelcomeMail(email)

	err := server.mailService.SendMail(
		"help@sebastijanzindl.me",
		email,
		"Welcome to Galore",
		mailTemplate,
	)
	if err != nil {
		return fmt.Errorf("there was an error sending the email to %s, reason: %s", email, err.Error())
	}
	return nil
}

func (server *Server) generateDailyFeatured(args map[string]interface{}) error {
	log.Println("JOB STARTED: Generate daily featured")
	numberOfFeatured := 10

	todaysFeatured, err := server.store.GetDailyFeatured(context.Background())
	if err != nil {
		return fmt.Errorf("there was an error getting the today's featured cocktails: %s", err.Error())
	}

	if len(todaysFeatured) >= numberOfFeatured {
		log.Println("Today's featured already generated")
		return nil
	}

	allCocktails, err := server.store.SearchCocktails(context.Background(), "")
	if err != nil {
		return fmt.Errorf("there was a problem with getting all the cocktails from the databse: %s", err.Error())
	}

	rand.Shuffle(len(allCocktails), func(i, j int) {
		allCocktails[i], allCocktails[j] = allCocktails[j], allCocktails[i]
	})

	var featuredCocktails []db.Cocktail

	for i := 0; i < numberOfFeatured; i++ {
		cocktail := allCocktails[i]
		_, err := server.store.CheckWasCocktailFeatured(context.Background(), cocktail.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				featuredCocktails = append(featuredCocktails, cocktail)
			} else {
				return fmt.Errorf("failed to check if cocktail %s was featured: %v", cocktail.ID, err)
			}
		} else {
			continue
		}
	}

	if len(featuredCocktails) < numberOfFeatured {
		return fmt.Errorf("unable to find enough unfeatured cocktails")
	}

	for _, cocktail := range featuredCocktails {
		_, err = server.store.CreateDailyFeatured(context.Background(), cocktail.ID)
		if err != nil {
			return fmt.Errorf("failed to mark cocktail %s as featured: %v", cocktail.ID, err)
		}
	}
	return nil
}

func (server *Server) migrateCocktails(args map[string]interface{}) error {
	log.Println("JOB STARTED: Migrate cocktails")
	req, err := http.NewRequest(
		"GET",
		server.config.MigrationServiceAddress,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", server.config.MigrationServiceKey)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 OK, got %d", response.StatusCode)
	}

	defer response.Body.Close()

	return nil
}

func (server *Server) createNotificationJob(args map[string]interface{}) error {
	log.Println("JOB STARTED: Create notification job")
	notificationTypeId, ok := args["notification_type_id"].(uuid.UUID)
	if !ok {
		return fmt.Errorf("missing required arguments: notification_type_id")
	}
	userId, ok := args["user_id"].(uuid.UUID)
	if !ok {
		return fmt.Errorf("missing required arguments: user_id")
	}

	notifArgs := db.CreateNotificationParams{
		UserID:             userId,
		NotificationTypeID: notificationTypeId,
	}

	notification, err := server.store.CreateNotification(context.Background(), notifArgs)
	if err != nil {
		return err
	}

	err = server.scheduler.EnqueueJob("send_notification", map[string]interface{}{
		"notification_type_id": notificationTypeId,
		"user_id":              notification.UserID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) sendNotification(args map[string]interface{}) error {
	notificationTypeId, ok := args["notification_type_id"].(uuid.UUID)
	if !ok {
		return fmt.Errorf("missing required arguments: notification_type")
	}
	userId, ok := args["user_id"].(uuid.UUID)
	if !ok {
		return fmt.Errorf("missing required arguments: user_id")
	}

	notificationType, err := server.store.GetNotificationType(context.Background(), notificationTypeId)
	if err != nil {
		return err
	}

	//Get the users FCM tokens
	fcmTokens, err := server.store.GetUserFCMTokens(context.Background(), userId)
	if err != nil {
		return err
	}
	var tokens []string
	for _, item := range fcmTokens {
		tokens = append(tokens, item.Token)
	}

	err = server.notificationService.SendNotification(notificationType.Title, notificationType.Content, tokens)
	if err != nil {
		return err
	}

	return nil
}
