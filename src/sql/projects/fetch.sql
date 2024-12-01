SELECT p.*,
(
	COALESCE(
		(SELECT
			jsonb_agg(
				json_build_object(
				'url', m.url, 
				'filename', m.filename
				)
			)
		FROM media m
		WHERE m.id = sws.document
		),
		'[]'
	)
) AS work_samples
FROM projects p
LEFT JOIN service_work_samples sws ON sws.service_id=p.id
WHERE p.id IN (?)