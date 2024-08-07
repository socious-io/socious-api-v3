package users

import (
	"socious/src/database"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        string    `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (*User) Columns() []string {
	return []string{"id", "created_at"}
}

func (*User) TableName() string {
	return "users"
}

func (*User) FetchQuery() string {
	return "users/fetch"
}

func (u *User) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(u)
}

func Get(id string) (*User, error) {
	u := new(User)
	if err := database.Get(u, id); err != nil {
		return nil, err
	}
	return u, nil
}
