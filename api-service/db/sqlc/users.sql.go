// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  id,
  email,
  name,
  password,
  avatar_url,
  birthday,
  hotp_secret
) VALUES (
    $1,
  $2,
  $3,
  $4,
  $5,
  $6,
   $7
) RETURNING id, name, email, avatar_url, birthday, enabled_push_notifications, enabled_email_notifications, created_at
`

type CreateUserParams struct {
	ID         uuid.UUID   `json:"id"`
	Email      string      `json:"email"`
	Name       string      `json:"name"`
	Password   string      `json:"password"`
	AvatarUrl  string      `json:"avatar_url"`
	Birthday   pgtype.Date `json:"birthday"`
	HotpSecret string      `json:"hotp_secret"`
}

type CreateUserRow struct {
	ID                        uuid.UUID   `json:"id"`
	Name                      string      `json:"name"`
	Email                     string      `json:"email"`
	AvatarUrl                 string      `json:"avatar_url"`
	Birthday                  pgtype.Date `json:"birthday"`
	EnabledPushNotifications  bool        `json:"enabled_push_notifications"`
	EnabledEmailNotifications bool        `json:"enabled_email_notifications"`
	CreatedAt                 time.Time   `json:"created_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Name,
		arg.Password,
		arg.AvatarUrl,
		arg.Birthday,
		arg.HotpSecret,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.AvatarUrl,
		&i.Birthday,
		&i.EnabledPushNotifications,
		&i.EnabledEmailNotifications,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1 
RETURNING id, name, email, avatar_url, birthday, enabled_push_notifications, enabled_email_notifications, created_at
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, email, name, password, avatar_url, hotp_secret, enabled_push_notifications, enabled_email_notifications, birthday, created_at FROM users 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.AvatarUrl,
		&i.HotpSecret,
		&i.EnabledPushNotifications,
		&i.EnabledEmailNotifications,
		&i.Birthday,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, name, password, avatar_url, hotp_secret, enabled_push_notifications, enabled_email_notifications, birthday, created_at FROM users 
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.AvatarUrl,
		&i.HotpSecret,
		&i.EnabledPushNotifications,
		&i.EnabledEmailNotifications,
		&i.Birthday,
		&i.CreatedAt,
	)
	return i, err
}

const getUserHOTPSecret = `-- name: GetUserHOTPSecret :one
SELECT hotp_secret FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserHOTPSecret(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRow(ctx, getUserHOTPSecret, id)
	var hotp_secret string
	err := row.Scan(&hotp_secret)
	return hotp_secret, err
}

const updateUserEmailNotifications = `-- name: UpdateUserEmailNotifications :one
UPDATE users
SET enabled_email_notifications = $2
WHERE id = $1
RETURNING enabled_email_notifications
`

type UpdateUserEmailNotificationsParams struct {
	ID                        uuid.UUID `json:"id"`
	EnabledEmailNotifications bool      `json:"enabled_email_notifications"`
}

func (q *Queries) UpdateUserEmailNotifications(ctx context.Context, arg UpdateUserEmailNotificationsParams) (bool, error) {
	row := q.db.QueryRow(ctx, updateUserEmailNotifications, arg.ID, arg.EnabledEmailNotifications)
	var enabled_email_notifications bool
	err := row.Scan(&enabled_email_notifications)
	return enabled_email_notifications, err
}

const updateUserInformation = `-- name: UpdateUserInformation :one
UPDATE users
SET email=$2, avatar_url=$3, name=$4, birthday=$5
WHERE id = $1
RETURNING id, name, email, avatar_url, birthday, enabled_push_notifications, enabled_email_notifications, created_at
`

type UpdateUserInformationParams struct {
	ID        uuid.UUID   `json:"id"`
	Email     string      `json:"email"`
	AvatarUrl string      `json:"avatar_url"`
	Name      string      `json:"name"`
	Birthday  pgtype.Date `json:"birthday"`
}

type UpdateUserInformationRow struct {
	ID                        uuid.UUID   `json:"id"`
	Name                      string      `json:"name"`
	Email                     string      `json:"email"`
	AvatarUrl                 string      `json:"avatar_url"`
	Birthday                  pgtype.Date `json:"birthday"`
	EnabledPushNotifications  bool        `json:"enabled_push_notifications"`
	EnabledEmailNotifications bool        `json:"enabled_email_notifications"`
	CreatedAt                 time.Time   `json:"created_at"`
}

func (q *Queries) UpdateUserInformation(ctx context.Context, arg UpdateUserInformationParams) (UpdateUserInformationRow, error) {
	row := q.db.QueryRow(ctx, updateUserInformation,
		arg.ID,
		arg.Email,
		arg.AvatarUrl,
		arg.Name,
		arg.Birthday,
	)
	var i UpdateUserInformationRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.AvatarUrl,
		&i.Birthday,
		&i.EnabledPushNotifications,
		&i.EnabledEmailNotifications,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1
`

type UpdateUserPasswordParams struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.ID, arg.Password)
	return err
}

const updateUserPushNotifications = `-- name: UpdateUserPushNotifications :one
UPDATE users 
SET enabled_push_notifications = $2
WHERE id = $1
RETURNING enabled_push_notifications
`

type UpdateUserPushNotificationsParams struct {
	ID                       uuid.UUID `json:"id"`
	EnabledPushNotifications bool      `json:"enabled_push_notifications"`
}

func (q *Queries) UpdateUserPushNotifications(ctx context.Context, arg UpdateUserPushNotificationsParams) (bool, error) {
	row := q.db.QueryRow(ctx, updateUserPushNotifications, arg.ID, arg.EnabledPushNotifications)
	var enabled_push_notifications bool
	err := row.Scan(&enabled_push_notifications)
	return enabled_push_notifications, err
}
