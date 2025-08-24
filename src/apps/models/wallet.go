package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

type Wallet struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.UUID     `json:"user_id"`
	Address   string        `json:"address"`
	Network   WalletNetwork `json:"network"`
	Testnet   bool          `json:"testnet"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (w *Wallet) Upsert(ctx context.Context) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}

	rows, err := database.Query(
		ctx,
		"wallets/upsert",
		w.ID,
		w.UserID,
		w.Address,
		w.Network,
		w.Testnet,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(w); err != nil {
			return err
		}
	}
	return nil
}
