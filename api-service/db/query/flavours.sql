-- name: CreateFlavour :one
INSERT INTO flavours (
  name
) VALUES (
  $1
) RETURNING *;

-- name: UpdateFlavour :one
UPDATE flavours
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteFlavour :exec
DELETE FROM flavours 
WHERE id = $1;

-- name: GetFlavourId :one
SELECT * FROM flavours
WHERE id = $1 LIMIT 1;

-- name: GetFlavourName :one
SELECT * FROM flavours
WHERE name = $1 LIMIT 1;

-- name: GetAllFlavours :many 
SELECT * FROM flavours;
