-- name: CreateCocktail :one
INSERT INTO cocktails (
  name,
  is_alcoholic,
  glass,
  image,
  instructions,
  ingredients,
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

-- name: GetCocktail :one
SELECT * FROM cocktails
WHERE id = $1 LIMIT 1;

-- name: GetCocktailAndSimilar :many
WITH target_cocktail AS (
    SELECT embedding
    FROM cocktails
    WHERE cocktails.id = $1
    LIMIT 1
)
SELECT c.*, c.embedding <=> t.embedding AS similarity_score
FROM cocktails c, target_cocktail t
ORDER BY similarity_score
LIMIT 10;


-- name: UpdateCocktail :one
UPDATE cocktails
SET name=$2, is_alcoholic=$3, glass=$4, image=$5, instructions=$6, ingredients=$7
WHERE id = $1
RETURNING *;

-- name: DeleteCocktail :exec 
DELETE FROM cocktails 
WHERE id = $1;
