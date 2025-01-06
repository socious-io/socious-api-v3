package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

type Contract struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`

	Type   ContractType   `db:"type" json:"type"`
	Status ContractStatus `db:"status" json:"status"`

	TotalAmount           float64                  `db:"total_amount" json:"total_amount"`
	Currency              *Currency                `db:"currency" json:"currency"`
	CryptoCurrency        *string                  `db:"crypto_currency" json:"crypto_currency"`
	CurrencyRate          float32                  `db:"currency_rate" json:"currency_rate"`
	Commitment            int                      `db:"commitment" json:"commitment"`
	CommitmentPeriod      ContractCommitmentPeriod `db:"commitment_period" json:"commitment_period"`
	CommitmentPeriodCount int                      `db:"commitment_period_count" json:"commitment_period_count"`
	PaymentType           *PaymentModeType         `db:"payment_type" json:"payment_type"`

	ProviderID uuid.UUID `db:"provider_id" json:"provider_id"`
	ClientID   uuid.UUID `db:"client_id" json:"client_id"`

	ApplicantID *uuid.UUID `db:"applicant_id" json:"applicant_id"`
	ProjectID   *uuid.UUID `db:"project_id" json:"project_id"`
	PaymentID   *uuid.UUID `db:"payment_id" json:"payment_id"`
	OfferID     *uuid.UUID `db:"offer_id" json:"-"`
	MissionID   *uuid.UUID `db:"mission_id" json:"-"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	PaymentJson   types.JSONText `db:"payment" json:"-"`
	ProviderJson  types.JSONText `db:"provider" json:"-"`
	ClientJson    types.JSONText `db:"client" json:"-"`
	ProjectJson   types.JSONText `db:"project" json:"-"`
	ApplicantJson types.JSONText `db:"applicant" json:"-"`
}

func (Contract) TableName() string {
	return "contracts"
}

func (Contract) FetchQuery() string {
	return "contracts/fetch"
}

func (c *Contract) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"contracts/create",
		c.Name,
		c.Description,
		c.Type,
		c.TotalAmount,
		c.Currency,
		c.CryptoCurrency,
		c.CurrencyRate,
		c.Commitment,
		c.CommitmentPeriod,
		c.CommitmentPeriodCount,
		c.PaymentType,
		c.ProjectID,
		c.ApplicantID,
		c.ProviderID,
		c.ClientID,
	)

	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(c); err != nil {
			return err
		}
	}

	return database.Fetch(c, c.ID)
}

func (c *Contract) Update(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"contracts/update",
		c.ID,
		c.Name,
		c.Description,
		c.TotalAmount,
		c.Currency,
		c.CryptoCurrency,
		c.CurrencyRate,
		c.Commitment,
		c.CommitmentPeriod,
		c.CommitmentPeriodCount,
		c.PaymentType,
		c.Status,
		c.PaymentID,
	)

	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(c); err != nil {
			return err
		}
	}
	return database.Fetch(c, c.ID)
}

func GetContract(id uuid.UUID) (*Contract, error) {
	c := new(Contract)
	if err := database.Fetch(c, id); err != nil {
		return nil, err
	}
	return c, nil
}

func GetContracts(identityId uuid.UUID, p database.Paginate) ([]Contract, int, error) {
	var (
		contracts = []Contract{}
		fetchList []database.FetchList
		ids       []interface{}
	)

	if err := database.QuerySelect("contracts/get", &fetchList, identityId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return contracts, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&contracts, ids...); err != nil {
		return nil, 0, err
	}
	return contracts, fetchList[0].TotalCount, nil
}
