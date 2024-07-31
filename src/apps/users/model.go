package users

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
}

func (*User) Columns() []string {
	return []string{"id", "created_at"}
}

func (*User) TableName() string {
	return "users"
}

func (u *User) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(u)
}
