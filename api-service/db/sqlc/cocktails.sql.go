// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: cocktails.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	dto "github.com/m1thrandir225/galore-services/dto"
	"github.com/pgvector/pgvector-go"
)

const createCocktail = `-- name: CreateCocktail :one
INSERT INTO cocktails (
  name,
  is_alcoholic,
  glass,
  image,
  instructions,
  ingredients,
  embedding
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7
) RETURNING id, name, is_alcoholic, glass, image, instructions, ingredients, embedding, created_at
`

type CreateCocktailParams struct {
	Name         string            `json:"name"`
	IsAlcoholic  pgtype.Bool       `json:"is_alcoholic"`
	Glass        string            `json:"glass"`
	Image        string            `json:"image"`
	Instructions string            `json:"instructions"`
	Ingredients  dto.IngredientDto `json:"ingredients"`
	Embedding    pgvector.Vector   `json:"embedding"`
}

func (q *Queries) CreateCocktail(ctx context.Context, arg CreateCocktailParams) (Cocktail, error) {
	row := q.db.QueryRow(ctx, createCocktail,
		arg.Name,
		arg.IsAlcoholic,
		arg.Glass,
		arg.Image,
		arg.Instructions,
		arg.Ingredients,
		arg.Embedding,
	)
	var i Cocktail
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
	)
	return i, err
}

const deleteCocktail = `-- name: DeleteCocktail :exec
DELETE FROM cocktails 
WHERE id = $1
`

func (q *Queries) DeleteCocktail(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteCocktail, id)
	return err
}

const getCocktail = `-- name: GetCocktail :one
SELECT id, name, is_alcoholic, glass, image, instructions, ingredients, embedding, created_at FROM cocktails
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCocktail(ctx context.Context, id uuid.UUID) (Cocktail, error) {
	row := q.db.QueryRow(ctx, getCocktail, id)
	var i Cocktail
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
	)
	return i, err
}

const getCocktailAndSimilar = `-- name: GetCocktailAndSimilar :many
WITH target_cocktail AS (
    SELECT embedding
    FROM cocktails
    WHERE cocktails.id = $1
    LIMIT 1
)
SELECT c.id, c.name, c.is_alcoholic, c.glass, c.image, c.instructions, c.ingredients, c.embedding, c.created_at
FROM cocktails c, target_cocktail t
WHERE c.id != $1
ORDER BY c.embedding <=> t.embedding
LIMIT 10
`

func (q *Queries) GetCocktailAndSimilar(ctx context.Context, id uuid.UUID) ([]Cocktail, error) {
	rows, err := q.db.Query(ctx, getCocktailAndSimilar, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Cocktail{}
	for rows.Next() {
		var i Cocktail
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

const getHomescreenForUser = `-- name: GetHomescreenForUser :many
WITH UserLikedCategories AS (
    SELECT cf.category_id
    FROM liked_flavours lf
             JOIN category_flavour cf ON lf.flavour_id = cf.flavour_id
    WHERE lf.user_id = $1
    GROUP BY cf.category_id
    ORDER BY RANDOM()
    LIMIT 4
),
     CategoryCocktails AS (
         SELECT cc.cocktail_id, cc.category_id
         FROM cocktail_categories cc
         WHERE cc.category_id IN (SELECT category_id FROM UserLikedCategories)
     ),
     RankedCocktails AS (
         SELECT
             cc.cocktail_id,
             cc.category_id,
             ROW_NUMBER() OVER (PARTITION BY cc.category_id ORDER BY RANDOM()) AS rank
         FROM CategoryCocktails cc
     )
SELECT
    rc.category_id,
    array_agg(rc.cocktail_id)::uuid[] AS cocktails
FROM RankedCocktails rc
WHERE rc.rank <= 5
GROUP BY rc.category_id
HAVING COUNT(rc.cocktail_id) >= 2
`

type GetHomescreenForUserRow struct {
	CategoryID uuid.UUID   `json:"category_id"`
	Cocktails  []uuid.UUID `json:"cocktails"`
}

func (q *Queries) GetHomescreenForUser(ctx context.Context, userID uuid.UUID) ([]GetHomescreenForUserRow, error) {
	rows, err := q.db.Query(ctx, getHomescreenForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetHomescreenForUserRow{}
	for rows.Next() {
		var i GetHomescreenForUserRow
		if err := rows.Scan(&i.CategoryID, &i.Cocktails); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchCocktails = `-- name: SearchCocktails :many
SELECT
    id,
    name,
    is_alcoholic,
    glass,
    image,
    instructions,
    ingredients,
    embedding,
    created_at
FROM cocktails
WHERE
    $1::TEXT IS NULL
   OR
    (name ILIKE '%' || $1::TEXT || '%')
   OR
    (ingredients::TEXT ILIKE '%' || $1::TEXT || '%')
`

func (q *Queries) SearchCocktails(ctx context.Context, dollar_1 string) ([]Cocktail, error) {
	rows, err := q.db.Query(ctx, searchCocktails, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Cocktail{}
	for rows.Next() {
		var i Cocktail
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

const updateCocktail = `-- name: UpdateCocktail :one
UPDATE cocktails
SET name=$2, is_alcoholic=$3, glass=$4, image=$5, instructions=$6, ingredients=$7
WHERE id = $1
RETURNING id, name, is_alcoholic, glass, image, instructions, ingredients, embedding, created_at
`

type UpdateCocktailParams struct {
	ID           uuid.UUID         `json:"id"`
	Name         string            `json:"name"`
	IsAlcoholic  pgtype.Bool       `json:"is_alcoholic"`
	Glass        string            `json:"glass"`
	Image        string            `json:"image"`
	Instructions string            `json:"instructions"`
	Ingredients  dto.IngredientDto `json:"ingredients"`
}

func (q *Queries) UpdateCocktail(ctx context.Context, arg UpdateCocktailParams) (Cocktail, error) {
	row := q.db.QueryRow(ctx, updateCocktail,
		arg.ID,
		arg.Name,
		arg.IsAlcoholic,
		arg.Glass,
		arg.Image,
		arg.Instructions,
		arg.Ingredients,
	)
	var i Cocktail
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
	)
	return i, err
}
