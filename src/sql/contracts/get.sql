SELECT id, COUNT(*) OVER () as total_count
FROM contracts c
WHERE c.provider_id = $1 OR c.client_id=$1
LIMIT $2 OFFSET $3