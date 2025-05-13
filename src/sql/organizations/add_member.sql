INSERT INTO org_members (org_id, user_id) VALUES ($1, $2)
ON CONFLICT (org_id, user_id) DO NOTHING