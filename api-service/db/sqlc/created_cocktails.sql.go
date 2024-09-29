// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: created_cocktails.sql

package db

import (
	"context"

	"github.com/google/uuid"
	dto "github.com/m1thrandir225/galore-services/dto"
	"github.com/pgvector/pgvector-go"
)

const createUserCocktail = `-- name: CreateUserCocktail :one
INSERT INTO created_cocktails (
  name,
  image, 
  ingredients,
  instructions, 
  description,
  user_id,
  embedding
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7
) RETURNING id, name, image, ingredients, instructions, description, user_id, embedding, created_at
`

type CreateUserCocktailParams struct {
	Name         string               `json:"name"`
	Image        string               `json:"image"`
	Ingredients  dto.IngredientDto    `json:"ingredients"`
	Instructions dto.AiInstructionDto `json:"instructions"`
	Description  string               `json:"description"`
	UserID       uuid.UUID            `json:"user_id"`
	Embedding    pgvector.Vector      `json:"embedding"`
}

func (q *Queries) CreateUserCocktail(ctx context.Context, arg CreateUserCocktailParams) (CreatedCocktail, error) {
	row := q.db.QueryRow(ctx, createUserCocktail,
		arg.Name,
		arg.Image,
		arg.Ingredients,
		arg.Instructions,
		arg.Description,
		arg.UserID,
		arg.Embedding,
	)
	var i CreatedCocktail
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.Ingredients,
		&i.Instructions,
		&i.Description,
		&i.UserID,
		&i.Embedding,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUserCocktail = `-- name: DeleteUserCocktail :exec
DELETE FROM created_cocktails 
WHERE id = $1
`

func (q *Queries) DeleteUserCocktail(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUserCocktail, id)
	return err
}

const getUserCocktail = `-- name: GetUserCocktail :one
SELECT id, name, image, ingredients, instructions, description, user_id, embedding, created_at from created_cocktails 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserCocktail(ctx context.Context, id uuid.UUID) (CreatedCocktail, error) {
	row := q.db.QueryRow(ctx, getUserCocktail, id)
	var i CreatedCocktail
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Image,
		&i.Ingredients,
		&i.Instructions,
		&i.Description,
		&i.UserID,
		&i.Embedding,
		&i.CreatedAt,
	)
	return i, err
}
