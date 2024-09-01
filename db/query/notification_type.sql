-- name: CreateNotifcationType :one 
INSERT INTO notification_types (
  title,
  content,
  tag
) VALUES (
  $1,
  $2,
  $3
) RETURNING *;

-- name: GetAllTypes :many
SELECT * from notification_types;

-- name: GetNotificationType :one 
SELECT * from notification_types
WHERE id = $1 LIMIT 1;

-- name: DeleteNotificationType :exec
DELETE from notification_types 
WHERE id = $1; 
