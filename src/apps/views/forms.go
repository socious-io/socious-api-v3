package views

import "github.com/google/uuid"

type ServiceForm struct {
	Title             string    `json:"title" validate:"required,min=3"`
	Description       string    `json:"description" validate:"required,min=3"`
	PaymentCurrency   string    `json:"payment_currency" validate:"required"`
	Skills            []string  `json:"skills" validate:"required"`
	JobCategoryId     uuid.UUID `json:"job_category_id" validate:"required"`
	ServiceTotalHours int       `json:"service_total_hours" validate:"required,min=3"`
	ServicePrice      int       `json:"service_price" validate:"required,min=3"`
	ServiceLength     string    `json:"service_length" validate:"required"`
	WorkSamples       []string  `json:"work_samples" validate:"required"`
}
