-- name: CreateNotification :one 
INSERT INTO notifications (
  user_id,
  notification_type
) VALUES (
  $1,
  $2
) RETURNING *;

-- name: GetUserNotifications :many
SELECT * FROM notifications 
WHERE user_id = $1 LIMIT 1;

-- name: UpdateUserNotification :one
UPDATE notifications
  SET opened = $2
  WHERE id = $1
  RETURNING *;
