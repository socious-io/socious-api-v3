UPDATE projects p SET
    title=$2,
    description=$3,
    payment_type=$4,
    payment_scheme=$5,
    payment_currency=$6,
    payment_range_lower=$7,
    payment_range_higher=$8,
    experience_level=$9,
    status=$10,
    remote_preference=$11,
    project_type=$12,
    project_length=$13,
    skills=$14,
    causes_tags=$15,
    country=$16,
    city=$17,
    geoname_id=$18,
	other_party_id=$19, 
    other_party_title=$20,
    other_party_url=$21,
    expires_at=$22,
    updated_at=$23,
    weekly_hours_lower=$24,
    weekly_hours_higher=$25,
    commitment_hours_lower=$26,
    commitment_hours_higher=$27,
    job_category_id=$28,
    payment_mode=$29
WHERE id=$1
RETURNING *