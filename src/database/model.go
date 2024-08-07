package database

import "github.com/jmoiron/sqlx"

type Model interface {
	Columns() []string
	TableName() string
	Scan(rows *sqlx.Rows) error
	FetchQuery() string
}
