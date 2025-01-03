package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type Identity struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	Type      string         `db:"type" json:"type"`
	Meta      types.JSONText `db:"meta" json:"meta"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
}

func (Identity) TableName() string {
	return "identities"
}

func (Identity) FetchQuery() string {
	return "identities/fetch"
}
