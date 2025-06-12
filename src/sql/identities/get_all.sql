SELECT DISTINCT ON (i.id)
    i.*,
    (CASE WHEN i.type='users' THEN true ELSE false END) AS primary,
    (CASE WHEN i.id=$2 THEN true ELSE false END) AS current
FROM identities i
LEFT JOIN org_members om ON om.user_id=$1
WHERE i.id=om.org_id OR i.id=om.user_id