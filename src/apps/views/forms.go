package views

import (
	"socious/src/apps/models"

	"github.com/google/uuid"
)

type ServiceForm struct {
	Title           string               `json:"title" validate:"required,min=3"`
	Description     string               `json:"description"`
	PaymentCurrency string               `json:"payment_currency" validate:"required"`
	Skills          []string             `json:"skills" validate:"required"`
	JobCategoryId   uuid.UUID            `json:"job_category_id" validate:"required"`
	TotalHours      string               `json:"total_hours" validate:"required"`
	Price           string               `json:"price" validate:"required"`
	ProjectLength   models.ProjectLength `json:"project_length" validate:"required"`
	WorkSamples     []uuid.UUID          `json:"work_samples" validate:"required"`
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
