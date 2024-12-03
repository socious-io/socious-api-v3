UPDATE projects p SET
    title=$2,
    description=$3,
    payment_currency=$4,
    skills=$5,
    job_category_id=$6,
    service_length=$7,
    service_total_hours=$8,
    service_price=$9
WHERE id=$1
RETURNING *