-- name: CreateGenerateCocktailDraft :one
INSERT INTO generate_cocktail_drafts(
         request_id,
         name,
         description,
         ingredients,
         instructions,
         main_image_prompt,
         steps_image_prompts
)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7
)
RETURNING *;