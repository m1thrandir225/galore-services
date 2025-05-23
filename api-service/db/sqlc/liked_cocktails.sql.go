// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
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
SELECT c.id,
        c.name,
        c.is_alcoholic,
        c.glass,
        c.image,
        c.embedding,
        c.instructions,
        c.ingredients,
        c.created_at
from cocktails c
JOIN liked_cocktails lc ON c.id = lc.cocktail_id
WHERE lc.user_id = $1 and lc.cocktail_id = $2
GROUP BY lc.user_id, lc.cocktail_id, c.id, lc.id
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
	Embedding    pgvector.Vector   `json:"embedding"`
	Instructions string            `json:"instructions"`
	Ingredients  dto.IngredientDto `json:"ingredients"`
	CreatedAt    time.Time         `json:"created_at"`
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
		&i.Embedding,
		&i.Instructions,
		&i.Ingredients,
		&i.CreatedAt,
	)
	return i, err
}

const getLikedCocktails = `-- name: GetLikedCocktails :many
SELECT c.id,
        c.name,
        c.is_alcoholic,
        c.glass,
        c.image,
        c.embedding,
        c.instructions,
        c.ingredients,
        c.created_at
from cocktails c
JOIN  liked_cocktails lc ON c.id = lc.cocktail_id
WHERE lc.user_id = $1
GROUP BY lc.user_id, c.id, lc.id
`

type GetLikedCocktailsRow struct {
	ID           uuid.UUID         `json:"id"`
	Name         string            `json:"name"`
	IsAlcoholic  pgtype.Bool       `json:"is_alcoholic"`
	Glass        string            `json:"glass"`
	Image        string            `json:"image"`
	Embedding    pgvector.Vector   `json:"embedding"`
	Instructions string            `json:"instructions"`
	Ingredients  dto.IngredientDto `json:"ingredients"`
	CreatedAt    time.Time         `json:"created_at"`
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
			&i.Embedding,
			&i.Instructions,
			&i.Ingredients,
			&i.CreatedAt,
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

const isCocktailLiked = `-- name: IsCocktailLiked :one
SELECT EXISTS (
    SELECT 1
    FROM liked_cocktails lc
    WHERE lc.user_id = $1 AND lc.cocktail_id = $2
) AS is_liked
`

type IsCocktailLikedParams struct {
	UserID     uuid.UUID `json:"user_id"`
	CocktailID uuid.UUID `json:"cocktail_id"`
}

func (q *Queries) IsCocktailLiked(ctx context.Context, arg IsCocktailLikedParams) (bool, error) {
	row := q.db.QueryRow(ctx, isCocktailLiked, arg.UserID, arg.CocktailID)
	var is_liked bool
	err := row.Scan(&is_liked)
	return is_liked, err
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
