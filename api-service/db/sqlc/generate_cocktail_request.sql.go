// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: generate_cocktail_request.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createGenerateCocktailRequest = `-- name: CreateGenerateCocktailRequest :one
INSERT INTO
generate_cocktail_requests(user_id, prompt, status)
VALUES ($1, $2, $3)
RETURNING id, user_id, prompt, status, error_message, updated_at, created_at
`

type CreateGenerateCocktailRequestParams struct {
	UserID uuid.UUID        `json:"user_id"`
	Prompt string           `json:"prompt"`
	Status GenerationStatus `json:"status"`
}

func (q *Queries) CreateGenerateCocktailRequest(ctx context.Context, arg CreateGenerateCocktailRequestParams) (GenerateCocktailRequest, error) {
	row := q.db.QueryRow(ctx, createGenerateCocktailRequest, arg.UserID, arg.Prompt, arg.Status)
	var i GenerateCocktailRequest
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Prompt,
		&i.Status,
		&i.ErrorMessage,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getGenerationRequest = `-- name: GetGenerationRequest :one
SELECT id, user_id, prompt, status, error_message, updated_at, created_at from generate_cocktail_requests
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetGenerationRequest(ctx context.Context, id uuid.UUID) (GenerateCocktailRequest, error) {
	row := q.db.QueryRow(ctx, getGenerationRequest, id)
	var i GenerateCocktailRequest
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Prompt,
		&i.Status,
		&i.ErrorMessage,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getIncompleteUserGenerationRequests = `-- name: GetIncompleteUserGenerationRequests :many
SELECT id, user_id, prompt, status, error_message, updated_at, created_at FROM generate_cocktail_requests
WHERE user_id = $1 AND status != 'success'
ORDER BY created_at DESC
`

func (q *Queries) GetIncompleteUserGenerationRequests(ctx context.Context, userID uuid.UUID) ([]GenerateCocktailRequest, error) {
	rows, err := q.db.Query(ctx, getIncompleteUserGenerationRequests, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GenerateCocktailRequest{}
	for rows.Next() {
		var i GenerateCocktailRequest
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Prompt,
			&i.Status,
			&i.ErrorMessage,
			&i.UpdatedAt,
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

const getUserGenerationRequests = `-- name: GetUserGenerationRequests :many
SELECT id, user_id, prompt, status, error_message, updated_at, created_at FROM generate_cocktail_requests
WHERE user_id = $1
ORDER BY created_at DESC
`

func (q *Queries) GetUserGenerationRequests(ctx context.Context, userID uuid.UUID) ([]GenerateCocktailRequest, error) {
	rows, err := q.db.Query(ctx, getUserGenerationRequests, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GenerateCocktailRequest{}
	for rows.Next() {
		var i GenerateCocktailRequest
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Prompt,
			&i.Status,
			&i.ErrorMessage,
			&i.UpdatedAt,
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

const updateGenerateCocktailRequest = `-- name: UpdateGenerateCocktailRequest :one
UPDATE generate_cocktail_requests
SET status = $2
WHERE id = $1
RETURNING id, user_id, prompt, status, error_message, updated_at, created_at
`

type UpdateGenerateCocktailRequestParams struct {
	ID     uuid.UUID        `json:"id"`
	Status GenerationStatus `json:"status"`
}

func (q *Queries) UpdateGenerateCocktailRequest(ctx context.Context, arg UpdateGenerateCocktailRequestParams) (GenerateCocktailRequest, error) {
	row := q.db.QueryRow(ctx, updateGenerateCocktailRequest, arg.ID, arg.Status)
	var i GenerateCocktailRequest
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Prompt,
		&i.Status,
		&i.ErrorMessage,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
