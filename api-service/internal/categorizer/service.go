package categorizer

import (
	"github.com/m1thrandir225/galore-services/internal/db/sqlc"
)

type Service interface {
	CategorizeCocktail(cocktail db.Cocktail) error
}
