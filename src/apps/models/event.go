package models

import (
	"time"

	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

type Event struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	EventAt     time.Time `json:"updated_at" db:"updated_at"`
}

func (Event) TableName() string {
	return "socious_events"
}

func (Event) FetchQuery() string {
	return "events/fetch"
}

func GetActiveEvent() (*Event, error) {
	e := new(Event)
	if err := database.Get(e, "events/get"); err != nil {
		return nil, err
	}
	return e, nil
}
