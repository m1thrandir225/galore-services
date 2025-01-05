-- name: CreateGenerateCocktailRequest :one
INSERT INTO
generate_cocktail_requests(user_id, prompt, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateGenerateCocktailRequest :one
UPDATE generate_cocktail_requests
SET status = $2
WHERE id = $1
RETURNING *;

-- name: GetGenerationRequest :one
SELECT * from generate_cocktail_requests
WHERE id = $1 LIMIT 1;

-- name: GetUserGenerationRequests :many
SELECT * FROM generate_cocktail_requests
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetIncompleteUserGenerationRequests :many
SELECT * FROM generate_cocktail_requests
WHERE user_id = $1 AND status != 'success'
ORDER BY created_at DESC;
