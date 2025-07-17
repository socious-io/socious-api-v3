INSERT INTO users (id, first_name, last_name, username, email, city, country, avatar, cover_image, language, impact_points, identity_verified) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT (id) DO UPDATE SET
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    username = EXCLUDED.username,
    city = COALESCE(EXCLUDED.city, users.city),
    country = COALESCE(EXCLUDED.country, users.country),
    avatar = EXCLUDED.avatar,
    cover_image = EXCLUDED.cover_image,
    language = EXCLUDED.language,
    impact_points = EXCLUDED.impact_points,
    identity_verified = EXCLUDED.identity_verified
RETURNING *;