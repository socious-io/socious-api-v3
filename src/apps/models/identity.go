package models

import (
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Type string    `db:"type" json:"type"` // 'users', 'organizations'
	Meta string    `db:"meta" json:"meta"` //jsonb

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (Identity) TableName() string {
	return "identities"
}

func (Identity) FetchQuery() string {
	return "identities/fetch"
}
