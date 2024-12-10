package views

import (
	"socious/src/apps/models"

	"github.com/google/uuid"
)

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

type ContractForm struct {
	Title                 string                          `json:"title" validate:"required,min=3"`
	Description           string                          `json:"description"`
	Type                  models.ContractType             `json:"type" validate:"required"`
	TotalAmount           float32                         `json:"total_amount"`
	Currency              models.Currency                 `json:"currency"`
	CryptoCurrency        string                          `json:"crypto_currency"`
	CurrencyRate          float32                         `json:"currency_rate"`
	Commitment            int                             `json:"commitment"`
	CommitmentPeriod      models.ContractCommitmentPeriod `json:"commitment_period"`
	CommitmentPeriodCount int                             `json:"commitment_period_count"`
	PaymentType           *models.PaymentModeType         `json:"payment_type"`
	ApplicantID           *uuid.UUID                      `json:"applicant_id"`
	ProjectID             *uuid.UUID                      `json:"project_id"`
	ClientID              uuid.UUID                       `json:"client_id" validate:"required"`
}
