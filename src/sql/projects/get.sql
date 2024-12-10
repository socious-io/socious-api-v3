SELECT id, COUNT(*) OVER () as total_count
FROM projects p
WHERE p.identity_id = $1
LIMIT $2 OFFSET $3