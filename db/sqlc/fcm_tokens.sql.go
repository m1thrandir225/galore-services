// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: fcm_tokens.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createFCMToken = `-- name: CreateFCMToken :one
INSERT INTO fcm_tokens (
  token,
  device_id,
  user_id
) VALUES (
  $1,
  $2,
  $3
) RETURNING id, token, device_id, user_id, created_at
`

type CreateFCMTokenParams struct {
	Token    string    `json:"token"`
	DeviceID string    `json:"device_id"`
	UserID   uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateFCMToken(ctx context.Context, arg CreateFCMTokenParams) (FcmToken, error) {
	row := q.db.QueryRow(ctx, createFCMToken, arg.Token, arg.DeviceID, arg.UserID)
	var i FcmToken
	err := row.Scan(
		&i.ID,
		&i.Token,
		&i.DeviceID,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteFCMToken = `-- name: DeleteFCMToken :exec
DELETE FROM fcm_tokens
WHERE id = $1
`

func (q *Queries) DeleteFCMToken(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteFCMToken, id)
	return err
}

const getFCMTokenById = `-- name: GetFCMTokenById :one
SELECT id, token, device_id, user_id, created_at FROM fcm_tokens
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetFCMTokenById(ctx context.Context, id uuid.UUID) (FcmToken, error) {
	row := q.db.QueryRow(ctx, getFCMTokenById, id)
	var i FcmToken
	err := row.Scan(
		&i.ID,
		&i.Token,
		&i.DeviceID,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}