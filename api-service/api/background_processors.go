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
	server.scheduler.RegisterJob("send_mail", false, server.sendMailJob)
	/**
	Generate Daily Featured Cocktails Cron
	*/
	server.scheduler.RegisterCronJob("generate_daily_featured", "0 0 * * *") //Cron for every day
	server.scheduler.RegisterJob("generate_daily_featured", true, server.generateDailyFeatured)

	/**
	Cocktail Migration Cron
	*/
	server.scheduler.RegisterCronJob("migrate_cocktails", "0 0 1 * *")
	server.scheduler.RegisterJob("migrate_cocktails", true, server.migrateCocktails)

	/**
	Send Notification Job
	*/
	server.scheduler.RegisterJob("send_notification", false, server.sendNotification)

	/**
	Generate Image Job
	*/
	//TODO: implement generate image background job
	server.scheduler.RegisterJob("generate_image", true, server.createImageGenerationRequest)

	/**
	Generate Cocktail Job
	*/
	//TODO: implement generate cocktail background job
	server.scheduler.RegisterJob("generate_cocktail_draft", true, server.createCocktailDraft)

	server.scheduler.RegisterJob("generate_cocktail_final", true, server.createGeneratedCocktail)

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
			"cocktail_draft_id":   draft.ID,
			"image_prompt":        promptInstruction.ImagePrompt,
			"is_main":             false,
			"user_id":             userId,
			"generate_request_id": draft.RequestID,
		})
	}

	server.scheduler.EnqueueJob("generate_image", map[string]interface{}{
		"cocktail_draft_id":   draft.ID,
		"image_prompt":        draft.MainImagePrompt,
		"is_main":             true,
		"user_id":             userId,
		"generate_request_id": draft.RequestID,
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

	generateRequestIdStr, ok := args["generate_request_id"].(string)
	if !ok {
		return fmt.Errorf("missing required arguments: generate_request_id")
	}
	generateRequestId, err := uuid.Parse(generateRequestIdStr)
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

	generateRequest, err := server.store.GetGenerationRequest(context.Background(), generateRequestId)
	if err != nil {
		return err
	}

	/**
	Check if the request has a status error (i.e one of the previous images had an error),
	if it does then just cancel the current image generation process
	*/
	if generateRequest.Status == db.GenerationStatusError {
		createArgs := db.CreateImageGenerationRequestParams{
			DraftID: cocktailDraftId,
			Prompt:  imagePrompt,
			Status:  db.ImageGenerationStatusCancelled,
			IsMain:  isMain,
		}

		_, imageReqErr := server.store.CreateImageGenerationRequest(context.Background(), createArgs)
		if imageReqErr != nil {
			log.Println("ERR: error while creating image generation request: " + err.Error())
			return imageReqErr
		}
		return nil
	}

	/*
		If the parent request for the draft doesn't have an error continue with the generation process
	*/
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
	httpClient := &http.Client{}

	/**
	Generate
	*/
	generatedImageData, err := server.imageGenerator.GenerateImage(imageRequest.Prompt, httpClient, "core")
	if err != nil {
		errorUpdateArgs := db.UpdateImageGenerationRequestParams{
			ID: imageRequest.ID,
			ImageUrl: pgtype.Text{
				String: "",
				Valid:  false,
			},
			ErrorMessage: pgtype.Text{
				String: err.Error(),
				Valid:  true,
			},
			Status: db.ImageGenerationStatusError,
		}
		_, updateImageReqErr := server.store.UpdateImageGenerationRequest(context.Background(), errorUpdateArgs)
		if updateImageReqErr != nil {
			log.Println("ERR: error while updating image generation request: " + err.Error())
			return err
		}
	}
	/**
	Upload the generated image to the current storage implementation
	If it errors then update the current image request
	*/

	uploadedFilePath, err := server.storage.UploadFile(
		generatedImageData.Content,
		"generated-images",
		fmt.Sprintf("%s%s", generatedImageData.FileName, generatedImageData.FileExt),
	)
	if err != nil {
		errorUpdateArgs := db.UpdateImageGenerationRequestParams{
			ID: imageRequest.ID,
			ImageUrl: pgtype.Text{
				String: "",
				Valid:  false,
			},
			ErrorMessage: pgtype.Text{
				String: err.Error(),
				Valid:  true,
			},
			Status: db.ImageGenerationStatusError,
		}
		_, updateImageReqErr := server.store.UpdateImageGenerationRequest(context.Background(), errorUpdateArgs)
		if updateImageReqErr != nil {
			log.Println("ERR: error while updating image generation request: " + err.Error())
			return err
		}
	}

	successUpdateArgs := db.UpdateImageGenerationRequestParams{
		ID: imageRequest.ID,
		ImageUrl: pgtype.Text{
			String: uploadedFilePath,
			Valid:  true,
		},
		ErrorMessage: pgtype.Text{
			String: "",
			Valid:  false,
		},
		Status: db.ImageGenerationStatusSuccess,
	}
	_, err = server.store.UpdateImageGenerationRequest(context.Background(), successUpdateArgs)

	if err != nil {
		log.Println("ERR: error while updating image generation request: " + err.Error())
		return err
	}

	/*
		Fetch the draft for the image to check if this is the last one that needs to be generated
	*/
	draft, err := server.store.GetCocktailDraft(context.Background(), cocktailDraftId)
	if err != nil {
		return err
	}

	/*
		Get the current status of all the image requests
	*/
	data, err := server.store.CheckImageGenerationProgress(context.Background(), draft.RequestID)
	if err != nil {
		log.Println("ERR: error while checking image generation: " + err.Error())
		return err
	}

	/*
		If there is one that is not successful, update the parent generate_request to status failed so the following ones can be cancelled
	*/
	if !data.AllSuccessful {
		log.Println("ERROR: not all images completed the generation process successfully")

		updateGenerateRequestArgs := db.UpdateGenerateCocktailRequestParams{
			ID:     generateRequestId,
			Status: db.GenerationStatusError,
		}
		_, err = server.store.UpdateGenerateCocktailRequest(context.Background(), updateGenerateRequestArgs)

		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	}
	/*
		Logic for comparing if this image_generation request is the last one, if it is then schedule the generation of the final cocktail
	*/

	var instructions []dto.PromptInstruction
	err = json.Unmarshal(draft.Instructions, &instructions)
	if err != nil {
		log.Println("ERR: error while unmarshalling instructions: " + err.Error())
		return err
	}

	if data.TotalImages == data.CompletedImages && data.CompletedImages == int64(len(instructions)+1) {
		log.Println("INFO: all images completed the generation process successfully")
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
