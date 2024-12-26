-- name: CreateUser :one
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
) RETURNING id, name, email, avatar_url, birthday, enabled_push_notifications, enabled_email_notifications, created_at;

-- name: GetUser :one 
SELECT * FROM users 
WHERE id = $1 LIMIT 1;

-- name: GetUserHOTPSecret :one
SELECT hotp_secret FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 LIMIT 1;

-- name: UpdateUserInformation :one
UPDATE users
SET email=$2, avatar_url=$3, name=$4, birthday=$5
WHERE id = $1
RETURNING id, name, email, avatar_url, birthday, enabled_push_notifications, enabled_email_notifications, created_at;

-- name: UpdateUserPushNotifications :one
UPDATE users 
SET enabled_push_notifications = $2
WHERE id = $1
RETURNING enabled_push_notifications;

-- name: UpdateUserEmailNotifications :one
UPDATE users
SET enabled_email_notifications = $2
WHERE id = $1
RETURNING enabled_email_notifications;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1 
RETURNING id, name, email, avatar_url, birthday, enabled_push_notifications, enabled_email_notifications, created_at;

