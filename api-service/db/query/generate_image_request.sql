-- name: CreateImageGenerationRequest :one
INSERT INTO generate_image_requests(draft_id, prompt, status, is_main)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateImageGenerationRequest :one
UPDATE generate_image_requests
SET image_url = $2,
    error_message = $3,
    status = $4
WHERE id = $1
RETURNING *;

-- name: CheckImageGenerationProgress :one
SELECT
    d.id as draft_id,
    d.request_id,
    COUNT(i.id) as total_images,
    COUNT(i.status = 'success') as completed_images,
    bool_and(i.error_message IS NULL) as all_successful
FROM generate_cocktail_drafts d
JOIN generate_image_requests i ON i.draft_id = d.id
WHERE d.request_id = $1
GROUP BY d.id, d.request_id;

-- name: GetImagesForDraft :many
SELECT *
FROM generate_image_requests i
WHERE i.draft_id = $1 AND i.status = 'success';