UPDATE contracts SET 
  name=$2,
  description=$3,
  total_amount=$4,
  currency=$5,
  crypto_currency=$6,
  currency_rate=$7,
  commitment=$8,
  commitment_period=$9,
  commitment_period_count=$10,
  payment_type=$11
WHERE id=$1
RETURNING *