package database

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Model interface {
	Columns() []string
	TableName() string
	Scan(rows *sqlx.Rows) error
	FetchQuery() string
}

// type Model interface {
// 	TableName() string
// 	// Scan(any) error
// 	FetchQuery() string
// }

type Paginate struct {
	Limit int
	Offet int
}

type FetchList struct {
	ID         uuid.UUID `db:"id"`
	TotalCount int       `db:"total_count"`
}
