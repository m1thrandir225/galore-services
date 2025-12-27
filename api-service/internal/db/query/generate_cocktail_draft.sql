-- name: CreateGenerateCocktailDraft :one
INSERT INTO generate_cocktail_drafts(
         request_id,
         name,
         description,
         ingredients,
         instructions,
         main_image_prompt
)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
)
RETURNING *;

-- name: GetCocktailDraft :one
SELECT *
FROM generate_cocktail_drafts
WHERE id = $1;