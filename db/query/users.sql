-- name: CreateUser :one
INSERT INTO users (
  email,
  name,
  password,
  avatar_url
) VALUES (
  $1,
  $2,
  $3,
  $4
) RETURNING *;

-- name: GetUser :one 
SELECT * FROM users 
WHERE id = $1 LIMIT 1; 

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1 
RETURNING *;
