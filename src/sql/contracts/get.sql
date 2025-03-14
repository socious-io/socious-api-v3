SELECT c.id, COUNT(*) OVER () as total_count
FROM contracts c
LEFT JOIN gopay_payments gp on gp.unique_ref::uuid=c.id
LEFT JOIN projects p ON p.id=c.project_id
WHERE
    (c.provider_id = $1 OR c.client_id=$1) AND
    (gp.transaction_status!='CANCELED' OR gp.transaction_status IS NULL OR p.kind!='SERVICE') AND
    ($4::contract_status[]='{}' OR c.status=ANY($4))
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3