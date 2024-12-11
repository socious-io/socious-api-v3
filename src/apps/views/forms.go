package views

import (
	"socious/src/apps/models"

	"github.com/google/uuid"
)

type ProjectForm struct {
	Title                 string                          `json:"title" validate:"required"`
	Description           string                          `json:"description" validate:"required"`
	ProjectType           *models.ProjectType             `json:"project_type"`
	ProjectLength         *models.ProjectLength           `json:"project_length"`
	PaymentCurrency       *string                         `json:"payment_currency"`
	PaymentRangeLower     *string                         `json:"payment_range_lower"`
	PaymentRangeHigher    *string                         `json:"payment_range_higher"`
	ExperienceLevel       *int                            `json:"experience_level"`
	Status                *models.ProjectStatus           `json:"status"`
	PaymentType           *models.PaymentType             `json:"payment_type"`
	PaymentScheme         *models.PaymentScheme           `json:"payment_scheme"`
	Country               *string                         `json:"country"`
	Skills                []string                        `json:"skills" validate:"required"`
	CausesTags            []string                        `json:"causes_tags"`
	RemotePreference      *models.ProjectRemotePreference `json:"remote_preference" validate:"required"`
	City                  *string                         `json:"city"`
	WeeklyHoursLower      *string                         `json:"weekly_hours_lower"`
	WeeklyHoursHigher     *string                         `json:"weekly_hours_higher"`
	CommitmentHoursLower  *string                         `json:"commitment_hours_lower"`
	CommitmentHoursHigher *string                         `json:"commitment_hours_higher"`
	GeonameId             *int                            `json:"geoname_id"`
	JobCategoryId         uuid.UUID                       `json:"job_category_id" validate:"required"`
	Kind                  models.ProjectKind              `json:"kind"`
	WorkSamples           []uuid.UUID                     `json:"work_samples" validate:"required"`
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
