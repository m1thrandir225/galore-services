package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m1thrandir225/galore-services/dto"
	"github.com/m1thrandir225/galore-services/util"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"golang.org/x/exp/rand"
)

func (server *Server) registerBackgroundHandlers() {
	server.scheduler.RegisterJob("send_mail", server.sendMailJob)
	/**
	Generate Daily Featured Cocktails Cron
	*/
	server.scheduler.RegisterCronJob("generate_daily_featured", "0 0 * * *") //Cron for every day
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
	server.scheduler.RegisterJob("generate_image", server.createImageGenerationRequest)

	/**
	Generate Cocktail Job
	*/
	//TODO: implement generate cocktail background job
	server.scheduler.RegisterJob("generate_cocktail_draft", server.createCocktailDraft)

	server.scheduler.RegisterJob("generate_cocktail_final", server.createGeneratedCocktail)

}

func (server *Server) createCocktailDraft(args map[string]interface{}) error {
	log.Println("JOB STARTED: Creating Generate Cocktail Draft")

	cocktailRequestIdStr, ok := args["cocktail_request_id"].(string)
	if !ok {
		return fmt.Errorf("missing required arguments: cocktail_request_id")
	}
	cocktailRequestId, err := uuid.Parse(cocktailRequestIdStr)
	if err != nil {
		return err
	}

	prompt, ok := args["prompt"].(string)
	if !ok {
		return fmt.Errorf("missing required arguments: prompt")
	}

	userIdStr, ok := args["user_id"].(string)
	if !ok {
		return fmt.Errorf("missing required arguments: user_id")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	promptCocktail, err := server.cocktailGenerator.GenerateRecipe(prompt)
	if err != nil {
		return err
	}
	jsonInstructions, err := json.Marshal(promptCocktail.Instructions)
	if err != nil {
		return err
	}

	jsonIngredients, err := json.Marshal(promptCocktail.Ingredients)
	if err != nil {
		return err
	}

	createArgs := db.CreateGenerateCocktailDraftParams{
		RequestID:       cocktailRequestId,
		Name:            promptCocktail.Name,
		Description:     promptCocktail.Description,
		Ingredients:     jsonIngredients,
		Instructions:    jsonInstructions,
		MainImagePrompt: promptCocktail.ImagePrompt,
	}

	draft, err := server.store.CreateGenerateCocktailDraft(context.Background(), createArgs)
	if err != nil {
		return err
	}

	updateGenerateRequestArgs := db.UpdateGenerateCocktailRequestParams{
		ID:     draft.RequestID,
		Status: db.GenerationStatusGeneratingImages,
	}

	_, err = server.store.UpdateGenerateCocktailRequest(context.Background(), updateGenerateRequestArgs)
	if err != nil {
		log.Println("ERR: error while updating generate request: " + err.Error())
		return err
	}

	for _, promptInstruction := range promptCocktail.Instructions {
		server.scheduler.EnqueueJob("generate_image", map[string]interface{}{
			"cocktail_draft_id": draft.ID,
			"image_prompt":      promptInstruction.ImagePrompt,
			"is_main":           false,
			"user_id":           userId,
		})
	}

	server.scheduler.EnqueueJob("generate_image", map[string]interface{}{
		"cocktail_draft_id": draft.ID,
		"image_prompt":      draft.MainImagePrompt,
		"is_main":           true,
		"user_id":           userId,
	})

	return nil
}

func (server *Server) createImageGenerationRequest(args map[string]interface{}) error {
	log.Println("JOB STARTED: Creating Images for Cocktail Draft")
	cocktailDraftIdStr, ok := args["cocktail_draft_id"].(string)
	if !ok {
		return fmt.Errorf("missing required arguments: cocktail_draft_id")
	}
	cocktailDraftId, err := uuid.Parse(cocktailDraftIdStr)
	if err != nil {
		return err
	}

	isMain, ok := args["is_main"].(bool)
	if !ok {
		log.Println("missing required arguments: is_main")
		return fmt.Errorf("missing required arguments: is_main")
	}

	imagePrompt, ok := args["image_prompt"].(string)
	if !ok {
		log.Println("missing required arguments: image_prompt")
		return fmt.Errorf("missing required arguments: image_prompt")
	}

	userIdStr, ok := args["user_id"].(string)
	if !ok {
		log.Println("ERROR: missing required arguments: user_id")
		return fmt.Errorf("missing required arguments: user_id")
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	createArgs := db.CreateImageGenerationRequestParams{
		DraftID: cocktailDraftId,
		Prompt:  imagePrompt,
		Status:  db.ImageGenerationStatusGenerating,
		IsMain:  isMain,
	}

	imageRequest, err := server.store.CreateImageGenerationRequest(context.Background(), createArgs)
	if err != nil {
		log.Println("ERR: error while creating image generation request: " + err.Error())
		return err
	}

	//TODO: add image generation logic

	//After the image is generated logic (if there is an image we update
	updateArgs := db.UpdateImageGenerationRequestParams{
		ID: imageRequest.ID,
		ImageUrl: pgtype.Text{
			String: "image_url",
			Valid:  true,
		},
		ErrorMessage: pgtype.Text{
			String: "",
			Valid:  false,
		},
		Status: db.ImageGenerationStatusSuccess,
	}
	_, err = server.store.UpdateImageGenerationRequest(context.Background(), updateArgs)
	if err != nil {
		log.Println("ERR: error while updating image generation request: " + err.Error())
		return err
	}
	draft, err := server.store.GetCocktailDraft(context.Background(), cocktailDraftId)
	if err != nil {
		return err
	}
	//Check if it's the last image that it's being generated from the draft request
	data, err := server.store.CheckImageGenerationProgress(context.Background(), draft.RequestID)
	if err != nil {
		log.Println("ERR: error while checking image generation: " + err.Error())
		return err
	}

	if !data.AllSuccessful {
		log.Println("ERROR: image generation progress check failed")
		return errors.New("there was a problem generating images")
	}

	var instructions []dto.PromptInstruction
	err = json.Unmarshal(draft.Instructions, &instructions)
	if err != nil {
		log.Println("ERR: error while unmarshalling instructions: " + err.Error())
		return err
	}
	if data.CompletedImages == int64(len(instructions)+1) {
		server.scheduler.EnqueueJob("generate_cocktail_final", map[string]interface{}{
			"cocktail_draft_id": cocktailDraftId,
			"user_id":           userId,
		})
	}
	return nil
}

func (server *Server) createGeneratedCocktail(args map[string]interface{}) error {
	cocktailDraftIdStr, ok := args["cocktail_draft_id"].(string)
	if !ok {
		log.Println("ERROR: missing required arguments - [cocktail_draft_id]")
		return errors.New("missing required arguments: cocktail_draft_id")
	}
	cocktailDraftId, err := uuid.Parse(cocktailDraftIdStr)
	if err != nil {
		return err
	}

	userIdStr, ok := args["user_id"].(string)
	if !ok {
		log.Println("ERROR: missing required arguments - [user_id]")
		return errors.New("missing required arguments: user_id")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	draft, err := server.store.GetCocktailDraft(context.Background(), cocktailDraftId)
	if err != nil {
		log.Println(err)
		return err
	}
	imagesForDraft, err := server.store.GetImagesForDraft(context.Background(), draft.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	var mainImageUrl string
	for _, img := range imagesForDraft {
		if img.IsMain {
			mainImageUrl = img.ImageUrl.String
			break

		}
	}

	var promptInstructions []dto.PromptInstruction
	err = json.Unmarshal([]byte(draft.Instructions), &promptInstructions)
	if err != nil {
		log.Println(err)
		return err
	}

	aiInstructions, err := util.ConvertPromptsToAiInstructionDto(promptInstructions, imagesForDraft)
	if err != nil {
		log.Println(err)
		return err
	}

	var ingredientData []dto.IngredientData
	err = json.Unmarshal([]byte(draft.Ingredients), &ingredientData)

	createFinalCocktailArgs := db.CreateGeneratedCocktailParams{
		Name:         draft.Name,
		UserID:       userId,
		RequestID:    draft.RequestID,
		DraftID:      draft.ID,
		Instructions: aiInstructions,
		Ingredients: dto.IngredientDto{
			Ingredients: ingredientData,
		},
		Description:  draft.Description,
		MainImageUrl: mainImageUrl,
	}

	finalCocktail, err := server.store.CreateGeneratedCocktail(context.Background(), createFinalCocktailArgs)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(finalCocktail)
	//TODO: SEND NOTIFICATION

	updateGenerateRequestArgs := db.UpdateGenerateCocktailRequestParams{
		ID:     draft.RequestID,
		Status: db.GenerationStatusSuccess,
	}

	_, err = server.store.UpdateGenerateCocktailRequest(context.Background(), updateGenerateRequestArgs)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (server *Server) sendMailJob(args map[string]interface{}) error {
	log.Println("JOB STARTED: Sending mail")

	email, ok := args["email"].(string)

	if !ok {
		return fmt.Errorf("missing required arguments")
	}

	fmt.Println("Started: sending mail job")

	mailTemplate, ok := args["mail_template"].(string)
	if !ok {
		return fmt.Errorf("missing required arguments: mail template")
	}

	subject, ok := args["subject"].(string)
	if !ok {
		return fmt.Errorf("missing required arguments: subject")
	}

	err := server.mailService.SendMail(
		"galore@sebastijanzindl.me",
		email,
		subject,
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
