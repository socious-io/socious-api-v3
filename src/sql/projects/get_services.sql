SELECT id, COUNT(*) OVER () as total_count
FROM projects p
WHERE p.identity_id = $1 AND p.kind='SERVICE'
LIMIT $2 OFFSET $3