-- name: CreateUser :one
INSERT INTO users (
  email,
  name,
  password,
  avatar_url,
  birthday
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
) RETURNING id, name, email, avatar_url, birthday, enabled_push_notifications, enabled_email_notifications, created_at;

-- name: GetUser :one 
SELECT * FROM users 
WHERE id = $1 LIMIT 1; 

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1 
RETURNING id, name, email, avatar_url, birthday, enabled_push_notifications, enabled_email_notifications, created_at;

-- name: GetUserFCMTokens :many 
SELECT * FROM fcm_tokens 
WHERE user_id = $1; 

-- name: GetCreatedCocktails :many 
SELECT * FROM created_cocktails 
WHERE user_id = $1;


