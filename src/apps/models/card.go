package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

type Card struct {
	ID         uuid.UUID      `db:"id" json:"id"`
	IdentityId uuid.UUID      `db:"identity_id" json:"identity_id"`
	HolderName *string        `db:"holder_name" json:"holder_name"`
	Brand      *string        `db:"brand" json:"brand"`
	Meta       types.JSONText `db:"meta" json:"meta"`
	Customer   *string        `db:"customer" json:"customer"`
	IsJp       bool           `db:"is_jp" json:"is_jp"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (Card) TableName() string {
	return "cards"
}

func (Card) FetchQuery() string {
	return "cards/fetch"
}

func GetCard(id uuid.UUID, identityId uuid.UUID) (*Card, error) {
	c := new(Card)
	if err := database.Get(c, "cards/fetch_by_identity", id, identityId); err != nil {
		return nil, err
	}
	return c, nil

}
