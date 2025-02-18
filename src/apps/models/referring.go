package models

import (
	"time"

	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

type Referring struct {
	ID uuid.UUID `db:"id" json:"id"`

	ReferredById  uuid.UUID `db:"referred_by_id" json:"referred_by_id"`
	WalletAddress *string   `db:"wallet_address" json:"wallet_address"`
	FeeDiscount   bool      `db:"fee_discount" json:"fee_discount"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (Referring) TableName() string {
	return "referrings"
}

func (Referring) FetchQuery() string {
	return "referrings/fetch"
}

func GetReferring(referredIdentityId uuid.UUID) (*Referring, error) {
	r := new(Referring)
	if err := database.Get(r, "referrings/get", referredIdentityId); err != nil {
		return nil, err
	}
	return r, nil
}
