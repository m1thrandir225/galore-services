-- name: CreateCocktail :one
INSERT INTO cocktails (
  name,
  is_alcoholic,
  glass,
  image,
  instructions,
  ingredients
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
) RETURNING *;

-- name: GetCocktail :one
SELECT * FROM cocktails
WHERE id = $1 LIMIT 1;

-- name: DeleteCocktail :exec 
DELETE FROM cocktails 
WHERE id = $1;
