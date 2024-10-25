package models

import (
	"context"
	"socious/src/database"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Media struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"-"`
	URL       string    `db:"url" json:"url"`
	Filename  string    `db:"filename" json:"filename"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (Media) TableName() string {
	return "media"
}

func (Media) FetchQuery() string {
	return "media/fetch"
}

func (m *Media) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(m)
}

func (m *Media) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"media/create",
		m.UserID, m.URL, m.Filename,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := m.Scan(rows); err != nil {
			return err
		}
	}
	return nil
}

func getAllMedia() ([]Media, error) {
	result := []Media{}
	return result, nil
}

func getManyMedia() ([]Media, error) {
	result := []Media{}
	return result, nil
}