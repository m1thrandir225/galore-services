-- name: CreateResetPasswordRequest :one
INSERT INTO reset_password_request(user_id, valid_until)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateResetPasswordRequest :one
UPDATE reset_password_request
SET password_reset = $2
WHERE id = $1
RETURNING *;

-- name: GetResetPasswordRequest :one
SELECT *
FROM reset_password_request
WHERE id = $1 LIMIT 1;

-- name: ClearExpiredRequests :exec
DELETE FROM reset_password_request
WHERE valid_until < NOW() - INTERVAL '30 days';