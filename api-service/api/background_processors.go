package api

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	db "github.com/m1thrandir225/galore-services/db/sqlc"
	"github.com/m1thrandir225/galore-services/mail"
	"golang.org/x/exp/rand"
)

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
