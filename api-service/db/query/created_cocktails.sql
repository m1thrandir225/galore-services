-- name: CreateUserCocktail :one
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
) RETURNING *;

-- name: GetUserCocktail :one
SELECT * from created_cocktails 
WHERE id = $1 LIMIT 1;

-- name: DeleteUserCocktail :exec
DELETE FROM created_cocktails 
WHERE id = $1;
