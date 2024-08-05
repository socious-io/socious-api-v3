package users

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Experience struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
}

func (*Experience) Columns() []string {
	return []string{"id", "title", "created_at"}
}

func (*Experience) TableName() string {
	return "experiences"
}

func (u *Experience) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(u)
}

func (u *Experience) Associations() []interface{} {
	return []interface{}{}
}
