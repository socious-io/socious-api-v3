UPDATE oauth_connects SET
  matrix_unique_id=$2,
  access_token=$3,
  refresh_token=$4,
  updated_at=NOW()
WHERE id=$1
RETURNING *