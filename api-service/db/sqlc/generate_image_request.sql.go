// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: generate_image_request.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const checkImageGenerationProgress = `-- name: CheckImageGenerationProgress :one
SELECT
    d.id as draft_id,
    d.request_id,
    COUNT(i.id) as total_images,
    COUNT(i.status = 'success') as completed_images,
    bool_and(i.error_message IS NULL) as all_successful
FROM generate_cocktail_drafts d
JOIN generate_image_requests i ON i.draft_id = d.id
WHERE d.request_id = $1
GROUP BY d.id, d.request_id
`

type CheckImageGenerationProgressRow struct {
	DraftID         uuid.UUID `json:"draft_id"`
	RequestID       uuid.UUID `json:"request_id"`
	TotalImages     int64     `json:"total_images"`
	CompletedImages int64     `json:"completed_images"`
	AllSuccessful   bool      `json:"all_successful"`
}

func (q *Queries) CheckImageGenerationProgress(ctx context.Context, requestID uuid.UUID) (CheckImageGenerationProgressRow, error) {
	row := q.db.QueryRow(ctx, checkImageGenerationProgress, requestID)
	var i CheckImageGenerationProgressRow
	err := row.Scan(
		&i.DraftID,
		&i.RequestID,
		&i.TotalImages,
		&i.CompletedImages,
		&i.AllSuccessful,
	)
	return i, err
}

const createImageGenerationRequest = `-- name: CreateImageGenerationRequest :one
INSERT INTO generate_image_requests(draft_id, prompt, status, is_main)
VALUES ($1, $2, $3, $4)
RETURNING id, draft_id, prompt, is_main, status, image_url, error_message, created_at, updated_at
`

type CreateImageGenerationRequestParams struct {
	DraftID uuid.UUID             `json:"draft_id"`
	Prompt  string                `json:"prompt"`
	Status  ImageGenerationStatus `json:"status"`
	IsMain  bool                  `json:"is_main"`
}

func (q *Queries) CreateImageGenerationRequest(ctx context.Context, arg CreateImageGenerationRequestParams) (GenerateImageRequest, error) {
	row := q.db.QueryRow(ctx, createImageGenerationRequest,
		arg.DraftID,
		arg.Prompt,
		arg.Status,
		arg.IsMain,
	)
	var i GenerateImageRequest
	err := row.Scan(
		&i.ID,
		&i.DraftID,
		&i.Prompt,
		&i.IsMain,
		&i.Status,
		&i.ImageUrl,
		&i.ErrorMessage,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getImagesForDraft = `-- name: GetImagesForDraft :many
SELECT id, draft_id, prompt, is_main, status, image_url, error_message, created_at, updated_at
FROM generate_image_requests i
WHERE i.draft_id = $1 AND i.status = 'success'
`

func (q *Queries) GetImagesForDraft(ctx context.Context, draftID uuid.UUID) ([]GenerateImageRequest, error) {
	rows, err := q.db.Query(ctx, getImagesForDraft, draftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GenerateImageRequest{}
	for rows.Next() {
		var i GenerateImageRequest
		if err := rows.Scan(
			&i.ID,
			&i.DraftID,
			&i.Prompt,
			&i.IsMain,
			&i.Status,
			&i.ImageUrl,
			&i.ErrorMessage,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateImageGenerationRequest = `-- name: UpdateImageGenerationRequest :one
UPDATE generate_image_requests
SET image_url = $2,
    error_message = $3,
    status = $4
WHERE id = $1
RETURNING id, draft_id, prompt, is_main, status, image_url, error_message, created_at, updated_at
`

type UpdateImageGenerationRequestParams struct {
	ID           uuid.UUID             `json:"id"`
	ImageUrl     pgtype.Text           `json:"image_url"`
	ErrorMessage pgtype.Text           `json:"error_message"`
	Status       ImageGenerationStatus `json:"status"`
}

func (q *Queries) UpdateImageGenerationRequest(ctx context.Context, arg UpdateImageGenerationRequestParams) (GenerateImageRequest, error) {
	row := q.db.QueryRow(ctx, updateImageGenerationRequest,
		arg.ID,
		arg.ImageUrl,
		arg.ErrorMessage,
		arg.Status,
	)
	var i GenerateImageRequest
	err := row.Scan(
		&i.ID,
		&i.DraftID,
		&i.Prompt,
		&i.IsMain,
		&i.Status,
		&i.ImageUrl,
		&i.ErrorMessage,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
