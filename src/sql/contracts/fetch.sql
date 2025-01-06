SELECT c.*,
  row_to_json(id1.*) as provider,
  row_to_json(id2.*) as client,
  row_to_json(a.*) as applicant,
  row_to_json(p.*) as project,
  row_to_json(pay.*) as payment
FROM contracts c
JOIN identities id1 ON id1.id = c.provider_id
JOIN identities id2 ON id2.id = c.client_id
LEFT JOIN projects p ON p.id = c.project_id
LEFT JOIN applicants a ON a.id = c.applicant_id
LEFT JOIN payments pay ON pay.id = c.payment_id
WHERE c.id IN (?)