-- name: GetCurrentCounter :one
SELECT counter
FROM hotp_counters
WHERE user_id = $1;

-- name: CreateHotpCounter :exec
INSERT INTO hotp_counters (user_id, counter)
VALUES ($1, $2)
ON CONFLICT (user_id) DO NOTHING;

-- name: CleanupExpiredCounters :exec
DELETE FROM hotp_counters
WHERE last_used < NOW() - INTERVAL '30 days';

-- name: IncreaseCounter :one
UPDATE hotp_counters
SET counter = counter + 1, last_used = NOW()
WHERE user_id = $1
RETURNING *;
