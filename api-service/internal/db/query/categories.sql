-- name: CreateCategory :one 
INSERT INTO categories (
  name,
  tag
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: GetAllCategories :many
SELECT * FROM categories
ORDER BY created_at DESC;

-- name: GetCategoryByTag :one
SELECT * FROM categories
WHERE tag = $1 LIMIT 1;

-- name: GetCategoryById :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, tag = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

