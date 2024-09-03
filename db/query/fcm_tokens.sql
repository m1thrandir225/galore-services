-- name: CreateFCMToken :one
INSERT INTO fcm_tokens (
  token,
  device_id,
  user_id
) VALUES (
  $1,
  $2,
  $3
) RETURNING *;


-- name: GetFCMTokenById :one
SELECT * FROM fcm_tokens
WHERE id = $1 LIMIT 1;

-- name: DeleteFCMToken :exec
DELETE FROM fcm_tokens
WHERE id = $1;
