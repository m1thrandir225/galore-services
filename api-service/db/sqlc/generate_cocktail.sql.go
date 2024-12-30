// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: generate_cocktail.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createGeneratedCocktail = `-- name: CreateGeneratedCocktail :one
INSERT INTO generated_cocktails (
user_id,
request_id,
draft_id,
instructions,
ingredients,
description,
main_image_url,
additional_images
)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
)
RETURNING id, user_id, request_id, draft_id, name, description, main_image_url, additional_images, instructions, ingredients, created_at
`

type CreateGeneratedCocktailParams struct {
	UserID           uuid.UUID `json:"user_id"`
	RequestID        uuid.UUID `json:"request_id"`
	DraftID          uuid.UUID `json:"draft_id"`
	Instructions     []byte    `json:"instructions"`
	Ingredients      []byte    `json:"ingredients"`
	Description      string    `json:"description"`
	MainImageUrl     string    `json:"main_image_url"`
	AdditionalImages []string  `json:"additional_images"`
}

func (q *Queries) CreateGeneratedCocktail(ctx context.Context, arg CreateGeneratedCocktailParams) (GeneratedCocktail, error) {
	row := q.db.QueryRow(ctx, createGeneratedCocktail,
		arg.UserID,
		arg.RequestID,
		arg.DraftID,
		arg.Instructions,
		arg.Ingredients,
		arg.Description,
		arg.MainImageUrl,
		arg.AdditionalImages,
	)
	var i GeneratedCocktail
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RequestID,
		&i.DraftID,
		&i.Name,
		&i.Description,
		&i.MainImageUrl,
		&i.AdditionalImages,
		&i.Instructions,
		&i.Ingredients,
		&i.CreatedAt,
	)
	return i, err
}

const getGeneratedCocktail = `-- name: GetGeneratedCocktail :one
SELECT id, user_id, request_id, draft_id, name, description, main_image_url, additional_images, instructions, ingredients, created_at
FROM generated_cocktails
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetGeneratedCocktail(ctx context.Context, id uuid.UUID) (GeneratedCocktail, error) {
	row := q.db.QueryRow(ctx, getGeneratedCocktail, id)
	var i GeneratedCocktail
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RequestID,
		&i.DraftID,
		&i.Name,
		&i.Description,
		&i.MainImageUrl,
		&i.AdditionalImages,
		&i.Instructions,
		&i.Ingredients,
		&i.CreatedAt,
	)
	return i, err
}

const getUserGeneratedCocktails = `-- name: GetUserGeneratedCocktails :many
SELECT id, user_id, request_id, draft_id, name, description, main_image_url, additional_images, instructions, ingredients, created_at
FROM generated_cocktails
WHERE user_id = $1
`

func (q *Queries) GetUserGeneratedCocktails(ctx context.Context, userID uuid.UUID) ([]GeneratedCocktail, error) {
	rows, err := q.db.Query(ctx, getUserGeneratedCocktails, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GeneratedCocktail{}
	for rows.Next() {
		var i GeneratedCocktail
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RequestID,
			&i.DraftID,
			&i.Name,
			&i.Description,
			&i.MainImageUrl,
			&i.AdditionalImages,
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
