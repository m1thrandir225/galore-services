-- name: CreateGeneratedCocktail :one
INSERT INTO generated_cocktails (
name,
user_id,
request_id,
draft_id,
instructions,
ingredients,
description,
main_image_url
)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
)
RETURNING *;

-- name: GetUserGeneratedCocktails :many
SELECT *
FROM generated_cocktails
WHERE user_id = $1;

-- name: GetGeneratedCocktail :one
SELECT *
FROM generated_cocktails
WHERE id = $1 LIMIT 1;
