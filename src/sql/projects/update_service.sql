UPDATE projects p SET
    title=$2,
    description=$3,
    payment_currency=$4,
    project_length=$5,
    skills=$6,
    job_category_id=$7,
    service_total_hours=$8,
    service_price=$9
WHERE id=$1
RETURNING *