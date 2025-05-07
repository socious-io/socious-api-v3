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
	PaymentMode           *models.PaymentModeType         `json:"payment_mode"`
	GeonameId             *int                            `json:"geoname_id"`
	JobCategoryId         uuid.UUID                       `json:"job_category_id" validate:"required"`
	Kind                  models.ProjectKind              `json:"kind"`
	WorkSamples           []uuid.UUID                     `json:"work_samples" validate:"required"`
}

type ContractForm struct {
	Name                  string                          `json:"name" validate:"required,min=3"`
	Description           string                          `json:"description"`
	Type                  models.ContractType             `json:"type" validate:"required"`
	TotalAmount           float64                         `json:"total_amount"`
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

type ContractDepositForm struct {
	CardID *uuid.UUID  `json:"card_id" validate:"required"`
	TxID   *string     `json:"txid" validate:"required"`
	Meta   interface{} `json:"meta" validate:"required"`
}

type ContractRequirementsForm struct {
	RequirementDescription string      `json:"requirement_description" validate:"required"`
	RequirementFiles       []uuid.UUID `json:"requirement_files" validate:"required"`
}

type ContractFeedbackForm struct {
	Content   string `json:"content" validate:"required"`
	Satisfied bool   `json:"satisfied" validate:"required"`
}

type UserUpdateForm struct {
	Username  *string    `json:"username" validate:"required,min=3,max=32"`
	Bio       *string    `json:"bio"`
	FirstName string     `json:"first_name" validate:"required,min=3,max=32"`
	LastName  string     `json:"last_name" validate:"required,min=3,max=32"`
	Phone     *string    `json:"phone"`
	AvatarID  *uuid.UUID `json:"avatar_id"`
}

type OrganizationUpdateForm struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	Name              *string    `db:"name" json:"name"`
	Bio               *string    `db:"bio" json:"bio"`
	Description       *string    `db:"description" json:"description"`
	Email             *string    `db:"email" json:"email"`
	Phone             *string    `db:"phone" json:"phone"`
	City              *string    `db:"city" json:"city"`
	Type              string     `db:"type" json:"type"` //type -> organization_type DEFAULT 'OTHER'
	Address           *string    `db:"address" json:"address"`
	Website           *string    `db:"website" json:"website"`
	SocialCauses      []string   `db:"social_causes" json:"social_causes"` //type -> social_causes_type[]
	Followers         int        `db:"followers" json:"followers"`
	Followings        int        `db:"followings" json:"followings"`
	Country           *string    `db:"country" json:"country"`
	WalletAddress     *string    `db:"wallet_address" json:"wallet_address"`
	ImpactPoints      float64    `db:"impact_points" json:"impact_points"`
	Mission           *string    `db:"mission" json:"mission"`
	Culture           *string    `db:"culture" json:"culture"`
	Logo              *uuid.UUID `db:"image" json:"image"`
	Avatar            *uuid.UUID `db:"cover_image" json:"cover_image"`
	MobileCountryCode *string    `db:"mobile_country_code" json:"mobile_country_code"`
	CreatedBy         *uuid.UUID `db:"created_by" json:"created_by"`
	Shortname         string     `db:"shortname" json:"shortname"`
	OtherPartyId      *string    `db:"other_party_id" json:"other_party_id"`
	OtherPartyTitle   *string    `db:"other_party_title" json:"other_party_title"`
	OtherPartyUrl     *string    `db:"other_party_url" json:"other_party_url"`
	GeonameId         *int       `db:"geoname_id" json:"geoname_id"`
	VerifiedImpact    bool       `db:"verified_impact" json:"verified_impact"`
	Hiring            bool       `db:"hiring" json:"hiring"`
	Size              *string    `db:"size" json:"size"`
	Industry          *string    `db:"industry" json:"industry"`
	Did               *string    `db:"did" json:"did"`
}

type SyncForm struct {
	Organizations []models.Organization `json:"organizations"`
	User          models.User           `json:"user" validate:"required"`
}
