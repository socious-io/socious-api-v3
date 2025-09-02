INSERT INTO deleted_users (user_id, username, reason, registered_at)
VALUES ($1, $2, $3, $4)
RETURNING *;