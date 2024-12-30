-- name: CreateGenerateCocktailRequest :one
INSERT INTO
generate_cocktail_requests(user_id, prompt, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserGenerationRequests :many
SELECT * FROM generate_cocktail_requests
WHERE user_id = $1
ORDER BY created_at DESC;
