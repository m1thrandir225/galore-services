package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/m1thrandir225/galore-services/background_jobs"
	"github.com/m1thrandir225/galore-services/cache"
	categorizer "github.com/m1thrandir225/galore-services/categorizer_service"
	embedding "github.com/m1thrandir225/galore-services/embedding_service"
	"github.com/m1thrandir225/galore-services/mail"
	"github.com/m1thrandir225/galore-services/storage"
	"golang.org/x/exp/rand"

	"github.com/gin-gonic/gin"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/token"
	"github.com/m1thrandir225/galore-services/util"
)

type Server struct {
	config      util.Config
	store       db.Store
	router      *gin.Engine
	tokenMaker  token.Maker
	storage     storage.FileService
	cache       cache.KeyValueStore
	embedding   embedding.EmbeddingService
	categorizer categorizer.CategorizerService
	scheduler   background_jobs.SchedulerService
	mailService mail.MailService
}

type HealthResponse struct {
	Status string `json:"status"`
}

func NewServer(
	config util.Config,
	store db.Store,
	storageService storage.FileService,
	cacheStore cache.KeyValueStore,
	embedding embedding.EmbeddingService,
	categorizer categorizer.CategorizerService,
	scheduler background_jobs.SchedulerService,
	mailService mail.MailService,
) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:      config,
		store:       store,
		tokenMaker:  tokenMaker,
		storage:     storageService,
		cache:       cacheStore,
		embedding:   embedding,
		categorizer: categorizer,
		scheduler:   scheduler,
		mailService: mailService,
	}

	server.setupRouter()
	server.registerBackgroundHandlers()

	return server, nil
}

func (server *Server) Start(address string) error {
	go server.scheduler.Start()
	defer server.scheduler.Stop()

	return server.router.Run(address)
}

func (server *Server) checkService(ctx *gin.Context) {
	message := HealthResponse{
		Status: "health",
	}
	ctx.JSON(http.StatusOK, message)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) registerBackgroundHandlers() {
	server.scheduler.RegisterJob("send_mail", server.sendMailJob)
	server.scheduler.RegisterCronJob("generate_daily_featured", "0 * * * * *")
	server.scheduler.RegisterJob("generate_daily_featured", server.generateDailyFeatured)
}

func (server *Server) sendMailJob(args map[string]interface{}) error {
	email, ok := args["email"].(string)

	if !ok {
		return fmt.Errorf("Missing required argument: email")
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
		return fmt.Errorf("There was an error sending the email to %s, reason: %s", email, err.Error())
	}

	return nil
}
func (server *Server) generateDailyFeatured(args map[string]interface{}) error {
	numberOfFeatured := 10

	todaysFeatured, err := server.store.GetDailyFeatured(context.Background())
	if err != nil {
		return fmt.Errorf("There was an error getting the today's featured cocktails: %s", err.Error())
	}

	if len(todaysFeatured) >= numberOfFeatured {
		log.Println("Today's featured already generated")
		return nil
	}

	allCocktails, err := server.store.SearchCocktails(context.Background(), "")
	if err != nil {
		return fmt.Errorf("There was a problem with getting all the cocktails from the databse: ", err.Error())
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
		_, err := server.store.CreateDailyFeatured(context.Background(), cocktail.ID)
		if err != nil {
			return fmt.Errorf("failed to mark cocktail %s as featured: %v", cocktail.ID, err)
		}
	}
	return nil
}
