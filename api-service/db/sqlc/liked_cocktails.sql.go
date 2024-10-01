// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: liked_cocktails.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	dto "github.com/m1thrandir225/galore-services/dto"
	"github.com/pgvector/pgvector-go"
)

const getLikedCocktail = `-- name: GetLikedCocktail :one
SELECT c.id, name, is_alcoholic, glass, image, instructions, ingredients, embedding, created_at, lc.id, cocktail_id, user_id from cocktails c
JOIN liked_cocktails lc ON c.id = lc.cocktail_id
WHERE lc.user_id = $1 and lc.cocktail_id = $2
GROUP BY lc.user_id, lc.cocktail_id, c.id
`

type GetLikedCocktailParams struct {
	UserID     uuid.UUID `json:"user_id"`
	CocktailID uuid.UUID `json:"cocktail_id"`
}

type GetLikedCocktailRow struct {
	ID           uuid.UUID         `json:"id"`
	Name         string            `json:"name"`
	IsAlcoholic  pgtype.Bool       `json:"is_alcoholic"`
	Glass        string            `json:"glass"`
	Image        string            `json:"image"`
	Instructions string            `json:"instructions"`
	Ingredients  dto.IngredientDto `json:"ingredients"`
	Embedding    pgvector.Vector   `json:"embedding"`
	CreatedAt    time.Time         `json:"created_at"`
	ID_2         uuid.UUID         `json:"id_2"`
	CocktailID   uuid.UUID         `json:"cocktail_id"`
	UserID       uuid.UUID         `json:"user_id"`
}

func (q *Queries) GetLikedCocktail(ctx context.Context, arg GetLikedCocktailParams) (GetLikedCocktailRow, error) {
	row := q.db.QueryRow(ctx, getLikedCocktail, arg.UserID, arg.CocktailID)
	var i GetLikedCocktailRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IsAlcoholic,
		&i.Glass,
		&i.Image,
		&i.Instructions,
		&i.Ingredients,
		&i.Embedding,
		&i.CreatedAt,
		&i.ID_2,
		&i.CocktailID,
		&i.UserID,
	)
	return i, err
}

const getLikedCocktails = `-- name: GetLikedCocktails :many
SELECt c.id, name, is_alcoholic, glass, image, instructions, ingredients, embedding, created_at, lc.id, cocktail_id, user_id from cocktails c
JOIN  liked_cocktails lc ON c.id = lc.cocktail_id
WHERE lc.user_id = $1
GROUP BY lc.user_id
`

type GetLikedCocktailsRow struct {
	ID           uuid.UUID         `json:"id"`
	Name         string            `json:"name"`
	IsAlcoholic  pgtype.Bool       `json:"is_alcoholic"`
	Glass        string            `json:"glass"`
	Image        string            `json:"image"`
	Instructions string            `json:"instructions"`
	Ingredients  dto.IngredientDto `json:"ingredients"`
	Embedding    pgvector.Vector   `json:"embedding"`
	CreatedAt    time.Time         `json:"created_at"`
	ID_2         uuid.UUID         `json:"id_2"`
	CocktailID   uuid.UUID         `json:"cocktail_id"`
	UserID       uuid.UUID         `json:"user_id"`
}

func (q *Queries) GetLikedCocktails(ctx context.Context, userID uuid.UUID) ([]GetLikedCocktailsRow, error) {
	rows, err := q.db.Query(ctx, getLikedCocktails, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetLikedCocktailsRow{}
	for rows.Next() {
		var i GetLikedCocktailsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IsAlcoholic,
			&i.Glass,
			&i.Image,
			&i.Instructions,
			&i.Ingredients,
			&i.Embedding,
			&i.CreatedAt,
			&i.ID_2,
			&i.CocktailID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const likeCocktail = `-- name: LikeCocktail :one
INSERT INTO liked_cocktails (
  cocktail_id,
  user_id
) VALUES (
  $1,
  $2
) RETURNING id, cocktail_id, user_id
`

type LikeCocktailParams struct {
	CocktailID uuid.UUID `json:"cocktail_id"`
	UserID     uuid.UUID `json:"user_id"`
}

func (q *Queries) LikeCocktail(ctx context.Context, arg LikeCocktailParams) (LikedCocktail, error) {
	row := q.db.QueryRow(ctx, likeCocktail, arg.CocktailID, arg.UserID)
	var i LikedCocktail
	err := row.Scan(&i.ID, &i.CocktailID, &i.UserID)
	return i, err
}

const unlikeCocktail = `-- name: UnlikeCocktail :exec
DELETE FROM liked_cocktails 
WHERE cocktail_id = $1 AND user_id = $2
`

type UnlikeCocktailParams struct {
	CocktailID uuid.UUID `json:"cocktail_id"`
	UserID     uuid.UUID `json:"user_id"`
}

func (q *Queries) UnlikeCocktail(ctx context.Context, arg UnlikeCocktailParams) error {
	_, err := q.db.Exec(ctx, unlikeCocktail, arg.CocktailID, arg.UserID)
	return err
}
