INSERT INTO feedbacks (
  content,
  is_contest,
  identity_id,
  project_id,
  mission_id,
  contract_id
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *