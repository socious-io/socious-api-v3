SELECT c.id, COUNT(*) OVER () as total_count
FROM contracts c
LEFT JOIN gopay_payments gp ON unique_ref=c.id::text
WHERE (c.provider_id = $1 OR c.client_id=$1) AND gp.status != 'CANCELED'
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3