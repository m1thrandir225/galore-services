// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: reset_password_request.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const clearExpiredRequests = `-- name: ClearExpiredRequests :exec
DELETE FROM reset_password_request
WHERE valid_until < NOW() - INTERVAL '30 days'
`

func (q *Queries) ClearExpiredRequests(ctx context.Context) error {
	_, err := q.db.Exec(ctx, clearExpiredRequests)
	return err
}

const createResetPasswordRequest = `-- name: CreateResetPasswordRequest :one
INSERT INTO reset_password_request(user_id, valid_until)
VALUES ($1, $2)
RETURNING id, user_id, password_reset, valid_until
`

type CreateResetPasswordRequestParams struct {
	UserID     uuid.UUID        `json:"user_id"`
	ValidUntil pgtype.Timestamp `json:"valid_until"`
}

func (q *Queries) CreateResetPasswordRequest(ctx context.Context, arg CreateResetPasswordRequestParams) (ResetPasswordRequest, error) {
	row := q.db.QueryRow(ctx, createResetPasswordRequest, arg.UserID, arg.ValidUntil)
	var i ResetPasswordRequest
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PasswordReset,
		&i.ValidUntil,
	)
	return i, err
}

const getResetPasswordRequest = `-- name: GetResetPasswordRequest :one
SELECT id, user_id, password_reset, valid_until
FROM reset_password_request
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetResetPasswordRequest(ctx context.Context, id uuid.UUID) (ResetPasswordRequest, error) {
	row := q.db.QueryRow(ctx, getResetPasswordRequest, id)
	var i ResetPasswordRequest
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PasswordReset,
		&i.ValidUntil,
	)
	return i, err
}

const updateResetPasswordRequest = `-- name: UpdateResetPasswordRequest :one
UPDATE reset_password_request
SET password_reset = $2, valid_until = $3
WHERE id = $1
RETURNING id, user_id, password_reset, valid_until
`

type UpdateResetPasswordRequestParams struct {
	ID            uuid.UUID        `json:"id"`
	PasswordReset bool             `json:"password_reset"`
	ValidUntil    pgtype.Timestamp `json:"valid_until"`
}

func (q *Queries) UpdateResetPasswordRequest(ctx context.Context, arg UpdateResetPasswordRequestParams) (ResetPasswordRequest, error) {
	row := q.db.QueryRow(ctx, updateResetPasswordRequest, arg.ID, arg.PasswordReset, arg.ValidUntil)
	var i ResetPasswordRequest
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PasswordReset,
		&i.ValidUntil,
	)
	return i, err
}