package categorizer

import db "github.com/m1thrandir225/galore-services/db/sqlc"

type CategorizerService interface {
	CategorizeCocktail(cocktail db.Cocktail) error
}
