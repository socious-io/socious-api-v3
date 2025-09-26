SELECT 
    u.*,
    row_to_json(m1.*) AS avatar,
    row_to_json(m2.*) AS cover,
    COALESCE(
        (
            SELECT json_agg(w)
            FROM wallets w
            WHERE w.user_id = u.id
        ),
        '[]'::json
    ) AS wallets,
    COALESCE(u.tags, '{}'::text[]) AS tags,
    COALESCE(u.events, '{}'::uuid[]) AS events
FROM users u
LEFT JOIN media m1 ON m1.id = u.avatar
LEFT JOIN media m2 ON m2.id = u.cover_image
WHERE u.id IN (?);