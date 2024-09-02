-- name: LikeFlavour :one
INSERT INTO liked_flavours (
  flavour_id,
  user_id
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: UnlikeFlavour :exec
DELETE FROM liked_flavours 
WHERE flavour_id = $1 AND user_id = $2; 
