INSERT INTO media(id, identity_id, url, filename)
VALUES($1, $2, $3, $4)
ON CONFLICT(id) DO UPDATE
SET
    identity_id=EXCLUDED.identity_id,
    url=EXCLUDED.url,
    filename=EXCLUDED.filename,
RETURNING *