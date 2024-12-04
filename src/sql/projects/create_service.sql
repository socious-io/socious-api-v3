INSERT
INTO projects(
    identity_id, title, description,
    payment_currency, skills, job_category_id,
    service_length, service_total_hours, service_price,
    payment_type, kind, status
)
VALUES (
    $1, $2, $3,
    $4, $5, $6, $7,
    $8, $9,
    'PAID', 'SERVICE', 'ACTIVE'
)
RETURNING *