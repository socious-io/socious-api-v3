package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

type Feedback struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	Content    *string    `db:"content" json:"content"`
	IsContest  *bool      `db:"is_contest" json:"is_contest"`
	IdentityID uuid.UUID  `db:"identity_id" json:"identity_id"`
	ProjectID  uuid.UUID  `db:"project_id" json:"project_id"`
	MissionID  *uuid.UUID `db:"mission_id" json:"mission_id"`
	ContractID *uuid.UUID `db:"contract_id" json:"contract_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (Feedback) TableName() string {
	return "feedbacks"
}

func (Feedback) FetchQuery() string {
	return "feedbacks/fetch"
}

func (c *Feedback) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"feedbacks/create",
		c.Content,
		c.IsContest,
		c.IdentityID,
		c.ProjectID,
		c.MissionID,
		c.ContractID,
	)

	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(c); err != nil {
			return err
		}
	}

	return database.Fetch(c, c.ID)
}
